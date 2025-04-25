package server

import "github.com/spf13/cobra"

func BuildServerCmd() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Start the web server",
		Long:  `Start the web server, which serves the registration app.`,
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}

	return serverCmd
}
