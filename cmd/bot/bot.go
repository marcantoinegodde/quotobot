package bot

import "github.com/spf13/cobra"

func BuildBotCmd() *cobra.Command {
	botCmd := &cobra.Command{
		Use:   "bot",
		Short: "Start the bot",
		Long:  `Start the telegram bot.`,
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}

	return botCmd
}
