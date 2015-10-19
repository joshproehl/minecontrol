// minecontrol is an application to interact with a minecraft server via the command line or HTTP
package main

import (
	"github.com/joshproehl/minecontrol/commands"
)

func main() {
	// commands/default.go is where the commands are set up.
	commands.GetGoing()
}
