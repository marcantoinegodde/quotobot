package bot

import (
	"context"
	"errors"
	"fmt"
	"quotobot/pkg/database"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

func (qb *QuotoBot) getHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	command := strings.Split(update.Message.Text, " ")

	if len(command) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		qb.Logger.Info.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	qid, err := strconv.Atoi(command[1])
	if err != nil || qid < 1 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		qb.Logger.Info.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	var quote database.Quote
	if err := qb.Database.Preload("Votes").Where("id = ?", qid).First(&quote).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Citation introuvable",
			})
			qb.Logger.Info.Printf("Citation introuvable de %s", update.Message.From.Username)
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Erreur lors de la récupération de la citation",
			})
			qb.Logger.Error.Printf("Erreur lors de la récupération de la citation: %v", err)
		}
		return
	}

	formattedQuote := fmt.Sprintf("\\#Q%d _\\(\\+%d\\)_\n*%s*\n\n_by %s_",
		quote.ID, len(quote.Votes), bot.EscapeMarkdown(quote.Content), bot.EscapeMarkdown(quote.Author))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      formattedQuote,
		ParseMode: models.ParseModeMarkdown,
	})
	qb.Logger.Info.Printf("Citation #%d envoyée à %s", quote.ID, update.Message.From.Username)
}
