package blackjack

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fogleman/gg"
	"github.com/t11230/ramenbot/lib/utils"
	"golang.org/x/image/font/inconsolata"
	"image"
)

var tableDrawingSlots = []DrawingSlot{
	DrawingSlot{
		CardX:       310,
		CardY:       40,
		NameX:       385,
		NameY:       30,
		NameAnchorX: 1.0,
		NameAnchorY: 0,
	},
	DrawingSlot{
		CardX:       230,
		CardY:       120,
		NameX:       260,
		NameY:       110,
		NameAnchorX: 0.5,
		NameAnchorY: 0,
	},
	DrawingSlot{
		CardX:       100,
		CardY:       120,
		NameX:       130,
		NameY:       110,
		NameAnchorX: 0.5,
		NameAnchorY: 0,
	},
	DrawingSlot{
		CardX:       10,
		CardY:       40,
		NameX:       5,
		NameY:       30,
		NameAnchorX: 0,
		NameAnchorY: 0,
	},
}

var tableDealerDrawingSlot = DrawingSlot{
	CardX:       170,
	CardY:       25,
	NameX:       195,
	NameY:       2,
	NameAnchorX: 0.5,
	NameAnchorY: 1.0,
}

func (r *Round) Render() (im image.Image, err error) {
	guild, _ := r.Session.Guild(r.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild": r.GuildID,
		}).Warning("Failed to grab guild")
		return
	}

	dc := gg.NewContext(390, 220)
	dc.DrawCircle(195, 0, 220)
	dc.SetRGB(0.32, 0.63, 0.20)
	dc.Fill()

	dc.SetFontFace(inconsolata.Regular8x16)
	dc.SetRGB(0, 0, 0)

	for idx, player := range r.Players {
		hand := player.Hands[0]
		slot := tableDrawingSlots[idx]

		username := utils.GetPreferredName(guild, player.UserID)

		if len(username) > 14 {
			username = username[:11] + "..."
		}

		dc.DrawStringAnchored(username, slot.NameX, slot.NameY,
			slot.NameAnchorX, slot.NameAnchorY)

		cardX, cardY := slot.CardX, slot.CardY
		for _, card := range hand.Pile.Cards {
			cardIm, err := card.GetImage()
			if err != nil {
				log.Errorf("Error getting card image %v", err)
				continue
			}

			dc.DrawImage(cardIm, cardX, cardY)
			cardX += 10
			cardY += 5
		}
	}

	dc.DrawStringAnchored("Dealer",
		tableDealerDrawingSlot.NameX, tableDealerDrawingSlot.NameY,
		tableDealerDrawingSlot.NameAnchorX, tableDealerDrawingSlot.NameAnchorY)

	cardX, cardY := tableDealerDrawingSlot.CardX, tableDealerDrawingSlot.CardY
	for _, card := range r.Dealer.Hands[0].Pile.Cards {
		cardIm, err := card.GetImage()
		if err != nil {
			log.Errorf("Error getting card image %v", err)
			continue
		}

		dc.DrawImage(cardIm, cardX, cardY)
		cardX += 10
		cardY += 5
	}

	return dc.Image(), nil
}
