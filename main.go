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
	_ = iota
	FIRE OrbAttribute = iota
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
}
