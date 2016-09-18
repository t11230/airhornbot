package cards

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"image"
	"io/ioutil"
	"net/http"
)

func NewDeck(deckCount int) (*Deck, error) {
	respData, err := apiRequest(fmt.Sprintf(endpointShuffle, deckCount))
	if err != nil {
		return nil, err
	}

	result := &Deck{}
	err = jsonUnmarshal(respData, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Deck) Draw(count int) (*DrawResult, error) {
	respData, err := apiRequest(fmt.Sprintf(endpointDraw, d.DeckID, count))
	if err != nil {
		return nil, err
	}

	result := &DrawResult{}
	err = jsonUnmarshal(respData, result)
	if err != nil {
		return nil, err
	}

	d.Remaining = result.Remaining

	return result, nil
}

func (d *Deck) Reshuffle() error {
	respData, err := apiRequest(fmt.Sprintf(endpointReshuffle, d.DeckID))
	if err != nil {
		return err
	}

	result := &Deck{}
	err = jsonUnmarshal(respData, result)
	if err != nil {
		return err
	}

	d.Remaining = result.Remaining

	return nil
}

func (c *Card) GetImage() (image.Image, error) {
	resp, err := http.Get(c.ImageURL)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return img, nil
}

func apiRequest(url string) ([]byte, error) {
	log.Debugf("Requesting: %v", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response code: %v", resp.Status)
	log.Debugf("Response body len: %v", len(respData))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Status code: %v", resp.StatusCode)
	}

	return respData, nil
}

func jsonUnmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return errors.New("Unmarshal Failure")
	}
	return nil
}
