package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Reading struct {
	Resoposnse_code int              `json:"response_code"`
	Questions       []QuestionFormat `json:"results"`
}

type QuestionFormat struct {
	Category   string   `json:"category"`
	Type       string   `json:"type"`
	Difficulty string   `json:"difficulty"`
	Question   string   `json:"question"`
	Correct    string   `json:"correct_answer"`
	Incorrect  []string `json:"incorrect_answers"`
}

type Question struct {
	Text           string
	Answers        map[string]string
	Points         int
	Correct_answer string
}

func (question Question) question_string() string {
	colors := [3]string{"\033[32m", "\033[34m", "\033[31m"}
	return colors[question.Points-1] + html.UnescapeString(question.Text+"\na)\t"+question.Answers["a"]+
		"\nb)\t"+question.Answers["b"]+
		"\nc)\t"+question.Answers["c"]+
		"\nd)\t"+question.Answers["d"]) + "\033[0m"
	// return html.UnescapeString(question.Text + "\na)\t" + question.Answers["a"] +
	// 	"\nb)\t" + question.Answers["b"] +
	// 	"\nc)\t" + question.Answers["c"] +
	// 	"\nd)\t" + question.Answers["d"])
}

func get_questions(n int) []Question {
	if n > 50 { // sa sajta moze da se dohvati 50 pitanja max
		n = 50
	}
	responses, err := http.Get("https://opentdb.com/api.php?amount=" + strconv.Itoa(n) + "&type=multiple")
	var data []byte
	if err != nil {
		fmt.Printf("HTTP request filed %s\n", err)
		return get_backup_questions(n)
	}
	data, _ = ioutil.ReadAll(responses.Body)

	questions := Reading{}
	err = json.Unmarshal(data, &questions)
	if err != nil {
		fmt.Println("Error, ", err)
		return make([]Question, 0)
	}

	rand.Seed(time.Now().Unix())
	options := []string{"a", "b", "c", "d"}
	var questions_ret = make([]Question, n)
	for i := 0; i < n; i++ {
		question := Question{}
		question.Text = questions.Questions[i].Question
		switch questions.Questions[i].Difficulty {
		case "easy":
			question.Points = 1
		case "medium":
			question.Points = 2
		case "hard":
			question.Points = 3
		default:
			question.Points = 1
		}

		rand.Shuffle(4, func(i, j int) {
			options[i], options[j] = options[j], options[i]
		})

		question.Correct_answer = options[0]
		question.Answers = make(map[string]string)
		question.Answers[options[0]] = questions.Questions[i].Correct
		for j := 1; j < 4; j++ {
			question.Answers[options[j]] = questions.Questions[i].Incorrect[j-1]
		}

		questions_ret[i] = question
	}
	return questions_ret
}

func get_backup_questions(n int) []Question {
	data, err := ioutil.ReadFile("backup_questions.json")
	if err != nil {
		var ret []Question
		return ret
	}
	var questions []Question
	json.Unmarshal(data, &questions)

	rand.Seed(time.Now().Unix())

	l := len(questions)
	if l < n {
		return questions
	}

	selected_questions := rand.Perm(l)[:n]

	ret := make([]Question, n)
	for i := 0; i < n; i++ {
		ret[i] = questions[selected_questions[i]]
	}
	return ret
}

func update_backup_questions(n int) {
	file, _ := json.Marshal(get_questions(n))
	err := ioutil.WriteFile("backup_questions.json", file, 0664)
	if err != nil {
		println("Error: write backup questions")
	}
}

// func main() {
// 	update_backup_questions(50)
// 	for i, q := range get_questions(10) {
// 		fmt.Println(i, q.question_string())
// 	}
// }
