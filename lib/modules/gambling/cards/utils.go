package cards

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
	"1": 1, "2": 2, "3": 3, "4": 4,
	"5": 5, "6": 6, "7": 7, "8": 8,
	"9": 9, "10": 10, "JACK": 10, "QUEEN": 10,
	"KING": 10, "ACE": 11,
}

var BlackjackAceLowMap = ValueMap{
	"ACE": 1, "1": 1, "2": 2, "3": 3,
	"4": 4, "5": 5, "6": 6, "7": 7,
	"8": 8, "9": 9, "10": 10, "JACK": 10,
	"QUEEN": 10, "KING": 10,
}
