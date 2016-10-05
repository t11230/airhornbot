package cards

func (s Suit) String() string {
	switch s {
	case Clubs:
		return "Clubs"
	case Diamonds:
		return "Diamonds"
	case Hearts:
		return "Hearts"
	case Spades:
		return "Spades"
	}
	return "Unknown"
}

func (v Value) String() string {
	if str, ok := valueStringMap[v]; ok {
		return str
	}
	return "Unknown"
}

var valueStringMap = map[Value]string{
	Ace: "Ace", Two: "Two", Three: "Three", Four: "Four",
	Five: "Five", Six: "Six", Seven: "Seven", Eight: "Eight",
	Nine: "Nine", Ten: "Ten", Jack: "Jack", Queen: "Queen",
	King: "King",
}

var suitImageNameMap = map[Suit]string{
	Clubs: "club", Diamonds: "diamond",
	Hearts: "heart", Spades: "spade",
}

var valueImageNameMap = map[Value]string{
	Ace: "Ace", Two: "2", Three: "3", Four: "4",
	Five: "5", Six: "6", Seven: "7", Eight: "8",
	Nine: "9", Ten: "10", Jack: "Jack", Queen: "Queen",
	King: "King",
}

// Cards useful functions
func (c *Card) NumericValue(m ValueMap) (value int) {
	return m[c.Value]
}

// Piles useful functions
func (p *Pile) Sum(m ValueMap) (sum int) {
	for _, c := range p.Cards {
		sum += c.NumericValue(m)
	}
	return
}

// Common card value maps
var BlackjackAceHighMap = ValueMap{
	Two: 2, Three: 3, Four: 4,
	Five: 5, Six: 6, Seven: 7, Eight: 8,
	Nine: 9, Ten: 10, Jack: 10, Queen: 10,
	King: 10, Ace: 11,
}

var BlackjackAceLowMap = ValueMap{
	Ace: 1, Two: 2, Three: 3, Four: 4,
	Five: 5, Six: 6, Seven: 7, Eight: 8,
	Nine: 9, Ten: 10, Jack: 10, Queen: 10,
	King: 10,
}
