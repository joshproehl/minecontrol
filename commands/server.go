package commands

import (
	"github.com/joshproehl/minecontrol/mcrcon/restServer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Create an HTTP server for the REST API and GUI",
	Long: `Create an HTTP server which will provide a JSON API to the connected Minecraft server.
By default the server will be available at http://127.0.0.0.1:7767`,
	Run: func(cmd *cobra.Command, args []string) {
		c := restServer.ServerConfig{
			RCON_address:  viper.GetString("rcon.address"),
			RCON_port:     viper.GetInt("rcon.port"),
			RCON_password: viper.GetString("rcon.password"),
			Username:      viper.GetString("server.username"),
			Password:      viper.GetString("server.password"),
			Port:          viper.GetInt("server.port"),
		}

		restServer.NewRestServer(&c)
	},
}

func init() {
	serverCmd.Flags().Int("serverPort", 7767, "Port to run the REST server on")
	serverCmd.Flags().String("serverUsername", "", "HTTP Basic auth username that the REST server will require")
	serverCmd.Flags().String("serverPassword", "", "HTTP Basic auth password that the REST server will require")
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("serverPort"))
	viper.BindPFlag("server.username", serverCmd.Flags().Lookup("serverUsername"))
	viper.BindPFlag("server.password", serverCmd.Flags().Lookup("serverPassword"))
}
