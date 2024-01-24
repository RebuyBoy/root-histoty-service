package dto

import "time"

type PlayerResultDTO struct {
	PlayerDTO
	Race          string `json:"race"`
	StartPosition int    `json:"start_position"`
	Rank          int    `json:"rank"`
	Point         int    `json:"point"`
}

type ResultDTO struct {
	GameId  int               `json:"game_id"`
	GameMap string            `json:"game_map"`
	Deck    string            `json:"deck"`
	Date    time.Time         `json:"date"`
	Results []PlayerResultDTO `json:"results"`
}
