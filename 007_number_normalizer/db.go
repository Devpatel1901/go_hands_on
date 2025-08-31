package main

import (
	"database/sql"
	"fmt"
	"regexp"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Don'tgiveup2easily"
	dbname   = "go_hands_on_demo"
)

type phone struct {
	id     int
	number string
}

func connectDatabase() *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable\n", host, port, user, password, dbname))

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func insertRecords(db *sql.DB, phoneNumbers []string) {
	sqlStatement := `INSERT INTO contacts (phone_number) VALUES ($1) RETURNING id`
	id := 0

	for _, number := range phoneNumbers {
		err := db.QueryRow(sqlStatement, number).Scan(&id)

		if err != nil {
			panic(err)
		}

		fmt.Println("New record ID is:", id)
	}
}

func allPhones(db *sql.DB) ([]phone, error) {
	rows, err := db.Query("SELECT id, phone_number FROM contacts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func findPhone(db *sql.DB, number string) (*phone, error) {
	var p phone
	row := db.QueryRow("SELECT id, phone_number FROM contacts WHERE phone_number=$1", number)
	err := row.Scan(&p.id, &p.number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func deletePhone(db *sql.DB, id int) error {
	statement := `DELETE FROM contacts WHERE id=$1`
	_, err := db.Exec(statement, id)
	return err
}

func updatePhone(db *sql.DB, p phone) error {
	statement := `UPDATE contacts SET phone_number=$2 WHERE id=$1`
	_, err := db.Exec(statement, p.id, p.number)
	return err
}

func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}
