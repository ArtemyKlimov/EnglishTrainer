package main

import (
	"english_trainer/internal/model"
	"english_trainer/internal/texts"
	"english_trainer/internal/utils"
	"english_trainer/pkg/trainerbot"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".enf file couldn't be loaded")
	}
	api_token := os.Getenv("TG_API_TOKEN")
	bot, err := tgbotapi.NewBotAPI(api_token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	config := trainerbot.NewConfig()

	b, err := trainerbot.New(config)
	if err != nil {
		log.Fatal(fmt.Errorf("Can not initialize config: %v", err))
	}

	apiUsers := make(map[string]*model.APIUser)

	for update := range updates {
		if update.Message == nil { // If we got a message
			continue
		}
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			apiUser, ok := apiUsers[update.Message.From.UserName]
			if !ok {
				apiUser = &model.APIUser{
					CurrentOperation: utils.StartOperation,
					UserName:         update.Message.From.UserName,
				}
				apiUsers[update.Message.From.UserName] = apiUser
			}
			if strings.HasPrefix(update.Message.Text, "/") {
				produceCommand(&update, apiUser, bot, b)
				continue
			}
			if apiUser.CurrentOperation == utils.AddOperation {
				processAdd(&update, apiUser, bot, b)
				continue
			} else if apiUser.CurrentOperation == utils.TrainOperation {
				processTrain(&update, apiUser, bot, b)
				continue
			} else if apiUser.CurrentOperation == utils.LearnOperation {
				processLearn(&update, apiUser, bot, b)
				continue
			} else {
				processUnknown(&update, apiUser, bot, b)
			}
		}
	}
}

func produceCommand(update *tgbotapi.Update, apiUser *model.APIUser, bot *tgbotapi.BotAPI, b *trainerbot.TrainerBot) {
	switch update.Message.Text {
	case "/start":
		apiUser.SetOperation(utils.StartOperation)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.StartText)
		bot.Send(msg)
	case "/help":
		apiUser.SetOperation(utils.StartOperation)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.HelpText)
		bot.Send(msg)
	case "/add":
		apiUser.SetOperation(utils.AddOperation)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.AddOperationText)
		bot.Send(msg)
	case "/train":
		apiUser.SetOperation(utils.TrainOperation)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.TrainInputText)
		apiUser.CurrentWord = ""
		apiUser.CurrentTranslation = nil
		bot.Send(msg)
		processTrain(update, apiUser, bot, b)
	case "/learn":
		apiUser.SetOperation(utils.LearnOperation)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.LearnOperationText)
		bot.Send(msg)
		processLearn(update, apiUser, bot, b)
	case "/exit":
		apiUser.SetOperation(utils.StartOperation)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.UnknownCommandText)
		apiUser.SetOperation(utils.StartOperation)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func processUnknown(update *tgbotapi.Update, apiUser *model.APIUser, bot *tgbotapi.BotAPI, b *trainerbot.TrainerBot) {
	apiUser.SetOperation(utils.StartOperation)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.UnknownInputText)
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func processAdd(update *tgbotapi.Update, apiUser *model.APIUser, bot *tgbotapi.BotAPI, b *trainerbot.TrainerBot) {
	result, err := b.AddNewPair(update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func processTrain(update *tgbotapi.Update, apiUser *model.APIUser, bot *tgbotapi.BotAPI, b *trainerbot.TrainerBot) {
	rusPhrase, err := b.StartTraining(apiUser, update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, rusPhrase)
	bot.Send(msg)
}

func processLearn(update *tgbotapi.Update, apiUser *model.APIUser, bot *tgbotapi.BotAPI, b *trainerbot.TrainerBot) {
	engPhrases, rusPhrase, err := b.StartLearning()
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		bot.Send(msg)
		return
	}
	result := ""
	for _, s := range engPhrases {
		result += s + ", "
	}
	result = strings.TrimRight(result, ", ")
	result += " : " + rusPhrase
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
	bot.Send(msg)
}
