package game

import "slices"

// game_check_line.go 文件实现了 GameCheckLine 接口，用于检查游戏中的连线情况。
// 核心功能包括：
func CheckLineWin(checkReel [][]int, lineData []int, odds map[int][]float64, wildSymbolID int, wildList []int, feverSymbolID int, lineBet float64, isFromLeft bool,
	skipFiveSymbol bool, wildMultiplier [][]float64, allMultiplier float64) float64 {
	// 方向决定列顺序
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
