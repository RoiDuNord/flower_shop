package s_order

// import (
// 	c_order "client/makeOrder"
// 	"encoding/json"
// )

// var id int

// type Order struct {
// 	ID         int               `json:"id"`
// 	Bouquets   []c_order.Bouquet `json:"bouquets"`
// 	TotalPrice int               `json:"price"`
// }

// func (ord *Order) Parse(data []byte) []byte {
// 	id = newID()

// 	var bouquets c_order.Bouquets
// 	json.Unmarshal(data, &bouquets)

// 	var totalPrice int
// 	for _, bouquet := range bouquets.Bouquets {
// 		totalPrice += bouquet.Price
// 	}

// 	ord = &Order{
// 		ID:         id,
// 		Bouquets:   bouquets.Bouquets,
// 		TotalPrice: totalPrice,
// 	}

// 	jd, _ := json.MarshalIndent(ord, "  ", "")
// 	return jd
// }

// func newID() int {
// 	id++
// 	return id
// }

// func Info() {
// 	id := newID()
// 	sum := price1 + price2

// 	order := c_order.Order{
// 		ID:         id,
// 		Bouquets:   bouquets,
// 		TotalPrice: sum,
// 	}

// 	fmt.Printf("Заказ №%d на сумму %d₽\n", order.ID, order.TotalPrice)
// 	for i, bouquet := range order.Bouquets {
// 		fmt.Printf("Букет %d\n", i+1)
// 		for _, flower := range bouquet.Flowers {
// 			quantity := quantity(flower.Quantity)
// 			fmt.Printf("- %s %s (%s)\n", flower.Name, flower.Color, quantity)
// 		}
// 	}

// }

// func quantity(num int) (pieceForm string) {
// 	switch {
// 	case num%100 >= 11 && num%100 <= 20:
// 		pieceForm = "штук"
// 	case num%10 == 1:
// 		pieceForm = "штука"
// 	case num%10 >= 2 && num%10 <= 4:
// 		pieceForm = "штуки"
// 	default:
// 		pieceForm = "штук"
// 	}

// 	return fmt.Sprintf("%d %s", num, pieceForm)
// }
