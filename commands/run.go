package commands

import (
	"fmt"
	"github.com/joshproehl/minecontrol/mcrcon"
	"github.com/spf13/cobra"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run [command]",
	Short: "Run the provided command on the server, then exit",
	Long:  `Sometimes you don't want a REPL, you just want to run a single command. This is how.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running command: ", args)
	},
}

// runCommand takes the options passed in from the command line, prints the output of the command, and terminates.
func runCommand(address string, port int, password string, command string) {

	client, err := mcrcon.NewClient(address, port, password)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Executing command: ", command)

	cmdResponse, rUserErr := client.SendCommand(command)

	if rUserErr != nil {
		fmt.Println("FATAL: ", rUserErr)
	}

	fmt.Println(cmdResponse)

	client.Close()
}
