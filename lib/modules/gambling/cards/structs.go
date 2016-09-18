package cards

type Deck struct {
	DeckID    string `json:"deck_id"`
	Remaining int    `json:"remaining"`
	Shuffled  bool   `json:"shuffled"`
}

type Card struct {
	ImageURL string `json:"image"`
	Value    string `json:"value"`
	Suit     string `json:"suit"`
	Code     string `json:"code"`
}

type DrawResult struct {
	Cards     []Card `json:"cards"`
	DeckID    string `json:"deck_id"`
	Remaining int    `json:"remaining"`
}
