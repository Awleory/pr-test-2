package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type orderRow struct {
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

type ordersTable struct {
}

type intIndex struct {
}

func (pf orderRow) PrintForm() string {
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

func (db *DB) selectOrdersStructByShelf(orders []int) (map[int][]orderRow, error) {
	if len(orders) == 0 {
		return nil, fmt.Errorf("orders is empty")
	}

	addIndex := func(mp *map[int][]orderRow, key int, printForm *orderRow) {
		if _, ok := (*mp)[key]; !ok {
			printForms := make([]orderRow, 0)
			(*mp)[key] = printForms
		}
		(*mp)[key] = append((*mp)[key], *printForm)
	}

	productIdIndex := make(map[int][]orderRow)
	query := `SELECT
	product_id,
	order_num,
	count
	FROM product_orders`
	if rows, err := db.sqlDB.Query(query); err == nil {
		defer rows.Close()

		ordersMap := make(map[int]int)
		for _, v := range orders {
			ordersMap[v] = 0
		}
		for rows.Next() {
			printForm := orderRow{}
			if err := rows.Scan(
				&printForm.productId,
				&printForm.orderId,
				&printForm.productCount,
			); err != nil {
				return nil, err
			}

			if _, ok := ordersMap[printForm.orderId]; !ok {
				continue
			}

			sub := make([]subShelf, 0)
			printForm.subShelfs = &sub

			addIndex(&productIdIndex, printForm.productId, &printForm)
		}
	} else {
		return nil, err
	}

	shelfIdIndex := make(map[int][]orderRow)
	subShelfsIndex := make(map[int][]orderRow)
	mainShelfsIndex := make(map[int][]orderRow)
	query = `SELECT
	product_id,
	shelf_id,
	is_main_shelf
	FROM product_shelf`
	if rows, err := db.sqlDB.Query(query); err == nil {
		defer rows.Close()

		var productId int
		var shelfId int
		var isMainShelf bool
		mapCount := len(productIdIndex)
		for rows.Next() {
			if mapCount == 0 {
				break
			}
			if err := rows.Scan(
				&productId,
				&shelfId,
				&isMainShelf,
			); err != nil {
				return nil, err
			}

			if _, ok := productIdIndex[productId]; ok {
				for i := 0; i < len(productIdIndex[productId]); i++ {
					productIdIndex[productId][i].shelfId = shelfId
					addIndex(&shelfIdIndex, shelfId, &(productIdIndex[productId][i]))

					if isMainShelf {
						addIndex(&subShelfsIndex, shelfId, &(productIdIndex[productId][i]))
					} else {
						addIndex(&mainShelfsIndex, shelfId, &(productIdIndex[productId][i]))
					}
				}
			} else {
				printForm := orderRow{
					productId: productId,
					shelfId:   shelfId,
				}

				addIndex(&shelfIdIndex, shelfId, &printForm)
				addIndex(&productIdIndex, productId, &printForm)
				if isMainShelf {
					addIndex(&subShelfsIndex, shelfId, &printForm)
				} else {
					addIndex(&mainShelfsIndex, shelfId, &printForm)
				}
			}

			mapCount--
		}
	}

	fmt.Println("222222")
	fmt.Println(printFormByIndex(productIdIndex))

	query = `SELECT
	id,
	title
	FROM product`
	if rows, err := db.sqlDB.Query(query); err == nil {
		defer rows.Close()

		var id int
		var title string
		mapCount := len(productIdIndex)
		for rows.Next() {
			if mapCount == 0 {
				break
			}
			if err := rows.Scan(
				&id,
				&title,
			); err != nil {
				return nil, err
			}

			if _, ok := productIdIndex[id]; ok {
				for i := 0; i < len(productIdIndex[id]); i++ {
					productIdIndex[id][i].productTitle = title
				}
			}

			mapCount--
		}
	} else {
		return nil, err
	}

	query = `SELECT
	id,
	title
	FROM shelf`
	if rows, err := db.sqlDB.Query(query); err == nil {
		defer rows.Close()

		var id int
		var title string
		mapCount := len(shelfIdIndex)
		for rows.Next() {
			if mapCount == 0 {
				break
			}
			if err := rows.Scan(
				&id,
				&title,
			); err != nil {
				return nil, err
			}

			if _, ok := shelfIdIndex[id]; ok {
				for i := 0; i < len(shelfIdIndex[id]); i++ {
					shelfIdIndex[id][i].shelfTitle = title
				}
			}

			mapCount--
		}
	} else {
		return nil, err
	}

	fmt.Println(len(subShelfsIndex))
	for index := range subShelfsIndex {
		for _, v := range subShelfsIndex[index] {
			if _, ok := mainShelfsIndex[v.productId]; ok {
				for i := 0; i < len(mainShelfsIndex[v.productId]); i++ {
					*mainShelfsIndex[v.productId][i].subShelfs =
						append(*mainShelfsIndex[v.productId][i].subShelfs,
							subShelf{v.shelfId, v.shelfTitle})
				}
			}
		}
	}

	return shelfIdIndex, nil
}

func (db *DB) OrdersPrintForm(orders []int) (string, error) {
	printForms, err := db.selectOrdersStructByShelf(orders)
	if err != nil {
		return "", err
	}

	result := "=+=+=+=\n"
	result += fmt.Sprintf("Страница сборки заказов %v\n\n", orders)
	if str, err := printFormByIndex(printForms); err != nil {
		return "", err
	} else {
		result += str
	}

	return result, nil
}

func printFormByIndex(printForms map[int][]orderRow) (string, error) {

	result := ""
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
