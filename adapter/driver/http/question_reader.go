package http

import (
	"encoding/json"
	"github.com/rwirdemann/questionmate/domain"
	"github.com/rwirdemann/questionmate/usecase"
	"io/ioutil"
	"log"
	"net/http"
)

func MakeNextQuestionHandler(questionReader usecase.QuestionReader) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("error: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var body struct {
			Answers []domain.Answer `json:"answers"`
		}
		if err := json.Unmarshal(b, &body); err != nil {
			log.Printf("error: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		q, hasNext := questionReader.NextQuestion(body.Answers)
		if hasNext {
			data, err := json.Marshal(q)
			if err != nil {
				log.Printf("error: %s", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = writer.Write(data)
			if err != nil {
				log.Printf("error: %s", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			writer.WriteHeader(http.StatusNoContent)
		}
	}
}

func MakeEvaluationsHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		evaluation := domain.Evaluation{}
		evaluation.Targets = append(evaluation.Targets, &domain.Target{Text: "changeability", Score: 12})
		evaluation.Targets = append(evaluation.Targets, &domain.Target{Text: "extendability", Score: 8})
		evaluation.Targets = append(evaluation.Targets, &domain.Target{Text: "robustness", Score: 3})
		data, err := json.Marshal(evaluation)
		if err != nil {
			log.Printf("error: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(data)
		if err != nil {
			log.Printf("error: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
