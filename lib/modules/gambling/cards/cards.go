package cards

import (
	log "github.com/Sirupsen/logrus"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
)

const (
	cardWidth      = 75
	cardXSpacing   = 2
	cardYSpacing   = 5
	cardSlotWidth  = cardWidth + cardXSpacing
	cardSlotHeight = 105
	cardsPerRow    = 5
)

// GenerateImage converts an array of Cards into an Image
func GenerateImage(cards []Card) (image.Image, error) {
	log.Debug("Generating image")

	nCards := len(cards)

	nCardsInFirstRow := nCards
	if nCardsInFirstRow > 5 {
		nCardsInFirstRow = 5
	}

	resultImg := image.NewRGBA(image.Rect(0, 0,
		cardSlotWidth*nCardsInFirstRow+cardXSpacing*2,
		cardSlotHeight*(((nCards-1)/cardsPerRow)+1)+cardYSpacing*2))

	log.Debugf("Image Size: %v", resultImg.Bounds())

	for index, c := range cards {
		img, err := c.GetImage()
		if err != nil {
			log.Error(err)
			return nil, err
		}

		rImg := resize.Thumbnail(cardWidth, uint(img.Bounds().Dy()), img, resize.Bilinear)

		pt := image.Point{cardSlotWidth*(index%cardsPerRow) + cardXSpacing,
			cardSlotHeight*(index/cardsPerRow) + cardYSpacing}
		rect := image.Rectangle{pt, pt.Add(rImg.Bounds().Size())}
		draw.Draw(resultImg, rect, rImg, image.Point{0, 0}, draw.Src)
	}

	return resultImg, nil
}
