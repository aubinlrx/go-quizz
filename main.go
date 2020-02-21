package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Question struct {
	Question string
	Answer   string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quizz in seconds")
	flag.Parse()
	questions, err := readQuestionsFromFile(csvFilename)
	if err != nil {
		log.Fatal(err)
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	goodAnswerCount := 0
	for _, question := range questions {
		fmt.Printf("what %s, sir?\n", question.Question)
		answerChan := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChan <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d.\n", goodAnswerCount, len(questions))
			return
		case answer := <-answerChan:
			if answer == question.Answer {
				goodAnswerCount++
			}
		}
	}

	fmt.Printf("Good answers: %d/%d\n", goodAnswerCount, len(questions))
}

func readQuestionsFromFile(filename *string) ([]Question, error) {
	csvFile, err := os.Open(*filename) // for read access
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	var questions []Question

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		questions = append(questions, Question{
			Question: line[0],
			Answer:   line[1],
		})
	}

	return questions, nil
}
