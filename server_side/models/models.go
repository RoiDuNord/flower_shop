package models

type Postcard struct {
	Message string `json:"message"`
	Cost    int    `json:"cost"`
}

type Pack struct {
	Material string `json:"material"`
	Cost     int    `json:"cost"`
}

type Decoration struct {
	Postcard       Postcard `json:"postcard"`
	Pack           Pack     `json:"pack"`
	DecorationCost int      `json:"decorationCost"`
}

type Flower struct {
	Name         string `json:"name"`
	Color        string `json:"color"`
	Quantity     int    `json:"quantity,omitempty"`
	Cost         int    `json:"cost,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type Bouquet struct {
	Position    int        `json:"position"`
	FlowerList  []Flower   `json:"bouquet"`
	Decoration  Decoration `json:"decoration"`
	BouquetCost int        `json:"bouquetCost"`
}

type Payment struct {
	OrderID   int  `json:"orderID,omitempty"`
	IsPaid    bool `json:"IsPaid"`
	PaymentID int  `json:"paymentID,omitempty"`
}

type Order struct {
	ID           int       `json:"orderID"`
	BouquetsList []Bouquet `json:"bouquetsList"`
	OrderCost    int       `json:"orderCost"`
	PaymentID    int       `json:"paymentID,omitempty"`
}
