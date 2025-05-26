// package tester

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"which/models"

// 	_ "github.com/lib/pq"
// )

// var db *sql.DB

// func main() {
// 	var err error

// 	connStr := "user=florist dbname=flower_shop sslmode=disable" // Замените 'yourusername' и 'yourdbname' на ваши значения
// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// 	Create()
// }

// func Create() {

// 	createFlowersTable := `
//     CREATE TABLE IF NOT EXISTS flowers (
//         id SERIAL PRIMARY KEY,
//         name VARCHAR(100),
//         color VARCHAR(50),
//         cost INT,
//         quantity INT
//     );`

// 	createPostcardsTable := `
//     CREATE TABLE IF NOT EXISTS postcards (
//         id SERIAL PRIMARY KEY,
//         note TEXT,
//         cost INT
//     );`

// 	createPacksTable := `
//     CREATE TABLE IF NOT EXISTS packs (
//         id SERIAL PRIMARY KEY,
//         material VARCHAR(100),
//         cost INT
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
// 		insertFlowerSQL := `INSERT INTO flowers (name, color, cost, quantity) VALUES ($1, $2, $3, $4)`
// 		_, err = db.Exec(insertFlowerSQL, flower.Name, flower.Color, flower.Cost, flower.Quantity)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	fmt.Println("Flowers data inserted successfully!")

// 	for _, postcard := range postcards {
// 		insertPostcardSQL := `INSERT INTO postcards (note, cost) VALUES ($1, $2)`
// 		_, err = db.Exec(insertPostcardSQL, postcard.Message, postcard.Cost)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	fmt.Println("Postcards data inserted successfully!")

// 	for _, pack := range packs {
// 		insertPackSQL := `INSERT INTO packs (material, cost) VALUES ($1, $2)`
// 		_, err = db.Exec(insertPackSQL, pack.Material, pack.Cost)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	fmt.Println("Packs data inserted successfully!")
// }

// package main

// import (
// 	"context"
// 	"time"
// 	order "which/orderToKafka"
// )

// func main() {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
// 	defer cancel()

// 	orderQty := 5

// 	getOrderData(ctx, orderQty)
// }

// func getOrderData(ctx context.Context, orderQty int) {
// 	orderChan := make(chan order.Order)

// 	order.OrderToKafka(ctx, orderQty, orderChan)
// }

// var id int
// var flowers, postcards, packs = initializeData()

// func makeOrder() ([]byte, error) {
// 	flowers1 := []models.Flower{flowers["yellowRose"], flowers["yellowLily"]}
// 	flowers1[0].Quantity, flowers1[1].Quantity = 3, 4
// 	decoration1 := models.Decoration{Postcard: postcards["womenDay"], Pack: packs["craft"]}
// 	cost1 := bouquetCost(flowers1, decoration1)

// 	flowers2 := []models.Flower{flowers["daisy"], flowers["whitePion"]}
// 	flowers2[0].Quantity, flowers2[1].Quantity = 50, 5
// 	decoration2 := models.Decoration{Postcard: postcards["valentineDay"], Pack: packs["tape"]}
// 	cost2 := bouquetCost(flowers2, decoration2)

// 	bouquet1 := models.Bouquet{
// 		Position:   nextID(),
// 		Flowers:    flowers1,
// 		BouquetCost:       cost1,
// 		Decoration: decoration1,
// 	}

// 	bouquet2 := models.Bouquet{
// 		Position:   nextID(),
// 		Flowers:    flowers2,
// 		BouquetCost:       cost2,
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

// func bouquetCost(flowersAr []models.Flower, decoration models.Decoration) (cost int) {
// 	for _, flower := range flowersAr {
// 		cost += flower.Cost * flower.Quantity
// 	}
// 	cost += decoration.Pack.Cost + decoration.Postcard.Cost
// 	return cost
// }

package db

import (
	"database/sql"
	"fmt"
	"log"
	"server/models"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error

	connStr := "user=florist dbname=flower_shop sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createDB()
}

func createDB() {
	createFlowersTable := `
    CREATE TABLE IF NOT EXISTS flowers (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        color VARCHAR(50),
        cost INT,
        quantity INT
    );`

	createPostcardsTable := `
    CREATE TABLE IF NOT EXISTS postcards (
        id SERIAL PRIMARY KEY,
        note TEXT,
        cost INT
    );`

	createPacksTable := `
    CREATE TABLE IF NOT EXISTS packs (
        id SERIAL PRIMARY KEY,
        material VARCHAR(100),
        cost INT
    );`

	var err error

	_, err = db.Exec(createFlowersTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createPostcardsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createPacksTable)
	if err != nil {
		log.Fatal(err)
	}

	flowers, postcards, packs := initializeData()

	for _, flower := range flowers {
		insertFlowerSQL := `INSERT INTO flowers (name, color, cost, quantity) VALUES ($1, $2, $3, $4)`
		_, err = db.Exec(insertFlowerSQL, flower.Name, flower.Color, flower.Cost, flower.Quantity)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Flowers data inserted successfully!")

	for _, postcard := range postcards {
		insertPostcardSQL := `INSERT INTO postcards (note, cost) VALUES ($1, $2)`
		_, err = db.Exec(insertPostcardSQL, postcard.Message, postcard.Cost)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Postcards data inserted successfully!")

	for _, pack := range packs {
		insertPackSQL := `INSERT INTO packs (material, cost) VALUES ($1, $2)`
		_, err = db.Exec(insertPackSQL, pack.Material, pack.Cost)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Packs data inserted successfully!")
}

func initializeData() (map[string]models.Flower, map[string]models.Postcard, map[string]models.Pack) {
	flowers := map[string]models.Flower{
		"redRose":    {Name: "Роза", Color: "Красная", Cost: 80, Quantity: 200},
		"whiteRose":  {Name: "Роза", Color: "Белая", Cost: 60, Quantity: 200},
		"yellowRose": {Name: "Роза", Color: "Жёлтая", Cost: 40, Quantity: 200},
		"whiteLily":  {Name: "Лилия", Color: "Белая", Cost: 100, Quantity: 50},
		"yellowLily": {Name: "Лилия", Color: "Жёлтая", Cost: 90, Quantity: 50},
		"pinkPion":   {Name: "Пион", Color: "Розовый", Cost: 120, Quantity: 100},
		"whitePion":  {Name: "Пион", Color: "Белый", Cost: 110, Quantity: 100},
		"lotus":      {Name: "Лотос", Color: "Белый", Cost: 200, Quantity: 50},
		"chamomile":  {Name: "Ромашка", Color: "Белая", Cost: 20, Quantity: 500},
	}

	postcards := map[string]models.Postcard{
		"birthday":         {Message: "С Днём рождения!", Cost: 5},
		"newYear":          {Message: "С Новым Годом!", Cost: 1},
		"happyWedding":     {Message: "Со свадьбой!", Cost: 2},
		"happyAnniversary": {Message: "С Юбилеем!", Cost: 3},
		"womenDay":         {Message: "С 8 марта!", Cost: 15},
		"valentineDay":     {Message: "С Днём Влюбленных!", Cost: 20},
	}

	packs := map[string]models.Pack{
		"craft": {Material: "Крафт", Cost: 100},
		"film":  {Material: "Плёнка", Cost: 50},
		"tape":  {Material: "Лента", Cost: 10},
	}

	return flowers, postcards, packs
}
