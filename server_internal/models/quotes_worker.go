package models

type QuotesWork struct {
	ID     uint
	Quote  string
	Error  error
	Finish chan struct{}
}
