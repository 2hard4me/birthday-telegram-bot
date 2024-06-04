package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type BirthdayInfo struct {
	ChatID int64
	Name   string
	Day    int
	Month  int
}

func (bd BirthdayInfo) Date() string {
	return time.Date(0, time.Month(bd.Month), bd.Day, 0, 0, 0, 0, time.UTC).Format("02-Jan")
}

func (bd BirthdayInfo) String() string {
	dateStr := bd.Date()
	return fmt.Sprintf("%s on %s", bd.Name, dateStr)
}

func (bd BirthdayInfo) DaysLeft() int {
	var nextBdYear int
	now := time.Now().In(tz)
	if now.Month() > time.Month(bd.Month) || (now.Month() == time.Month(bd.Month) && now.Day() > bd.Day) {
		nextBdYear = now.Year() + 1
	} else {
		nextBdYear = now.Year()
	}
	nextBdDate := time.Date(nextBdYear, time.Month(bd.Month), bd.Day, 0, 0, 0, 0, tz)
	diff := nextBdDate.Sub(now)
	return int(math.Ceil((diff.Hours() / 24)))
}

func BirthdayListStr(bdList []BirthdayInfo) string {
	var bdStrList []string
	prevMonth := -1
	for _, bd := range bdList {
		if prevMonth != bd.Month {
			prevMonth = bd.Month
			bdStrList = append(bdStrList, fmt.Sprintf("\n<i>%v</i>", time.Month(bd.Month)))
		}
		bdStrList = append(bdStrList, fmt.Sprintf("  • %s (%d days left)", bd.String(), bd.DaysLeft()))
	}
	return strings.Join(bdStrList, "\n")
}

const (
	InvalidAddReply = "Please send the name of the person along with the command\neg. /add John Doe"
	InvalidLoginReply = "Incorrect login"
	InvalidPasswordReply = "Incorrect password"
)

func GetMsgData(msg string) (string, bool) {
	_, data, isfound := strings.Cut(msg, " ")
	data = strings.TrimSpace(data)
	if data == "" {
		return "", false
	}
	return data, isfound
}

func RetryReply(command string) string {
	return fmt.Sprintf("Something went wrong!!\nPlease retry %s again.", command)
}