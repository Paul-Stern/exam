package main

import (
	"net/url"
	"strings"
	"unicode/utf8"
)

type Option struct {
	Id   int
	Text string
}

type Card struct {
	Id       int
	Question string
	Appendix []string
	Options  []Option
}

type Test struct {
	User  string
	Cards []Card
}

type Result struct {
	QuestionId string   `json:"questionId"`
	AnswerIds  []string `json:"answersIds"`
}

type TestResult struct {
	UserId  int      `json:"userId"`
	Results []Result `json:"results"`
}

type CardsResult []struct {
	QuestionId int   `json:"QuestionId"`
	AnswerIds  []int `json:"AnswerId"`
}

type restOpt struct {
	Id          int    `json:"ID"`
	Answer_text string `json:"ANSWER_TEXT"`
}

type restBlock struct {
	Id            int       `json:"ID"`
	Task_text     string    `json:"TASK_TEXT"`
	Task_appendix []string  `json:"TASK_APPENDIX"`
	Answers       []restOpt `json:"ANSWERS"`
}

type restBlocks []restBlock

func newCardsResult(vals url.Values) TestResult {
	var tr TestResult
	tr.UserId = 1
	for k := range vals {
		var r Result
		// cr[k] = vals[k]
		qId, _ := strings.CutPrefix(k, "question_")
		var aIds []string
		for _, aId := range vals[k] {
			aIds = append(aIds, aId)
		}
		r.QuestionId = qId
		r.AnswerIds = aIds
		tr.Results = append(tr.Results, r)
	}
	return tr
}

func newTest(u User, c []Card) Test {
	return Test{
		User:  getFullName(u),
		Cards: c,
	}
}

func getCards(rbs restBlocks) (cards []Card) {
	for _, block := range rbs {
		var c Card
		c.Id = block.Id
		c.Question = block.Task_text
		c.Appendix = removeAppendixPrefix(block.Task_appendix)
		for _, o := range block.Answers {
			c.Options = append(
				c.Options,
				newOption(o.Id, o.Answer_text),
			)
		}
		cards = append(cards, c)
	}
	return cards
}

func removeAppendixPrefix(ap []string) []string {
	if ap == nil {
		return ap
	}
	var result []string
	for _, s := range ap {
		_, j := utf8.DecodeRuneInString(s)
		// result[i] = s[(j * 3):]
		result = append(result, s[(j*3):])
	}
	return result
}
