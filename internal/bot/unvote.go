package bot

import (
	"context"
	"errors"
	"quotobot/pkg/database"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

func (qb *QuotoBot) unvoteHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	if err := qb.Database.Where("id = ?", qid).First(&database.Quote{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Citation introuvable",
			})
			qb.Logger.Info.Printf("Citation introuvable de %s", update.Message.From.Username)
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Erreur lors de la recherche de la citation",
			})
			qb.Logger.Error.Printf("Erreur lors de la recherche de la citation: %v", err)
		}
		return
	}

	query := qb.Database.Where("person_id = ? AND quote_id = ?", update.Message.From.ID, qid).Delete(&database.Vote{})
	if err := query.Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur lors de la suppression du vote",
		})
		qb.Logger.Error.Printf("Erreur lors de la suppression du vote: %v", err)
		return
	}

	if query.RowsAffected == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "T'as pas encore voté pour cette citation crétin !",
		})
		qb.Logger.Info.Printf("Vote introuvable de %s", update.Message.From.Username)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "J'ai supprimé ton vote",
	})

	qb.Logger.Info.Printf("Vote supprimé de %s pour la citation #%d", update.Message.From.Username, qid)
}
