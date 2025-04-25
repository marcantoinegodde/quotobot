package cmd

import (
	"log"
	"os"
	"quotobot/cmd/bot"
	"quotobot/cmd/server"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quotobot",
	Short: "QuotoBot - ViaRézo Telegram quotes bot",
	Long:  "QuotoBot - A Telegram bot for sharing quotes from ViaRézo.",
}

func init() {
	rootCmd.AddCommand(bot.BuildBotCmd())
	rootCmd.AddCommand(server.BuildServerCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
