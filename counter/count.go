package counter

import (
	"sort"
	"strconv"
)

var (
	cnt map[string]int
)

func Init() {
	cnt = make(map[string]int)
}

func Add(UserName string) {
	cnt[UserName]++
}

type Rating struct {
	CntMessages int
	UserName    string
}

func Status() (ans string) {
	if len(cnt) == 0 {
		return "No messages yet"
	}
	var rating []Rating
	for key, val := range cnt {
		rating = append(rating, Rating{CntMessages: val, UserName: key})
	}
	sort.Slice(rating, func(i, j int) bool {
		return rating[i].CntMessages > rating[j].CntMessages
	})
	for _, val := range rating {
		ans += val.UserName + ": " + strconv.Itoa(val.CntMessages) + "\n"
	}
	return
}
