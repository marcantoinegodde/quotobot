package main

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
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

	if err := qb.Database.Where("id = ?", qid).First(&Quote{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Citation introuvable",
			})
			log.Printf("Citation introuvable de %s", update.Message.From.Username)
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Erreur de base de données",
			})
			log.Printf("Erreur de base de données : %v", err)
		}
		return
	}

	vote := Vote{
		PersonID: update.Message.From.ID,
		QuoteID:  uint(qid),
	}

	if err := qb.Database.Create(&vote).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "T'as déjà voté crétin !",
			})
			log.Printf("Vote déjà enregistré de %s", update.Message.From.Username)
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Erreur de base de données",
			})
			log.Printf("Erreur de base de données : %v", err)
		}
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "A voté !",
	})

	log.Printf("Vote enregistré de %s pour la citation %d", update.Message.From.Username, qid)
}
