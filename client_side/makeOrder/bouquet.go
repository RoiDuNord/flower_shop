package c_order

type Postcard struct {
	Note  string `json:"note"`
	Price int    `json:"price"`
}

type Pack struct {
	Material string `json:"material"`
	Price    int    `json:"price"`
}

type Decoration struct {
	Postcard Postcard `json:"postcard"`
	Pack     Pack     `json:"pack"`
}

type Flower struct {
	Name     string `json:"name"`
	Color    string `json:"color"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type Bouquet struct {
	Position   int        `json:"position"`
	Bouquet    []Flower   `json:"bouquet"`
	Decoration Decoration `json:"decoration"`
	Price      int        `json:"price"`
}

type Bouquets struct {
	List []Bouquet `json:"list"`
}

func BouquetPrice(flowers []Flower, decoration Decoration) (price int) {
	for _, flower := range flowers {
		price += flower.Price * flower.Quantity
	}
	price += decoration.Pack.Price + decoration.Postcard.Price
	return price
}
