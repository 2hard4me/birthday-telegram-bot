package utils

import (
	"database/sql"
	"fmt"
	"strings"
	"github.com/jmoiron/sqlx"
)

var (
	DBConn *sqlx.DB
	accessGranted bool
)

const (
	insertSQL         = "INSERT INTO birthdays (chat_id, name, day, month) VALUES ($1, $2, $3, $4)"
	selectOnSQL       = "SELECT chat_id, name, day, month FROM birthdays WHERE day = $1 AND month = $2"
	selectAllSQL      = "SELECT name, day, month FROM birthdays WHERE chat_id = $1"
	selectNameSQL     = "SELECT name, day, month FROM birthdays WHERE chat_id = $1 AND name = $2"
	selectNameLikeSQL = "SELECT name, day, month FROM birthdays WHERE chat_id = $1 AND LOWER(name) LIKE $2"
	selectMonthSQL    = "SELECT name, day, month FROM birthdays WHERE chat_id = $1 AND month = $2"
	selectDaySQL      = "SELECT name, day, month FROM birthdays WHERE chat_id = $1 AND day = $2"
	selectDateSQL     = "SELECT name, day, month FROM birthdays WHERE chat_id = $1 AND day = $2 AND month = $3"
	updateDateSQL     = "UPDATE birthdays SET day = $1, month = $2 WHERE chat_id = $3 AND name = $4"
	updateNameSQL     = "UPDATE birthdays SET name = $1 WHERE chat_id = $2 AND name = $3"
	deleteSQL         = "DELETE FROM birthdays WHERE chat_id = $1 AND name = $2"
	orderbySQL        = " ORDER BY month, day"
	authSQL 		  = "SELECT password from accounts WHERE login = &1"
)

func GetPassword(login string) (string, error) {
	var pass string

	err := DBConn.QueryRowx(authSQL, login).Scan(&pass)
	if err != nil {
		return "", err
	}
	return pass, nil
}

func GetBirthday(chatID int64, name string) (BirthdayInfo, error) {
	var day, month int

	err := DBConn.QueryRowx(selectNameSQL, chatID, name).Scan(&name, &day, &month)
	if err != nil {
		if err == sql.ErrNoRows {
			return BirthdayInfo{}, nil
		} else {
			return BirthdayInfo{}, err
		}
	}

	bd := BirthdayInfo{
		Name:  name,
		Day:   day,
		Month: month,
	}
	return bd, nil
}

func GetBirthdaysOn(day int, month int) ([]BirthdayInfo, error) {
	rows, err := DBConn.Queryx(selectOnSQL, day, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var info BirthdayInfo
	var infoList []BirthdayInfo
	for rows.Next() {
		err := rows.Scan(&info.ChatID, &info.Name, &info.Day, &info.Month)
		if err != nil {
			return nil, err
		}
		infoList = append(infoList, info)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return infoList, nil
}

func AddBirthday(bd *BirthdayInfo) error {
	_, err := DBConn.Exec(insertSQL, bd.ChatID, bd.Name, bd.Day, bd.Month)
	return err
}

func SearchBirthday(bd *BirthdayInfo, searchby string) ([]BirthdayInfo, error) {
	var rows *sqlx.Rows
	var err error

	switch searchby {
	case "all":
		rows, err = DBConn.Queryx(selectAllSQL+orderbySQL, bd.ChatID)
	case "name":
		likeExp := "%" + strings.ToLower(bd.Name) + "%"
		rows, err = DBConn.Queryx(selectNameLikeSQL+orderbySQL, bd.ChatID, likeExp)
	case "month":
		rows, err = DBConn.Queryx(selectMonthSQL+orderbySQL, bd.ChatID, bd.Month)
	case "day":
		rows, err = DBConn.Queryx(selectDaySQL+orderbySQL, bd.ChatID, bd.Day)
	case "date":
		rows, err = DBConn.Queryx(selectDateSQL+orderbySQL, bd.ChatID, bd.Day, bd.Month)
	default:
		return nil, fmt.Errorf("invalid searchby: %s", searchby)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	info := BirthdayInfo{ChatID: bd.ChatID}
	var infoList []BirthdayInfo
	for rows.Next() {
		err := rows.Scan(&info.Name, &info.Day, &info.Month)
		if err != nil {
			return nil, err
		}
		infoList = append(infoList, info)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return infoList, nil
}

func UpdateDate(bd *BirthdayInfo) error {
	_, err := DBConn.Exec(updateDateSQL, bd.Day, bd.Month, bd.ChatID, bd.Name)
	return err
}

func UpdateName(bd *BirthdayInfo, newName string) error {
	_, err := DBConn.Exec(updateNameSQL, newName, bd.ChatID, bd.Name)
	return err
}

func RemoveBirthday(bd *BirthdayInfo) error {
	_, err := DBConn.Exec(deleteSQL, bd.ChatID, bd.Name)
	return err
}