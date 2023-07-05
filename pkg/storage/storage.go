package storage

type RusEngFrase struct {
	rusFrase []string
	engFrase []string
}

func New() *RusEngFrase {
	return &RusEngFrase{}
}
