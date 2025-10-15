package game

// BaseGame 提供基础结构体，供具体游戏组合使用。
type BaseGame struct {
	GameID   string
	GameName string
}

func NewBaseGame(id, name string) *BaseGame {
	return &BaseGame{GameID: id, GameName: name}
}

func (b *BaseGame) GetID() string   { return b.GameID }
func (b *BaseGame) GetName() string { return b.GameName }

type Game interface {
	GetID() string
	GetName() string
	// 获取游戏配置
	GetGameInfo() *GameInfo
	Spin()
}
