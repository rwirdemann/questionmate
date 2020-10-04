package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpadapter "github.com/rwirdemann/questionmate/adapter/driver/http"
	"github.com/rwirdemann/questionmate/adapter/repositories/file"
	"github.com/rwirdemann/questionmate/adapter/repositories/parser"
	"github.com/rwirdemann/questionmate/usecase"
	"log"
	"net/http"
	"os"
)

func main() {
	fn := fmt.Sprintf("%s/src/github.com/rwirdemann/questionmate/config/legacylab", os.Getenv("GOPATH"))
	directoryPtr := flag.String("directory", fn, "the questions directory")
	flag.Parse()

	// 1. Instantiate the "I need to go out httpadapter"
	repositoryAdapter := file.NewQuestionRepository(*directoryPtr, parser.QMParser{})

	// 2. Instantiate the hexagons
	hexagon := usecase.NextQuestion{QuestionRepository: repositoryAdapter}
	evaluator := usecase.Evaluator{QuestionRepository: repositoryAdapter}

	// 3. Instantiate the "I need to go in adapter"
	nextQuestoionHttpAdapter := httpadapter.MakeNextQuestionHandler(hexagon)
	evaluatorHttpAdapter := httpadapter.MakeEvaluationsHandler(evaluator)

	r := mux.NewRouter()
	r.HandleFunc("/{questionaire}/questions", nextQuestoionHttpAdapter).Methods("POST")
	r.HandleFunc("/{questionaire}/evaluations", evaluatorHttpAdapter).Methods("POST")
	log.Printf("Service listening on http://localhost:8080...")
	handler := cors.AllowAll().Handler(r)
	_ = http.ListenAndServe(":8080", handler)
}
