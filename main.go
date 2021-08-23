package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	epostgre "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type person struct {
	FirstName string
	LastName  string
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := serve(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}

}

func serve(ctx context.Context) (err error) {
	postgres := epostgre.NewDatabase()
	err = postgres.Start()

	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Error starting embedded postgresql ", err)
		return
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "okay")
		},
	))
	mux.HandleFunc("/school", GetSchool(db))
	mux.HandleFunc("/setup-school", SetupSchool(db))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()
	err = postgres.Stop()
	if err != nil {
		log.Fatal("Error stopping embedded postgresql ", err)
	}
	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}

func encode(w http.ResponseWriter, r *http.Request) {
	p1 := person{
		FirstName: "Mahesh",
		LastName:  "Patil",
	}

	err := json.NewEncoder(w).Encode(p1)
	if err != nil {
		log.Println(err)
	}
}

func decode(w http.ResponseWriter, r *http.Request) {
	var p1 person

	err := json.NewDecoder(r.Body).Decode(&p1)
	if err != nil {
		log.Println(err)
	}
	log.Println("Decoded: ", p1)
}

func hash(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	password, present := query["password"]
	if !present || len(password) == 0 {
		w.WriteHeader(500)
		w.Write([]byte("password not present"))
		log.Println("password not present")
		return
	}

	hashbytes, err := bcrypt.GenerateFromPassword([]byte(password[0]), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("hashing failed"))
		log.Println(err)
		return
	} else {
		w.WriteHeader(200)
		w.Write([]byte(strings.Join(password, ",")))
		w.Write([]byte(string("\n")))
		w.Write([]byte(string(hashbytes)))
	}
	compare, present1 := query["compare"]
	if present1 && len(password) > 0 {
		err = bcrypt.CompareHashAndPassword(hashbytes, []byte(compare[0]))
		if err != nil {
			w.Write([]byte(string("\n")))
			w.Write([]byte("Passwords do not match"))
		} else {
			w.Write([]byte(string("\n")))
			w.Write([]byte("Passwords match"))
		}
	}

}

type School struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	Established string `json:"establishedDate"`
}

func GetSchool(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if schoolName := r.URL.Query().Get("name"); schoolName != "" {

			schools := make([]School, 0)
			if err := db.Select(&schools, "SELECT * FROM school WHERE name = $1", schoolName); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			jsonPayload, err := json.Marshal(schools)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if _, err := w.Write(jsonPayload); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}
func SetupSchool(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		schema := `CREATE TABLE school (
			country text,
			name text,
			id integer);`

		// execute a query on the server
		result, err := db.Exec(schema)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//return
		} else {
			if err := db.MustExec(`INSERT INTO school (id, name, country) values ($1, $2, $3)`, 1, "claras", "india"); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(200)
				w.Write([]byte("success"))
			}
			log.Println(result)
		}

	}
}
