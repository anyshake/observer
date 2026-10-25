package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	licenseToolModule = "github.com/google/go-licenses/v2@v2.0.1"
	reportName        = "third-party-go-licenses.md"
)

//go:embed report.tpl
var reportTemplate []byte

//go:embed collector.tpl
var collectorTemplate []byte

type options struct {
	goCommand  string
	targetOS   string
	targetArch string
	targetARM  string
	targetMIPS string
}

func main() {
	var opts options
	flag.StringVar(&opts.goCommand, "go", "go", "Go executable used to install the license collector")
	flag.StringVar(&opts.targetOS, "goos", os.Getenv("GOOS"), "Target operating system (GOOS)")
	flag.StringVar(&opts.targetArch, "goarch", os.Getenv("GOARCH"), "Target architecture (GOARCH)")
	flag.StringVar(&opts.targetARM, "goarm", os.Getenv("GOARM"), "Target ARM version (GOARM)")
	flag.StringVar(&opts.targetMIPS, "gomips", os.Getenv("GOMIPS"), "Target MIPS floating-point mode (GOMIPS)")
	flag.Parse()

	if err := generate(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generate(opts options) error {
	distDir := filepath.Join("web", "dist")
	info, err := os.Stat(distDir)
	if err != nil {
		return fmt.Errorf("build the web frontend before generating Go licenses: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("web distribution path is not a directory: %s", distDir)
	}

	workDir, err := os.MkdirTemp("", "observer-licensegen-")
	if err != nil {
		return fmt.Errorf("failed to create license collector directory: %w", err)
	}
	defer os.RemoveAll(workDir)

	// Build the collector for the host, independently of the application's target.
	install := exec.Command(opts.goCommand, "install", licenseToolModule)
	install.Env = append(os.Environ(),
		"GOBIN="+workDir,
		"GOOS="+runtime.GOOS,
		"GOARCH="+runtime.GOARCH,
		"GOARM=",
		"GOMIPS=",
		"CGO_ENABLED=0",
	)
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr
	if err := install.Run(); err != nil {
		return fmt.Errorf("failed to install license collector: %w", err)
	}

	templatePath := filepath.Join(workDir, "report.tpl")
	if err := os.WriteFile(templatePath, collectorTemplate, 0o644); err != nil {
		return fmt.Errorf("failed to write license template: %w", err)
	}

	toolName := "go-licenses"
	if runtime.GOOS == "windows" {
		toolName += ".exe"
	}
	collect := exec.Command(filepath.Join(workDir, toolName), "report", "./cmd/observer",
		"--ignore", "github.com/anyshake/observer", "--template", templatePath)
	collect.Env = append(os.Environ(),
		"GOOS="+opts.targetOS,
		"GOARCH="+opts.targetArch,
		"GOARM="+opts.targetARM,
		"GOMIPS="+opts.targetMIPS,
		"CGO_ENABLED=0",
	)
	collect.Stderr = os.Stderr
	metadata, err := collect.Output()
	if err != nil {
		return fmt.Errorf("failed to collect Go licenses: %w", err)
	}
	if len(bytes.TrimSpace(metadata)) == 0 {
		return fmt.Errorf("license collector returned an empty report")
	}
	content, err := renderReport(metadata, opts.goCommand, collect.Env)
	if err != nil {
		return fmt.Errorf("failed to render Go licenses: %w", err)
	}

	for _, directory := range []string{distDir, filepath.Join("web", "src", "public")} {
		if err := writeReport(filepath.Join(directory, reportName), content); err != nil {
			return fmt.Errorf("failed to write license report in %s: %w", directory, err)
		}
	}
	fmt.Fprintln(os.Stdout, "[Info] Go license reports generated in web/dist and web/src/public.")
	return nil
}

func writeReport(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".licensegen-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return err
	}
	if err := file.Chmod(0o644); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}
