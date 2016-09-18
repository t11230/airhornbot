package cards

import (
	log "github.com/Sirupsen/logrus"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
)

func GenerateImage(cards []Card) (image.Image, error) {
	log.Debug("Generating image")

	nCards := len(cards)

	resultImg := image.NewRGBA(image.Rect(0, 0, 80*nCards+20*2, 104+10*2))

	log.Debugf("Image Size: %v", resultImg.Bounds())

	for index, c := range cards {
		img, err := c.GetImage()
		if err != nil {
			log.Error(err)
			return nil, err
		}

		rImg := resize.Thumbnail(75, uint(img.Bounds().Dy()), img, resize.Bilinear)

		pt := image.Point{10 + 80*index, 10}

		log.Debugf("Drawing at %v", pt)

		rect := image.Rectangle{pt, pt.Add(rImg.Bounds().Size())}

		log.Debugf("Rect is %v", rect)

		draw.Draw(resultImg, rect, rImg, image.Point{0, 0}, draw.Src)
	}

	return resultImg, nil
}
