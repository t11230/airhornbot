package cards

var (
	endpointRoot = "https://deckofcardsapi.com/api/deck/"

	endpointShuffle        = endpointRoot + "new/shuffle/?deck_count=%d"
	endpointDraw           = endpointRoot + "%s/draw/?count=%d"
	endpointReshuffle      = endpointRoot + "%s/shuffle/"
	endpointNewDeck        = endpointRoot + "new/"
	endpointNewPartialDeck = endpointRoot + "new/shuffle/?cards=%s"
	endpointAddToPile      = endpointRoot + "%s/pile/%s/add/?cards=%s"
	endpointDrawFromPile   = endpointRoot + "%s/pile/%s/draw/?cards=%s"
)
