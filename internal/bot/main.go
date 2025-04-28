package bot

import (
	"context"
	"os"
	"os/signal"
	"quotobot/pkg/config"
	"quotobot/pkg/database"
	"quotobot/pkg/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

type QuotoBot struct {
	Logger   *logger.Logger
	Config   *config.Config
	Database *gorm.DB
}

func NewQuotoBot() *QuotoBot {
	l := logger.NewLogger()
	c := config.LoadConfig(l)
	db := database.LoadDatabase(l)

	return &QuotoBot{
		Logger:   l,
		Config:   c,
		Database: db,
	}
}

func (qb *QuotoBot) Start() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(qb.defaultHandler),
	}

	b, err := bot.New(qb.Config.Bot.Token, opts...)
	if err != nil {
		panic(err)
	}

	botUser, err := b.GetMe(ctx)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "register", bot.MatchTypeCommandStartOnly, qb.registerHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "add", bot.MatchTypeCommandStartOnly, qb.addHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "last", bot.MatchTypeCommandStartOnly, qb.lastHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "random", bot.MatchTypeCommandStartOnly, qb.randomHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "get", bot.MatchTypeCommandStartOnly, qb.getHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "search", bot.MatchTypeCommandStartOnly, qb.searchHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "vote", bot.MatchTypeCommandStartOnly, qb.voteHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "unvote", bot.MatchTypeCommandStartOnly, qb.unvoteHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "score", bot.MatchTypeCommandStartOnly, qb.scoreHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "top", bot.MatchTypeCommandStartOnly, qb.topHandler, qb.requireAuth)

	// Register commands with group-like syntax in the absence of a better solution
	b.RegisterHandler(bot.HandlerTypeMessageText, "register@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.registerHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "add@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.addHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "last@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.lastHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "random@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.randomHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "get@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.getHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "search@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.searchHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "vote@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.voteHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "unvote@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.unvoteHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "score@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.scoreHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "top@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.topHandler, qb.requireAuth)
	b.RegisterHandler(bot.HandlerTypeMessageText, "top@"+botUser.Username, bot.MatchTypeCommandStartOnly, qb.topHandler, qb.requireAuth)

	qb.Logger.Info.Println("QuotoBot started")

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
