package main

import (
	"fmt"
	"os"

	"github.com/rwirdemann/questionmate/adapter/repositories/file"
	"github.com/rwirdemann/questionmate/adapter/repositories/parser"
)

func main() {
	fn := fmt.Sprintf("%s/src/github.com/rwirdemann/questionmate/config/legacylab", os.Getenv("GOPATH"))
	r := file.NewQuestionRepository(fn, parser.QMParser{})
	fmt.Printf("%s", r)
}
