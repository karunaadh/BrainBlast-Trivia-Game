package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Structure for JSON response
type Data struct {
	ResponseCode int `json:"response_code"`
	Results      []struct {
		Type            string   `json:"type"`
		Difficulty      string   `json:"difficulty"`
		Category        string   `json:"category"`
		Question        string   `json:"question"`
		CorrectAnswer   string   `json:"correct_answer"`
		IncorrectAnswer []string `json:"incorrect_answers"`
	} `json:"results"`
}

func main() {
	//Track scores
	var numCorrect int = 0

	// Intro
	fmt.Println("Welcome to BrainBlast!")
	fmt.Print("Enter your name: ")
	var name string
	fmt.Scan(&name)

	fmt.Printf("Hello %v!\n", name)

	// Get number of questions
	var numQuestions int = -1
	for numQuestions < 0 {
		fmt.Print("Choose a number of questions: ")
		fmt.Scan(&numQuestions)
		//Check for negative inputs
		if numQuestions < 0 {
			fmt.Println("Please enter a valid number.")
			continue
		}
	}

	// Get difficulty
	var difficulty string

	for {
		fmt.Print("Choose difficulty level (1 - Easy, 2 - Medium, 3 - Hard): ")
		fmt.Scan(&difficulty)

		// Set difficulty
		if difficulty == "1" {
			difficulty = "easy"
			break
		} else if difficulty == "2" {
			difficulty = "medium"
			break
		} else if difficulty == "3" {
			difficulty = "hard"
			break
		} else {
			fmt.Println("Invalid difficulty. Please choose 1, 2, or 3.")
		}
	}

	// Set API URL
	apiURL := fmt.Sprintf("https://opentdb.com/api.php?amount=%d&category=9&difficulty=%s&type=boolean", numQuestions, difficulty)

	//Send get request
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error making GET request: ", err)
		return
	}

	defer response.Body.Close()

	// Read the response body
	var data Data
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return
	}

	// Print the parsed data
	if data.ResponseCode == 0 {
		fmt.Println("-----------------------------")
		fmt.Println("Get ready!")
		fmt.Println("-----------------------------")
		fmt.Println("")

		var answer string
		for i, question := range data.Results {
			//Prompt Question
			fmt.Printf("  Question %d: ", i+1)
			fmt.Println(question.Question)

			//Prompt Answer
			fmt.Print("  Answer (true/false): ")
			fmt.Scan(&answer)
			answer = strings.ToLower(answer)

			//Correct Answer
			correctAnswer := question.CorrectAnswer
			correctAnswer = strings.ToLower(correctAnswer)

			//Check for correction
			if answer == correctAnswer {
				fmt.Println("Correct! +1 points ")
				numCorrect++
			} else {
				fmt.Println("Incorrect! :( ")
			}
			fmt.Println("-----------------------------")
			fmt.Println("")
		}

		fmt.Println("STATS")
		fmt.Println("-----------------------------")
		fmt.Println("Total Score:", numCorrect)
		fmt.Printf("Success rate: %.1f%%\n", (float64(numCorrect)/float64(numQuestions))*100)
		fmt.Println("-----------------------------")

	} else if data.ResponseCode == 1 {
		fmt.Printf("Sorry, %d questions of your specified difficulty are not available.", numQuestions)
	} else {
		fmt.Println("We are currently experiencing an error. Try again later.")
	}

	fmt.Printf("\nCiao %v!", name)
}
