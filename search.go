package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (qb *QuotoBot) searchHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	command := strings.Split(update.Message.Text, " ")

	if len(command) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		log.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	search := strings.TrimSpace(command[1])
	if search == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		log.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	n := 1

	if len(command) > 2 {
		num, err := strconv.Atoi(command[2])
		if err != nil || num < 1 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Mauvais format",
			})
			log.Printf("Mauvais format de %s", update.Message.From.Username)
			return
		}
		n = num
	}

	const maxQuotes = 10
	if n > maxQuotes {
		n = maxQuotes
	}

	var quotes []Quote
	if err := qb.Database.Where("content LIKE ?", "%"+search+"%").Or("author LIKE ?", "%"+search+"%").Order("RANDOM()").Limit(n).Find(&quotes).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur lors de la récupération des citations",
		})
		log.Printf("Erreur lors de la récupération des citations: %v", err)
		return
	}

	if len(quotes) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Aucune citation trouvée",
		})
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

	log.Printf("%d quotes envoyées à %s\n", len(quotes), update.Message.From.Username)
}
