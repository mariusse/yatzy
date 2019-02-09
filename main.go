package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type game struct {
	categories map[int]int
	score      int
}

var (
	score int
	wg    sync.WaitGroup
)

// GetGame returns a game that can be *game.Play()'ed
func GetGame() *game {
	return &game{
		categories: map[int]int{
			1: 0,
			2: 0,
			3: 0,
			4: 0,
			5: 0,
			6: 0,
		},
		score: 0,
	}
}

// pickCategory returns highest scoring category in roll
// returns 0 if roll only contains taken categories
func (g *game) pickCategory(diceroll []int) int {
	numberOfCategories := map[int]int{}
	for _, v := range diceroll {
		numberOfCategories[v]++
		if numberOfCategories[v] == 5 && g.categories[v] == 0 {
			return v
		}
		if numberOfCategories[v] > 4 && g.categories[v] == 0 {
			return v
		}
	}

	highestScoreCat := 0
	for k, v := range numberOfCategories {
		if k*v > highestScoreCat && g.categories[k] == 0 {
			highestScoreCat = k
		}
	}

	return highestScoreCat
}

//Play plays a round of yatzy
func (g *game) Play() int {
	for i := 0; i < len(g.categories); i++ {
		diceKept := 0
		chosenCategory := 0

		for rollNum := 0; rollNum < 3; rollNum++ {
			r, _ := rollDice(5 - diceKept)

			if chosenCategory == 0 {
				chosenCategory = g.pickCategory(r)
			}

			for _, v := range r {
				if v == chosenCategory {
					g.categories[v]++
				}
			}
		}

	}

	g.tallyScore()
	return g.score
}

func (g *game) tallyScore() {
	for k, v := range g.categories {
		g.score += k * v
	}
}

func playMany() {
	wg.Add(1)
	g := GetGame()
	score += g.Play()
	wg.Done()
}

func main() {
	numGames := 100
	for i := 0; i < numGames; i++ {
		go playMany()
	}
	wg.Wait()
	fmt.Println("Average score of", numGames, "games is: ", score/numGames)
}

func rollDice(n int) ([]int, error) {
	if n == 0 || n > 5 {
		return nil, fmt.Errorf("invalid numer of dice. Must be between 1 and 5")
	}
	url := "http://nav-yatzy.herokuapp.com/throw?n=" + strconv.Itoa(n)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r := []int{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
