package main

import (
	"errors"
	"fmt"
	"unsafe"
	"strings"
	"math/rand" // Runs deterministically unless we set a seed.
)

type OrbAttribute uint8

const (
	EMPTY OrbAttribute = iota
	FIRE
	WATER
	WOOD
	LIGHT
	DARK
	HEART
	JAMMER
	POISON
	MORTAL_POISON
	BOMB
)

var AttributeToName map[OrbAttribute]string = map[OrbAttribute]string{
	EMPTY: "EMPTY",
	FIRE: "Fire",
	WATER: "Water",
	WOOD: "Wood",
	LIGHT: "Light",
	DARK: "Dark",
	HEART: "Heart",
	JAMMER: "Jammer",
	POISON: "Poison",
	MORTAL_POISON: "Mortal Poison",
	BOMB: "Bomb",
}

var AttributeToLetter map[OrbAttribute]string = map[OrbAttribute]string{
	EMPTY: " ",
	FIRE: "R",
	WATER: "B",
	WOOD: "G",
	LIGHT: "L",
	DARK: "D",
	HEART: "H",
	JAMMER: "J",
	POISON: "P",
	MORTAL_POISON: "M",
	BOMB: "o",
}

var NormalOrbs [6]OrbAttribute = [6]OrbAttribute{
	FIRE, WATER, WOOD, LIGHT, DARK, HEART}
var HazardOrbs [4]OrbAttribute = [4]OrbAttribute{
	JAMMER, POISON, MORTAL_POISON, BOMB}

type OrbStateFlag uint8

const (
  ENHANCED OrbStateFlag = 1 << iota
	LOCKED
	BLIND
	STICKY_BLIND
	UNMATCHABLE
)

type Orb struct {
	Attribute OrbAttribute
	State OrbStateFlag
}

func (self Orb) Clone() Orb {
	return Orb{self.Attribute, self.State}
}

func (self Orb) IsBlinded() bool {
	return self.State & (BLIND | STICKY_BLIND) != 0
}

func (self Orb) String() string {
	locked := " "
	if self.State & LOCKED != 0 {
		locked = "&"
	}
	char := AttributeToLetter[self.Attribute]
	if self.IsBlinded() {
		char = "?"
	}
  enhanced := " "
	if self.State & ENHANCED != 0 {
		enhanced = "+"
	}
	return locked + char + enhanced
}

type BoardSpaceStateFlag uint8

const (
	TAPE BoardSpaceStateFlag = 1 << iota
	CLOUD
	SPINNER_1S
	SPINNER_2S
)

type BoardSpace struct {
	Orb Orb
	State BoardSpaceStateFlag
}

func (self BoardSpace) Clone() BoardSpace {
	return BoardSpace{self.Orb.Clone(), self.State}
}

func (self BoardSpace) String() string {
	if self.State & CLOUD != 0 {
		return "{%}"
	}
	if self.State & TAPE != 0 {
		return fmt.Sprintf("<%s>", string(self.Orb.String()[1]))
	}
	return self.Orb.String()
}

type Board struct {
	Slots []BoardSpace
	Height uint8
	Width uint8
}

func (self Board) Clone() Board {
	new_slots := make([]BoardSpace, len(self.Slots))
	for i, board_space := range self.Slots {
		new_slots[i] = board_space.Clone()
	}
	return Board{new_slots, self.Height, self.Width}
}

func (self Board) String() string {
	border := "++" + strings.Repeat("-", 3 * int(self.Width)) + "++\n"
	body := ""
	for y := uint8(0); y < self.Height; y++ {
		body += "||"
		for x := uint8(0); x < self.Width; x++ {
			body += self.Slots[Placement{y, x}.ToPos(self)].String()
		}
		body += "||\n"
	}
	return fmt.Sprintf("%s%s%s", border, body, border)
}

func (self Board) Print() {
	fmt.Println(self)
}

type Direction uint8

const (
	_ = iota
	UP Direction = iota
	UP_RIGHT
	RIGHT
	DOWN_RIGHT
	DOWN
	DOWN_LEFT
	LEFT
	UP_LEFT
)

