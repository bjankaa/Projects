package websockets

type Card struct {
	ID    int
	Color string
	Hand  *Client
}

type CardList []Card

var CardDeck = CardList{
	{ID: 7, Color: "tok", Hand: nil},  // hetes
	{ID: 8, Color: "tok", Hand: nil},  // nyolcas
	{ID: 9, Color: "tok", Hand: nil},  // kilences
	{ID: 10, Color: "tok", Hand: nil}, // tizes
	{ID: 11, Color: "tok", Hand: nil}, // alsó
	{ID: 12, Color: "tok", Hand: nil}, // felső
	{ID: 13, Color: "tok", Hand: nil}, // király
	{ID: 14, Color: "tok", Hand: nil}, // ász
	{ID: 17, Color: "makk", Hand: nil},
	{ID: 18, Color: "makk", Hand: nil},
	{ID: 19, Color: "makk", Hand: nil},
	{ID: 20, Color: "makk", Hand: nil},
	{ID: 21, Color: "makk", Hand: nil},
	{ID: 22, Color: "makk", Hand: nil},
	{ID: 23, Color: "makk", Hand: nil},
	{ID: 24, Color: "makk", Hand: nil},
	{ID: 27, Color: "zold", Hand: nil},
	{ID: 28, Color: "zold", Hand: nil},
	{ID: 29, Color: "zold", Hand: nil},
	{ID: 30, Color: "zold", Hand: nil},
	{ID: 31, Color: "zold", Hand: nil},
	{ID: 32, Color: "zold", Hand: nil},
	{ID: 33, Color: "zold", Hand: nil},
	{ID: 34, Color: "zold", Hand: nil},
	{ID: 37, Color: "piros", Hand: nil},
	{ID: 38, Color: "piros", Hand: nil},
	{ID: 39, Color: "piros", Hand: nil},
	{ID: 40, Color: "piros", Hand: nil},
	{ID: 41, Color: "piros", Hand: nil},
	{ID: 42, Color: "piros", Hand: nil},
	{ID: 43, Color: "piros", Hand: nil},
	{ID: 44, Color: "piros", Hand: nil},
}
