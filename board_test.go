package main

import (
	// "errors"
	// "fmt"
	// "unsafe"
	// "strings"
	// "math/rand" // Runs deterministically unless we set a seed.
	"testing"
)

// func TestPrinting(t *testing.T) {
// 	slots := make([]BoardSpace, 30)
// 	for i := 0; i < len(slots); i++ {
// 		slots[i].Orb.Attribute = NormalOrbs[rand.Intn(len(NormalOrbs))]
// 		// slots[i].Orb.State |= BLIND
// 		if i % 6 == 0 {
// 			slots[i].State |= TAPE
// 		}
// 		if i >= 24 {
// 			slots[i].State |= CLOUD
// 		}
// 		if i % 6 == 5 {
// 			slots[i].Orb.State |= BLIND
// 		}
// 		if i >= 18 && i < 24 {
// 			slots[i].Orb.State |= ENHANCED
// 		}
// 		if i % 6 == 1 {
// 			slots[i].Orb.State |= LOCKED
// 		}
// 	}
// 	board := Board{slots, 5, 6}
// 	board.Print()
// 	fmt.Printf("Size of board: %d\n", unsafe.Sizeof(board))
//
// 	fmt.Println()
// 	board_swapped, err := board.Swap(Placement{1, 1}, RIGHT)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Printf("Old Board\n%s\nNew Board\n%s", board, board_swapped)
// 	}
// }
//
func TestDios(t *testing.T) {
	dios_box_board := CreateEmptyBoard(6)
	for i := 0; i < len(dios_box_board.Slots); i++ {
		dios_box_board.Slots[i].Orb.Attribute = WOOD
	}
	// G R R R G G
	// G R G R G G
	// G R R R G G
	// G G G G G G
	// G G G G G G
	dios_box_board.Slots[1].Orb.Attribute = FIRE
	dios_box_board.Slots[2].Orb.Attribute = FIRE
	dios_box_board.Slots[3].Orb.Attribute = FIRE
	dios_box_board.Slots[7].Orb.Attribute = FIRE
	dios_box_board.Slots[9].Orb.Attribute = FIRE
	dios_box_board.Slots[13].Orb.Attribute = FIRE
	dios_box_board.Slots[14].Orb.Attribute = FIRE
	dios_box_board.Slots[15].Orb.Attribute = FIRE
	// fmt.Println(dios_box_board)
	combos, _ := dios_box_board.GetCombos()
	if len(combos) != 2 {
		t.Error("This should have two combos.")
	}
	if len(combos[0].Positions) != 21 || combos[0].Attribute != WOOD {
		combos[0].Print(dios_box_board.Width)
		t.Errorf("The first combo should be wood with 21 orbs: %s", combos[0])
	}
	if len(combos[1].Positions) != 8 || combos[1].Attribute != FIRE {
		t.Errorf("The second combo should be fire with 8 orbs: %s", combos[1])
	}
}

var swng_board Board = CreateBoard("RRHHHDGDGGGDGDDRRRRDRGGDGRRDDG", 6)

func TestSwngBoard(t *testing.T) {
	swng_board := CreateBoard("RRHHHDGDGGGDGDDRRRRDRGGDGRRDDG", 6)
	combos := swng_board.GetAllCombos()

	if len(combos) != 10 {
		t.Error("This should be a 10c board.")
	}
}

func BenchmarkSwngBoard(b *testing.B) {
	swng_board.GetAllCombos()
}

func TestBombDios(t *testing.T) {
	// Should become:
	// G G G - G -
	// G G G - G -
	// G G G - G -
	// - - - o - o
	// - - - o - -
	board := CreateBoard("GGGGGGGGGGGGGGGGGGGGGoGoGGGoGG", 6)

	combos, new_board := board.GetCombos()
	if len(combos) != 2 {
		t.Errorf("This should have two combos:\n%s\n%s", new_board, combos)
	}
	if len(combos[0].Positions) != 9 {
		t.Errorf("This combo should be a VDP\n%s", combos[0])
	}
	if len(combos[1].Positions) != 3 {
		t.Errorf("The second combo should be a single 3-match\n%s", combos[1])
	}
}

func TestBicolor(t *testing.T) {
	board := CreateBoard("LLLRLLRRLLRLLRRRLLLLRLLRLLLRRR", 6)

	combos := board.GetAllCombos()

	// for _, combo := range combos {
	// 	combo.Print(6)
	// }

	if len(combos) != 6 {
		t.Errorf("This should be a 6 combo board got %d", len(combos))
	}
}
