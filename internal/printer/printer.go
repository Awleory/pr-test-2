package printer

import (
	"fmt"
	"strings"
	"test/internal/storage"
)

type MainShelf struct {
	ShelfName   string
	ProductName string
	ProductID   int
	OrderNum    int
	Count       int
}

type SubShelf struct {
	ProductID int
	ShelfName string
}

func Printer(shelfs []storage.Shelf, productShelf []storage.ProductShelf, orders []storage.Order, nums []int) {
	var mainShelfs []MainShelf
	var subShelfs []SubShelf

	for _, shelf := range shelfs {
		for _, conn := range productShelf {
			for _, order := range orders {
				for _, n := range nums {
					if shelf.ShelfID == conn.ShelfID && order.ProductID == conn.ProductID && conn.IsMainShelf && n == order.OrderNum {
						mainShelfs = append(mainShelfs, MainShelf{
							ShelfName:   shelf.ShelfName,
							ProductName: order.ProductName,
							ProductID:   order.ProductID,
							OrderNum:    order.OrderNum,
							Count:       order.Count,
						})
					} else if shelf.ShelfID == conn.ShelfID && order.ProductID == conn.ProductID && conn.IsMainShelf == false && n == order.OrderNum {

						subShelfs = append(subShelfs, SubShelf{
							ProductID: order.ProductID,
							ShelfName: shelf.ShelfName,
						})
					}
				}
			}
		}
	}

	for k, m := range mainShelfs {
		if k == 0 {
			fmt.Printf("===Стеллаж %s\n", m.ShelfName)
			Output(m, subShelfs)
		} else {
			if m.ShelfName == mainShelfs[k-1].ShelfName {
				Output(m, subShelfs)
			} else {
				fmt.Printf("===Стеллаж %s\n", m.ShelfName)
				Output(m, subShelfs)
			}
		}
	}
}

func Output(mSh MainShelf, adSh []SubShelf) {
	fmt.Printf("%s (id=%d)\n", mSh.ProductName, mSh.ProductID)
	fmt.Printf("заказ %d, %d шт\n", mSh.OrderNum, mSh.Count)
	var addShStr []string
	for _, a := range adSh {
		if a.ProductID == mSh.ProductID {
			addShStr = append(addShStr, a.ShelfName)
		}
	}
	if addShStr != nil {

		fmt.Printf("доп стеллаж:%v\n", strings.Join(addShStr, ","))

	}
	fmt.Println()
}
