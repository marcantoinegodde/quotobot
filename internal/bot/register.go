package bot

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (qb *QuotoBot) registerHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Chat.Type != models.ChatTypePrivate {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Pour t'enregistrer, viens me parler en privé.",
		})
		return
	}

	id := update.Message.From.ID
	username := update.Message.From.Username

	params := fmt.Sprintf("id=%d&username=%s", id, username)
	mac := hmac.New(sha256.New, []byte(qb.Config.Bot.HMACSecret))
	mac.Write([]byte(params))

	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	url := fmt.Sprintf("https://%s/register?id=%d&username=%s&signature=%s", qb.Config.Bot.BaseURL, id, username, signature)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Salut ! Pour t'enregistrer, clique sur ce lien : %s", url),
	})

	qb.Logger.Info.Printf("Lien d'enregistrement envoyé à %s", update.Message.From.Username)
}
