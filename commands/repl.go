package commands

import (
	"bufio"
	"fmt"
	"github.com/joshproehl/minecontrol/mcrcon"
	"github.com/spf13/cobra"
	"os"
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Open a REPL for your Minecraft server",
	Long: `Read
	Evaluate
	Print
	Loop`,
	Run: func(cmd *cobra.Command, args []string) {
		//runREPL()
	},
}

// runREPL takes an address and password, then sest up a connection to the RCON server and presents the user with a
// read-evaluate-print-loop command prompt for the connected RCON server.
func runREPL(address string, port int, password string) {
	client, err := mcrcon.NewClient(address, port, password)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Type \"exit\" to quit")

	inputreader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := inputreader.ReadString('\n') // this will prompt the user for input

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if input == "exit\n" {
			client.Close()
			return
		}

		cmdResponse, rUserErr := client.SendCommand(input)

		if rUserErr != nil {
			fmt.Println("FATAL: ", rUserErr)
		}

		fmt.Println(cmdResponse)
	}

}
