package router

import (
	"fmt"
	"gowizard/commands"
	"gowizard/responses"
	"strings"
)

func Run(args []string) {
	// Handle args count
	if len(args) == 0 {
		fmt.Println(responses.Help)
		return
	}

	var response = ""
	switch strings.ToLower(args[0]) {
	case g, gen, generate:
		response = handleCommand(args, commands.NewGenerateCommand)
		break

	case help:
		responses.PrintHelp(responses.Help)
		return
	default:
		fmt.Println("Invalid command", args)
		responses.PrintHelp(responses.Help)
		return
	}

	if response != "" {
		responses.PrintResp(response)
	}
}

// Command list
const (
	help = "help"

	// Generate command
	g        = "g"
	gen      = "gen"
	generate = "generate"
)

func handleCommand(args []string, commandGet func(args []string) (commands.Command, error)) string {
	command, err := commandGet(args)
	if err != nil {
		responses.PrintErr(fmt.Errorf("invalid command usage: %w", err))
		return ""
	}

	if args[len(args)-1] == help {
		responses.PrintHelp(command.GetHelp())
		return ""
	}

	response, err := command.Run()
	if err != nil {
		responses.PrintErr(fmt.Errorf("unable to run command: %w", err))
		return ""
	}

	return response
}