type Placement struct {
	Y uint8
	X uint8
}

func (self Placement) ToPos(board Board) uint8 {
	return self.Y * board.Width + self.X
}

func (self Placement) String() string {
	return fmt.Sprintf("(%d,%d)", self.X, self.Y)
}

func (self Board) Swap(placement Placement, direction Direction) (Board, error) {
	new_placement := Placement{placement.Y, placement.X}
	// TODO: Consider changing this to a hashmap of placement, direction: placement.
	// Will need to profile if faster or not.
	switch direction {
	case UP:
		new_placement.Y--
  case UP_RIGHT:
		new_placement.Y--
		new_placement.X++
	case RIGHT:
		new_placement.X++
	case DOWN_RIGHT:
		new_placement.Y++
		new_placement.X++
	case DOWN:
		new_placement.Y++
	case DOWN_LEFT:
		new_placement.Y++
		new_placement.X--
	case LEFT:
		new_placement.X--
	case UP_LEFT:
		new_placement.Y--
		new_placement.X--
	}

	new_board := self.Clone()

	if new_placement.X >= self.Width || new_placement.Y >= self.Height {
		return new_board, errors.New("Invalid move. Off board.")
	}
	new_pos := new_placement.ToPos(self)
	old_pos := placement.ToPos(self)

  if self.Slots[new_pos].State & TAPE != 0 {
		return new_board, errors.New("Invalid move. New Position is Immobilized.")
	}
	if self.Slots[old_pos].State & TAPE != 0 {
		return new_board, errors.New("Invalid move. Original Position is Immobilized.")
	}
	temp := new_board.Slots[new_pos].Orb.Clone()
	new_board.Slots[new_pos].Orb = new_board.Slots[old_pos].Orb.Clone()
	new_board.Slots[old_pos].Orb = temp
	new_board.Slots[new_pos].Orb.State &= ^BLIND
	new_board.Slots[old_pos].Orb.State &= ^BLIND
	return new_board, nil
}

type ComboType uint8

const (
	MATCH_3 ComboType = iota
	MATCH_TPA
  MATCH_CROSS
	MATCH_L
	MATCH_COLUMN
	MATCH_VDP
)

type BoardCombo struct {
	Attribute OrbAttribute
	Positions []Placement
	IsEnhanced bool
}

func (self BoardCombo) String() string {
	return fmt.Sprintf("%s: [%d]", AttributeToName[self.Attribute], len(self.Positions))
}

func (self BoardCombo) Print() {
	fmt.Println(self)
	board := CreateEmptyBoard()
	for _, placement := range self.Positions {
		board.Slots[placement.ToPos(board)].Orb.Attribute = self.Attribute
	}
	fmt.Println(board)
}

type BoardCombos map[uint8]uint8

// type BoardRestriction uint16
//
// const (
// 	RESTRICT_4PLUS BoardRestriction = 1 << iota
// 	RESTRICT_5PLUS
// 	RESTRICT_FIRE
// 	RESTRICT_WATER
// 	RESTRICT_LIGHT
// 	RESTRICT_DARK
// 	RESTRICT_HEART
// 	RESTRICT_POISON
// )

