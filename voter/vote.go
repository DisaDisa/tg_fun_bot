package voter

import (
	"log"
	"sort"
	"strconv"
	"time"
)

var (
	cnt        map[string]int
	top_active map[int64]bool
)

func Init() {
	cnt = make(map[string]int)
	top_active = make(map[int64]bool)
}

func Add(UserName string) {
	cnt[UserName]++
}

type Rating struct {
	CntVotes int
	UserName string
}

func Status() (ans string) {
	if len(cnt) == 0 {
		return "No messages yet"
	}
	var rating []Rating
	for user, votes := range cnt {
		rating = append(rating, Rating{CntVotes: votes, UserName: user})
	}
	sort.Slice(rating, func(i, j int) bool {
		return rating[i].CntVotes > rating[j].CntVotes
	})
	for _, val := range rating {
		ans += val.UserName + ": " + strconv.Itoa(val.CntVotes) + "\n"
	}
	return
}

func Top() (topUser string, count int) {
	for user, rating := range cnt {
		if count < rating {
			count = rating
			topUser = user
		}
	}
	return
}

func ActivateDayTop(ID int64, bot *tgbotapi.BotAPI) {
	if top_active[ID] == true {
		return
	}
	top_active[ID] = true
	go func() {
		if top_active[ID] == false {
			return
		}
		ticker := time.NewTicker(time.Hour * 24)
		for tickTime := range ticker.C {
			log.Println("Ticker: ", tickTime)
			user, cnt_r := Top()
			cnt = make(map[string]int)
			reply := "Выбор дня: " + user + "\nC этим согласны " + strconv.Itoa(cnt_r) + " человек."
			msg := tgbotapi.NewMessage(ID, reply)
			//msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}()
}

func DisactivateDayTop(ID int64) {
	top_active[ID] = false
}
