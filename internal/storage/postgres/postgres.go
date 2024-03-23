package postgres

import (
	"database/sql"
	"test/internal/storage"

	_ "github.com/lib/pq"
)

type DB struct {
	sqlDB *sql.DB
}

func NewConn() (*DB, error) {
	connStr := "user=postgres password=so2037456va dbname=Shop sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &DB{db}, err
}

func (db *DB) Close() error {
	return db.sqlDB.Close()
}

func (db *DB) GetOrder() ([]storage.Order, error) {

	query := `SELECT 
	product.id, 
	title, 
	count, 
	order_num
	FROM product
	JOIN product_orders
	ON product.id=product_orders.product_id;`
	rows, err := db.sqlDB.Query(query)
	if err != nil {

		return nil, err
	}

	var OrderData []storage.Order
	defer rows.Close()

	for rows.Next() {
		var oData storage.Order

		err = rows.Scan(&oData.ProductID, &oData.ProductName, &oData.Count, &oData.OrderNum)
		if err != nil {

			return nil, err
		}

		OrderData = append(OrderData, oData)

	}
	return OrderData, nil

}

func (db *DB) GetShelf() ([]storage.Shelf, error) {
	query := `SELECT 
	title, 
	id
	FROM shelf;`
	rows, err := db.sqlDB.Query(query)
	if err != nil {

		return nil, err
	}

	var ShelfData []storage.Shelf
	defer rows.Close()

	for rows.Next() {
		var sData storage.Shelf

		err = rows.Scan(&sData.ShelfName, &sData.ShelfID)
		if err != nil {
			return nil, err
		}

		ShelfData = append(ShelfData, sData)
	}
	return ShelfData, nil
}

func (db *DB) GetConn() ([]storage.ProductShelf, error) {
	query := `SELECT 
	product_id, 
	shelf_id, 
	is_main_shelf 
	FROM product_shelf;`
	rows, err := db.sqlDB.Query(query)
	if err != nil {

		return nil, err
	}

	var ConnData []storage.ProductShelf
	defer rows.Close()

	for rows.Next() {
		var cData storage.ProductShelf

		err = rows.Scan(&cData.ProductID, &cData.ShelfID, &cData.IsMainShelf)
		if err != nil {

			return nil, err
		}

		ConnData = append(ConnData, cData)
	}
	return ConnData, nil
}
