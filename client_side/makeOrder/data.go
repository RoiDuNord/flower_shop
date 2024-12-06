package c_order

func initializeData() (map[string]Flower, map[string]Postcard, map[string]Pack) {
	flowers := map[string]Flower{
		"redRose":    {Name: "Роза", Color: "Красная", Price: 80},
		"whiteRose":  {Name: "Роза", Color: "Белая", Price: 60},
		"yellowRose": {Name: "Роза", Color: "Жёлтая", Price: 40},
		"whiteLily":  {Name: "Лилия", Color: "Белая", Price: 100},
		"yellowLily": {Name: "Лилия", Color: "Жёлтая", Price: 90},
		"pinkPion":   {Name: "Пион", Color: "Розовый", Price: 120},
		"whitePion":  {Name: "Пион", Color: "Белый", Price: 110},
		"lotus":      {Name: "Лотос", Color: "Белый", Price: 200},
		"daisy":      {Name: "Ромашка", Color: "Белая", Price: 20},
	}

	postcards := map[string]Postcard{
		"birthday":         {Note: "С Днём рождения!", Price: 5},
		"newYear":          {Note: "С Новым Годом!", Price: 1},
		"happyWedding":     {Note: "Со свадьбой!", Price: 2},
		"happyAnniversary": {Note: "С Юбилеем!", Price: 3},
		"womenDay":         {Note: "С 8 марта!", Price: 15},
		"valentineDay":     {Note: "С Днём Влюбленных!", Price: 20},
	}

	packs := map[string]Pack{
		"craft": {Material: "Крафт", Price: 100},
		"film":  {Material: "Плёнка", Price: 50},
		"tape":  {Material: "Лента", Price: 10},
	}

	return flowers, postcards, packs
}
