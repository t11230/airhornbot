package cards

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fogleman/gg"
	"image"
	"math/rand"
	"path"
)

func NewDeck(shuffle bool) (d *Deck) {
	log.Debug("Creating a new deck")
	d = &Deck{}

	d.Cards = make([]Card, len(DefaultSuits)*len(DefaultValues))

	// Generate all combinations of the default suits and values
	for suitIdx, suit := range DefaultSuits {
		for valueIdx, value := range DefaultValues {
			deckIdx := suitIdx*len(DefaultValues) + valueIdx
			d.Cards[deckIdx] = Card{
				Value:      value,
				Suit:       suit,
				IsFaceDown: false,
			}
		}
	}

	// Shuffle if necessary
	if shuffle {
		d.Shuffle()
	}
	return
}

func (d *Deck) Shuffle() {
	// https://gist.github.com/quux00/8258425
	N := len(d.Cards)
	for i := 0; i < N; i++ {
		r := i + rand.Intn(N-i)
		d.Cards[r], d.Cards[i] = d.Cards[i], d.Cards[r]
	}

	d.Shuffled = true
}

func (d *Deck) Draw(nCards int) (result *Pile) {
	result = &Pile{d.Cards[:nCards]}
	d.Cards = d.Cards[nCards:]
	return
}

func (p *Pile) AddCards(cards ...Card) {
	p.Cards = append(p.Cards, cards...)
}

func (p *Pile) AddPile(pile *Pile) {
	p.Cards = append(p.Cards, pile.Cards...)
}

func (c *Card) GetImage() (image.Image, error) {
	return gg.LoadImage(c.GetFilepath())
}

func (c *Card) GetFilename() (str string) {
	if c.IsFaceDown {
		return "redBack.png"
	}

	suitFname, ok := suitImageNameMap[c.Suit]
	if !ok {
		log.Errorf("Unknown suit %v", c.Suit)
		return "unknown"
	}
	str += suitFname

	valueFname, ok := valueImageNameMap[c.Value]
	if !ok {
		log.Errorf("Unknown value %v", c.Value)
		return "unknown"
	}
	str += valueFname

	str += ".png"
	return
}

func (c *Card) GetFilepath() string {
	return path.Join("assets", "png-42px", c.GetFilename())
}
