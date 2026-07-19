package modules

import (
	"math/rand"
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
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		// 1. Trigger a Random Auto-Reaction
		reactions := []string{"❤️", "🔥", "⚡", "😍", "🎉", "🥰", "✨", "🦋", "🌸", "💋", "💖"}
		randomReaction := reactions[r.Intn(len(reactions))]
		_ = m.React(randomReaction)

		// 2. "hi honnnney type write" Frame-by-Frame Typing Animation
		animText := []string{
			"❁",
			"❁ ℎ",
			"❁ ℎ𝑖",
			"❁ ℎ𝑖 ℎ",
			"❁ ℎ𝑖 ℎ||𝑜𝑛",
			"❁ ℎ𝑖 ℎ||𝑜𝑛𝑒",
			"❁ ℎ𝑖 ℎ||𝑜𝑛𝑒ee",
			"❁ ℎ𝑖 ℎ||𝑜𝑛𝑒eeʏ",
			"❁ ℎ𝑖 ℎ||𝑜𝑛𝑒eeʏ ❁",
			"✨ ℎ𝑖 ℎ||i ℎ||𝑜𝑛𝑒eeʏ ✨",
			"🌸 ℎ𝑖 ℎ||𝑜𝑛𝑒eeʏ 🌸",
		}

		// Send initial typing frame
		animMsg, err := m.Respond(animText[0], &tg.SendOptions{})
		if err == nil {
			// Loop frames with slight delay
			for _, text := range animText[1:] {
				time.Sleep(40 * time.Millisecond)
				_, _ = animMsg.Edit(text, &tg.SendOptions{})
			}
			// Clean up animation frame
			time.Sleep(100 * time.Millisecond)
			_, _ = animMsg.Delete()
		}

		// 3. Select and Send Random Sticker (Fixed to use ReplySticker)
		stickers := []string{
			"CAACAgUAAxkBAAERSZ5qFrUovFMtksurKhQTv45yVUrOfQAC8x0AAui3IVY8DSpAuqVR7jsE",
			"CAACAgIAAxkBAAERSaBqFrWCmOjc6nrqWKMTiZE0FpFXjwACup8AArLXgUgE5umHBy9ewzsE",
			"CAACAgIAAxkBAAERSaJqFrXAyyTxU1YAAS36RxVxwHO921AAAgUuAAIQlDlKtx2PXUs8Y307BA",
			"CAACAgUAAxkBAAERSaRqFrXlFmoKB9fr-yrgb-4XzqmDtwACfwgAAnVOGFavmNzCq-QSnjsE",
			"CAACAgIAAxkBAAERSaZqFrYzK553Zc_hl86IYI5UiBhPvgAC-3EAAhdc2UoveorfYh-18zsE",
		}
		randomSticker := stickers[r.Intn(len(stickers))]
		
		// Use ReplySticker or SendSticker so Telegram renders the graphic instead of text
		_, _ = m.ReplySticker(randomSticker)
		time.Sleep(300 * time.Millisecond)

		// 4. Send Main Menu Layout with Spoiler Image
		caption := F(m.ChannelID(), "start_private", locales.Arg{
			"user": utils.MentionHTML(m.Sender),
			"bot":  utils.MentionHTML(m.Client.Me()),
		})
		
		sendOpt := &tg.SendOptions{
			Caption:     caption,
			NoForwards:  true,
			ReplyMarkup: core.GetStartMarkup(m.ChannelID()),
			Media:       config.StartImage,
			Spoiler:     true, 
		}

		_, err = m.Respond(caption, sendOpt)
		if err != nil {
			gologging.Error(
				"[start] Media send with spoiler failed: " + err.Error(),
			)
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
