package cards

type Deck struct {
	DeckID    string `json:"deck_id"`
	Remaining int    `json:"remaining"`
	Shuffled  bool   `json:"shuffled"`
}

type Card struct {
	ImageURL   string `json:"image"`
	Value      string `json:"value"`
	Suit       string `json:"suit"`
	Code       string `json:"code"`
	IsFaceDown bool
}

type DrawResult struct {
	Cards     []Card `json:"cards"`
	DeckID    string `json:"deck_id"`
	Remaining int    `json:"remaining"`
}

type Pile struct {
	Cards     []Card
	DeckID    string
	Remaining int
}

type ValueMap map[string]int
