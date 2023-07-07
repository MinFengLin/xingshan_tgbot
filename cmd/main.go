package main

import (
	"fmt"
	"strconv"
	"time"

	/*
	 * env
	 */
	"log"
	"os"

	bot_service "github.com/MinFengLin/xingshan_tgbot/bot"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	yourToken := os.Getenv("YOUR_TOKEN")
	crontab_time_package := os.Getenv("Crontime_Package")
	crontab_time_rent := os.Getenv("Crontime_Rent")

	Package_count := 0
	Rent := map[string]int{
		"Rent_remind": 1,
		"Rent_done":   0,
	}

	go bot_service.Telegram_reply_run(&chatID, &yourToken, &Package_count, &Rent)

	fmt.Printf("chatID: %d, yourToken: %s, Use crontab: %s \n", chatID, yourToken, crontab_time_package)
	fmt.Printf("chatID: %d, yourToken: %s, Use crontab: %s \n", chatID, yourToken, crontab_time_rent)
	p := cron.New()
	r := cron.New()
	go func(pct *int) {
		_, _ = p.AddFunc(crontab_time_package, func() {
			if *pct > 0 {
				bot_service.Telegram_bot_run(&chatID, &yourToken, "package_service", *pct)
			}
		})
	}(&Package_count)

	go func(rent *map[string]int) {
		_, _ = r.AddFunc(crontab_time_rent, func() {
			now := time.Now()
			date := now.Format("02")
			if date == "15" && (*rent)["Rent_done"] != 1 {
				(*rent)["Rent_remind"] = 1
			}
			if date == "21" {
				(*rent)["Rent_remind"] = 1
				(*rent)["Rent_done"] = 0
			}
			if (*rent)["Rent_remind"] != 0 && date != "21" {
				bot_service.Telegram_bot_run(&chatID, &yourToken, "rent_service", 0)
			}
		})
	}(&Rent)

	start_info := "xingshan service start, Time setting: Pacakge:" + crontab_time_package + ", Rent:" + crontab_time_rent

	bot_service.Telegram_bot_run(&chatID, &yourToken, start_info, 0)
	p.Start()
	r.Start()
	select {}
}
