package tester

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
