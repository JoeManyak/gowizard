package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Builder struct {
	ProjectName string            `yaml:"project_name"`
	Layers      map[string]string `yaml:"layers"`
	Path        string            `yaml:"-"`
}

func (b *Builder) setDefaultsIfEmpty() string {
	if b.Path == "" {
		b.Path = "./gowizard/"
	}
	if b.Path[len(b.Path)-1] != '/' {
		b.Path = b.Path + "/"
	}

	if b.ProjectName == "" {
		b.ProjectName = "gowizard"
	}

	if b.Layers == nil {
		b.Layers = map[string]string{
			"handler":    "http",
			"service":    "",
			"repository": "postgres",
		}

	}
	return b.ProjectName
}

func (b *Builder) CodeGenerate() error {
	b.setDefaultsIfEmpty()

	err := b.generateDirectories()
	if err != nil {
		return fmt.Errorf("unable to generate directories: %w", err)
	}

	err = b.mainGenerate()
	if err != nil {
		return fmt.Errorf("unable to generate main.go: %w", err)
	}

	return nil
}

const templateMain = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}
`

func (b *Builder) generateDirectories() error {
	if !checkIfExist(b.Path) {
		err := os.Mkdir(b.Path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to create main directory: %w", err)
		}
	}

	return nil
}

func checkIfExist(fp string) bool {
	dir := filepath.Dir(fp)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}

	return true
}

func (b *Builder) mainGenerate() error {
	f, err := os.Create(filepath.Join(b.Path, "main.go"))
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}

	_, err = f.Write([]byte(templateMain))
	if err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("unable to close file: %w", err)
	}

	// If there is no go.mod file or go.sum file, there will be an expected error, so we ignore it
	_ = os.Remove(filepath.Join(b.Path, "go.mod"))
	_ = os.Remove(filepath.Join(b.Path, "go.sum"))

	cmd := exec.Command("go", "mod", "init", b.ProjectName)
	cmd.Dir = b.Path
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to run go mod init: %w", err)
	}

	return nil
}
