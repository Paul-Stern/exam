package main

import (
	"net/url"
	"strconv"
	"strings"
	"time"
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
	User    User        `json:"PERSINFO"`
	Profile TestProfile `json:"TASK_PROFILE"`
	Time    Time        `json:"TIMING"`
	Cards   []Card      `json:"RESULTS"`
}

type ResultStore struct {
	Id           int  `json:"TESTING_RESULT_ID"`
	Total        int  `json:"TOTAL_QUESTIONS"`
	CorrectCount int  `json:"RIGHT_ANSWERS"`
	Percent      int  `json:"PERCENT"`
	Certified    bool `json:"CERTIFIED"`
}

type Time struct {
	Start time.Time `json:"TESTING_START"`
	End   time.Time `json:"TESTING_END"`
}

type Result struct {
	QuestionId int   `json:"questionId"`
	AnswerIds  []int `json:"answersIds"`
}

type TestResult struct {
	User    User        `json:"PERSINFO"`
	Profile TestProfile `json:"TASK_PROFILE"`
	Time    Time        `json:"TIMING"`
	Results []Result    `json:"RESULTS"`
}

type CardsResult []struct {
	QuestionId int   `json:"QuestionId"`
	AnswerIds  []int `json:"AnswerId"`
}

type TestPlan struct {
	Id int
}

type AvailableTestProfiles struct {
	User     User          `json:"USER"`
	Profiles []TestProfile `json:"TASK_PROFILES"`
}

type TestProfile struct {
	Id   int    `json:"TASK_PROFILE_ID"`
	Text string `json:"TASK_PROFILE_NAME"`
}

type Profiles struct {
	Profiles []TestProfile `json:"TASK_PROFILES"`
}

type TaskOption struct {
	Id          int    `json:"ID"`
	Answer_text string `json:"ANSWER_TEXT"`
}

type Task struct {
	Id            int          `json:"ID"`
	Task_text     string       `json:"TASK_TEXT"`
	Task_appendix []string     `json:"TASK_APPENDIX"`
	Answers       []TaskOption `json:"ANSWERS"`
}

type Tasks struct {
	Tasks []Task `json:"QUESTIONS"`
}

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
	tr.User.Id = m["userId"]
	tr.Time.End = time.Now()
	tr.Profile.Id = m["profile_id"]
	if err != nil {
		return tr, err
	}
	for k, v := range m {
		// Create Result instance
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

func newTest(u User, c []Card, pr TestProfile) (t Test) {
	return Test{
		User:    u,
		Cards:   c,
		Profile: pr,
		Time: Time{
			Start: time.Now(),
		},
	}
}

func getCards(tasks Tasks) (cards []Card) {
	for _, task := range tasks.Tasks {
		var c Card
		c.Id = task.Id
		c.Question = task.Task_text
		c.Appendix = removeAppendixPrefix(task.Task_appendix)
		for _, o := range task.Answers {
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
