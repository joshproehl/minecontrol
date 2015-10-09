package main

import (
	"bufio"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/joshproehl/minecontrol-go/mcrcon"
	"github.com/voxelbrain/goptions"
	"os"
	"strings"
)

func main() {
	options := struct {
		Address  string        `goptions:"-a, --address, description='The domain name or IP address to connect to'"`
		Port     int           `goptions:"-p, --port, description='The port to connect to'"`
		Password string        `goptions:"--password, description='Supply a password at the command line rather than be prompted'"`
		Help     goptions.Help `goptions:"-h, --help, description='Show this help output'"`
		goptions.Remainder

		goptions.Verbs
		Command struct {
		} `goptions:"command"`
		Repl struct {
		} `goptions:"repl"`
		Server struct {
			Username string `goptions:"--server-username, description='Require this username to connect to the server'"`
			Password string `goptions:"--server-password, description='Require this password to connect to the server'"`
			Port     int    `goptions:"--server-port, description='Run the server on this port'"`
		} `goptions:"server"`
	}{ // Default values
		Address:  "127.0.0.1",
		Port:     25575, //Have to use this to be able to detect if they passed in an empty string via -p... Don't like.
		Password: "\"\"",
	}
	// Set the nested default values
	options.Server.Port = 7767

	goptions.ParseAndFail(&options)

	// Catch this before the switch so we can get the password if needed EXCEPT for thi scase.
	if options.Verbs == "" {
		goptions.PrintHelp()
		return
	}

	// If they haven't passed in a password, prompt for one.
	if options.Password == "\"\"" {
		fmt.Printf("Enter RCON password: ")
		options.Password = string(gopass.GetPasswd())
	}
	switch options.Verbs {
	case "repl":
		runREPL(options.Address, options.Password)
	case "command":
		runCommand(options.Address, options.Password, strings.Join(options.Remainder, " "))
	case "server":
		mcrcon.NewRestServer(options.Address, options.Password, options.Port, options.Server.Username, options.Server.Password, options.Server.Port)
	}
}

// runREPL takes an address and password, then sest up a connection to the RCON server and presents the user with a
// read-evaluate-print-loop command prompt for the connected RCON server.
func runREPL(address string, password string) {
	client := mcrcon.NewClient(address, password)

	if client.Connected != true {
		fmt.Println("FATAL: Client could not connect")
		return
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

// runCommand takes the options passed in from the command line, prints the output of the command, and terminates.
func runCommand(address string, password string, command string) {

	client := mcrcon.NewClient(address, password)

	if client.Connected != true {
		fmt.Println("FATAL: Client could not connect")
		return
	}

	fmt.Println("Executing command: ", command)

	cmdResponse, rUserErr := client.SendCommand(command)

	if rUserErr != nil {
		fmt.Println("FATAL: ", rUserErr)
	}

	fmt.Println(cmdResponse)

	client.Close()
}
