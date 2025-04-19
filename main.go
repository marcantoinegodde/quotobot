package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

type QuotoBot struct {
	Config   *Config
	Database *gorm.DB
}

func NewQuotoBot() *QuotoBot {
	c := loadConfig()
	db := loadDatabase()

	return &QuotoBot{
		Config:   c,
		Database: db,
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	qb := NewQuotoBot()

	opts := []bot.Option{
		bot.WithDefaultHandler(qb.defaultHandler),
	}

	b, err := bot.New(qb.Config.Token, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "add", bot.MatchTypeCommandStartOnly, qb.addHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "last", bot.MatchTypeCommandStartOnly, qb.lastHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "random", bot.MatchTypeCommandStartOnly, qb.randomHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "search", bot.MatchTypeCommandStartOnly, qb.searchHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "vote", bot.MatchTypeCommandStartOnly, qb.voteHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "unvote", bot.MatchTypeCommandStartOnly, qb.unvoteHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "score", bot.MatchTypeCommandStartOnly, qb.scoreHandler)

	b.Start(ctx)
}

func (qb *QuotoBot) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Ignore message edits, necessary to avoid panics
	if update.Message == nil {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Commande inconnue",
	})
}
