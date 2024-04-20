package responses

import "fmt"

// todo finish
const (
	Help = `SomeHelpCommands`
)

var (
	WrongArgs = fmt.Errorf("wrong set of arguments, run \"gowizard [command] help\" for usage")
)

func PrintErr(str error) {
	fmt.Printf(">> ERROR <<\n%s\n", str.Error())
}

func PrintHelp(str string) {
	fmt.Printf(">> HELP <<\n%s\n", str)
}

func PrintResp(str string) {
	fmt.Printf(">> SUCCESSFUL <<\n%s\n", str)
}
