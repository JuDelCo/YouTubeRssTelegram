package lib

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

type TelegramChatInfo struct {
	chatType string
	chatID   int64
}

// Private global variables (for lib package)
var telegramBot *tb.Bot
var telegramChats = make([]TelegramChatInfo, 0)

func TelegramLoadChatIDs(telegramIDsPath string) error {
	records, err := ReadCsvFile(telegramIDsPath)

	if err != nil {
		return err
	}

	for _, line := range records {
		chatInfo := TelegramChatInfo{}
		chatInfo.chatType = strings.TrimSpace(line[0])
		chatInfo.chatID, err = strconv.ParseInt(line[1], 10, 0)

		telegramChats = append(telegramChats, chatInfo)
	}

	return err
}

func TelegramGetChatID(chatType string) (int64, error) {
	for _, chatInfo := range telegramChats {
		if chatInfo.chatType == chatType {
			return chatInfo.chatID, nil
		}
	}

	return 0, errors.New("Telegram Chat ID not found: " + chatType)
}

func InitializeTelegramBot(telegramApiTokenPath string) error {
	telegramApiTokenFile, err := os.Open(telegramApiTokenPath)

	if err != nil {
		return errors.New("Unable to open Telegram token file: " + telegramApiTokenPath + " \n" + err.Error())
	}

	telegramApiToken, err := ioutil.ReadAll(telegramApiTokenFile)

	if err != nil {
		return errors.New("Unable to read Telegram token file: " + telegramApiTokenPath + " \n" + err.Error())
	}

	b, err := tb.NewBot(tb.Settings{Token: string(telegramApiToken), Synchronous: true})

	if err != nil {
		return errors.New("Unable to initialize Telegram Bot\n" + err.Error())
	}

	telegramBot = b

	return nil
}

func TelegramSendMessage(chatId int64, msg string, disablePreview bool) error {
	chat := tb.ChatID(chatId)

	var err error

	if disablePreview {
		_, err = telegramBot.Send(chat, msg, tb.NoPreview)
	} else {
		_, err = telegramBot.Send(chat, msg)
	}

	return err
}

func TelegramSendImageMessage(chatId int64, imgUrl string, msg string) error {
	chat := tb.ChatID(chatId)

	img := &tb.Photo{File: tb.FromURL(imgUrl)}
	img.Caption = msg

	_, err := telegramBot.Send(chat, img, tb.NoPreview)

	return err
}
