package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

type library struct {
	Name        string
	Version     string
	LicenseName string
	LicensePath string
	LicenseURL  string
	LicenseText string
}

func renderReport(metadata []byte, goCommand string, environment []string) ([]byte, error) {
	var libraries []library
	if err := json.Unmarshal(metadata, &libraries); err != nil {
		return nil, fmt.Errorf("failed to decode license metadata: %w", err)
	}
	if len(libraries) == 0 {
		return nil, fmt.Errorf("license collector returned no libraries")
	}
	slices.SortStableFunc(libraries, func(left, right library) int {
		if order := strings.Compare(strings.ToLower(left.Name), strings.ToLower(right.Name)); order != 0 {
			return order
		}
		return strings.Compare(left.Version, right.Version)
	})
	for index := range libraries {
		lib := &libraries[index]
		override, allowed := licenseAllowlist[libraryKey{Name: lib.Name, Version: lib.Version}]
		if allowed && lib.LicenseName == "Unknown" {
			if err := supplementLicense(lib, override, goCommand, environment); err != nil {
				return nil, fmt.Errorf("failed to supplement %s: %w", lib.Name, err)
			}
			fmt.Fprintf(os.Stderr, "[Info] Supplemented %s using the license allowlist.\n", lib.Name)
			continue
		}
		if lib.LicensePath == "" || lib.LicensePath == "Unknown" {
			continue
		}
		text, err := os.ReadFile(lib.LicensePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read license for %s: %w", lib.Name, err)
		}
		lib.LicenseText = string(text)
	}

	report, err := template.New("report").Funcs(template.FuncMap{
		"markdownBlock": markdownBlock,
	}).Parse(string(reportTemplate))
	if err != nil {
		return nil, fmt.Errorf("failed to parse report template: %w", err)
	}
	var output bytes.Buffer
	if err := report.Execute(&output, libraries); err != nil {
		return nil, fmt.Errorf("failed to execute report template: %w", err)
	}
	return output.Bytes(), nil
}

func markdownBlock(text string) string {
	fence := "~~~"
	for strings.Contains(text, fence) {
		fence += "~"
	}
	return fence + "\n" + text + "\n" + fence
}

func supplementLicense(lib *library, override licenseOverride, goCommand string, environment []string) error {
	if len(override.Files) == 0 {
		return fmt.Errorf("license allowlist entry has no documents")
	}
	command := exec.Command(goCommand, "list", "-m", "-json", override.Module)
	command.Env = environment
	command.Stderr = os.Stderr
	metadata, err := command.Output()
	if err != nil {
		return fmt.Errorf("failed to locate module: %w", err)
	}
	var module struct {
		Dir     string
		Version string
	}
	if err := json.Unmarshal(metadata, &module); err != nil {
		return fmt.Errorf("failed to decode module metadata: %w", err)
	}
	if module.Dir == "" || module.Version != lib.Version {
		return fmt.Errorf("module directory or version does not match the collected license metadata")
	}

	var parts []string
	if override.Notice != "" {
		parts = append(parts, override.Notice)
	}
	for _, name := range override.Files {
		text, err := os.ReadFile(filepath.Join(module.Dir, filepath.FromSlash(name)))
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", name, err)
		}
		if len(bytes.TrimSpace(text)) == 0 {
			return fmt.Errorf("upstream license document is empty: %s", name)
		}
		parts = append(parts, name+"\n\n"+string(text))
	}
	lib.LicenseName = override.LicenseName
	lib.LicensePath = filepath.Join(module.Dir, filepath.FromSlash(override.Files[0]))
	lib.LicenseURL = override.LicenseURL
	lib.LicenseText = strings.Join(parts, "\n\n")
	return nil
}
