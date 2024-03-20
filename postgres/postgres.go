package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type orderPrintForm struct {
	shelfTitle   string
	shelfId      int
	orderId      int
	productTitle string
	productId    int
	productCount int
	subShelfs    *[]subShelf
}

type subShelf struct {
	id    int
	title string
}

func (pf orderPrintForm) PrintForm() string {
	result := fmt.Sprintf("%s (id=%d)\n", pf.productTitle, pf.productId)
	result += fmt.Sprintf("заказ %d, %d шт\n", pf.orderId, pf.productCount)
	if len(*pf.subShelfs) > 0 {
		var subShelfsStr []string
		for _, v := range *pf.subShelfs {
			subShelfsStr = append(subShelfsStr, v.title)
		}
		result += fmt.Sprintf("доп стеллаж: %v\n", strings.Join(subShelfsStr, ","))
	}
	result += "\n"
	return result
}

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

func (db *DB) ordersPrintStruct(orders []int) (map[int][]orderPrintForm, error) {

	query :=
		`SELECT 
	order_num as orderId,
	product_orders.product_id as productId,
	product.title as productTitle,
	shelf.id as shelfId,
	shelf.title as shelfTitle,
	is_main_shelf as isMainShelf,
	count as productCount
	FROM product_orders
	LEFT JOIN product 
	ON product.id=product_orders.product_id
	LEFT JOIN product_shelf 
	ON product_shelf.product_id = product_orders.product_id
	LEFT JOIN shelf 
	ON shelf.id = shelf_id
	Where order_num = ANY($1)
	Order by is_main_shelf desc`

	rows, err := db.sqlDB.Query(query, pq.Array(orders))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shelfIdPrintForm := make(map[int][]orderPrintForm)
	productIdPrintForm := make(map[int][]orderPrintForm)
	for rows.Next() {
		printForm := orderPrintForm{}
		isMainShelf := false
		if err := rows.Scan(
			&printForm.orderId,
			&printForm.productId,
			&printForm.productTitle,
			&printForm.shelfId,
			&printForm.shelfTitle,
			&isMainShelf,
			&printForm.productCount,
		); err != nil {
			return nil, err
		}

		sub := make([]subShelf, 0)
		printForm.subShelfs = &sub

		if !isMainShelf {

			if array, ok := productIdPrintForm[printForm.productId]; ok {
				for i := 0; i < len(array); i++ {
					*productIdPrintForm[printForm.productId][i].subShelfs =
						append(*productIdPrintForm[printForm.productId][i].subShelfs, subShelf{printForm.shelfId, printForm.shelfTitle})
				}
			}
			continue
		}

		if _, ok := shelfIdPrintForm[printForm.shelfId]; !ok {
			printForms := make([]orderPrintForm, 0)
			shelfIdPrintForm[printForm.shelfId] = printForms
		}

		if _, ok := productIdPrintForm[printForm.productId]; !ok {
			printForms := make([]orderPrintForm, 0)
			productIdPrintForm[printForm.productId] = printForms
		}

		shelfIdPrintForm[printForm.shelfId] = append(shelfIdPrintForm[printForm.shelfId], printForm)
		productIdPrintForm[printForm.productId] = append(productIdPrintForm[printForm.productId], printForm)
	}

	return shelfIdPrintForm, nil
}

func (db *DB) OrdersPrintForm(orders []int) (string, error) {
	printForms, err := db.ordersPrintStruct(orders)
	if err != nil {
		return "", err
	}

	result := "=+=+=+=\n"
	result += fmt.Sprintf("Страница сборки заказов %v\n\n", orders)

	for _, pfs := range printForms {
		if len(pfs) > 0 {
			result += fmt.Sprintf("===Стеллаж %s\n", pfs[0].shelfTitle)
		}
		for _, pf := range pfs {
			result += pf.PrintForm()
		}
	}

	return result, nil
}
