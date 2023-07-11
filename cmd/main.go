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

var trainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("не знаю 🤷", "dontknow"),
		tgbotapi.NewInlineKeyboardButtonData("подсказка🤔", "givemeahint"),
		tgbotapi.NewInlineKeyboardButtonData(nextString, "nextquestionpls"),
		tgbotapi.NewInlineKeyboardButtonData(exitString, "exittomainmenu"),
	),
)

var learnMyWordsKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(nextString),
		tgbotapi.NewKeyboardButton(backString),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(LearnMyWordsString),
	),
)

var learnAllWordsKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(nextString),
		tgbotapi.NewKeyboardButton(backString),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(LearnAllWordsString),
	),
)

var startKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("тренироваться"),
		tgbotapi.NewKeyboardButton("учить"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("добавить фразу в словарь"),
	),
)

var addKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(backString, backString),
	),
)

var startKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("тренироваться 🧠", "/train"),
		tgbotapi.NewInlineKeyboardButtonData("учить ✍", "/learn"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("добавить фразу в словарь 📒", "/add"),
	),
)
var (
	exitString          string = "выход ❌"
	backString          string = "назад 🔙"
	nextString          string = "дальше ✅"
	LearnMyWordsString  string = "Только мой словарь"
	LearnAllWordsString string = "Общий словарь"
	apiUsers            map[string]*model.APIUser
	bot                 *tgbotapi.BotAPI
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

	apiUsers = make(map[string]*model.APIUser)

	for update := range updates {
		apiUser, err := getApiUser(&update, b)
		if err != nil {
			processError(&update, bot, err, "")
			continue
		}
		if update.Message != nil { // If we got a message
			if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
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
		} else if update.CallbackQuery != nil {
			if apiUser.CurrentOperation == utils.TrainOperation || update.CallbackQuery.Data == "/train" {
				processTrain(&update, apiUser, bot, b)
				continue
			}
			if apiUser.CurrentOperation == utils.LearnOperation || update.CallbackQuery.Data == "/learn" {
				processLearn(&update, apiUser, bot, b)
				continue
			}
			if update.CallbackQuery.Data == "/add" {
				apiUser.SetOperation(utils.AddOperation)
				if apiUser.PreviousMessage != nil {
					removePreviousDialogButtons(apiUser.PreviousMessage.ChatId, apiUser.PreviousMessage.MessageId, bot)
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, texts.AddOperationText)
				msg.ReplyMarkup = addKeyBoard
				sentMsg, _ := bot.Send(msg)
				if apiUser.PreviousMessage == nil {
					apiUser.PreviousMessage = &model.TgMessage{}
				}
				apiUser.PreviousMessage.ChatId = sentMsg.Chat.ID
				apiUser.PreviousMessage.MessageId = sentMsg.MessageID
				continue
			}
			if update.CallbackQuery.Data == backString {
				if apiUser.PreviousMessage != nil {
					removePreviousDialogButtons(apiUser.PreviousMessage.ChatId, apiUser.PreviousMessage.MessageId, bot)
				}
				goToMainMenu(update.CallbackQuery.Message.Chat.ID, apiUser, bot)
				continue
			}

		}
	}
}

func processError(update *tgbotapi.Update, bot *tgbotapi.BotAPI, err error, additionalText string) {
	var chatId int64
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
	} else {
		chatId = update.Message.Chat.ID
	}
	msg := tgbotapi.NewMessage(chatId, err.Error()+additionalText)
	bot.Send(msg)
}

func getApiUser(update *tgbotapi.Update, b *trainerbot.TrainerBot) (*model.APIUser, error) {
	userName, err := getUserName(update)
	if err != nil {
		return nil, err
	}
	apiUser, ok := apiUsers[userName]
	if ok {
		return apiUser, nil
	}
	apiUser = &model.APIUser{
		CurrentOperation: utils.StartOperation,
		UserName:         userName,
		CurrentAttempt:   0,
		OnlyMyDictionary: false,
	}

	err = b.FindUserByUserName(apiUser)
	if err != nil {
		return nil, err
	}
	apiUsers[userName] = apiUser
	return apiUser, nil
}

