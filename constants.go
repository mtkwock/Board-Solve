package main

import (
	"fmt"
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

var ALL_ATTRIBUTES = []OrbAttribute{FIRE, WATER, WOOD, LIGHT, DARK, HEART, JAMMER, POISON, MORTAL_POISON, BOMB}

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

func (self OrbAttribute) String() string {
	return AttributeToName[self]
}

func OrbAttributeArrayToString(attrs []OrbAttribute) string {
	result := "["
	for _, attr := range attrs {
		result += attr.String() + ", "
	}
	if len(result) == 1 {
		return "[]"
	}
	return result[:len(result) - 2] + "]"
}

func AllOrbAttributesExcept(to_ignore []OrbAttribute) []OrbAttribute {
	possibilities := make([]OrbAttribute, 0)
	for _, attribute := range ALL_ATTRIBUTES {
		will_add := true
		for _, ignored_attribute := range to_ignore {
			if attribute == ignored_attribute {
				will_add = false
				break
			}
		}
		if will_add {
			possibilities = append(possibilities, attribute)
		}
	}
	return possibilities
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
	COMBO
)

type BoardSpaceStateFlag uint8

const (
	TAPE BoardSpaceStateFlag = 1 << iota
	CLOUD
	SPINNER_1S
	SPINNER_2S
)

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

var DirectionToLetter map[Direction]string = map[Direction]string {
	RIGHT: "R",
	DOWN_RIGHT: "DR",
	DOWN: "D",
	DOWN_LEFT: "DL",
	LEFT: "L",
	UP_LEFT: "UL",
	UP: "U",
	UP_RIGHT: "UR",
}

func (self Direction) String() string {
	return DirectionToLetter[self]
}

func DirectionsToString(directions []Direction) string {
	result := ""
	for _, direction := range directions {
		result += fmt.Sprintf("%s, ", direction)
	}
	if len(result) <= 1 {
		return ""
	}
	return result[:len(result) - 2]
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
