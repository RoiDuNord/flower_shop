package models

type Postcard struct {
	Message string `json:"message"`
	Cost    int    `json:"cost,omitempty"`
}

type Pack struct {
	Material string `json:"material"`
	Cost     int    `json:"cost,omitempty"`
}

type Decoration struct {
	Postcard       Postcard `json:"postcard"`
	Pack           Pack     `json:"pack"`
	DecorationCost int      `json:"decorationCost,omitempty"`
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
	BouquetCost int        `json:"bouquetCost,omitempty"`
}

type Order struct {
	ID            int       `json:"orderID,omitempty"`
	BouquetsList  []Bouquet `json:"bouquetsList"`
	OrderCost     int       `json:"orderCost,omitempty"`
	PaymentStatus string    `json:"status,omitempty"`
	PaymentID     int       `json:"paymentID,omitempty"`
}

type Payment struct {
	OrderID   int  `json:"orderID,omitempty"`
	IsPaid    bool `json:"IsPaid"`
	PaymentID int  `json:"paymentID,omitempty"`
}
