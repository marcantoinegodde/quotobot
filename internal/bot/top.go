package bot

import (
	"context"
	"fmt"
	"quotobot/pkg/database"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (qb *QuotoBot) topHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	command := strings.Split(update.Message.Text, " ")
	n := 1

	if len(command) > 1 {
		num, err := strconv.Atoi(command[1])
		if err != nil || num < 1 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Mauvais format",
			})
			qb.Logger.Info.Printf("Mauvais format de %s", update.Message.From.Username)
			return
		}
		n = num
	}

	const maxQuotes = 10
	if n > maxQuotes {
		n = maxQuotes
	}

	var quotes []database.Quote
	if err := qb.Database.Model(&database.Vote{}).Preload("Votes").
		Select("quotes.content, quotes.author, quotes.id, COUNT(votes.person_id) AS votes").
		Joins("JOIN quotes ON quotes.id = votes.quote_id").
		Group("quotes.id").
		Order("votes DESC").
		Limit(n).
		Find(&quotes).Error; err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erreur lors de la récupération des citations",
		})
		qb.Logger.Error.Printf("Erreur lors de la récupération des citations: %v", err)
		return
	}

	if len(quotes) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Aucune citation trouvée",
		})
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

	qb.Logger.Info.Printf("%d quote(s) envoyée(s) à %s", len(quotes), update.Message.From.Username)
}