// var potential_combos_6x5 = [...][3]Placement{
// 	// Horizontal top row.
// 	[3]Placement{Placement{0, 0}, Placement{0, 1}, Placement{0, 2}},
// 	[3]Placement{Placement{0, 1}, Placement{0, 2}, Placement{0, 3}},
// 	[3]Placement{Placement{0, 2}, Placement{0, 3}, Placement{0, 4}},
// 	[3]Placement{Placement{0, 3}, Placement{0, 4}, Placement{0, 5}},
// 	// Horizontal 2nd row.
// 	[3]Placement{Placement{1, 0}, Placement{1, 1}, Placement{1, 2}},
// 	[3]Placement{Placement{1, 1}, Placement{1, 2}, Placement{1, 3}},
// 	[3]Placement{Placement{1, 2}, Placement{1, 3}, Placement{1, 4}},
// 	[3]Placement{Placement{1, 3}, Placement{1, 4}, Placement{1, 5}},
// 	// Horizontal 3rd row.
// 	[3]Placement{Placement{2, 0}, Placement{2, 1}, Placement{2, 2}},
// 	[3]Placement{Placement{2, 1}, Placement{2, 2}, Placement{2, 3}},
// 	[3]Placement{Placement{2, 2}, Placement{2, 3}, Placement{2, 4}},
// 	[3]Placement{Placement{2, 3}, Placement{2, 4}, Placement{2, 5}},
// 	// Horizontal 4th row.
// 	[3]Placement{Placement{3, 0}, Placement{3, 1}, Placement{3, 2}},
// 	[3]Placement{Placement{3, 1}, Placement{3, 2}, Placement{3, 3}},
// 	[3]Placement{Placement{3, 2}, Placement{3, 3}, Placement{3, 4}},
// 	[3]Placement{Placement{3, 3}, Placement{3, 4}, Placement{3, 5}},
// 	// Horizontal bottom row.
// 	[3]Placement{Placement{4, 0}, Placement{4, 1}, Placement{4, 2}},
// 	[3]Placement{Placement{4, 1}, Placement{4, 2}, Placement{4, 3}},
// 	[3]Placement{Placement{4, 2}, Placement{4, 3}, Placement{4, 4}},
// 	[3]Placement{Placement{4, 3}, Placement{4, 4}, Placement{4, 5}},
//
// 	// Vertical left column
// 	[3]Placement{Placement{0, 0}, Placement{1, 0}, Placement{2, 0}},
// 	[3]Placement{Placement{1, 0}, Placement{2, 0}, Placement{3, 0}},
// 	[3]Placement{Placement{2, 0}, Placement{3, 0}, Placement{4, 0}},
// 	// Vertical 2nd column
// 	[3]Placement{Placement{0, 1}, Placement{1, 1}, Placement{2, 1}},
// 	[3]Placement{Placement{1, 1}, Placement{2, 1}, Placement{3, 1}},
// 	[3]Placement{Placement{2, 1}, Placement{3, 1}, Placement{4, 1}},
// 	// Vertical 3rd column
// 	[3]Placement{Placement{0, 2}, Placement{1, 2}, Placement{2, 2}},
// 	[3]Placement{Placement{1, 2}, Placement{2, 2}, Placement{3, 2}},
// 	[3]Placement{Placement{2, 2}, Placement{3, 2}, Placement{4, 2}},
// 	// Vertical 4th column
// 	[3]Placement{Placement{0, 3}, Placement{1, 3}, Placement{2, 3}},
// 	[3]Placement{Placement{1, 3}, Placement{2, 3}, Placement{3, 3}},
// 	[3]Placement{Placement{2, 3}, Placement{3, 3}, Placement{4, 3}},
// 	// Vertical 5th column
// 	[3]Placement{Placement{0, 4}, Placement{1, 4}, Placement{2, 4}},
// 	[3]Placement{Placement{1, 4}, Placement{2, 4}, Placement{3, 4}},
// 	[3]Placement{Placement{2, 4}, Placement{3, 4}, Placement{4, 4}},
// 	// Vertical right column
// 	[3]Placement{Placement{0, 5}, Placement{1, 5}, Placement{2, 5}},
// 	[3]Placement{Placement{1, 5}, Placement{2, 5}, Placement{3, 5}},
// 	[3]Placement{Placement{2, 5}, Placement{3, 5}, Placement{4, 5}},
// }
//
// touching_combos = map[int]([]int) {
// 	0: [0]int{}
// 	1: [1]int{0},
// 	2: [2]int{1, 0},
// 	3: [3]int{2, 1, 0},
// 	4: [3]int{2, 1, 0},
// 	5: [5]int{4, 3, 2, 1, 0},
// 	6: [6]int{5, 4, 3, 2, 1, 0},
// 	7: [6]int{6, 5, 4, 3, 2, 1},
// }

