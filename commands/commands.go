package commands

type Command interface {
	Run() (string, error)
	GetHelp() string
}
