package main

import (
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	d := newDeck()
	if len(d) != 52 {
		t.Errorf("Expected deck length of 52, but got %d", len(d))
	}
	if d[0] != "Ace of Spades" {
		t.Errorf("Expected first card value to be 'Ace of Spades', but got %v", d[0])
	}
	if d[len(d)-1] != "King of Clubs" {
		t.Errorf("Expected last card value to be 'King of Clubs', but got %v", d[len(d)-1])
	}
}

func TestSaveToFileAndNewDeckFromFile(t *testing.T) {

	os.Remove("_decktesting")

	d := newDeck()
	d.shuffle()
	d.saveToFile("_decktesting")

	loadedDeck := newDeckFromFile("_decktesting")
	if len(loadedDeck) != len(d) {
		t.Errorf("Expected deck length of %d, but got %d", len(d), len(loadedDeck))
	}
	if loadedDeck[0] != d[0] {
		t.Errorf("Expected first card value to be '%v', but got '%v'", d[0], loadedDeck[0])
	}
	if loadedDeck[len(loadedDeck)-1] != d[len(d)-1] {
		t.Errorf("Expected last card value to be '%v', but got '%v'", d[len(d)-1], loadedDeck[len(loadedDeck)-1])
	}

	os.Remove("_decktesting")

}
