package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpadapter "github.com/rwirdemann/questionmate/adapter/driver/http"
	"github.com/rwirdemann/questionmate/adapter/repositories/file"
	"github.com/rwirdemann/questionmate/adapter/repositories/parser"
	"github.com/rwirdemann/questionmate/usecase"
)

func main() {
	fn := fmt.Sprintf("%s/src/github.com/rwirdemann/questionmate/config", os.Getenv("GOPATH"))
	directoryPtr := flag.String("directory", fn, "the config directory")
	flag.Parse()

	// 1. Instantiate the "I need to go out adapter"
	comaRepositoryAdapter := file.NewQuestionRepository(*directoryPtr+"/coma", parser.YAMLParser{})

	// 2. Instantiate the hexagons
	getQuestionnaire := usecase.NewGetQuestionnaire()
	getQuestionnaire.Repositories["coma"] = comaRepositoryAdapter

	hexagon := usecase.NextQuestion{QuestionRepository: comaRepositoryAdapter}
	evaluator := usecase.Assessment{QuestionRepository: comaRepositoryAdapter}

	// 3. Instantiate the "I need to go in adapter"
	getQuestionnaireHTTPAdapter := httpadapter.MakeGetQuestionnaireHandler(getQuestionnaire)
	nextQuestionHTTPAdapter := httpadapter.MakeNextQuestionHandler(hexagon)
	evaluatorHTTPAdapter := httpadapter.MakeAssessmentHandler(evaluator)

	r := mux.NewRouter()
	r.HandleFunc("/{questionnaire}", getQuestionnaireHTTPAdapter).Methods("GET")
	r.HandleFunc("/{questionnaire}/questions", nextQuestionHTTPAdapter).Methods("POST")
	r.HandleFunc("/{questionnaire}/assessment", evaluatorHTTPAdapter).Methods("POST")
	log.Printf("Service listening on http://localhost:8080...")
	handler := cors.AllowAll().Handler(r)
	_ = http.ListenAndServe(":8080", handler)
}
