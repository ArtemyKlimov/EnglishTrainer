package main

import (
	"english_trainer/internal/model"
	"english_trainer/internal/texts"
	"english_trainer/internal/utils"
	"english_trainer/pkg/trainerbot"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	bot *tgbotapi.BotAPI
)

func initTelegram(botToken, baseURL, pemFile string) {
	var err error

	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Println(err)
		return
	}

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
	url := baseURL + "/" + bot.Token
	pemFilePath := tgbotapi.FilePath(pemFile)

	wh, _ := tgbotapi.NewWebhookWithCert(url, pemFilePath)
	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	log.Printf("IPAddress = %s", info.IPAddress)
	log.Printf("URL = %s", info.URL)
	if err != nil {
		log.Println(err)
	}

}

func main() {
	f, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	err = godotenv.Load()
	if err != nil {
		log.Fatal(".enf file couldn't be loaded")
	}
	api_token := os.Getenv("TG_API_TOKEN")
	keyPem := os.Getenv("PEM_KEY")
	certPem := os.Getenv("PEM_CERT")
	baseURL := os.Getenv("BASE_URL")
	config := trainerbot.NewConfig()
	b, err := trainerbot.New(config)
	if err != nil {
		log.Fatal(fmt.Errorf("Can not initialize config: %v", err))
	}
	bot, err := tgbotapi.NewBotAPI(api_token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	initTelegram(api_token, baseURL, certPem)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook("/" + bot.Token)

	go http.ListenAndServeTLS(":8443", certPem, keyPem, nil)
	log.Printf("BOT_TOKEN = %s", bot.Token)

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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.ExitText)
		bot.Send(msg)
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
