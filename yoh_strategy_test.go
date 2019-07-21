package main

import (
	"fmt"
	"testing"
)

func BenchmarkMakeYohRowStrategy_WithPerfectMatch(b *testing.B) {
	MakeYohRowStrategy(nice_yoh_board.Clone())
}

func TestYohBoard_RandomBoard(t *testing.T) {
	// Board found by entering Endless Corridors.
	// G H D B D L
	// D G B L G G
	// H L H L R G
	// L D R H L R
	// G L R L B B
	board := CreateBoard("GHDBDLDGBLGGHLHLRGLDRHLRGLRLBB", 6)

	setup := MakeYohRowStrategy(board)

	fmt.Println(YohAnalyze(board))
	panic(setup.String())
}
