package websockets

type Card struct {
	ID    int
	Color string
	Hand  *Client
}

type CardList []Card

var CardDeck = CardList{
	// tok 1-8
	{ID: 1, Color: "tok", Hand: nil},  // hetes
	{ID: 2, Color: "tok", Hand: nil},  // nyolcas
	{ID: 3, Color: "tok", Hand: nil},  // kilences
	{ID: 4, Color: "tok", Hand: nil},  // tizes
	{ID: 5, Color: "tok", Hand: nil},  // alsó
	{ID: 6, Color: "tok", Hand: nil},  // felső
	{ID: 7, Color: "tok", Hand: nil},  // király
	{ID: 8, Color: "tok", Hand: nil},  // ász
	// makk 9-16
	{ID: 9, Color: "makk", Hand: nil},
	{ID: 10, Color: "makk", Hand: nil},
	{ID: 11, Color: "makk", Hand: nil},
	{ID: 12, Color: "makk", Hand: nil},
	{ID: 13, Color: "makk", Hand: nil},
	{ID: 14, Color: "makk", Hand: nil},
	{ID: 15, Color: "makk", Hand: nil},
	{ID: 16, Color: "makk", Hand: nil},
	// zold 17-24
	{ID: 17, Color: "zold", Hand: nil},
	{ID: 18, Color: "zold", Hand: nil},
	{ID: 19, Color: "zold", Hand: nil},
	{ID: 20, Color: "zold", Hand: nil},
	{ID: 21, Color: "zold", Hand: nil},
	{ID: 22, Color: "zold", Hand: nil},
	{ID: 23, Color: "zold", Hand: nil},
	{ID: 24, Color: "zold", Hand: nil},
	// piros 25-32
	{ID: 25, Color: "piros", Hand: nil},
	{ID: 26, Color: "piros", Hand: nil},
	{ID: 27, Color: "piros", Hand: nil},
	{ID: 28, Color: "piros", Hand: nil},
	{ID: 29, Color: "piros", Hand: nil},
	{ID: 30, Color: "piros", Hand: nil},
	{ID: 31, Color: "piros", Hand: nil},
	{ID: 32, Color: "piros", Hand: nil},
}
