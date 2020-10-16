package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseYAMLQuestionannaire(t *testing.T) {
	fn := fmt.Sprintf("%s/src/github.com/rwirdemann/questionmate/config/%s", os.Getenv("GOPATH"), "coma/questions.yaml")
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal(err)
	}

	parser := YAMLParser{}
	questionnaire := parser.ParseQuestionnaire(data)
	assert.Equal(t, "Dieser Fragebogen enthält sechs Fragen zum Zustand Ihres Softwaresystems. Nach Beantwortung der Fragen erhalten Sie eine erste Einschätzung sowie eine Bewertung Ihres Systems.", questionnaire.Abstract)
}

func TestParseYAMLQuestions(t *testing.T) {
	fn := fmt.Sprintf("%s/src/github.com/rwirdemann/questionmate/config/%s", os.Getenv("GOPATH"), "coma/questions.yaml")
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal(err)
	}

	parser := YAMLParser{}
	questions := parser.ParseQuestions(data)
	assert.Equal(t, 1, questions[0].ID)
	assert.Equal(t, "Die aktuellen Sourcen sind...", questions[0].Text)
	assert.Equal(t, 1, questions[0].Options[0].Value)
	assert.Equal(t, "vollständig vorhanden und übersetzbar", questions[0].Options[0].Text)
	assert.Equal(t, 2, questions[0].Options[0].Targets["fitness"].Value)
}

func TestParseYAMLRatings(t *testing.T) {
	fn := fmt.Sprintf("%s/src/github.com/rwirdemann/questionmate/config/%s", os.Getenv("GOPATH"), "coma/ratings.yaml")
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal(err)
	}

	parser := YAMLParser{}
	ratings := parser.ParseRatings(data)
	assert.Equal(t, 0, ratings["fitness"][0].Min)
	assert.Equal(t, 5, ratings["fitness"][0].Max)
	assert.Equal(t, "Akuter Handlungsbedarf, da in einem Fehlerfall nicht das Problem behoben werden kann, sondern erst aufwendige Vorarbeiten erledigt werden müssen. Das System ist akut gefährdet.", ratings["fitness"][0].Description)

	assert.Equal(t, 6, ratings["fitness"][1].Min)
	assert.Equal(t, 11, ratings["fitness"][1].Max)
	assert.Equal(t, "Hier besteht Handlungsbedarf, um im Fehlerfall vorbereitet zu sein.", ratings["fitness"][1].Description)
}
