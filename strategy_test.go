package main

import (
	"fmt"
	"testing"
)

// G G G G G G
// L L L L L H
// R B G L D H
// R B G L D H
// R B G L D H
var nice_yoh_board Board = CreateBoard("GGGGGGLLLLLHRBGLDHRBGLDHRBGLDH", 6)

func TestYohAnalyze_PerfectBoard_ReturnsValues(t *testing.T) {
	analysis := YohAnalyze(nice_yoh_board)

	if analysis.wood_count != 9 {
		fmt.Println(analysis.wood_count)
		panic("Should have 9 wood orbs.")
	}
	if len(analysis.priority_five_match) != 1 || analysis.priority_five_match[0] != LIGHT {
		fmt.Println(analysis.priority_five_match)
		panic("Light should be five matched priority.")
	}
	if len(analysis.fallback_five_match) != 0 {
		fmt.Println(analysis.fallback_five_match)
		panic("Shouldn't have fallback five match set.")
	}
	if len(analysis.priority_extra) != 1 {
		fmt.Println(analysis.priority_extra)
		panic("Should have one priority extra.")
	}
	if len(analysis.three_matches) != 4 {
		fmt.Println(analysis.three_matches)
		panic("Should have 4x three_matches, Fire, Water, Wood, Dark")
	}
}

func TestBoardSetup_StringCall_GetsReasonableString(t *testing.T) {
	board_setup := BoardSetup {
		Combos: []SetupCombo {
			SetupCombo {
				WOOD,
				[]Pair {Pair{0, 0}, Pair{0, 1}, Pair{0, 2}, Pair{0, 3}, Pair{0, 4}, Pair{0, 5}},
			},
			SetupCombo {
				LIGHT,
				[]Pair {Pair{1, 0}, Pair{1, 1}, Pair{1, 2}, Pair{1, 3}, Pair{1, 4}},
			},
		},
	}
	board_setup.Init(6)
	g_count := 0
	l_count := 0

	setup_string := board_setup.String()
	for i := 0; i < len(setup_string); i++ {
		if setup_string[i] == 'G' {
			g_count++
		}
		if setup_string[i] == 'L' {
			l_count++
		}
	}

  if g_count != 6 || l_count != 5 {
		panic(setup_string)
	}

	g_count = 0
	l_count = 0
	setup_string = board_setup.MirrorHorizontal(6).String()
	for i := 0; i < len(setup_string); i++ {
		if setup_string[i] == 'G' {
			g_count++
		}
		if setup_string[i] == 'L' {
			l_count++
		}
	}
	if g_count != 6 || l_count != 5 {
		panic(setup_string)
	}
}

func TestManhattanDistance_AlreadyInPlace_Returns0(t *testing.T) {
	board_setup := BoardSetup {
		Combos: []SetupCombo {
			SetupCombo {
				WOOD,
				[]Pair {Pair{0, 0}, Pair{0, 1}, Pair{0, 2}, Pair{0, 3}, Pair{0, 4}, Pair{0, 5}},
			},
			SetupCombo {
				LIGHT,
				[]Pair {Pair{1, 0}, Pair{1, 1}, Pair{1, 2}, Pair{1, 3}, Pair{1, 4}},
			},
		},
	}
	board_setup.Init(6)

	distance := board_setup.ManhattanDistanceGreedyEdges(nice_yoh_board)

	if distance > 0 {
		panic("Distance should be 0!")
	}
}

func TestManhattanDistance_OneShifted_Returns6(t *testing.T) {
	board_setup := BoardSetup {
		Combos: []SetupCombo {
			SetupCombo {
				WOOD,
				[]Pair {Pair{1, 0}, Pair{1, 1}, Pair{1, 2}, Pair{1, 3}, Pair{1, 4}, Pair{1, 5}},
			},
		},
	}

	distance := board_setup.ManhattanDistanceGreedyEdges(nice_yoh_board)

	if distance != 6 {
		panic("Distance should be 6.")
	}
}

func TestUnusedOrbIdxs(t *testing.T) {
	board_setup := BoardSetup {
		Combos: []SetupCombo {
			SetupCombo {
				WOOD,
				[]Pair {Pair{0, 0}, Pair{0, 1}, Pair{0, 2}, Pair{0, 3}, Pair{0, 4}, Pair{0, 5}},
			},
			SetupCombo {
				LIGHT,
				[]Pair {Pair{1, 0}, Pair{1, 1}, Pair{1, 2}, Pair{1, 3}, Pair{1, 4}},
			},
		},
	}
	board_setup.Init(6)
	board := CreateBoard("GGGGGGLLLLLLRRRRRRBBBBBBDDDDDD", 6)

	unused_idx := board_setup.UnusedOrbIdxs(board)

  if len(unused_idx) != 24 {
		fmt.Println(len(unused_idx))
		s := ""
		for _, idx := range unused_idx {
			s += idx.String() + ", "
		}
		panic(s[:len(s) - 2])
	}
}
