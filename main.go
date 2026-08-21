package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var db *sql.DB

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	users := getUsersFromDB(db)
	html, err := template.ParseFiles("views/view.html")
	check(err)

	err = html.Execute(writer, GuestbookInfo{len(users), users})
	check(err)
}

func newHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("views/new.html")
	check(err)

	err = html.Execute(writer, nil)
	check(err)
}

func profileHandler(writer http.ResponseWriter, request *http.Request) {
	urlPath := strings.Split(request.URL.Path, "/")[1:]

	if len(urlPath) != 3 || urlPath[2] == "" {
		http.NotFound(writer, request)
		return
	}

	idUrlSegment := urlPath[2]
	u, userItems, err := getUserAndUserItemsFromDB(db, idUrlSegment, writer)
	if err != nil {
		return
	}

	print(userItems)

	html, err := template.ParseFiles("views/profile.html")
	check(err)
	err = html.Execute(writer, UserProfile{User: u, Items: userItems})
	check(err)
}

func createItemHandler(writer http.ResponseWriter, request *http.Request) {
	urlQuery := request.URL.Query()
	userIdStr := urlQuery.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Println("err is ", err)
		http.NotFound(writer, request)
		return
	}

	print(userId)
	err = insertItemToDB(userId, request)

	if err != nil {
		check(err)
	}

	http.Redirect(writer, request, "/guestbook/profile/"+userIdStr, http.StatusFound)
}

func newItemHandler(writer http.ResponseWriter, request *http.Request) {
	urlQuery := request.URL.Query()
	userIdStr := urlQuery.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Println("err is ", err)
		http.NotFound(writer, request)
		return
	}

	print(userId)

	html, err := template.ParseFiles("views/newItem.html")
	check(err)
	err = html.Execute(writer, userId)
	check(err)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	userName := request.FormValue("user_name")
	userAge := request.FormValue("age")
	userCity := request.FormValue("city")
	fieldErrors := validateUserCreateFormFields(userName, userAge, userCity)

	if len(fieldErrors) != 0 {
		html, err := template.ParseFiles("views/new.html")
		check(err)

		err = html.Execute(writer, fieldErrors)
		check(err)
		return
	}

	err := insertUserToDB(userName, userAge, userCity)

	if err != nil {
		log.Printf("INSERT error: %#v", err)
		check(err)
	}

	http.Redirect(writer, request, "/guestbook", http.StatusFound)
}

func main() {
	var err error
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	db, err = runDB(db, err)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}

	defer db.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/guestbook", viewHandler)
	http.HandleFunc("/guestbook/new", newHandler)
	http.HandleFunc("/guestbook/create", createHandler)
	http.HandleFunc("/guestbook/profile/", profileHandler)
	http.HandleFunc("/guestbook/items/new", newItemHandler)
	http.HandleFunc("/guestbook/items/create", createItemHandler)

	log.Println("running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
