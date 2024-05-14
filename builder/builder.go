package builder

import (
	"fmt"
	"gowizard/builder/model/system"
	"gowizard/consts"
	"os"
	"os/exec"
	"path/filepath"
)

type Builder struct {
	ProjectName string     `yaml:"project_name"`
	Layers      []LayerDTO `yaml:"layers"`
	Unsafe      bool       `yaml:"unsafe"`
	Path        string     `yaml:"path"`

	Models []*system.Model `yaml:"models"`

	LayerController *LayerController `yaml:"-"`
}

type LayerDTO struct {
	Layer string `yaml:"layer"`
	Tag   string `yaml:"tag"`
}

func (b *Builder) setDefaultsIfEmpty() string {
	if b.ProjectName == "" {
		b.ProjectName = "gowizard"
	}

	if b.Path == "" {
		b.Path = "gowizard"
	}
	if !b.Unsafe {
		b.Path = filepath.Join("magic", b.Path)
	}
	if b.Path[len(b.Path)-1] != '/' {
		b.Path = b.Path + "/"
	}

	if b.Layers == nil {
		b.Layers = []LayerDTO{
			{
				Layer: "controller",
				Tag:   "http",
			},
			{
				Layer: "service",
			},
			{
				Layer: "repository",
				Tag:   "postgres",
			},
		}
	}
	return b.ProjectName
}

func (b *Builder) CodeGenerate() error {
	b.setDefaultsIfEmpty()

	err := b.initStructure()
	if err != nil {
		return fmt.Errorf("unable to generate directories: %w", err)
	}

	err = b.mainGenerate()
	if err != nil {
		return fmt.Errorf("unable to generate main.go: %w", err)
	}

	err = b.LayerController.Generate()
	if err != nil {
		return fmt.Errorf("unable to generate layers: %w", err)
	}

	err = b.goModTidy()
	if err != nil {
		return err
	}

	err = b.gofmt()
	if err != nil {
		return err
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

func (b *Builder) initStructure() error {
	if _, err := createIfNoExist(b.Path); err != nil {
		return fmt.Errorf("unable to create main directory: %w", err)
	}

	b.LayerController = NewLayerController(b, b.Layers, b.Models)

	for i, layer := range b.LayerController.Layers {
		path, err := createIfNoExist(filepath.Join(b.Path, layer.Name))
		if err != nil {
			return fmt.Errorf("unable to create %s directory: %w", layer.Name, err)
		}

		b.LayerController.Layers[i].Path = path
	}

	_, err := createIfNoExist(filepath.Join(b.Path, consts.DefaultModelsFolder))
	if err != nil {
		return fmt.Errorf("unable to create models directory: %w", err)
	}

	_, err = createIfNoExist(filepath.Join(b.Path, consts.DefaultConfigFolder))
	if err != nil {
		return fmt.Errorf("unable to create config directory: %w", err)
	}

	return nil
}

func createIfNoExist(fp string) (string, error) {
	if fp[len(fp)-1] != '/' {
		fp = fp + "/"
	}

	if !checkIfExist(fp) {
		err := os.MkdirAll(fp, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("unable to create directory %s: %w", fp, err)
		}
	}

	return fp, nil
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

func (b *Builder) goModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = b.Path
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to run go mod tidy: %w", err)
	}

	return nil
}

func (b *Builder) gofmt() error {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Dir = b.Path
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to run gofmt: %w", err)
	}

	return nil
}
