package main

import (
	"fmt"
	"testing"
)

// func TestBfsBasic(t *testing.T) {
// 	board := CreateBoard("GGGGGGGGGGGGGGGGGGGGGGGGGGGGGG", 6)
// 	requirement := SolveRequirement{1}
//
// 	solver := BfsFourDirectionSolver{}
//
// 	moves := solver.Solve(board, requirement)
//
// 	if len(moves.Directions) > 1 {
// 		t.Errorf("Starting: %s, Moves: %d", moves.StartingPosition, len(moves.Directions))
// 	}
// }
//
// func TestBfsOneMove(t *testing.T) {
// 	board := CreateBoard("GGGGGGGGGGGGGGGGGGGGGGRRGGGRGG", 6)
// 	// fmt.Println(board)
// 	requirement := SolveRequirement{2}
// 	solver := BfsFourDirectionSolver{}
//
// 	moves := solver.Solve(board, requirement)
//
// 	if len(moves.Directions) > 1 {
// 		t.Errorf("Starting: %s, Moves: %d", moves.StartingPosition, len(moves.Directions))
// 	}
// }
//
// var one_move_board Board = CreateBoard("GGGGGGGGGGGGGGGGGGGGGGRRGGGRGG", 6)
// var two_combo_requirement SolveRequirement = SolveRequirement{2}
//
// func BenchmarkBfsOneMove(b *testing.B) {
// 	BfsFourDirectionSolver{}.Solve(one_move_board, two_combo_requirement)
// }

func TestRandomBoard(t *testing.T) {
	board := CreateBoard("BRHRDBBDDBHLBHGLDBHRRBBDRLLRRL", 6)

	fmt.Println(board)
	requirement := SolveRequirement{4}

	moves := BfsFourDirectionSolver{}.Solve(board, requirement)

	t.Errorf("%s\n", moves)
}
