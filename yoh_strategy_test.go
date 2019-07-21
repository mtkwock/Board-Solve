package main

import (
	"fmt"
	"testing"
)

func BenchmarkMakeYohRowStrategy_WithPerfectMatch(b *testing.B) {
	YohFindSetup(nice_yoh_board.Clone())
}

func TestYohBoard_RandomBoard(t *testing.T) {
	// Board found by entering Endless Corridors.
	// G H D B D L
	// D G B L G G
	// H L H L R G
	// L D R H L R
	// G L R L B B
	board := CreateBoard("GHDBDLDGBLGGHLHLRGLDRHLRGLRLBB", 6)

	_ = YohFindSetup(board)

	// fmt.Println(YohAnalyze(board))
	// panic(setup.String())
}

func TestYohRowSolve_RandomBoard(t *testing.T) {
	// Board found by entering Endless Corridors.
	// G H D B D L
	// D G B L G G
	// H L H L R G
	// L D R H L R
	// G L R L B B
	board := CreateBoard("GHDBDLDGBLGGHLHLRGLDRHLRGLRLBB", 6)
	setup := YohFindSetup(board)

	moves := StrategySolve(board, setup, SolveRequirement{
		AllowDiagonals: false,
		FinishedFn: func(state AStarState) bool {
			return len(state.board.GetAllCombos()) >= 7
		},
		ScoreState: func(state AStarState) int {
			return len(state.board.GetAllCombos()) * 17 - len(state.moves)
		},
		RejectionFn: MakeRejectionFunction(),
	})
	fmt.Println(moves)
	fmt.Println(len(moves.Directions))
	fmt.Println(ToDawnglare(board, moves))
	panic("Let's see what we have here...")

	// fmt.Println(YohAnalyze(board))
	// panic(setup.String())
}
