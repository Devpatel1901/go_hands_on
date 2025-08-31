package main

import (
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	db := connectDatabase()

	defer db.Close()

	phoneNumbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	insertRecords(db, phoneNumbers)

	phones, err := allPhones(db)

	if err != nil {
		panic(err)
	}

	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.number)
		if number != p.number {
			fmt.Println("Updating or removing...", number)
			existing, err := findPhone(db, number)

			if err != nil {
				panic(err)
			}

			if existing != nil {
				err = deletePhone(db, p.id)

				if err != nil {
					panic(err)
				}
			} else {
				p.number = number
				err = updatePhone(db, p)

				if err != nil {
					panic(err)
				}
			}
		} else {
			fmt.Println("No changes required")
		}
	}

}