func produceCommand(update *tgbotapi.Update, apiUser *model.APIUser, bot *tgbotapi.BotAPI, b *trainerbot.TrainerBot) {
	switch update.Message.Text {
	case "/start":
		apiUser.SetOperation(utils.StartOperation)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.StartText)
		msg.ReplyMarkup = startKeyboard
		if apiUser.PreviousMessage == nil {
			apiUser.PreviousMessage = &model.TgMessage{}
		} else {
			removePreviousDialogButtons(update.Message.Chat.ID, update.Message.MessageID, bot)
		}
		sentMsg, _ := bot.Send(msg)
		apiUser.PreviousMessage.ChatId = sentMsg.Chat.ID
		apiUser.PreviousMessage.MessageId = sentMsg.MessageID

	case "/help":
		apiUser.SetOperation(utils.StartOperation)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.HelpText)
		bot.Send(msg)
	case "/add":
		apiUser.SetOperation(utils.AddOperation)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.AddOperationText)
		msg.ReplyMarkup = addKeyBoard
		bot.Send(msg)
	case "/train":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.TrainInputText)
		apiUser.CurrentWord = ""
		apiUser.CurrentTranslation = nil
		bot.Send(msg)
		processTrain(update, apiUser, bot, b)
	case "/learn":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.LearnOperationText)
		bot.Send(msg)
		processLearn(update, apiUser, bot, b)
	case "/exit":
		goToMainMenu(update.Message.Chat.ID, apiUser, bot)
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

func getChatId(update *tgbotapi.Update) int64 {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}
	return update.Message.Chat.ID
}

func getMessageId(update *tgbotapi.Update) int {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.MessageID
	}
	return update.Message.MessageID
}

func processAdd(update *tgbotapi.Update, apiUser *model.APIUser, bot *tgbotapi.BotAPI, b *trainerbot.TrainerBot) {
	if apiUser.PreviousMessage != nil {
		removePreviousDialogButtons(apiUser.PreviousMessage.ChatId, apiUser.PreviousMessage.MessageId, bot)
	}
	result, err := b.AddNewPair(update.Message.Text, apiUser)
	var msg tgbotapi.MessageConfig
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, result)
		msg.ReplyToMessageID = update.Message.MessageID
	}
	msg.ReplyMarkup = addKeyBoard
	sentMsg, _ := bot.Send(msg)
	apiUser.PreviousMessage.ChatId = sentMsg.Chat.ID
	apiUser.PreviousMessage.MessageId = sentMsg.MessageID
}
func processLearn(update *tgbotapi.Update, apiUser *model.APIUser, bot *tgbotapi.BotAPI, b *trainerbot.TrainerBot) {
	apiUser.SetOperation(utils.LearnOperation)
	chatId := getChatId(update)
	if apiUser.PreviousMessage != nil {
		removePreviousDialogButtons(apiUser.PreviousMessage.ChatId, apiUser.PreviousMessage.MessageId, bot)
	}

	if update.Message != nil {
		switch update.Message.Text {
		case LearnMyWordsString:
			apiUser.OnlyMyDictionary = true
			msg := tgbotapi.NewMessage(chatId, "Теперь изучаем слова и фразы из моего словаря")
			msg.ReplyMarkup = learnAllWordsKeyBoard
			bot.Send(msg)
			return
		case LearnAllWordsString:
			apiUser.OnlyMyDictionary = false
			msg := tgbotapi.NewMessage(chatId, "Теперь изучаем слова и фразы из общего словаря")
			msg.ReplyMarkup = learnMyWordsKeyBoard
			bot.Send(msg)
			return
		case backString:
			goToMainMenu(chatId, apiUser, bot)
			return
		}
	}
	var currentMarkUp tgbotapi.ReplyKeyboardMarkup
	if apiUser.OnlyMyDictionary {
		currentMarkUp = learnAllWordsKeyBoard
	} else {
		currentMarkUp = learnMyWordsKeyBoard
	}
	engPhrases, rusPhrase, err := b.StartLearning(apiUser)
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, err.Error())
		msg.ReplyMarkup = currentMarkUp
		bot.Send(msg)
		return
	}
	result := ""
	for _, s := range engPhrases {
		result += s + ", "
	}
	result = strings.TrimRight(result, ", ")
	result += " : " + rusPhrase
	msg := tgbotapi.NewMessage(chatId, result)
	msg.ReplyMarkup = currentMarkUp
	sentMsg, _ := bot.Send(msg)
	apiUser.PreviousMessage.ChatId = sentMsg.Chat.ID
	apiUser.PreviousMessage.MessageId = sentMsg.MessageID
}

