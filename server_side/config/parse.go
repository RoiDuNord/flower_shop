package config

import (
	"fmt"
	"os"
)

// func ParseConfig() (Config, error) {
// 	cfg, err := loadFromFile("config.yaml")
// 	if err != nil {
// 		return Config{}, err
// 	}

// 	if err := cfg.validate(); err != nil {
// 		return Config{}, err
// 	}

// 	return cfg, nil
// }

func ParseConfig() (Config, error) {
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config.yaml not found")
	}

	cfg, err := loadFromFile("config.yaml")
	if err != nil {
		return Config{}, err
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Error parsing config: open config.yaml: no such file or directory

// ├── README.md
// ├── server_side
// │   ├── businessLogic
// │   │   └── order
// │   │       ├── decorationCost.go
// │   │       ├── flowerQuantityAndCost.go
// │   │       ├── parseOrder.go
// │   │       ├── payment1.json
// │   │       └── requestToHandledBouquets.go
// │   ├── cmd
// │   │   ├── main
// │   │   └── main.go
// │   ├── config
// │   │   ├── getDBparams.go
// │   │   ├── loadFromFile.go
// │   │   ├── parse.go
// │   │   └── validate.go
// │   ├── config.yaml
// │   ├── db
// │   │   ├── db.go
// │   │   ├── getDecorElPrice.go
// │   │   ├── getFlowersQtyAndPrice.go
// │   │   └── updateQty.go
// │   ├── go.mod
// │   ├── go.sum
// │   ├── handlers
// │   │   └── info.go
// │   ├── main
// │   ├── models
// │   │   └── models.go
// │   └── server
// │       └── server.go
// └── test
//     ├── createOrder.go
//     ├── order1.json
//     └── testReq.go
