package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
)

var generatedFiles = []string{
	"internal/server/router/graph/generated.go",
	"internal/server/router/graph/model/models_gen.go",
	"internal/server/router/graph/schema.resolvers.go",
}

func main() {
	if err := generate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generate() error {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		return fmt.Errorf("failed to load gqlgen config: %w", err)
	}

	backupDir, err := os.MkdirTemp("", "graphqlgen-backup-")
	if err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	defer os.RemoveAll(backupDir)

	if err := copyFiles(".", backupDir, generatedFiles); err != nil {
		return fmt.Errorf("failed to back up generated files: %w", err)
	}

	if err := api.Generate(cfg); err != nil {
		if restoreErr := copyFiles(backupDir, ".", generatedFiles); restoreErr != nil {
			return fmt.Errorf("failed to generate GraphQL code: %w (also failed to restore generated files: %v)", err, restoreErr)
		}
		return fmt.Errorf("failed to generate GraphQL code: %w", err)
	}

	return nil
}

func copyFiles(sourceDir, destinationDir string, files []string) error {
	for _, name := range files {
		source := filepath.Join(sourceDir, name)
		destination := filepath.Join(destinationDir, name)

		info, err := os.Stat(source)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		content, err := os.ReadFile(source)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(destination, content, info.Mode().Perm()); err != nil {
			return err
		}
	}
	return nil
}
