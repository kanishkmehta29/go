package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type problem struct {
	a string
	b string
	c string
}

func ProblemPuller(filename string) ([]problem, error) {

	fObj, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening the CSV file: %w", err)
	}
	csvR := csv.NewReader(fObj)
	cLines, err := csvR.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading the CSV file: %w", err)
	}
	return ParseProblem(cLines), err
}

func ParseProblem(lines [][]string) []problem {
	r := make([]problem, len(lines))

	for i := 0; i < len(lines); i++ {
		r[i] = problem{
			a: lines[i][0],
			b: lines[i][1],
			c: lines[i][2],
		}
	}
	return r

}

func main() {
	fname := flag.String("f", "quiz.csv", "path of the csv file")
	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()

	problems, err := ProblemPuller(*fname)
	if err != nil {
		log.Fatalf("Some error occured during problem pulling %v", err)
	}

	correctcnt := 0

	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)

problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %v: %v+%v = ", i+1, p.a, p.b)
		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()

		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop
		case iAns := <-ansC:
			if iAns == p.c {
				correctcnt++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}

	}
	fmt.Printf("Your score is %d/%d", correctcnt, len(problems))
}
