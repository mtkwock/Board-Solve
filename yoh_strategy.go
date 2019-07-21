package main

import (
	"fmt"
)

type YohAnalysis struct {
	wood_count int
	priority_five_match []OrbAttribute
	fallback_five_match []OrbAttribute
	priority_extra []OrbAttribute
	fallback_extra []OrbAttribute
	three_matches []OrbAttribute
}

func (self YohAnalysis) String() string {
	return fmt.Sprintf(
		"Woods: %d\nFive-Match: %s\nFallback Five-Match: %s\nExtra Orb: %s\nMaybe Extra: %s\nThrees: %s",
		self.wood_count,
		OrbAttributeArrayToString(self.priority_five_match),
		OrbAttributeArrayToString(self.fallback_five_match),
		OrbAttributeArrayToString(self.priority_extra),
		OrbAttributeArrayToString(self.fallback_extra),
		OrbAttributeArrayToString(self.three_matches))
}

func YohAnalyze(board Board) YohAnalysis {
	// Analyze the attribute counts to determine which boards are viable.
	orb_to_count := board.GetCounts()
	wood_count := 0
	if v, exists := orb_to_count[WOOD]; exists {
		wood_count = v
	}
	priority_five_match := make([]OrbAttribute, 0)
	fallback_five_match := make([]OrbAttribute, 0)
	priority_extra := make([]OrbAttribute, 0)
	fallback_extra := make([]OrbAttribute, 0)
	three_matches := make([]OrbAttribute, 0)
	for attribute, count := range orb_to_count {
		if count % 3 == 1 {
			priority_extra = append(priority_extra, attribute)
		} else if count % 3 == 2 {
			fallback_extra = append(fallback_extra, attribute)
		} else {
			if attribute == WOOD {
				if count >= 9 {
					three_matches = append(three_matches, attribute)
				}
			} else {
				three_matches = append(three_matches, attribute)
			}
		}
		if attribute == WOOD {
			continue
		}
		if count >= 5 {
			if count % 3 == 2 {
				priority_five_match = append(priority_five_match, attribute)
			} else {
				fallback_five_match = append(fallback_five_match, attribute)
			}
		}
	}
	return YohAnalysis {
		wood_count,
		priority_five_match,
		fallback_five_match,
		priority_extra,
		fallback_extra,
		three_matches,
	}
}

