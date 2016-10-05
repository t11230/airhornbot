package blackjack

const (
	blackjackHelpString = `**BETROLL**
**usage:** !!$ betroll *ante* *<dietype>*
    This command initiates a bet on a dice roll. The second argument is the ante that all participants must pay into the pool.
    The third optional argument specifies a type of die for the roll.
**Die Types**
    **d6 (default):** 6-sided die.
    **d4:** 4-sided die.
    **d8:** 8-sided die.
    **d10:** 10-sided die.
    **d12:** 12-sided die.
    **d20:** 20-sided die.
    **other:** random integer generator between 1 and input.`

	bidHelpString = `**BID**
**usage:** !!$ bid *result*
    This command initiates bids on an in progress betroll. The second argument is the result that the user is bidding on.
    The number of bits placed on the bid is determined by the ante of the betroll.`
)
