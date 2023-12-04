package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/db_contact")
	if err != nil {
		return nil, err
	}
	return db, nil
}

type contact struct {
	Id    string
	Nama  string
	Phone string
}

type Rumus struct {
	Total float32
}

func routerIndexGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		rows, err := db.Query("select * from tb_contact")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer rows.Close()

		var contacts []contact
		for rows.Next() {
			var each = contact{}
			var err = rows.Scan(&each.Id, &each.Nama, &each.Phone)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			contacts = append(contacts, each)
		}
		if err = rows.Err(); err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, each := range contacts {
			fmt.Println(each.Nama, each.Id, each.Phone)
		}
		var tmpl = template.Must(template.New("home").ParseFiles("index.html"))
		err = tmpl.Execute(w, contacts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "method not found", http.StatusBadRequest)
}

func routerCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("create").ParseFiles("index.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if r.Method == "POST" {
		db, err := connect()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		var nama = r.FormValue("nama")
		var phone = r.Form.Get("phone")
		_, err = db.Exec("insert into tb_contact (nama,phone)values (?, ?)", nama, phone)
		if err != nil {
			http.Redirect(w, r, "/create", http.StatusMovedPermanently)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	http.Error(w, "Bad Request", http.StatusBadRequest)
}
func routeDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")

		db, err := connect()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		_, err = db.Exec("delete from tb_contact where id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func routerEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		db, err := connect()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var contact = contact{}
		err = db.
			QueryRow("select * from tb_contact where id = ?", id).
			Scan(&contact.Id, &contact.Nama, &contact.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var tmpl = template.Must(template.New("edit").ParseFiles("index.html"))
		err = tmpl.Execute(w, contact)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if r.Method == "POST" {
		id := r.URL.Query().Get("id")
		db, err := connect()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		var nama = r.FormValue("nama")
		var phone = r.Form.Get("phone")
		_, err = db.Exec("update tb_contact set nama = ? , phone = ? where id =?", nama, phone, id)
		if err != nil {
			http.Redirect(w, r, "/edit", http.StatusMovedPermanently)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	http.Error(w, "Bad Request", http.StatusBadRequest)
}
func main() {
	http.HandleFunc("/home", routerIndexGet)
	http.HandleFunc("/", routerIndexGet)
	http.HandleFunc("/delete", routeDelete)
	http.HandleFunc("/create", routerCreate)
	http.HandleFunc("/edit", routerEdit)
	fmt.Println("sukses")
	http.ListenAndServe(":9000", nil)

}
