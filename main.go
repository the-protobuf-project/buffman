// package main is the entry point for the Buffman CLI application.
package main

import "github.com/machanirobotics/buffman/cmd"

// main is the primary entry point for the executable. It calls the Execute
// function from the cmd package, which initializes and runs the Cobra
// command-line interface. This handles parsing command-line arguments and
// flags, and executes the corresponding command logic.
func main() {
	cmd.Execute()
}