func goToMainMenu(chatId int64, apiUser *model.APIUser, bot *tgbotapi.BotAPI) {
	apiUser.SetOperation(utils.StartOperation)
	apiUser.CurrentAttempt = 0
	apiUser.CurrentWord = ""
	apiUser.CurrentTranslation = nil
	if apiUser.PreviousMessage != nil {
		removePreviousDialogButtons(apiUser.PreviousMessage.ChatId, apiUser.PreviousMessage.MessageId, bot)
	}
	msg := tgbotapi.NewMessage(chatId, texts.ExitText)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
	msg = tgbotapi.NewMessage(chatId, "Что будем делать:")
	msg.ReplyMarkup = startKeyboard
	sentMsg, _ := bot.Send(msg)
	if apiUser.PreviousMessage == nil {
		apiUser.PreviousMessage = &model.TgMessage{}
	}
	apiUser.PreviousMessage.ChatId = sentMsg.Chat.ID
	apiUser.PreviousMessage.MessageId = sentMsg.MessageID
}

func removePreviousDialogButtons(chatId int64, messageId int, bot *tgbotapi.BotAPI) {
	editedMsg := tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0, 0)})
	bot.Send(editedMsg)
}

func processTrain(update *tgbotapi.Update, apiUser *model.APIUser, bot *tgbotapi.BotAPI, b *trainerbot.TrainerBot) {
	apiUser.SetOperation(utils.TrainOperation)
	var rusPhrase string
	var err error
	var msg tgbotapi.MessageConfig
	var chatId int64
	var messageId int
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
		messageId = update.CallbackQuery.Message.MessageID
		removePreviousDialogButtons(chatId, messageId, bot)
		if update.CallbackQuery.Data == "exittomainmenu" {
			goToMainMenu(chatId, apiUser, bot)
			return
		}
		rusPhrase, err = b.StartTraining(apiUser, update.CallbackQuery.Data)
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, rusPhrase)
	} else {
		if apiUser.PreviousMessage != nil {
			removePreviousDialogButtons(update.Message.Chat.ID, update.Message.MessageID, bot)
		}

		fmt.Println("TEXT = ", update.Message.Text)
		if update.Message.Text == "exittomainmenu" {
			goToMainMenu(update.Message.Chat.ID, apiUser, bot)
			return
		}

		rusPhrase, err = b.StartTraining(apiUser, update.Message.Text)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, rusPhrase)

	}
	if err != nil {
		processError(update, bot, err, ". Нажмите /exit для выхода в главное меню")
		return
	}

	msg.ReplyMarkup = trainKeyboard
	sendedMsg, _ := bot.Send(msg)
	if apiUser.PreviousMessage == nil {
		apiUser.PreviousMessage = &model.TgMessage{}
	}
	apiUser.PreviousMessage.ChatId = sendedMsg.Chat.ID
	apiUser.PreviousMessage.MessageId = sendedMsg.MessageID
}

func getUserName(update *tgbotapi.Update) (string, error) {
	if update.Message != nil {
		return update.Message.From.UserName, nil
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.From.UserName, nil
	} else if update.ChannelPost != nil {
		return update.ChannelPost.From.UserName, nil
	} else if update.EditedMessage != nil {
		return update.EditedMessage.From.UserName, nil
	} else {
		return "", fmt.Errorf("Unsopported message type")
	}
}
