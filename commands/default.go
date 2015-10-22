// commands defines all the subcommands and command line options available in minecontrol
package commands

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"os"
)

var mcCmd = &cobra.Command{
	Use: "minecontrol",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// This will set up the app

		// Note at this point only WARN or above is actually logged to file, and ERROR or above to console.
		jww.SetLogFile("minecontrol.log")

		if viper.GetBool("verbose") {
			fmt.Println("Enabling verbose output...")
			jww.SetLogThreshold(jww.LevelTrace)
			jww.SetStdoutThreshold(jww.LevelInfo)
		}

		if fvVersion {
			// TODO: Get version numbers dynamically
			fmt.Println(" ")
			fmt.Println("Minecontrol version 0.0.1")
			os.Exit(0)
		}

		if viper.GetString("rcon.password") == "" { // Should detect if we have a password via config or flag, and only execute this if NOT
			fmt.Printf("Enter RCON password: ")
			passwd := string(gopass.GetPasswd())
			viper.Set("rcon.password", passwd)
		}
	},
}

// Flag values
var fvAddress, fvPassword string
var fvPort int
var fvVerbose, fvVersion bool

// GetGoing is what sets up the app, and then runs Execute() on whichever command was called.
func GetGoing() {
	// Note that it's critically important to add commands and flags BEFORE you do the config file
	// stuff, otherwise Viper just silently gives up and doesn't bind the two.
	addCommands()
	addFlags()
	getConfigFile()

	if err := mcCmd.Execute(); err != nil {
		// Cobra already spat out any errors, but...
		jww.FATAL.Println("Command failure:", err)
		os.Exit(1)
	}
}

// Checks for a config file, sets up sensible default options
func getConfigFile() {
	viper.SetConfigName("minecontrol")
	viper.AddConfigPath(".")
	configErr := viper.ReadInConfig()

	if configErr != nil {
		jww.WARN.Println("No config file found, using default values.")
	}

	// Bind config file values to the command line options passed in
	viper.BindPFlag("verbose", mcCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("rcon.address", mcCmd.PersistentFlags().Lookup("address"))
	viper.BindPFlag("rcon.port", mcCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("rcon.password", mcCmd.PersistentFlags().Lookup("password"))
}

func addFlags() {
	mcCmd.PersistentFlags().StringVarP(&fvAddress, "address", "a", "127.0.0.1", "The IP address or domain name of the server to connect to")
	mcCmd.PersistentFlags().IntVarP(&fvPort, "port", "p", 25566, "The port number that minecraft is running on at the provided address")
	mcCmd.PersistentFlags().StringVarP(&fvPassword, "password", "P", "", "The RCON Password needed to connect to the server")
	mcCmd.PersistentFlags().BoolVar(&fvVersion, "version", false, "Print the version number and exit")
	mcCmd.PersistentFlags().BoolVar(&fvVerbose, "verbose", false, "Set verbose mode. (Logs even more to the logfile)")
}

func addCommands() {
	mcCmd.AddCommand(runCmd)
	mcCmd.AddCommand(replCmd)
	mcCmd.AddCommand(serverCmd)
}
