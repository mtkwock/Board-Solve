package main

import (
	"errors"
	"fmt"
	// "unsafe"
	"strings"
	"math/rand" // Runs deterministically unless we set a seed.
	// "testing"
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

var LetterToAttribute map[string]OrbAttribute = invertLetters()

func invertLetters() map[string]OrbAttribute {
	result := make(map[string]OrbAttribute)
	for attribute, repr := range AttributeToLetter {
		result[repr] = attribute
	}
	return result
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

func (self Board) SimpleString() string {
	result := make([]byte, len(self.Slots))
	for i, slot := range self.Slots {
		result[i] = byte(AttributeToLetter[slot.Orb.Attribute][0])
	}
	return string(result)
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
	return fmt.Sprintf("(%d,%d)", self.Y, self.X)
}

func (self Placement) Swap(direction Direction) Placement {
	new_placement := Placement{self.Y, self.X}
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
	return new_placement
}

func (self Board) Swap(placement Placement, direction Direction) (Board, error) {
	new_placement := placement.Swap(direction)

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

func (self BoardCombo) Print(width uint8) {
	fmt.Println(self)
	board := CreateEmptyBoard(width)
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
func (self Board) GetCombos() ([]BoardCombo, Board) {
	// Determine which orbs will be comboed out. Do not group them.
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
			if (x > 1) {
				orb_next := self.GetOrbAt(Placement{placement.Y, placement.X - 1})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Placement{placement.Y, placement.X - 2})
					if orb.Attribute == orb_next_next.Attribute {
						fmt.Println(y, x)
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
			if (y > 1) {
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

	// Determine if there are any unmatched bombs.  If so, clear orbs first then
	// get combos from that new board.
	unmatched_bombs := make([]Placement, 0)
	for i := uint8(0); i < uint8(len(self.Slots)); i++ {
		if self.Slots[i].Orb.Attribute == BOMB && !marked_combos[i] {
			unmatched_bombs = append(unmatched_bombs, Placement{uint8(i / self.Width), i % self.Width})
		}
	}
	if len(unmatched_bombs) != 0 {
		new_board := self.Clone()
		for _, placement := range unmatched_bombs {
			for x2 := uint8(0); x2 < self.Width; x2++ {
				pos := Placement{placement.Y, x2}.ToPos(new_board)
				if new_board.Slots[pos].Orb.Attribute != BOMB {
					new_board.Slots[pos].Orb.Attribute = EMPTY
				}
			}
			for y2 := uint8(0); y2 < self.Height; y2++ {
				pos := Placement{y2, placement.X}.ToPos(new_board)
				if new_board.Slots[pos].Orb.Attribute != BOMB {
					new_board.Slots[pos].Orb.Attribute = EMPTY
				}
			}
			new_board.Slots[placement.ToPos(new_board)].Orb.Attribute = EMPTY
		}
		return new_board.GetCombos()
	}

  // Group orbs to combo out.
	is_used := make([]bool, len(self.Slots))
	combos := make([]BoardCombo, 0)

  // For each orb, use a DFS to
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

	return combos, self
}

func (self Board) DropOrbs() {
	// new_to_old := map[uint8]uint8
	for y := self.Height - 1; y > 0; y-- {
		for x := uint8(0); x < self.Width; x++ {
			if self.GetOrbAt(Placement{y, x}).Attribute != EMPTY {
				continue
			}
			pos := Placement{y, x}.ToPos(self)
			for yo := y - 1; yo < self.Height; yo-- {
				new_pos := Placement{yo, x}.ToPos(self)
				moved_orb := self.Slots[new_pos].Orb
				if moved_orb.Attribute != EMPTY {
					self.Slots[pos].Orb = moved_orb
					self.Slots[new_pos].Orb = Orb{EMPTY, 0}
					// new_to_old[new_pos] = pos
					break
				}
			}
		}
	}
	// return new_to_old
}

func (self Board) GetAllCombos() []BoardCombo {
	// fmt.Println(self)
	all_combos := make([]BoardCombo, 0)
	current_board := self.Clone()
	// To determine the positions of the originally removed orbs.
	// new_to_old := map[uint8]uint8
	// for i := uint8(0); i < uint8(len(self.Slots)); i++ {
	// 	new_to_old[i] = i
	// }
	for new_combos, current_board := current_board.GetCombos(); len(new_combos) > 0; new_combos, current_board = current_board.GetCombos() {
		all_combos = append(all_combos, new_combos...)
		// fmt.Println(new_combos)
		for _, combo := range new_combos {
			for _, placement := range combo.Positions {
				current_board.Slots[placement.ToPos(current_board)].Orb.Attribute = EMPTY
			}
		}
		current_board.DropOrbs()
		// for new, old := range new_to_old {
		// 	if val, exists := newer_map[old]; exists {
		// 		new_to_old[val]
		// 	}
		// }
		// for new, old := range newer_map {
		// 	 new_to_old[]
		// }
		// fmt.Println(current_board)
	}
	return all_combos
}

func CreateEmptyBoard(width uint8) Board {
	return Board{make([]BoardSpace, width * (width - 1)), width - 1, width}
}

func CreateRandomBoard(width uint8) Board {
	size := width * (width - 1)
	slots := make([]BoardSpace, size)
	for i := uint8(0); i < size; i++ {
		slots[i].Orb.Attribute = OrbAttribute(uint8(rand.Intn(6) + 1))
	}
	return Board{slots, width - 1, width}
}

func CreateBoard(s string, width int) Board {
  slots := make([]BoardSpace, len(s))
	for i, rune := range s {
		slots[i].Orb.Attribute = LetterToAttribute[string(rune)]
	}
	return Board{slots, uint8(len(s) / width), uint8(width)}
}
