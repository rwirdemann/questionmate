package main

import (
	. "fmt"
	"github.com/rwirdemann/questionmate/adapter/repositories/file"
	"github.com/rwirdemann/questionmate/domain"
	"github.com/rwirdemann/questionmate/usecase/nextquestion"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	fn := Sprintf("%s/src/github.com/rwirdemann/questionmate/config/%s", os.Getenv("GOPATH"), "legacylab.file")
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal(err)
	}

	var store file.QuestionRepository
	uc := nextquestion.NewUseCase(store, data)
	var answers []domain.Answer
	hasNext := true
	for hasNext {
		var nextQuestion domain.Question
		nextQuestion, hasNext = uc.NextQuestion(answers)
		if hasNext {
			Printf("%s\n", nextQuestion.Text)
			for _, option := range nextQuestion.Options {
				Printf("%d: %s\n", option.Value, option.Text)
			}
			Print("Your answer: ")
			var answer string
			_, err := Scanln(&answer)
			if err != nil {
				log.Fatal(err)
			}

			// loop until the user has selected a valid option
			var option *domain.Option
			for option == nil {
				i, err := strconv.Atoi(answer)
				if err == nil {
					option = nextQuestion.GetOption(i)
					if option == nil {
						Print("Try again: ")
						_, _ = Scanln(&answer)
					} else {
						a := domain.Answer{QuestionID: nextQuestion.ID}
						answers = append(answers, a)
					}
				} else {
					Print("Try again: ")
					_, _ = Scanln(&answer)
				}
			}
		}
	}
}
