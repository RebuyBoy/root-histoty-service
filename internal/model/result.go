package model

import "time"

type PlayerResult struct {
	Player
	Race  string
	Rank  int
	Point int
}

type Result struct {
	GameId  int
	GameMap string
	Deck    string
	Date    time.Time
	Results []PlayerResult
}
