package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"kv-store/main/db"
	"net/http"
)

func GetStorage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	fmt.Println(id)
	kvs, err := db.Open()
	if err != nil {
		fmt.Println("error")
		return
	}

	var iWasStored = kvs.Get(id)
	fmt.Printf("%+v\n", iWasStored)
	_ = kvs.Close()

	json.NewEncoder(w).Encode(iWasStored)

}

func CreateStorage(w http.ResponseWriter, r *http.Request) {

	u, err := uuid.NewV4()
	var value string
	_ = json.NewDecoder(r.Body).Decode(&value)

	fmt.Printf("%+v\n", r.Body)
	fmt.Printf("%+v\n", value)

	kvs, err := db.Open()
	if err != nil {
		return
	}
	kvs.Put(u.String(), value)

	_ = kvs.Close()

	json.NewEncoder(w).Encode(u.String())

}
