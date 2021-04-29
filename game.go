package main

type game struct {
	question_list     []Question
	points            map[*client]int
	attempted_answers map[*client]bool // svi koji su pokusali da odg na trenutno pitanje
	br_pitanja        int              // Zbog testiranja
}

func newGame(members []*client) *game {
	m := make(map[*client]int)
	a := make(map[*client]bool)

	for _, ptr := range members {
		m[ptr] = 0
		a[ptr] = false
	}

	return &game{
		question_list:     get_questions(11),
		points:            m,
		attempted_answers: a,
		br_pitanja:        0,
	}

}

func (g *game) getNextQuestion() (string, bool) {
	if g.br_pitanja == len(g.question_list)-1 {
		return "", true
	}
	for k := range g.attempted_answers {
		g.attempted_answers[k] = false
	}
	g.br_pitanja += 1
	return g.question_list[g.br_pitanja].question_string(), false
}

func (g *game) attemptAnswer(c *client, ans string) int {
	if g.attempted_answers[c] {
		return -1
	}
	g.attempted_answers[c] = true
	if g.question_list[g.br_pitanja].Correct_answer == ans {
		g.points[c] = g.points[c] + g.question_list[g.br_pitanja].Points
		return 1
	}
	return 0
}

// bool moveToNextQuestion() -> ako su svi odg onda true
func (g *game) moveToNextQuestion() bool {
	for _, v := range g.attempted_answers {
		if !v {
			return false
		}
	}
	return true
}

func (g *game) getPoints(c *client) int {
	return g.points[c]
}

func (g *game) leaveGame(c *client) {
	_, ok := g.attempted_answers[c]
	if ok {
		delete(g.attempted_answers, c)
		delete(g.points, c)
	}
}
