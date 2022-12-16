package model

type PlayerModel struct {
	ID     int64  `db:"id" json:"id"`
	Name   string `db:"name" json:"name,omitempty"`
	TeamID int64  `db:"team_id" json:"teamId,omitempty"`
}