// 6x5 Yoh Row Strategies
// TODO: Row + SFua?  VDP + SFua?  VDP + Fua? Fua + Green Blob?
// Possible Yoh strategies:
//  * 5-match 1c
//  * 5-match max c
//  * Row + 5-match (This)
//  * Row + SFua
//  * Green Blob + Fua
//  * VDP + 5-match
//  * VDP + Fua
//  * VDP + SFua
func YohFindSetup(board Board) BoardSetup {
	analysis := YohAnalyze(board)

	potential_board_setups := make([]BoardSetup, 0)

	// This is for wood row strategies, impossible to do.
	if analysis.wood_count < 6 {
		return BoardSetup{}
	}
	five_match_attrs := analysis.priority_five_match
	if len(five_match_attrs) == 0 {
		five_match_attrs = analysis.fallback_five_match
	}
	// Can't match 5.  No point in trying.
	if len(five_match_attrs) == 0 {
		return BoardSetup{}
	}
	for _, five_match_attr := range five_match_attrs {
		// G G G G G G
		// 1 1 1 1 1 .
		// . . . . . .
		// . . . . . .
		// . . . . . .
		board_setup_1 := BoardSetup {
			Combos: []SetupCombo {
				SetupCombo {
					WOOD,
					[]Pair {Pair{0, 0}, Pair{0, 1}, Pair{0, 2}, Pair{0, 3}, Pair{0, 4}, Pair{0, 5}},
				},
				SetupCombo {
					five_match_attr,
					[]Pair {Pair{1, 0}, Pair{1, 1}, Pair{1, 2}, Pair{1, 3}, Pair{1, 4}},
				},
			},
		}
		board_setup_1.Init(board.Width)

		// G G G G G G
		// . . . . . .
		// . . . . . .
		// . . . . . .
		// 1 1 1 1 1 .
		board_setup_2 := BoardSetup {
			Combos: []SetupCombo {
				SetupCombo {
					WOOD,
					[]Pair {Pair{0, 0}, Pair{0, 1}, Pair{0, 2}, Pair{0, 3}, Pair{0, 4}, Pair{0, 5}},
				},
				SetupCombo {
					five_match_attr,
					[]Pair {Pair{4, 0}, Pair{4, 1}, Pair{4, 2}, Pair{4, 3}, Pair{4, 4}},
				},
			},
		}
		board_setup_2.Init(board.Width)

		potential_board_setups = append(potential_board_setups, board_setup_1, board_setup_2)

		for _, extra_match_attr := range analysis.priority_extra {
			if extra_match_attr == five_match_attr {
				continue
			}
			// G G G G G G
			// . . . . . .
			// . . . . . .
			// . . . . . .
			// 1 1 1 1 1 2
			board_setup_2b := board_setup_2.Clone()
			board_setup_2b.Combos = []SetupCombo{
				board_setup_2b.Combos[0],
				SetupCombo{extra_match_attr, []Pair{Pair{4, 5}}},
				board_setup_2b.Combos[1],
			}
			board_setup_2b.Init(board.Width)

			// 1 1 1 1 1 2
			// G G G G G G
			// . . . . . .
			// . . . . . .
			// . . . . . .
			board_setup_3 := BoardSetup {
				Combos: []SetupCombo {
					SetupCombo {
						five_match_attr,
						[]Pair {Pair{0, 0}, Pair{0, 1}, Pair{0, 2}, Pair{0, 3}, Pair{0, 4}},
					},
					SetupCombo {
						extra_match_attr,
						[]Pair{Pair{0, 5}},
					},
					SetupCombo {
						WOOD,
						[]Pair {
							Pair{1, 0}, Pair{1, 1}, Pair{1, 2}, Pair{1, 3}, Pair{1, 4}, Pair{1, 5},
						},
					},
				},
			}
			board_setup_3.Init(board.Width)
			potential_board_setups = append(potential_board_setups, board_setup_2b, board_setup_3)
		}
		// G G G G G G
		// 1 1 1 . . .
		// 2 3 1 . . .
		// 2 3 1 . . .
		// 2 3 . . . .
		if len(analysis.three_matches) >= 2 {
			for i := 0; i < len(analysis.three_matches) - 1; i++ {
				off_color_1 := analysis.three_matches[i]
				for j := i + 1; j < len(analysis.three_matches); j++ {
					off_color_2 := analysis.three_matches[j]
					board_setup_4a := BoardSetup {
						Combos: []SetupCombo {
							SetupCombo {
								WOOD,
								[]Pair{Pair{0, 0}, Pair{0, 1}, Pair{0, 2}, Pair{0, 3}, Pair{0, 4}, Pair{0, 5}},
							},
							SetupCombo {
								five_match_attr,
								[]Pair{Pair{1, 0}, Pair{1, 1}, Pair{1, 2}},
							},
							SetupCombo {
								off_color_1,
								[]Pair{Pair{2, 0}, Pair{3, 0}, Pair{4, 0}},
							},
							SetupCombo {
								off_color_2,
								[]Pair{Pair{2, 1}, Pair{3, 1}, Pair{4, 1}},
							},
							SetupCombo {
								five_match_attr,
								[]Pair{Pair{2, 2}, Pair{3, 2}},
							},
						},
					}
					board_setup_4a.Init(board.Width)
					board_setup_4b := board_setup_4a.Clone()
					board_setup_4b.Combos[2].Attribute = off_color_2
					board_setup_4b.Combos[3].Attribute = off_color_1
					board_setup_4b.Init(board.Width)
					potential_board_setups = append(potential_board_setups, board_setup_4a, board_setup_4b)
				}
			}
		}
	}

	best_setup := BoardSetup{}
	var best_average float32 = 100.0 // Worse than any possible case.

	for _, setup := range potential_board_setups {
		distance := setup.ManhattanDistanceAverage(board)
		// fmt.Printf("%s, %f\n", setup, distance)
		if distance < best_average {
			best_setup = setup
			best_average = distance
		}
		horizontal_flip := setup.MirrorHorizontal(board.Width)
		distance = horizontal_flip.ManhattanDistanceAverage(board)
		// fmt.Printf("%s, %f\n", horizontal_flip, distance)
		if distance < best_average {
			best_setup = horizontal_flip
			best_average = distance
		}
		horizontal_vertical_flip := horizontal_flip.MirrorVertical(board.Height)
		distance = horizontal_vertical_flip.ManhattanDistanceAverage(board)
		// fmt.Printf("%s, %f\n", horizontal_vertical_flip, distance)
		if distance < best_average {
			best_setup = horizontal_vertical_flip
			best_average = distance
		}

		vertical_flip := setup.MirrorVertical(board.Height)
		distance = vertical_flip.ManhattanDistanceAverage(board)
		// fmt.Printf("%s, %f\n", vertical_flip, distance)
		if distance < best_average {
			best_setup = vertical_flip
			best_average = distance
		}
	}

	return best_setup
}
