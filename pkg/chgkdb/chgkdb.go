package chgkdb

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type Question struct {
	Championship string
	Tour         string
	Date         string
	Number       int
	Question     string
	Answer       string
	Author       string
	Source       string
	Comment      string
}

func LoadQuestions(directory string) ([]Question, error) {
	var questions []Question

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			fileQuestions, err := parseFile(path)
			if err != nil {
				return err
			}
			questions = append(questions, fileQuestions...)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return questions, nil
}

func parseFile(filename string) ([]Question, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var questions []Question
	var currentQuestion Question
	var currentField *string

	questionNumberPattern := regexp.MustCompile(`Вопрос\s+(\d+):`)

	koi8rReader := transform.NewReader(file, charmap.KOI8R.NewDecoder())
	scanner := bufio.NewScanner(koi8rReader)

	fmt.Printf("parseFile %s\n", filename)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "Чемпионат:"):
			currentQuestion.Championship = strings.TrimSpace(strings.TrimPrefix(line, "Чемпионат:"))
			currentField = &currentQuestion.Championship
		case strings.HasPrefix(line, "Тур:"):
			currentQuestion.Tour = strings.TrimSpace(strings.TrimPrefix(line, "Тур:"))
			currentField = &currentQuestion.Tour
		case strings.HasPrefix(line, "Дата:"):
			currentQuestion.Date = strings.TrimSpace(strings.TrimPrefix(line, "Дата:"))
			currentField = &currentQuestion.Date
		case questionNumberPattern.MatchString(line):
			if currentQuestion.Question != "" {
				questions = append(questions, currentQuestion)
				currentQuestion = Question{Championship: currentQuestion.Championship, Tour: currentQuestion.Tour, Date: currentQuestion.Date}
			}
			matches := questionNumberPattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				currentQuestion.Number, _ = strconv.Atoi(matches[1])
				currentQuestion.Question = strings.TrimSpace(strings.TrimPrefix(line, matches[0]))
			}
			currentField = &currentQuestion.Question
		case strings.HasPrefix(line, "Ответ:"):
			currentQuestion.Answer = strings.TrimSpace(strings.TrimPrefix(line, "Ответ:"))
			currentField = &currentQuestion.Answer
		case strings.HasPrefix(line, "Автор:"):
			currentQuestion.Author = strings.TrimSpace(strings.TrimPrefix(line, "Автор:"))
			currentField = &currentQuestion.Author
		case strings.HasPrefix(line, "Источник:"):
			currentQuestion.Source = strings.TrimSpace(strings.TrimPrefix(line, "Источник:"))
			currentField = &currentQuestion.Source
		case strings.HasPrefix(line, "Комментарий:"):
			currentQuestion.Comment = strings.TrimSpace(strings.TrimPrefix(line, "Комментарий:"))
			currentField = &currentQuestion.Comment
		default:
			if currentField != nil {
				*currentField += "\n" + line
			}
		}
	}

	if currentQuestion.Question != "" {
		questions = append(questions, currentQuestion)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}
