package main

import (
	"fmt"
	"gowizard/router"
	"os"
)

func main() {
	// Should ignore the first argument that contains the program name
	if len(os.Args) == 0 {
		fmt.Println("No arguments provided")
		return
	}

	router.Run(os.Args[1:])
}
