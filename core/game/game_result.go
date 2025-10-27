package game

// TumbleResult 旋转一次盘面结果
// 游戏可自行定义扩展
type TumbleResult struct {
	TumbleSymbol [][]int `json:"TumbleSymbol"` // 盘面符号
	LineSymbol   []int   `json:"LineSymbol"`   // 连线中奖符号      [0,0,0,1,0,0,0,0]  -表示下标3的连线，中奖符号为1
	LineCount    []int   `json:"LineCount"`    // 连线中奖符号數量  [0,0,0,5,0,0,0,0]  -表示下标3的连线，中奖符号數量为5
	LineWin      []int64 `json:"LineWin"`      // 连线中奖金额      [0,0,0,500,0,0,0,0]-	表示下标3的连线，中奖金额为500
	Win          int64   `json:"Win"`          // 盘面总中奖
}

// MGResult 主游戏结果，游戏可自行定义扩展
type MGResult[T any] struct {
	MGTumbleList []T   `json:"MGTumbleList"` // 主遊戲盤面列表
	MainWin      int64 `json:"MainWin"`      // 主遊戲贏分
}

// FGResult 免费游戏结果,游戏可自行定义扩展
type FGResult[T any] struct {
	FGTumbleList []T   `json:"FGTumbleList"` // 免費遊戲 Spin 結果列表
	FreeWin      int64 `json:"FreeWin"`      // 免費遊戲贏分

}

type SlotResult[M any, F any] struct {
	MGResult MGResult[M] `json:"MGResult"` // MAIN GAME
	FGResult FGResult[F] `json:"FGResult"` // FREE GAME
	TotalWin int64       `json:"TotalWin"` // 总赢FG+MG
	Balance  int64       `json:"Balance"`  // 玩家余额
	Extra    any         `json:"Extra"`    // 额外扩展数据
}
