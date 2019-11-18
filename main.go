package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type problem struct {
	ques string
	ans  string
}

func main() {
	csvfile := flag.String("csv", "problems.csv", "A CSV file that has questions and answers")
	timelimit := flag.Int("limit", 30, "Time limit for the quiz")
	flag.Parse()

	file, error := os.Open(*csvfile)

	if error != nil {
		fmt.Printf("There was an error opening the CSV file %s\n", *csvfile)
		os.Exit(1)
	}

	r := csv.NewReader(file)
	input_lines, err := r.ReadAll()

	if err != nil {
		fmt.Print("There was an error while reading the CSV file")
		os.Exit(1)
	}
	problems := line_parser(input_lines)

	timer := time.NewTimer(time.Duration(*timelimit) * time.Second)

	correct := 0
	incorrect := 0
problem_loop:
	for i, prob := range problems {
		fmt.Printf("Problem %d: %s = \n", i+1, prob.ques)
		ans_channel := make(chan int)
		go func() {
			var answer int
			fmt.Scanf("%d\n", &answer)
			ans_channel <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nCorrect answers : %d out of %d\n", correct, len(input_lines))
			fmt.Printf("\nIncorrect answers : %d out of %d\n", incorrect, len(input_lines))
			break problem_loop
		case answer := <-ans_channel:
			real_answer, _ := strconv.Atoi(prob.ans)
			if answer == real_answer {
				correct++
			} else {
				incorrect++
			}

		}

	}
	fmt.Printf("Correct answers : %d out of %d\n", correct, len(input_lines))
	fmt.Printf("Incorrect answers : %d out of %d\n", incorrect, len(input_lines))
}

func line_parser(lines [][]string) []problem {
	parsed_result := make([]problem, len(lines))

	for i, v := range lines {
		parsed_result[i] = problem{
			ques: v[0],
			ans:  v[1],
		}
	}
	return parsed_result
}
