package model

import (
	player_model "github.com/tesarwijaya/ouroboros/internal/domain/player/model"
)

type TeamModel struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type TeamPlayerRespModel struct {
	TeamModel
	Players []player_model.PlayerModel
}
