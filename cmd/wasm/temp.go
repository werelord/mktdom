package main


func temploadStoredPlanetData() {
	tb := baseProductMap["transport_bot"]
	tb.productTypeExport = productTypeExport{4, 6}
	mr := baseProductMap["miner_robot"]
	mr.productTypeExport = productTypeExport{5, 6}

	var p1 = planet{"planet 331", "sector 1", 15415, 14.0, 134, make(map[string]productType, 0)}
	p1.productList[tb.ID()] = tb
	p1.productList[mr.ID()] = mr

	var p2 = planet{"planet 2", "sector 2", 31337, 57.0, 42, make(map[string]productType, 0)}
	p2.productList[tb.ID()] = tb
	p2.productList[mr.ID()] = mr

	var p3 = planet{"planet3", "sector 3", 95136, 0.0, 1337, make(map[string]productType, 0)}
	p3.productList[mr.ID()] = mr
	p3.productList[tb.ID()] = tb
	var p4 = planet{"planet 4", "sector 1", 15415, 14.0, 134, make(map[string]productType, 0)}
	p4.productList[tb.ID()] = tb
	p4.productList[mr.ID()] = mr
	var p5 = planet{"planet 5", "sector 2", 31337, 57.0, 42, make(map[string]productType, 0)}
	p5.productList[tb.ID()] = tb
	p5.productList[mr.ID()] = mr
	var p6 = planet{"planet6", "sector 3", 95136, 0.0, 1337, make(map[string]productType, 0)}
	p6.productList[tb.ID()] = tb
	p6.productList[mr.ID()] = mr

	planetList[p1.name] = p1
	planetList[p2.name] = p2
	planetList[p3.name] = p1
	planetList[p4.name] = p4
	planetList[p5.name] = p5
	planetList[p6.name] = p6


}