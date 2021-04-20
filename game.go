// TODO struktura game
// lista pitanja -> question_list
// mapa poeni - inicijalizovati na 0 -> points
// mapa attempted_answers -> ko je pokusao da da odg

package main

type game struct {
	question_list []Question
	points        map[*client]int
	br_pitanja    int // Zbog testiranja
}

func newGame(members []*client) *game {
	m := make(map[*client]int)

	for _, ptr := range members {
		m[ptr] = 0
	}

	return &game{
		question_list: get_questions(10),
		points:        m,
		br_pitanja:    0,
	}

}

// get next question -> string getNextQuestion() {question string fja}
func (g *game) getNextQuestion() (string, bool) {
	if g.br_pitanja == 4 {
		return "", true
	}
	g.br_pitanja += 1
	return "Pitanje", false
}

// bool attemptAnswer(client, string answer)
func (g *game) attemptAnswer(c *client, ans string) bool {
	return true
}

// bool moveToNextQuestion() -> ako su svi odg onda true
func (g *game) moveToNextQuestion() bool {
	return true
}

// int getPoints(client)
func (g *game) getPoints(c *client) int {
	return 0
}
