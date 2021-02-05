package file

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/rwirdemann/questionmate/adapter/repositories/parser"
	"github.com/rwirdemann/questionmate/domain"
)

// todo: rename to QuestionnaireRepository
type QuestionRepository struct {
	questionnaire domain.Questionnaire
	Questions     []domain.Question
	Descriptions  map[int]string
	Targets       map[string]string
	Ratings       map[string][]domain.Rating
}

func NewQuestionRepository(filename string, parser parser.Parser) QuestionRepository {
	var questionnaire domain.Questionnaire
	var questions []domain.Question
	var ratings map[string][]domain.Rating
	var descriptions map[int]string
	var targets map[string]string

	if bytes, ok := readFile(filename + "." + parser.Suffix()); ok {
		questions = parser.ParseQuestions(bytes)
		questionnaire = parser.ParseQuestionnaire(bytes)
		ratings = parser.ParseRatings(bytes)
	}

	return QuestionRepository{questionnaire: questionnaire,
		Questions: questions, Ratings: ratings, Descriptions: descriptions, Targets: targets}
}

func readFile(file string) ([]byte, bool) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	return data, true
}

func (q QuestionRepository) GetQuestions() []domain.Question {
	return q.Questions
}

func (q QuestionRepository) GetRatings() map[string][]domain.Rating {
	return q.Ratings
}

func (q QuestionRepository) GetDescriptions() map[int]string {
	return q.Descriptions
}

func (q QuestionRepository) GetQuestionnaire() domain.Questionnaire {
	return q.questionnaire
}

func (q QuestionRepository) String() string {
	var s string
	for _, question := range q.Questions {
		s = fmt.Sprintf("%s---------------------------------------------------------------------------------------\n", s)
		s = fmt.Sprintf("%sFrage: %s\n", s, question.Text)
		if desc, ok := q.Descriptions[question.ID]; ok {
			s = fmt.Sprintf("%sBeschreibung: %s\n", s, desc)
		}
		for _, option := range question.Options {
			s = fmt.Sprintf("%s- %s\n", s, option.Text)
			for k, v := range option.Targets {
				s = fmt.Sprintf("%s  - %s: %d\n", s, k, v.Value)
			}
		}
	}
	return s
}
