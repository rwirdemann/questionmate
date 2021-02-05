# Questionmate

Questionmate is a [Golang](https://www.golang.org) library and API for building configurable online
surveys, assessments and recommendations. 

## Build

```shell
make
```

## Test

```shell
make test
```

## Run

```shell
./bin/questionsmate
```

## Usage

Surveys are modeled in a YAML based DSL:

```
questions:
  - id: 1
    text: Which programming languages does the system use?
    type: single
    options:
      - value: 1
        text: modern (e.g. go, Rust, Swift)
        targets:
          - sustainability: 3
      - value: 2
        text: standard (e.g. Java, C/C++, JavaScript, Php, Ruby)
        targets:
          - sustainability: 2
      - value: 4
        text: expiring (Cobol, Delphi)
        targets:
          - sustainability: 1
```

A survey consists of a list of questions. A question consists of its id (1), it's text, a type
(single or multi) and a list of options (modern, ...). An option can have one or more targets
(sustainability) associated with a credit value (3). 

## Credits

Users can select one or more of the question's options, depending on the question type. A selected
option pays into its associated targets. For instance, option 1 *modern* pays 3 credits into the
target *sustainability*.

## Question order

Questions are served in a natural order defined by their ids, thus question 2 follows question 1. 

## Assessment

At the end of answering a survey an assessment of the given answers can be requested. An assessment consists of a list
of targets each associated with its resulting score, modeled by the following go types:

```go
type Target struct {
	Text   string 
	Score  int    
	Rating string 
}

type Assessment struct {
	Targets []Target 
}
```

A targets score is calculated by adding up all credits configured for the target based on the given answer set, e.g. if
question 1 pays 3 credits into *sustainability* and question 5 pays 1 credits into *sustainability* the target will get
an overall score of 4. 

A target might contain an optional rating description that explains or assesses the meaning of a targets score.

## API

**POST http://questionmate.ralfwirdemann.de/{questionnaire}/questions**

Returns the next question according to the answers the user has given so far contained in the request body as follows:

```json
{
  "answers": [
    {
      "question_id": 1,
      "value": 1
    }
  ]
}
```

HTTP/1.1 200 Ok
```json
{
  "id": 2,
  "text": "How often do incidents occur without first changing the software?",
  "type": "single",
  "options": [
    {
      "value": 1,
      "text": "almost never",
      "targets": {
        "robustness": {
          "value": 10
        }
      }
    }
  ]
}
```

**POST http://questionmate.ralfwirdemann.de/{questionnaire}/assessment**

Returns the assessment of the answers contained in the request body as follows:

```json
{
  "answers": [
    {
      "question_id": 1,
      "value": 1
    }
  ]
}
```

HTTP/1.1 200 Ok
```json
{
  "targets": [
    {
      "text": "robustness",
      "score": 3,
      "rating": "Your software isn't very robust. You have to fix it soon."
    },
    {
      "text": "sustainability",
      "score": 7,
      "rating": "Your software is sustainable enough."
    }
  ]
}
```

## Architecture

Questionmate's design is based on hexagonal architecture ideas. Thus, it focuses on its domain model and rules but makes no assumptions about how this model is used within a specific context, e.g. the web.

Instead of forcing programmers to use *Questionmate* within a specific context examples are given about how to use the library wrapped as REST services deployed to AWS lambda.

### Domain Model

```
Questionnaire --> *Question --> *Option --> *Target
                      1            1           *
                      |            |           |
                      |            |           |
                    Answer --------/       Assessment
```

A Questionnaire consists of many Questions. A question has many options, whereby each options represents a possibe answer. Each option has many targets to which this option pays. 

An answer models a user's answer to a question, thus the answer's values must match one of the option values of its question.  


## License

* [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)