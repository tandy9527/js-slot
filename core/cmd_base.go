package core

import (
	"github.com/tandy9527/js-slot/core/game"
)

func GetBalance(user *User, gameinfo *game.GameInfo, msg Message) GameResult {
	result := map[string]any{}
	result["balance"] = user.GetBalance()
	return GameResult{Data: result}
}
