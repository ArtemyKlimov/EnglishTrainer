package model

type DBUser struct {
	id               int
	firstName        string
	userName         string
	isStudyOwnFrases bool
	addedWordsToday  int
}

type DBFrase struct {
	id      int
	value   string
	addedBy int
}

type Operation int

type APIUser struct {
	CurrentOperation   Operation
	UserName           string
	CurrentWord        string
	CurrentAttempt     int
	CurrentTranslation []string
}

func (u *APIUser) SetOperation(o Operation) {
	u.CurrentOperation = o
}