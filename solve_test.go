package main
//
// import (
// 	"fmt"
// 	"testing"
// )
//
// func TestBfsBasic(t *testing.T) {
// 	board := CreateBoard("GGGGGGGGGGGGGGGGGGGGGGGGGGGGGG", 6)
// 	requirement := SolveRequirement{1}
//
// 	solver := AStarSolve{}
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
//
// func TestRandomBoard(t *testing.T) {
// 	board := CreateBoard("BRHRDBBDDBHLBHGLDBHRRBBDRLLRRL", 6)
//
// 	fmt.Println(board)
// 	requirement := SolveRequirement{7}
//
// 	moves := BfsFourDirectionSolver{}.Solve(board, requirement)
//
// 	t.Errorf("%s\n", moves)
// }
//
// func TestNextState_CenterOfBoard_ReturnsFourDirections(t *testing.T) {
// 	board := CreateEmptyBoard()
// 	placement := Placement{0,0}
// 	direction := DOWN
// 	board.Slots[placement.ToPos(board)].Orb.Attribute = FIRE
// 	new_board, _ := board.Swap(placement, direction)
// 	state := BfsState{
// 		new_board,
// 		placement,
// 		placement.Swap(direction),
// 		[]Direction{direction},
// 		nil,
// 	}
//
// 	next_states := state.NextStates()
// 	if len(next_states) != 2 {
// 		t.Error("Expecting 2 states to be generated.")
// 	}
// 	if next_states[0].moves[1] != RIGHT || next_states[1].moves[1] != DOWN {
// 		t.Error("Expected the next two movements to be right and down.")
// 	}
// }
