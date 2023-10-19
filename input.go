package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

type Input interface {}

type ShortAnswerField struct {
	textinput textinput.Model
}
// textinput
func NewLShortAnswerField() *ShortAnswerField {
	ta := textinput.New()
	return &ShortAnswerField{ta}
}

type LongAnswerField struct {
	textinput textarea.Model
}

// textarea
func NewLongAnswerField() *LongAnswerField {
	ta := textarea.New()
	return &LongAnswerField{ta}
}