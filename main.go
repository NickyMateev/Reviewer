package main

import (
	"database/sql"
	"fmt"
	"github.com/NickyMateev/Reviewer/api"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/NickyMateev/Reviewer/web"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open(storage.DbType, storage.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = storage.UpdateSchema(db)
	fmt.Println(err)

	router := buildRouter(api.Default(db))
	log.Fatal(http.ListenAndServe(":8888", router))
}

func buildRouter(api web.API) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	controllers := api.Controllers()
	for _, controller := range controllers {
		routes := controller.Routes()
		for _, route := range routes {
			router.HandleFunc(route.Path, route.HandleFunc).Methods(route.Method)
		}
	}

	return router
}
