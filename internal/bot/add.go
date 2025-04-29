package bot

import (
	"context"
	"fmt"
	"quotobot/pkg/database"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (qb *QuotoBot) addHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Chat.Type != models.ChatTypePrivate {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Arrête de faire chier les autres et viens me voir en privé",
		})

		return
	}

	command := strings.Split(update.Message.Text, "|")
	if len(command) != 5 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		qb.Logger.Info.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	content := strings.TrimSpace(command[1])
	author := strings.TrimSpace(command[3])

	if content == "" || author == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		qb.Logger.Info.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	quote := database.Quote{
		Content: content,
		Author:  author,
	}
	if err := qb.Database.Create(&quote).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur lors de l'ajout de la citation",
		})
		qb.Logger.Error.Printf("Erreur lors de l'ajout de la citation: %v", err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("Voici la citation ajoutée :\n\n*%s*\n\n_by %s_", bot.EscapeMarkdown(quote.Content), bot.EscapeMarkdown(quote.Author)),
		ParseMode: models.ParseModeMarkdown,
	})

	qb.Logger.Info.Printf("Quote #%d ajoutée par %s", quote.ID, update.Message.From.Username)
}
