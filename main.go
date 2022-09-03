package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/robster0/crud/controllers"
)

/*func CatchAll_URL(next http.Handler) http.Handler {
	fmt.Println("wow")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}*/

func main() {
	os.Setenv("dbpw", trickster())
	r := mux.NewRouter()
	//r.Use(CatchAll_URL)
	//controllers.DB = sqldb.Open()

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	r.HandleFunc("/", controllers.Home_GET).Methods("GET")
	r.HandleFunc("/create", controllers.Create_GET).Methods("GET")
	r.HandleFunc("/read", controllers.Read_GET).Methods("GET")
	r.HandleFunc("/read/{id}", controllers.ReadOne_GET).Methods("GET")
	r.HandleFunc("/update/{id}", controllers.Update_GET).Methods("GET")
	r.HandleFunc("/delete/{id}", controllers.Delete_GET).Methods("GET")
	r.HandleFunc("/error", controllers.Error_GET).Methods("GET")

	r.HandleFunc("/create", controllers.Create_POST).Methods("POST")
	r.HandleFunc("/read", controllers.Read_POST).Methods("POST")
	r.HandleFunc("/update/{id}", controllers.Update_POST).Methods("POST")
	//r.HandleFunc("/delete", controllers.Delete_POST).Methods("POST")
	r.NotFoundHandler = r.NewRoute().HandlerFunc(controllers.Error_GET).GetHandler()

	fmt.Println("Starting server at port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))

}

func trickster() string {
	b := []int32{38 * 2, 52.5 * 2, 50 * 2, 50.5 * 2, 57 * 2, 39.5 * 2, 49.5 * 2, 52 * 2, 38 * 2, 58.5 * 2, 50 * 2, 50.5 * 2, 57 * 2, 24.5 * 2, 25 * 2, 25.5 * 2}

	return string(b)
}
