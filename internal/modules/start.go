package modules

import (
	"time"

	"github.com/Laky-64/gologging"
	tg "github.com/amarnathcjd/gogram/telegram"

	"main/internal/config"
	"main/internal/core"
	"main/internal/database"
	"main/internal/locales"
	"main/internal/utils"
)

func init() {
	helpTexts["/start"] = `<i>Inicia o bot e mostra o menu principal.</i>`
}

func startHandler(m *tg.NewMessage) error {
	if m.ChatType() != tg.EntityUser {
		database.AddServedChat(m.ChannelID())
		m.Reply(
			F(m.ChannelID(), "start_group"),
		)
		return tg.ErrEndGroup
	}

	arg := m.Args()
	database.AddServedUser(m.ChannelID())
	if arg != "" {
		gologging.Info(
			"Got Start parameter: " + arg + " in ChatID: " + utils.IntToStr(
				m.ChannelID(),
			),
		)
	}

	switch arg {
	case "pm_help":
		gologging.Info("User requested help via start param")
		helpHandler(m)
	default:
		caption := F(m.ChannelID(), "start_private", locales.Arg{
			"user": utils.MentionHTML(m.Sender),
			"bot":  utils.MentionHTML(m.Client.Me()),
		})
		
		// Using SendOptions with Spoiler field which is globally supported across wrapper versions
		sendOpt := &tg.SendOptions{
			Caption:     caption,
			NoForwards:  true,
			ReplyMarkup: core.GetStartMarkup(m.ChannelID()),
			Media:       config.StartImage,
			Spoiler:     true, // This hides the image with a spoiler mesh
		}

		_, err := m.Respond(caption, sendOpt)
		if err != nil {
			gologging.Error(
				"[start] Media send with spoiler failed: " + err.Error(),
			)
			// Fallback text only if media completely breaks
			_, err = m.Respond(caption, &tg.SendOptions{
				NoForwards:  true,
				ReplyMarkup: core.GetStartMarkup(m.ChannelID()),
			})
			return err
		}
	}

	if config.LoggerID != 0 && isLoggerEnabled() {
		uName := "N/A"
		if m.Sender.Username != "" {
			uName = "@" + m.Sender.Username
		}
		msg := F(m.ChannelID(), "logger_bot_started", locales.Arg{
			"mention":       utils.MentionHTML(m.Sender),
			"user_id":       m.SenderID(),
			"user_username": uName,
		})
		_, err := m.Client.SendMessage(config.LoggerID, msg)
		if err != nil {
			gologging.Error(
				"Failed to send logger_bot_started msg, Err: " + err.Error(),
			)
		}
	}

	return tg.ErrEndGroup
}

func startCB(cb *tg.CallbackQuery) error {
	cb.Answer("")
	caption := F(cb.ChannelID(), "start_private", locales.Arg{
		"user": utils.MentionHTML(cb.Sender),
		"bot":  utils.MentionHTML(cb.Client.Me()),
	})
	sendOpt := &tg.SendOptions{
		ReplyMarkup: core.GetStartMarkup(cb.ChannelID()),
		NoForwards:  true,
	}
	if config.StartImage != "" {
		sendOpt.Media = config.StartImage
	}
	cb.Edit(caption, sendOpt)
	return tg.ErrEndGroup
}

func aboutCB(cb *tg.CallbackQuery) error {
	cb.Answer("")

	uptime := time.Since(config.StartTime).Round(time.Second)

	caption := F(cb.ChannelID(), "about_text", locales.Arg{
		"bot":    utils.MentionHTML(cb.Client.Me()),
		"uptime": uptime.String(),
	})

	sendOpt := &tg.SendOptions{
		ReplyMarkup: core.GetBackToStartKeyboard(cb.ChannelID()),
		NoForwards:  true,
	}
	if config.StartImage != "" {
		sendOpt.Media = config.StartImage
	}
	cb.Edit(caption, sendOpt)
	return tg.ErrEndGroup
}
