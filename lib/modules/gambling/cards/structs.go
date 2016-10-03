package cards

type Suit int
type Value int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

const (
	Ace Value = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

var (
	DefaultSuits  = []Suit{Clubs, Diamonds, Hearts, Spades}
	DefaultValues = []Value{
		Ace, Two, Three, Four, Five,
		Six, Seven, Eight, Nine, Ten,
		Jack, Queen, King,
	}
)

type Deck struct {
	Cards    []Card
	Shuffled bool
}

type Card struct {
	Value      Value
	Suit       Suit
	IsFaceDown bool
}

type Pile struct {
	Cards []Card
}

type ValueMap map[Value]int
