package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

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

	b.Start(ctx)
}

func (qb *QuotoBot) addHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Chat.ID == qb.Config.ChatID {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Arrête de faire chier les autres et viens me voir en privé",
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
	author := strings.TrimSpace(command[3])

	if content == "" || author == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Mauvais format",
		})
		log.Printf("Mauvais format de %s", update.Message.From.Username)
		return
	}

	quote := Quote{
		Content: content,
		Author:  author,
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
		Text:      fmt.Sprintf("Voici la citation ajoutée :\n\n*%s*\n\n_by %s_", bot.EscapeMarkdown(quote.Content), bot.EscapeMarkdown(quote.Author)),
		ParseMode: models.ParseModeMarkdown,
	})

	log.Println("Quote ajoutée par", update.Message.From.Username)
}

func (qb *QuotoBot) lastHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	n := 1

	command := strings.Split(update.Message.Text, " ")
	if len(command) > 1 {
		num, err := strconv.Atoi(command[1])
		if err != nil || num < 1 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Mauvais format",
			})
			log.Printf("Mauvais format de %s", update.Message.From.Username)
			return
		}
		n = num
	}

	const maxQuotes = 10
	if n > maxQuotes {
		n = maxQuotes
	}

	var quotes []Quote
	if err := qb.Database.Order("id desc").Limit(n).Find(&quotes).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur lors de la récupération des citations",
		})
		log.Printf("Erreur lors de la récupération des citations: %v", err)
		return
	}

	if len(quotes) == 0 {
		return
	}

	separator := "\n" + strings.Repeat("_", 20) + "\n\n"
	var formattedQuotes []string

	for _, quote := range quotes {
		formattedQuote := fmt.Sprintf("\\#Q%d _\\(\\+%d\\)_\n*%s*\n\n_by %s_",
			quote.ID, len(quote.Votes), bot.EscapeMarkdown(quote.Content), bot.EscapeMarkdown(quote.Author))
		formattedQuotes = append(formattedQuotes, formattedQuote)
	}

	response := strings.Join(formattedQuotes, bot.EscapeMarkdown(separator))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      response,
		ParseMode: models.ParseModeMarkdown,
	})

	log.Printf("%d quotes envoyées à %s\n", len(quotes), update.Message.From.Username)
}

func (qb *QuotoBot) randomHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	n := 1

	command := strings.Split(update.Message.Text, " ")
	if len(command) > 1 {
		num, err := strconv.Atoi(command[1])
		if err != nil || num < 1 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Mauvais format",
			})
			log.Printf("Mauvais format de %s", update.Message.From.Username)
			return
		}
		n = num
	}

	const maxQuotes = 10
	if n > maxQuotes {
		n = maxQuotes
	}

	var quotes []Quote
	if err := qb.Database.Order("RANDOM()").Limit(n).Find(&quotes).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur lors de la récupération des citations",
		})
		log.Printf("Erreur lors de la récupération des citations: %v", err)
		return
	}

	if len(quotes) == 0 {
		return
	}

	separator := "\n" + strings.Repeat("_", 20) + "\n\n"
	var formattedQuotes []string

	for _, quote := range quotes {
		formattedQuote := fmt.Sprintf("\\#Q%d _\\(\\+%d\\)_\n*%s*\n\n_by %s_",
			quote.ID, len(quote.Votes), bot.EscapeMarkdown(quote.Content), bot.EscapeMarkdown(quote.Author))
		formattedQuotes = append(formattedQuotes, formattedQuote)
	}

	response := strings.Join(formattedQuotes, bot.EscapeMarkdown(separator))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      response,
		ParseMode: models.ParseModeMarkdown,
	})

	log.Printf("%d quotes envoyées à %s\n", len(quotes), update.Message.From.Username)
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
