package trainerbot

import (
	"english_trainer/internal/model"
	"english_trainer/internal/store"
	"english_trainer/internal/utils"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type TrainerBot struct {
	config *Config
	store  *store.Store
}

func New(c *Config) (*TrainerBot, error) {
	t := &TrainerBot{
		config: c,
	}
	err := t.configureStore()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (bot *TrainerBot) AddNewPair(text string, user *model.APIUser) (string, error) {
	splitted_text := strings.Split(text, ":")
	if len(splitted_text) == 1 {
		return "", errors.New("Неправильно переданы значения. Должен быть разделитель \":\"")
	}
	if strings.ContainsAny(text, ".!?*()#@") {
		return "", errors.New("Переданы недопустимые символы в строке")
	}
	splitted := strings.Split(text, ":")
	if len(splitted) != 2 {
		return "", errors.New("Разделитель \":\" должен быть только один")
	}
	if strings.ContainsAny(splitted[1], "абвгдеёжзийклмнопрстуфхцчшщъыьэюя") {
		return "", errors.New("Перевод не должен содержать кириллических символов")
	}
	rusPhrases := strings.Split(splitted[0], ",")
	engPhrases := strings.Split(splitted[1], ",")
	var totalResult string
	for _, rs := range rusPhrases {
		for _, es := range engPhrases {
			var result string
			rusV := strings.Replace(rs, "`", "'", -1)
			engV := strings.Replace(es, "`", "'", -1)
			if err := bot.store.Db.QueryRow("SELECT EnglishTrainer.AddNewPair ($1,$2,$3)", rusV, engV, user.Id).
				Scan(&result); err != nil {
				return "", err
			}
			if result != "" {
				totalResult += "Пара " + rusV + ":" + engV + " уже присутствует в словаре.\n"
			} else {
				totalResult += "Пара " + rusV + ":" + engV + " успешно добавлена в словарь.\n"
			}
		}
	}

	return totalResult, nil
}

func (bot *TrainerBot) FindUserByUserName(user *model.APIUser) error {
	var userId int
	var istrainOwnFrases bool
	var addedWordsToday int
	if err := bot.store.Db.QueryRow("SELECT id, istrainOwnFrases, addedWordsToday FROM englishTrainer.Users WHERE userName = $1", user.UserName).
		Scan(&userId, &istrainOwnFrases, &addedWordsToday); err == nil {
		user.Id = userId
		user.OnlyMyDictionary = istrainOwnFrases
		user.AddedWordsToday = addedWordsToday
		return nil
	} else if !strings.Contains(err.Error(), "no rows in result set") {
		return err
	}
	if err := bot.store.Db.QueryRow("INSERT INTO englishTrainer.Users (userName, istrainOwnFrases, addedWordsToday)"+
		" VALUES ($1, $2, $3) RETURNING id", user.UserName, false, 0).Scan(&userId); err != nil {
		return err
	}
	user.Id = userId
	return nil
}

func (bot *TrainerBot) StartLearning(user *model.APIUser) ([]string, string, error) {
	results := make([]string, 0)
	var totalNum int
	var rusPhraseId int
	var random int
	var rusPhrase string
	rand.Seed(time.Now().UnixNano())
	if user.OnlyMyDictionary {
		idMap := make(map[int]int)
		rows, err := bot.store.Db.Query("SELECT id FROM englishTrainer.rusPhrase WHERE addedById = $1", user.Id)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				return nil, "", fmt.Errorf("Ваш словарь пока пустой")
			}
			return nil, "", err
		}
		i := 0
		var id int
		for rows.Next() {
			rows.Scan(&id)
			idMap[i] = id
			i++
		}
		if len(idMap) == 0 {
			return nil, "", fmt.Errorf("Вы не добавили еще слов в Ваш словарь")
		}
		num := rand.Intn(len(idMap))
		random = idMap[num]
	} else {
		if err := bot.store.Db.QueryRow("SELECT COUNT(*) FROM englishTrainer.rusPhrase").Scan(&totalNum); err != nil {
			return nil, "", err
		}
		random = 1 + rand.Intn(totalNum)

	}
	if err := bot.store.Db.QueryRow("SELECT id, value FROM englishTrainer.rusPhrase WHERE id = $1", random).Scan(&rusPhraseId, &rusPhrase); err != nil {
		fmt.Println("ERROR is HERE!")
		return nil, "", err
	}
	rows, err := bot.store.Db.Query(
		"SELECT value FROM englishTrainer.engPhrase WHERE id IN "+
			"(SELECT engId FROM englishTrainer.RusEngPhrase WHERE rusId = $1)", rusPhraseId)
	if err != nil {
		return nil, "", err

	}
	var scanString string
	for rows.Next() {
		rows.Scan(&scanString)
		results = append(results, utils.Normalizetext(scanString))
	}
	return results, rusPhrase, nil
}

