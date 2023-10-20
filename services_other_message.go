package main

import (
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func SendMessageToFinanceGroup(msg string) error {
	godotenv.Load(".env")
	bot, _ := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	bot.Debug = true
	chat_id, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHATID_FINANCE"), 10, 64)
	msg_send := tgbotapi.NewMessage(chat_id, msg)
	if _, err := bot.Send(msg_send); err != nil {
		return err
	}
	return nil
}

func SendMessageToContestGroup(msg string) error {
	godotenv.Load(".env")
	bot, _ := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	bot.Debug = true
	chat_id, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHATID_CONTEST"), 10, 64)
	msg_send := tgbotapi.NewMessage(chat_id, msg)
	if _, err := bot.Send(msg_send); err != nil {
		return err
	}
	return nil
}

func SaveToMessages(type_id int, msg string) error {
	message := CpsMessages{
		TypeID:  type_id,
		Message: msg,
	}
	if err := db_ksc.Create(&message).Error; err != nil {
		return err
	}
	switch type_id {
	case 1:
		{
			if err := SendMessageToContestGroup(msg); err != nil {
				return err
			}
			message.IsSent = 1
		}
	case 2:
		{
			if err := SendMessageToFinanceGroup(msg); err != nil {
				return err
			}
			message.IsSent = 1
		}
	}
	if err := db_ksc.Save(&message).Error; err != nil {
		return err
	}
	return nil
}

func GetAndSendMessageFromDb() error {
	messages := []CpsMessages{}
	if err := db_ksc.Model(messages).Where("is_sent = 0").Find(&messages).Error; err != nil {
		return err
	}
	if len(messages) == 0 {
		return nil
	}
	for i, v := range messages {
		switch v.TypeID {
		case 1:
			{
				if err := SendMessageToContestGroup(v.Message); err != nil {
					return err
				}
				messages[i].IsSent = 1
			}
		case 2:
			{
				if err := SendMessageToFinanceGroup(v.Message); err != nil {
					return err
				}
				messages[i].IsSent = 1
			}
		}
	}
	if err := db_ksc.Save(&messages).Error; err != nil {
		return err
	}
	return nil
}
