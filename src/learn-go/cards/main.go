package main

func main() {
	//cards := newDeck()
	// hand, remainingCards := deal(cards, 5)
	// hand.print()
	// remainingCards.print()

	//cards.saveToFile("my_cards")
	newCards := newDeckFromFile("my_cards")
	newCards.shuffle()
	newCards.print()
}
