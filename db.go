package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func getUsersFromDB(db *sql.DB) []User {
	rows, err := db.Query("SELECT id, name, city, age FROM users")
	check(err)
	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var name string
		var city string
		var age int
		var id int

		err = rows.Scan(&id, &name, &city, &age)
		check(err)
		users = append(
			users,
			User{Id: id, Name: name, City: city, Age: age},
		)
	}

	return users
}

func getUserAndUserItemsFromDB(
	db *sql.DB,
	userIdSegment string,
	writer http.ResponseWriter,
) (User, []Item, error) {
	var u User

	err := db.QueryRow(
		"SELECT id, name, city, age FROM users WHERE id = $1", userIdSegment).Scan(
		&u.Id,
		&u.Name,
		&u.City,
		&u.Age,
	)

	if err == sql.ErrNoRows {
		http.Error(writer, "Пользователь не найден", http.StatusNotFound)
		return u, []Item{}, err
	} else if err != nil {
		log.Printf("error: %#v", err)
		http.Error(writer, "Ошибка сервера", http.StatusInternalServerError)
		return u, []Item{}, err
	}

	var rows *sql.Rows

	rows, err = db.Query(`
		SELECT title, description, price, status
		FROM items
		INNER JOIN users ON items.user_id = users.id
		WHERE users.id = $1`,
		u.Id,
	)

	check(err)
	defer rows.Close()

	userItems := make([]Item, 0)

	for rows.Next() {
		var id int
		var title string
		var description string
		var price string
		var status string

		err = rows.Scan(&title, &description, &price, &status)
		check(err)

		userItems = append(
			userItems,
			Item{
				Id:          id,
				Title:       title,
				UserId:      u.Id,
				Description: description,
				Price:       price,
				Status:      status,
			},
		)
	}

	return u, userItems, nil
}

func runDB(db *sql.DB, err error) (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	log.Println(host, user, password, dbname)

	// connStr := "host=localhost port=5432 user=postgres password=123456 dbname=owner sslmode=disable"
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		user,
		password,
		dbname,
	)

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Println("Ошибка открытия соединения:", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println("Ошибка подключения:", err)
		return nil, err
	}

	log.Println("Подключились к БД!")
	return db, nil
}

func insertItemToDB(userId int, request *http.Request) error {
	title := request.FormValue("title")
	description := request.FormValue("description")
	price := request.FormValue("price")

	var id int

	insertQuery := (`
		INSERT INTO items
		(user_id, title, description, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`)

	err := db.QueryRow(insertQuery, userId, title, description, price).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

func insertUserToDB(userName string, age string, city string) error {
	var id int

	insertQuery := "INSERT INTO users (name, age, city) VALUES ($1, $2, $3) RETURNING id"
	err := db.QueryRow(insertQuery, userName, age, city).Scan(&id)
	return err
}
