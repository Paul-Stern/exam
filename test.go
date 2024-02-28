package main

import (
	"net/url"
	"strconv"
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
	User  User
	Cards []Card
}

type Result struct {
	QuestionId int   `json:"questionId"`
	AnswerIds  []int `json:"answersIds"`
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

type finishTest struct {
	RetCode int `json:"RetCode"`
}

func flattenMap(vals url.Values) (m map[string]int) {
	m = make(map[string]int)
	for k, v := range vals {
		m[k], _ = strconv.Atoi(v[0])
	}
	return m
}

func (tr TestResult) indexOf(id int) (index int, found bool) {
	index = -1
	for i, v := range tr.Results {
		if id == v.QuestionId {
			found = true
			index = i
			break
		}
	}
	return
}

func newTestResult(vals url.Values) (tr TestResult, err error) {
	m := flattenMap(vals)
	tr.UserId = m["userId"]
	if err != nil {
		return tr, err
	}
	for k, v := range m {
		var r Result
		if strings.Contains(k, "answer_on_question_") {
			idString, _ := strings.CutPrefix(k, "answer_on_question_")
			r.QuestionId, err = strconv.Atoi(idString)
			r.AnswerIds = append(r.AnswerIds, v)
			i, found := tr.indexOf(r.QuestionId)
			if found {
				tr.Results[i].AnswerIds = append(tr.Results[i].AnswerIds, v)
			} else {
				tr.Results = append(tr.Results, r)
			}
		} else if strings.Contains(k, "_id") {
			r.QuestionId = v
			_, found := tr.indexOf(r.QuestionId)
			if !found {
				tr.Results = append(tr.Results, r)
			}
		}
	}
	return tr, err
}

func newTest(u User, c []Card) Test {
	return Test{
		User:  u,
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
