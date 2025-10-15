package game

// game_check_line.go 文件实现了 GameCheckLine 接口，用于检查游戏中的连线情况。
import "slices"

// CheckLineWin 检查连线是否中奖
// checkReel: 检查的转盘符号:[row][col]->[1,2,3][1,2,3,4,5][1,2,3,4,5]
// lineData: 连线数据 ps: [0,0,0]  第一列第一行，第二列第一行，第三列第一行符号是否相同
// odds: 赔率- 连线后该符号的赔率
// wildSymbolID: 百变符号
// feverSymbolID: 特殊符号,不参与连线
// lineBet: 连线下注
func CheckLineWin(checkReel [][]int, lineData []int, odds map[int][]float64, wildSymbolID int, wildList []int, feverSymbolID int, lineBet float64, isFromLeft bool,
	skipFiveSymbol bool, wildMultiplier [][]float64, allMultiplier float64) float64 {

	checkReelList := make([]int, len(checkReel))
	for i := range checkReel {
		if isFromLeft {
			checkReelList[i] = i
		} else {
			checkReelList[i] = len(checkReel) - 1 - i
		}
	}

	checkSymbol := checkReel[checkReelList[0]][lineData[checkReelList[0]]]
	lineSymbolCount, lineWildCount := 0, 0

	for _, col := range checkReelList {
		sym := checkReel[col][lineData[col]]

		if sym == feverSymbolID { // 特殊符号中断
			break
		}

		if slices.Contains(wildList, checkSymbol) {
			if slices.Contains(wildList, sym) {
				lineWildCount++
				lineSymbolCount++
			} else {
				lineSymbolCount++
				checkSymbol = sym
			}
		} else if sym == checkSymbol || slices.Contains(wildList, sym) {
			lineSymbolCount++
		} else {
			break
		}
	}

	// 计算赢分
	var wildWin, symbolWin float64
	if lineWildCount > 0 {
		list := odds[wildSymbolID]
		if lineWildCount-1 < len(list) {
			wildWin = list[lineWildCount-1] * lineBet
		}
	}
	if lineSymbolCount > 0 {
		list := odds[checkSymbol]
		if lineSymbolCount-1 < len(list) {
			symbolWin = list[lineSymbolCount-1] * lineBet
		}
	}

	extraMul := 1.0
	if symbolWin > 0 && len(wildMultiplier) > 0 {
		for col := range checkReel {
			extraMul *= wildMultiplier[col][lineData[col]]
		}
		symbolWin *= extraMul
	}

	win := symbolWin
	if wildWin > win {
		win = wildWin
	}

	if allMultiplier > 1 {
		win *= allMultiplier
		extraMul *= allMultiplier
	}

	if skipFiveSymbol && lineSymbolCount == len(lineData) {
		win = 0
	}

	return win
}
