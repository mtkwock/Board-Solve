package main

type Orb struct {
	Attribute OrbAttribute // See constants
	State OrbStateFlag // See constants
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
