package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

type QuotoBot struct {
	Database *gorm.DB
}

func NewQuotoBot() *QuotoBot {
	db := loadDatabase()

	return &QuotoBot{
		Database: db,
	}
}

const (
	CHAT_ID = 123456789
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	qb := NewQuotoBot()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(os.Getenv("TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "add", bot.MatchTypeCommandStartOnly, qb.addHandler)

	b.Start(ctx)
}

func (qb *QuotoBot) addHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Chat.ID == CHAT_ID {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Arrete de faire chier les autres et viens me voir en privé",
		})

		return
	}

	command := strings.Split(update.Message.Text, "|")
	if len(command) != 5 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		log.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	content := strings.TrimSpace(command[1])
	header := strings.TrimSpace(command[3])

	if content == "" || header == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		log.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	quote := Quote{
		Header:  header,
		Content: content,
	}
	if err := qb.Database.Create(&quote).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur lors de l'ajout de la citation",
		})
		log.Printf("Erreur lors de l'ajout de la citation: %v", err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("Voici la citation ajoutée :\n\n*%s*\n\n_by %s_", bot.EscapeMarkdown(quote.Content), bot.EscapeMarkdown(quote.Header)),
		ParseMode: models.ParseModeMarkdown,
	})

	log.Println("Quote ajoutée par", update.Message.From.Username)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Commande inconnue",
	})
}
