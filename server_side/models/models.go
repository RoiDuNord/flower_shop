package models

type Postcard struct {
	Message string `json:"message"`
	Price   int    `json:"price"`
}

type Pack struct {
	Material string `json:"material"`
	Price    int    `json:"price"`
}

type Decoration struct {
	Postcard Postcard `json:"postcard"`
	Pack     Pack     `json:"pack"`
	Cost     int      `json:"decorationCost"`
}

type Flower struct {
	Name         string `json:"name"`
	Color        string `json:"color"`
	Quantity     int    `json:"quantity,omitempty"`
	Cost         int    `json:"cost,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type Bouquet struct {
	Position   int        `json:"position"`
	Flowers    []Flower   `json:"bouquet"`
	Decoration Decoration `json:"decoration"`
	Cost       int        `json:"bouquetCost"`
}

type Order struct {
	ID        int       `json:"orderID"`
	List      []Bouquet `json:"bouquetsList"`
	OrderCost int       `json:"orderCost"`
}