func (bot *TrainerBot) StartTraining(user *model.APIUser, text string) (string, error) {
	results := make([]string, 0)
	resultRusPhrase := "как сказать: "
	if user.CurrentWord != "" && text != "nextquestionpls" {
		if strings.ToLower(text) == "не знаю" || strings.ToLower(text) == "dontknow" {
			resultRusPhrase = user.CurrentWord + " - это " + user.CurrentTranslation[0] + ". А как сказать: "
		} else if text == "givemeahint" {
			return "Вот подсказка: " + utils.ReplaceEverySecondSymbol(user.CurrentTranslation[0]), nil
		} else {
			var isMatch bool
			for _, s := range user.CurrentTranslation {
				if utils.Normalizetext(text) == utils.Normalizetext(s) {
					resultRusPhrase = "Молодец! А " + resultRusPhrase
					isMatch = true
					break
				}
			}
			if !isMatch {
				user.CurrentAttempt += 1
				return "Неверно, попробуй еще раз. Если не знаешь, напиши 'не знаю'", nil
			}

		}
	}
	var totalNum int
	var random int
	rand.Seed(time.Now().UnixNano())
	if user.OnlyMyDictionary {
		idMap := make(map[int]int)
		rows, err := bot.store.Db.Query("SELECT id FROM englishTrainer.rusPhrase WHERE addedById = $1", user.Id)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return "", fmt.Errorf("Ваш словарь пока пустой")
			}
			return "", err
		}
		i := 0
		var id int
		for rows.Next() {
			rows.Scan(&id)
			idMap[i] = id
			i++
		}
		if len(idMap) == 0 {
			return "", fmt.Errorf("Вы не добавили еще слов в Ваш словарь")
		}
		num := rand.Intn(len(idMap))
		random = idMap[num]
	} else {

		if err := bot.store.Db.QueryRow("SELECT COUNT(*) FROM englishTrainer.rusPhrase").Scan(&totalNum); err != nil {
			return "", err
		}
		random = 1 + rand.Intn(totalNum)
	}
	var rusPhraseId int
	var rusPhrase string
	if err := bot.store.Db.QueryRow("SELECT id, value FROM englishTrainer.rusPhrase WHERE id = $1", random).Scan(&rusPhraseId, &rusPhrase); err != nil {
		return "", err
	}

	rows, err := bot.store.Db.Query(
		"SELECT value FROM englishTrainer.engPhrase WHERE id IN "+
			"(SELECT engId FROM englishTrainer.RusEngPhrase WHERE rusId = $1)", rusPhraseId)
	if err != nil {
		return "", err

	}
	var scanString string
	for rows.Next() {
		rows.Scan(&scanString)
		results = append(results, utils.Normalizetext(scanString))
	}
	user.CurrentWord = rusPhrase
	user.CurrentTranslation = results
	user.CurrentAttempt = 1

	return resultRusPhrase + "'" + rusPhrase + "'" + "?", nil
}

func (bot *TrainerBot) configureStore() error {
	st := store.New(bot.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	bot.store = st
	return nil
}
