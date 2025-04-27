package bot

import (
	"context"
	"quotobot/pkg/database"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (qb *QuotoBot) requireAuth(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message.Chat.ID == qb.Config.Bot.ChatID {
			next(ctx, b, update)
			return
		}

		var count int64
		if err := qb.Database.Model(&database.User{}).Where("telegram_id = ?", update.Message.From.ID).Count(&count).Error; err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Erreur lors de la vérification de l'authentification",
			})
			qb.Logger.Error.Printf("Erreur lors de la vérification de l'authentification: %v", err)
			return
		}

		if count == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Vous devez vous enregister avant de pouvoir utiliser cette commande. Utilisez la commande /register pour vous enregistrer.",
			})
			qb.Logger.Info.Printf("Utilisateur non authentifié: %s", update.Message.From.Username)
			return
		}

		next(ctx, b, update)
	}
}
