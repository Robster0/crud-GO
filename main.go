package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/robster0/crud/controllers"
)

/*func CatchAll_URL(next http.Handler) http.Handler {
	fmt.Println("wow")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}*/

func main() {
	godotenv.Load(".env")
	r := mux.NewRouter()
	//r.Use(CatchAll_URL)
	//controllers.DB = sqldb.Open()

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.PathPrefix("/static/").Handler(s)

	r.HandleFunc("/", controllers.Home_GET).Methods("GET")
	r.HandleFunc("/create/p", controllers.CreateP_GET).Methods("GET")
	r.HandleFunc("/create/c/{id}", controllers.CreateC_GET).Methods("GET")
	r.HandleFunc("/read", controllers.Read_GET).Methods("GET")
	r.HandleFunc("/read/{id}", controllers.ReadOne_GET).Methods("GET")
	r.HandleFunc("/update/{id}", controllers.Update_GET).Methods("GET")
	r.HandleFunc("/delete/{id}", controllers.Delete_GET).Methods("GET")
	r.HandleFunc("/error", controllers.Error_GET).Methods("GET")

	r.HandleFunc("/create/p", controllers.CreateP_POST).Methods("POST")
	r.HandleFunc("/create/c/{id}", controllers.CreateC_POST).Methods("POST")
	r.NotFoundHandler = r.NewRoute().HandlerFunc(controllers.Error_GET).GetHandler()

	fmt.Println("Starting server at port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))

}
