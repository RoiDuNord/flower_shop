package test

func t() {

}

// import (
// 	"encoding/json"
// 	"models"
// )

// var id int
// var flowers, postcards, packs = initializeData()

// func makeOrder() ([]byte, error) {
// 	flowers1 := []models.Flower{flowers["yellowRose"], flowers["yellowLily"]}
// 	flowers1[0].Quantity, flowers1[1].Quantity = 3, 4
// 	decoration1 := models.Decoration{Postcard: postcards["womenDay"], Pack: packs["craft"]}
// 	price1 := bouquetPrice(flowers1, decoration1)

// 	flowers2 := []models.Flower{flowers["daisy"], flowers["whitePion"]}
// 	flowers2[0].Quantity, flowers2[1].Quantity = 50, 5
// 	decoration2 := models.Decoration{Postcard: postcards["valentineDay"], Pack: packs["tape"]}
// 	price2 := bouquetPrice(flowers2, decoration2)

// 	bouquet1 := models.Bouquet{
// 		Position:   nextID(),
// 		Flowers:    flowers1,
// 		Cost:       price1,
// 		Decoration: decoration1,
// 	}

// 	bouquet2 := models.Bouquet{
// 		Position:   nextID(),
// 		Flowers:    flowers2,
// 		Cost:       price2,
// 		Decoration: decoration2,
// 	}

// 	bouquets := make([]models.Bouquet, 0)
// 	bouquets = append(bouquets, bouquet1, bouquet2)

// 	jsonData, err := json.MarshalIndent(bouquets, "", "   ")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return jsonData, nil
// }

// func nextID() int {
// 	id++
// 	return id
// }

// func bouquetPrice(flowersAr []models.Flower, decoration models.Decoration) (price int) {
// 	for _, flower := range flowersAr {
// 		price += flower.Cost * flower.Quantity
// 	}
// 	price += decoration.Pack.Price + decoration.Postcard.Price
// 	return price
// }

// func initializeData() (map[string]models.Flower, map[string]models.Postcard, map[string]models.Pack) {
// 	flowers := map[string]models.Flower{
// 		"redRose":    {Name: "Роза", Color: "Красная", Cost: 80, Quantity: 20},
// 		"whiteRose":  {Name: "Роза", Color: "Белая", Cost: 60, Quantity: 20},
// 		"yellowRose": {Name: "Роза", Color: "Жёлтая", Cost: 40, Quantity: 20},
// 		"whiteLily":  {Name: "Лилия", Color: "Белая", Cost: 100, Quantity: 5},
// 		"yellowLily": {Name: "Лилия", Color: "Жёлтая", Cost: 90, Quantity: 5},
// 		"pinkPion":   {Name: "Пион", Color: "Розовый", Cost: 120, Quantity: 10},
// 		"whitePion":  {Name: "Пион", Color: "Белый", Cost: 110, Quantity: 10},
// 		"lotus":      {Name: "Лотос", Color: "Белый", Cost: 200, Quantity: 5},
// 		"daisy":      {Name: "Ромашка", Color: "Белая", Cost: 20, Quantity: 50},
// 	}

// 	postcards := map[string]models.Postcard{
// 		"birthday":         {Message: "С Днём рождения!", Price: 5},
// 		"newYear":          {Message: "С Новым Годом!", Price: 1},
// 		"happyWedding":     {Message: "Со свадьбой!", Price: 2},
// 		"happyAnniversary": {Message: "С Юбилеем!", Price: 3},
// 		"womenDay":         {Message: "С 8 марта!", Price: 15},
// 		"valentineDay":     {Message: "С Днём Влюбленных!", Price: 20},
// 	}

// 	packs := map[string]models.Pack{
// 		"craft": {Material: "Крафт", Price: 100},
// 		"film":  {Material: "Плёнка", Price: 50},
// 		"tape":  {Material: "Лента", Price: 10},
// 	}

// 	return flowers, postcards, packs
// }

// func main() {
// 	Create()
// }

// var db *sql.DB

// func Create() {

// 	createFlowersTable := `
//     CREATE TABLE IF NOT EXISTS flowers (
//         id SERIAL PRIMARY KEY,
//         name VARCHAR(100),
//         color VARCHAR(50),
//         price INT,
//         quantity INT
//     );`

// 	createPostcardsTable := `
//     CREATE TABLE IF NOT EXISTS postcards (
//         id SERIAL PRIMARY KEY,
//         note TEXT,
//         price INT
//     );`

// 	createPacksTable := `
//     CREATE TABLE IF NOT EXISTS packs (
//         id SERIAL PRIMARY KEY,
//         material VARCHAR(100),
//         price INT
//     );`

// 	var err error

// 	_, err = db.Exec(createFlowersTable)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = db.Exec(createPostcardsTable)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = db.Exec(createPacksTable)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	flowers, postcards, packs := initializeData()

// 	for _, flower := range flowers {
// 		insertFlowerSQL := `INSERT INTO flowers (name, color, price, quantity) VALUES ($1, $2, $3, $4)`
// 		_, err = db.Exec(insertFlowerSQL, flower.Name, flower.Color, flower.Cost, flower.ReqQty)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	fmt.Println("Flowers data inserted successfully!")

// 	for _, postcard := range postcards {
// 		insertPostcardSQL := `INSERT INTO postcards (note, price) VALUES ($1, $2)`
// 		_, err = db.Exec(insertPostcardSQL, postcard.Message, postcard.Price)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	fmt.Println("Postcards data inserted successfully!")

// 	for _, pack := range packs {
// 		insertPackSQL := `INSERT INTO packs (material, price) VALUES ($1, $2)`
// 		_, err = db.Exec(insertPackSQL, pack.Material, pack.Price)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	fmt.Println("Packs data inserted successfully!")
// }
