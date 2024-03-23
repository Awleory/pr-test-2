package storage

type Order struct {
	ProductID   int
	ProductName string
	Count       int
	OrderNum    int
}

type Shelf struct {
	ShelfID   int
	ShelfName string
}

type ProductShelf struct {
	ProductID   int
	ShelfID     int
	IsMainShelf bool
}
