package file

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/rwirdemann/questionmate/domain"
)

type QuestionRepository struct {
}

func (q QuestionRepository) LoadQuestions(data []byte) domain.Questionaire {
	var questionaire domain.Questionaire
	questionaire.Questions = make(map[int]*domain.Question)
	lines := strings.Split(string(data), "\n")

	var question *domain.Question
	var option *domain.Option

	for _, l := range lines {
		if isQuestion(l) {
			q := strings.Split(l, ":")
			id := toInt(q[0])
			question = domain.NewQuestion(id, strings.Trim(q[1], " "))
			questionaire.Questions[id] = question
		}

		if question != nil && isType(l) {
			t := strings.Split(l, ":")
			question.Type = strings.Trim(t[1], " ")
		}

		if question != nil && isDependency(l) {
			d := strings.Split(l, "=>")
			questionID := toInt(d[0])
			optionID := toInt(d[1])
			question.Dependencies[questionID] = optionID
		}

		if question != nil && isOption(l) {
			o := strings.Split(l, ":")
			id := toInt(o[0])
			option = domain.NewOption(id, strings.Trim(o[1], " "))
			question.Options = append(question.Options, option)
		}

		if option != nil && isTarget(l) {
			t := strings.Split(l, ":")
			target := strings.Trim(t[0], " ")
			value := toInt(t[1])
			option.Targets[target[2:]] = domain.Score{Value: value}
		}
	}

	return questionaire
}

func isTarget(s string) bool {
	match, _ := regexp.MatchString("(^ {2}- [a-z]+): \\d+", s)
	return match
}

func toInt(s string) int {
	i, err := strconv.Atoi(strings.Trim(s, " "))
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func isOption(s string) bool {
	match, _ := regexp.MatchString("(^ {2}\\d+):", s)
	return match
}

func isType(s string) bool {
	return strings.HasPrefix(s, "type:")
}

func isQuestion(s string) bool {
	match, _ := regexp.MatchString("(^\\d+):", s)
	return match
}

func isDependency(s string) bool {
	match, _ := regexp.MatchString("(^ {2}\\d+) => \\d+", s)
	return match
}