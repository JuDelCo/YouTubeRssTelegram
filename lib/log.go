package lib

import (
	"log"
	"runtime/debug"
)

func LogDebug(msg string) {
	log.Println("[DEBUG] " + msg)
}

func LogError(msg string) {
	log.Println("[ERROR] " + msg)

	telegramChatID, err := TelegramGetChatID("logging")

	if err != nil {
		log.Println("[ERROR] " + err.Error())
		log.Println("[ERROR] Error sending error log to Telegram")
	} else {
		err = TelegramSendMessage(telegramChatID, "[ERROR] "+msg, true)

		if err != nil {
			log.Println("[ERROR] Error sending error log to Telegram")
		}
	}
}

func LogErrorFatal(msg string) {
	telegramChatID, err := TelegramGetChatID("logging")

	if err != nil {
		log.Println("[ERROR] " + err.Error())
		log.Println("[ERROR] Error sending FATAL error log to Telegram")
	} else {
		err = TelegramSendMessage(telegramChatID, "[FATAL] "+msg, true)

		if err != nil {
			log.Println("[ERROR] Error sending FATAL error log to Telegram")
		} else {
			TelegramSendMessage(telegramChatID, string(debug.Stack()), true)
		}
	}

	log.Println("[FATAL] " + msg)
	log.Fatalln(string(debug.Stack()))
}