// Ideas:
// Simple iteration
// 1. Iterate through board y and x-2, finding all three-match groups horizontally
// 2 .Iterate through board y-2 and x, finding all three-match groups vertically
// 3. Group these up iteratively by finding combos that border or intersect.

// 0.1. Build a lookup table of all possible combos
// 0.2. Build a lookup table of all combos that map to each other.
// 1. Using (0.1)'s table,

// 1) Iterate through all orbs and find if they're going to combo out.
// 2) Iterate through all comboed orbs and blob together orbs that are same
//    color and comboing.

var EmptyOrb Orb = Orb{EMPTY, 0}

func (self Board) GetOrbAt(placement Placement) Orb {
	// Guard Clause
	if placement.Y >= self.Height || placement.X >= self.Width {
		return EmptyOrb
	}
	return self.Slots[placement.ToPos(self)].Orb
}

// TODO: Add BoardRestriction capabilities.
func (self Board) GetCombos() []BoardCombo {
  marked_combos := make([]bool, len(self.Slots))
	horizontal_max := self.Width - 2
	vertical_max := self.Height - 2
  for y := uint8(0); y < self.Height; y++ {
		for x := uint8(0); x < self.Width; x++ {
			placement := Placement{y, x}
			position := placement.ToPos(self)
			// Already known to be matchable, ignore.
			if marked_combos[position] {
				continue
			}

			// Unmatchable orb, ignore
			orb := self.GetOrbAt(placement)
			if orb.Attribute == EMPTY || orb.State & UNMATCHABLE != 0 {
				continue
			}

			// Determine matches to the right.
			if (x < horizontal_max) {
        orb_next := self.GetOrbAt(Placement{placement.Y, placement.X + 1})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Placement{placement.Y, placement.X + 2})
					if orb.Attribute == orb_next_next.Attribute {
						marked_combos[position] = true
						marked_combos[position + 1] = true
						marked_combos[position + 2] = true
						continue
					}
				}
			}

			// Determine matches to the left.
			if (x > 2) {
				orb_next := self.GetOrbAt(Placement{placement.Y, placement.X - 1})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Placement{placement.Y, placement.X - 2})
					if orb.Attribute == orb_next_next.Attribute {
						marked_combos[position] = true
						marked_combos[position - 1] = true
						marked_combos[position - 2] = true
						continue
					}
				}
			}

			// Determine matches below.
			if (y < vertical_max) {
				orb_next := self.GetOrbAt(Placement{placement.Y + 1, placement.X})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Placement{placement.Y + 2, placement.X})
					if orb.Attribute == orb_next_next.Attribute {
						marked_combos[position] = true
						marked_combos[position + self.Width] = true
						marked_combos[position + 2 * self.Width] = true
						continue
					}
				}
			}

			// Determine matches above.
			if (y > 2) {
				orb_next := self.GetOrbAt(Placement{placement.Y - 1, placement.X})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Placement{placement.Y - 2, placement.X})
					if orb.Attribute == orb_next_next.Attribute {
						marked_combos[position] = true
						marked_combos[position - self.Width] = true
						marked_combos[position - 2 * self.Width] = true
						continue
					}
				}
			}

			if (x > 0 && x < self.Width - 1) {
				orb_next := self.GetOrbAt(Placement{placement.Y, placement.X - 1})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Placement{placement.Y, placement.X + 1})
					if orb.Attribute == orb_next_next.Attribute {
						marked_combos[position] = true
						continue
					}
				}
			}

			if (y > 0 && y < self.Height - 1) {
				orb_next := self.GetOrbAt(Placement{placement.Y - 1, placement.X})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Placement{placement.Y + 1, placement.X})
					if orb.Attribute == orb_next_next.Attribute {
						marked_combos[position] = true
						continue
					}
				}
			}
		}
	}

	is_used := make([]bool, len(self.Slots))
	combos := make([]BoardCombo, 0)

	for i, is_comboed := range marked_combos {
		if !is_comboed || is_used[i] {
			continue
		}
		is_used[i] = true
		placement := Placement{uint8(i) / self.Width, uint8(i) % self.Width}
		attribute := self.GetOrbAt(placement).Attribute
		placements := []Placement{placement}
		for j := 0; j < len(placements); j++ {
			current := placements[j]
			pos := current.ToPos(self) + 1
			right := Placement{current.Y, current.X + 1}
			if self.GetOrbAt(right).Attribute == attribute && marked_combos[pos] && !is_used[pos] {
				placements = append(placements, right)
				is_used[pos] = true
			}
			pos -= 2
			left := Placement{current.Y, current.X - 1}
			if self.GetOrbAt(left).Attribute == attribute && marked_combos[pos] && !is_used[pos]{
				placements = append(placements, left)
				is_used[pos] = true
			}
			pos += self.Width + 1
			down := Placement{current.Y + 1, current.X}
			if self.GetOrbAt(down).Attribute == attribute && marked_combos[pos] && !is_used[pos]{
				placements = append(placements, down)
				is_used[pos] = true
			}
			pos -= 2 * self.Width
			up := Placement{current.Y - 1, current.X}
			if self.GetOrbAt(up).Attribute == attribute && marked_combos[pos] && !is_used[pos]{
				placements = append(placements, up)
				is_used[pos] = true
			}
		}
		combo := BoardCombo{attribute, placements, false}
		// TODO: Mark the combos with special values.
		combos = append(combos, combo)
	}

	return combos
}

