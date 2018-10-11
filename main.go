package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"telegram_bot/counter"
	"telegram_bot/voter"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
)

func main() {
	args := os.Args[1:]
	proxyAddres := args[0]
	proxyUser := args[1]
	proxyPassword := args[2]
	botApi := args[3]

	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	auth := &proxy.Auth{User: proxyUser, Password: proxyPassword}
	dialer, err := proxy.SOCKS5("tcp", proxyAddres, auth, proxy.Direct)
	httpTransport.Dial = dialer.Dial
	bot, err := tgbotapi.NewBotAPIWithClient(botApi, httpClient)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	counter.Init()
	voter.Init()

	for update := range updates {

		go func() {
			for {
				ticker := time.NewTicker(time.Minute)
				select {
				case <-ticker.C:
					user, cnt := voter.Top()
					reply := "Пидор дня: " + user + "\nC этим согласны " + strconv.Itoa(cnt) + " человек."
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					//msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
			}
		}()

		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Command() == "" {
			counter.Add(update.Message.From.UserName)
		}

		var reply string

		switch update.Message.Command() {
		case "status":
			log.Printf("status command")
			reply = counter.Status()
		case "reset":
			log.Printf("reset command")
			counter.Init()
			reply = "Done"
		case "vote_pidor":
			log.Printf("votePidor command")
			text := update.Message.Text
			for i := len(text); i > 0; i-- {
				if text[i-1:i] == " " {
					voter.Add(text[i:])
					break
				}
			}
		case "vote_status":
			log.Printf("voteStatus command")
			reply = voter.Status()
		case "activate_day_top":
			log.Printf("activate day top")
			voter.ActivateDayTop(update.Message.Chat.ID, bot)
		case "disactivate_day_top":
			log.Printf("disactivate day top")
			voter.DisactivateDayTop(update.Message.Chat.ID)
		}
		if reply != "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			//msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
