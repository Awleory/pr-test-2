package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"test/internal/printer"
	"test/internal/storage/postgres"
)

func main() {
	arguments := os.Args

	if len(arguments) < 2 {
		log.Fatal("enter order number")
	}

	var argsInt []int

	for _, v := range arguments[1:] {
		num, err := strconv.Atoi(v)
		if err != nil {
			log.Println("unable to convert to int")
		}
		argsInt = append(argsInt, num)
	}

	db, err := postgres.NewConn()
	if err != nil {
		log.Fatal("unnable to connect to db")
	}
	fmt.Println("=+=+=+=")
	fmt.Printf("Страница сборки заказов %v\n\n", argsInt)

	orderData, err := db.GetOrder()
	if err != nil {
		log.Fatal("unnable to get order data")
	}

	shelfData, err := db.GetShelf()
	if err != nil {
		log.Fatal("unnable to get shelf data")
	}

	connData, err := db.GetConn()
	if err != nil {
		log.Fatal("unnable to get conn data")
	}

	printer.Printer(shelfData, connData, orderData, argsInt)
}