func CreateEmptyBoard() Board {
	return Board{make([]BoardSpace, 30), 5, 6}
}

func main() {
	// fmt.Println("Hello World")

	// for i := 0; i < len(NormalOrbs); i++ {
	// 	fmt.Printf("%d: %s\n", NormalOrbs[i], AttributeToName[NormalOrbs[i]])
	// }

	// fmt.Printf("Size of NormalOrbs: %d\n", unsafe.Sizeof(NormalOrbs))
	slots := make([]BoardSpace, 30)
	for i := 0; i < len(slots); i++ {
		slots[i].Orb.Attribute = NormalOrbs[rand.Intn(len(NormalOrbs))]
		// slots[i].Orb.State |= BLIND
		if i % 6 == 0 {
			slots[i].State |= TAPE
		}
		if i >= 24 {
			slots[i].State |= CLOUD
		}
		if i % 6 == 5 {
			slots[i].Orb.State |= BLIND
		}
		if i >= 18 && i < 24 {
			slots[i].Orb.State |= ENHANCED
		}
		if i % 6 == 1 {
			slots[i].Orb.State |= LOCKED
		}
	}
	board := Board{slots, 5, 6}
	board.Print()
	fmt.Printf("Size of board: %d\n", unsafe.Sizeof(board))

	fmt.Println()
	board_swapped, err := board.Swap(Placement{1, 1}, RIGHT)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Old Board\n%s\nNew Board\n%s", board, board_swapped)
	}

	dios_board := CreateEmptyBoard()
	for i := 0; i < len(dios_board.Slots); i++ {
		dios_board.Slots[i].Orb.Attribute = WOOD
	}
	dios_board.Slots[1].Orb.Attribute = FIRE
	dios_board.Slots[2].Orb.Attribute = FIRE
	dios_board.Slots[3].Orb.Attribute = FIRE
	dios_board.Slots[7].Orb.Attribute = FIRE
	// dios_board.Slots[8].Orb.Attribute = FIRE
	dios_board.Slots[9].Orb.Attribute = FIRE
	dios_board.Slots[13].Orb.Attribute = FIRE
	dios_board.Slots[14].Orb.Attribute = FIRE
	dios_board.Slots[15].Orb.Attribute = FIRE
	fmt.Println(dios_board)
	combos := dios_board.GetCombos()
	combos[0].Print()
	combos[1].Print()
	// fmt.Println(dios_board.GetCombos())

	// fmt.Println(len(dios_board.GetCombos()[0].Positions))
}
