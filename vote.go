package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (qb *QuotoBot) voteHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	command := strings.Split(update.Message.Text, " ")

	if len(command) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		log.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	qid, err := strconv.Atoi(command[1])
	if err != nil || qid < 1 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		log.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	var count int64
	if err := qb.Database.Model(&Vote{}).Where("person_id = ? AND quote_id = ?", update.Message.From.ID, qid).Count(&count).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur de base de données",
		})
		log.Printf("Erreur de base de données : %v", err)
		return
	}

	if count > 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "T'as déjà voté crétin !",
		})
		log.Printf("Vote déjà enregistré de %s", update.Message.From.Username)
		return
	}

	vote := Vote{
		PersonID: update.Message.From.ID,
		QuoteID:  uint(qid),
	}

	if err := qb.Database.Create(&vote).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur de base de données",
		})
		log.Printf("Erreur de base de données : %v", err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "A voté !",
	})

	log.Printf("Vote enregistré de %s pour la citation %d", update.Message.From.Username, qid)
}
