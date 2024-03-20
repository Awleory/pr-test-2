package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"test/postgres"
)

func main() {
	arguments := os.Args

	if len(arguments) < 2 {
		log.Fatal("enter order number")
	}

	var orders []int
	for _, v := range arguments[1:] {
		num, err := strconv.Atoi(v)
		if err != nil {
			log.Println("unable to convert to int")
		}
		orders = append(orders, num)
	}

	db, err := postgres.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if result, err := db.OrdersPrintForm(orders); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}
