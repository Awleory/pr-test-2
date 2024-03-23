package model

import (
	"fmt"
	"strings"
)

type orderRow struct {
	orderId       int
	shelfId       int
	shelfTitle    string
	productId     int
	productTitle  string
	productCount  int
	subShelfTable subShelfTable
}

func NewOrderRow() orderRow {
	return orderRow{
		subShelfTable: NewSubShelfTable(),
	}
}

type subShelf struct {
	productId int
	shelfId   int
	title     string
}

type subShelfTable struct {
	rows *[]subShelf
}

func (table subShelfTable) Len() int {
	return len(*table.rows)
}

func NewSubShelfTable() subShelfTable {
	rows := make([]subShelf, 0)
	return subShelfTable{&rows}
}

type ordersTable struct {
	rows *[]orderRow
}

func NewOrdersTable() ordersTable {
	table := ordersTable{}
	rows := make([]orderRow, 0)
	table.rows = &rows
	return table
}

func (table ordersTable) AddRow() *orderRow {
	newRow := NewOrderRow()
	*table.rows = append(*table.rows, newRow)
	return &newRow
}

func (row orderRow) PrintForm() string {
	result := fmt.Sprintf("%s (id=%d)\n", row.productTitle, row.productId)
	result += fmt.Sprintf("заказ %d, %d шт\n", row.orderId, row.productCount)
	if row.subShelfTable.Len() > 0 {
		var subShelfsStr []string
		for _, v := range *row.subShelfTable.rows {
			subShelfsStr = append(subShelfsStr, v.title)
		}
		result += fmt.Sprintf("доп стеллаж: %v\n", strings.Join(subShelfsStr, ","))
	}
	result += "\n"
	return result
}
