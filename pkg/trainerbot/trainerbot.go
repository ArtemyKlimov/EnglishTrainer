package trainerbot

import (
	"english_trainer/internal/model"
	"english_trainer/internal/store"
	"english_trainer/internal/utils"
	"errors"
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

func (bot *TrainerBot) AddNewPair(text string) (string, error) {
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
			if err := bot.store.Db.QueryRow("SELECT EnglishTrainer.AddNewPair ($1,$2)", rusV, engV).
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

func (bot *TrainerBot) StartLearning() ([]string, string, error) {
	results := make([]string, 0)
	var totalNum int
	if err := bot.store.Db.QueryRow("SELECT COUNT(*) FROM englishTrainer.rusPhrase").Scan(&totalNum); err != nil {
		return nil, "", err
	}
	rand.Seed(time.Now().UnixNano())
	rand := 1 + rand.Intn(totalNum)
	var rusPhraseId int
	var rusPhrase string
	if err := bot.store.Db.QueryRow("SELECT id, value FROM englishTrainer.rusPhrase WHERE id = $1", rand).Scan(&rusPhraseId, &rusPhrase); err != nil {
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
	if err := bot.store.Db.QueryRow("SELECT COUNT(*) FROM englishTrainer.rusPhrase").Scan(&totalNum); err != nil {
		return "", err
	}
	rand.Seed(time.Now().UnixNano())
	rand := 1 + rand.Intn(totalNum)
	var rusPhraseId int
	var rusPhrase string
	if err := bot.store.Db.QueryRow("SELECT id, value FROM englishTrainer.rusPhrase WHERE id = $1", rand).Scan(&rusPhraseId, &rusPhrase); err != nil {
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
