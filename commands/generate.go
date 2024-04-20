package commands

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"gowizard/builder"
	"gowizard/responses"
	"io"
	"os"
)

type GenerateCommand struct {
	filepath string
}

var _ Command = &GenerateCommand{}

func NewGenerateCommand(args []string) (Command, error) {
	command := &GenerateCommand{}
	if len(args) < 2 {
		return command, responses.WrongArgs
	}

	command.filepath = args[1]

	return command, nil
}

func (cmd *GenerateCommand) Run() (string, error) {
	cfgFile, err := os.Open(cmd.filepath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}
	defer cfgFile.Close()

	cfgRaw, err := io.ReadAll(cfgFile)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}

	var b builder.Builder
	err = yaml.Unmarshal(cfgRaw, &b)
	if err != nil {
		return "", fmt.Errorf("could not parse file: %w", err)
	}

	r, _ := json.Marshal(b)

	err = b.CodeGenerate()
	if err != nil {
		return "", fmt.Errorf("could not generate code: %w", err)
	}

	return string(r), nil
}

func (cmd *GenerateCommand) GetHelp() string {
	return `Usage: gowizard generate <filepath>`
}
