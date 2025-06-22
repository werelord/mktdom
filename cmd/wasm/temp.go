package main

/* func temploadStoredPlanetData() {
	tb := baseProductMap["transport_bot"]
	tb.Supply, tb.Demand = 4, 6
	mr := baseProductMap["miner_robot"]
	mr.Supply, mr.Demand = 5, 6

	var p1 = planet{Name: "planet 331", Sector: "sector 1", market: marketVolume{total: 15415}, DomPoints: 134,
		ProductList: make(map[string]productType, 0)}
	p1.ProductList[tb.ID()] = tb
	p1.ProductList[mr.ID()] = mr

	var p2 = planet{Name: "planet 2", Sector: "sector 2", market: marketVolume{total: 31337}, DomPoints: 42,
		ProductList: make(map[string]productType, 0)}
	p2.ProductList[tb.ID()] = tb
	p2.ProductList[mr.ID()] = mr

	var p3 = planet{Name: "planet 32", Sector: "sector 3", market: marketVolume{total: 95136}, DomPoints: 1337,
		ProductList: make(map[string]productType, 0)}
	p3.ProductList[mr.ID()] = mr
	p3.ProductList[tb.ID()] = tb

	var p4 = planet{Name: "planet 4", Sector: "sector 1", market: marketVolume{total: 15415}, DomPoints: 134,
		ProductList: make(map[string]productType, 0)}
	p4.ProductList[tb.ID()] = tb
	p4.ProductList[mr.ID()] = mr

	var p5 = planet{Name: "planet 5", Sector: "sector 2", market: marketVolume{total: 31337}, DomPoints: 42,
		ProductList: make(map[string]productType, 0)}
	p5.ProductList[tb.ID()] = tb
	p5.ProductList[mr.ID()] = mr

	var p6 = planet{Name: "planet6", Sector: "sector 3", market: marketVolume{total: 95136}, DomPoints: 1337,
		ProductList: make(map[string]productType, 0)}
	p6.ProductList[tb.ID()] = tb
	p6.ProductList[mr.ID()] = mr

	planetMap[p1.Name] = &p1
	planetMap[p2.Name] = &p2
	planetMap[p3.Name] = &p3
	planetMap[p4.Name] = &p4
	planetMap[p5.Name] = &p5
	planetMap[p6.Name] = &p6

}
*/