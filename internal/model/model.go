package model

type Operation int

type APIUser struct {
	Id                 int
	CurrentOperation   Operation
	UserName           string
	CurrentWord        string
	CurrentAttempt     int
	PreviousMessage    *TgMessage
	CurrentTranslation []string
	OnlyMyDictionary   bool
	AddedWordsToday    int
}

func (u *APIUser) SetOperation(o Operation) {
	u.CurrentOperation = o
}

type TgMessage struct {
	MessageId int
	ChatId    int64
}
