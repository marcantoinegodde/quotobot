package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (qb *QuotoBot) unvoteHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	if count == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "T'as pas encore voté pour cette citation crétin !",
		})
		log.Printf("Vote déjà enregistré de %s", update.Message.From.Username)
		return
	}

	if err := qb.Database.Model(&Vote{}).Where("person_id = ? AND quote_id = ?", update.Message.From.ID, qid).Delete(&Vote{}).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur lors de la suppression du vote",
		})
		log.Printf("Erreur lors de la suppression du vote de %s", update.Message.From.Username)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "J'ai supprimé ton vote",
	})

	log.Printf("Vote supprimé de %s", update.Message.From.Username)
}
