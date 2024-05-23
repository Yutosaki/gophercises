package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of question,answer")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle",false,"shuffle the problems, type true or false")
	flag.Parse()

	//open the file
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("faild to open:%s\n", *csvFilename))
	}
	defer file.Close()

	//parse the file
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("failed to parse"))
	}

	//parse the line
	problems := parseLines(lines)

	//shuffle
	if *shuffle {
		rand.NewSource(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) {
            problems[i], problems[j] = problems[j], problems[i]
        })
	}

	//set the timer and counter
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	counter := 0
problemLoop:
	for i, p := range problems {
		fmt.Printf("problem: %d: %s = ", i+1, p.Q)
		answerCh := make(chan string)
		go func(){
			answer := ""
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			break problemLoop
		case answer := <- answerCh:
			if p.A == answer {
				fmt.Println("correct!")
				counter++
			}
		}
		
	}
	fmt.Printf("\nyou scored %d out of %d\n", counter, len(problems))
}

func parseLines(lines [][]string) []problem {
	res := make([]problem, len(lines))
	for i, line := range lines {
		res[i] = problem{
			Q: line[0],
			A: line[1],
		}
	}
	return res
}

type problem struct {
	Q string
	A string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
