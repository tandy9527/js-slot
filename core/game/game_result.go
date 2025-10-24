package game

// ReelSymbol 轉輪帶
type ReelSymbol []Symboler

// GameSymbol 獎圖盤面
type GameSymbol []ReelSymbol

// TumbleResult 旋转一次盘面结果
type TumbleResult struct {
	TumbleSymbol GameSymbol `json:"TumbleSymbol"` // 盘面符号
	LineSymbol   []Symboler `json:"LineSymbol"`   // 连线中奖符号      [0,0,0,1,0,0,0,0]  -表示下标3的连线，中奖符号为1
	LineCount    []int      `json:"LineCount"`    // 连线中奖符号數量  [0,0,0,5,0,0,0,0]  -表示下标3的连线，中奖符号數量为5
	LineWin      []uint64   `json:"LineWin"`      // 连线中奖金额      [0,0,0,500,0,0,0,0]-	表示下标3的连线，中奖金额为500
	Win          uint64     `json:"Win"`          // 盘面总中奖
}

// MGResult 主遊戲結果
type MGResult struct {
	MGTumbleList []TumbleResult `json:"MGTumbleList"` // 主遊戲盤面列表
	MainWin      uint64         `json:"MainWin"`      // 主遊戲贏分
}

// FGResult 免費遊戲結果
type FGResult struct {
	FGTumbleList []TumbleResult `json:"FGTumbleList"` // 免費遊戲 Spin 結果列表
	FreeWin      uint64         `json:"FreeWin"`      // 免費遊戲贏分

}

type SlotResult struct {
	MGResult
	FGResult
	TotalBet uint64 `json:"-"`        // 總投注額
	TotalWin uint64 `json:"TotalWin"` // 總贏分
}
