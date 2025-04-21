package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (qb *QuotoBot) randomHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	command := strings.Split(update.Message.Text, " ")
	n := 1

	if len(command) > 1 {
		num, err := strconv.Atoi(command[1])
		if err != nil || num < 1 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Mauvais format",
			})
			qb.Logger.Info.Printf("Mauvais format de %s", update.Message.From.Username)
			return
		}
		n = num
	}

	const maxQuotes = 10
	if n > maxQuotes {
		n = maxQuotes
	}

	var quotes []Quote
	if err := qb.Database.Model(&Quote{}).Preload("Votes").Order("RANDOM()").Limit(n).Find(&quotes).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur lors de la récupération des citations",
		})
		qb.Logger.Error.Printf("Erreur lors de la récupération des citations: %v", err)
		return
	}

	if len(quotes) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Aucune citation trouvée",
		})
		qb.Logger.Info.Printf("Aucune citation trouvée pour %s", update.Message.From.Username)
		return
	}

	separator := "\n" + strings.Repeat("_", 20) + "\n\n"
	var formattedQuotes []string

	for _, quote := range quotes {
		formattedQuote := fmt.Sprintf("\\#Q%d _\\(\\+%d\\)_\n*%s*\n\n_by %s_",
			quote.ID, len(quote.Votes), bot.EscapeMarkdown(quote.Content), bot.EscapeMarkdown(quote.Author))
		formattedQuotes = append(formattedQuotes, formattedQuote)
	}

	response := strings.Join(formattedQuotes, bot.EscapeMarkdown(separator))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      response,
		ParseMode: models.ParseModeMarkdown,
	})

	qb.Logger.Info.Printf("%d quote(s) envoyée(s) à %s", len(quotes), update.Message.From.Username)
}
