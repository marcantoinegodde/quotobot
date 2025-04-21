package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (qb *QuotoBot) scoreHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	var quote Quote
	if err := qb.Database.Preload("Votes").Where("id = ?", qid).First(&quote).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Citation introuvable",
		})
		qb.Logger.Info.Printf("Citation introuvable pour %s", update.Message.From.Username)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Il y a %d vote(s) pour la citation n° %d", len(quote.Votes), qid),
	})

	qb.Logger.Info.Println("Score envoyé à", update.Message.From.Username)
}
