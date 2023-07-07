package bot

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	bot_r, bot_d *tgbotapi.BotAPI
)

func sendMsg(msg string, chatID int64) {
	NewMsg := tgbotapi.NewMessage(chatID, msg)
	// NewMsg.ParseMode = tgbotapi.ModeHTML   //傳送html格式的訊息
	_, err := bot_d.Send(NewMsg)
	if err == nil {
		log.Printf("Send telegram message success")
	} else {
		log.Printf("Send telegram message error")
	}
}

func replyMsg(chatID **int64, package_count **int, rent **map[string]int) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, _ := bot_r.GetUpdatesChan(updateConfig)
	for update_i := range updates {
		update := update_i
		cmd_text := update.Message.Text
		chatID := update.Message.Chat.ID
		replyMsg := tgbotapi.NewMessage(chatID, cmd_text)
		replyMsg.ReplyToMessageID = update.Message.MessageID
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "pa":
				**package_count = **package_count + 1
				replyMsg.Text = "尚有 " + strconv.Itoa(**package_count) + " 個📦待領取"
			case "pd":
				if **package_count == 0 {
					replyMsg.Text = "尚有 " + strconv.Itoa(**package_count) + " 個📦待領取"
				} else {
					**package_count = **package_count - 1
					replyMsg.Text = "尚有 " + strconv.Itoa(**package_count) + " 個📦待領取"
				}
			case "rd":
				(**rent)["Rent_remind"] = 0
				(**rent)["Rent_done"] = 1
				replyMsg.Text = "✔️房租確認完成, 已關閉提醒"
			case "rs":
				(**rent)["Rent_remind"] = 1
				(**rent)["Rent_done"] = 0
				replyMsg.Text = "✔️房租提醒已開啟"
			case "help":
				replyMsg.Text = "/pa /pd /rd /rs"
			default:
				replyMsg.Text = ""
			}
		} else {
			replyMsg.Text = ""
		}
		if len(replyMsg.Text) > 0 {
			_, _ = bot_r.Send(replyMsg)
		}
	}
}

func Telegram_reply_run(chatID *int64, TOKEN *string, package_count *int, rent *map[string]int) {
	var err error
	bot_r, err = tgbotapi.NewBotAPI(*TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	bot_r.Debug = false

	replyMsg(&chatID, &package_count, &rent)
}

func Telegram_bot_run(chatID *int64, TOKEN *string, service string, value int) {
	var err error
	var msg string

	bot_d, err = tgbotapi.NewBotAPI(*TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	bot_d.Debug = false

	switch service {
	case "package_service":
		msg = "尚有 " + strconv.Itoa(value) + " 個📦待領取"
	case "rent_service":
		msg = "!!!📣📣📣!!! 該繳房租了, 請確認戶頭的錢夠不夠扣款"
	default:
		msg = service
	}

	sendMsg(msg, *chatID)
}
