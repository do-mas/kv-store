package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kv-store/main/db"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/nu7hatch/gouuid"
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
	_ = kvs.Close()

	json.NewEncoder(w).Encode(iWasStored)

}

func GetFile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Disposition", "attachment; filename=WHATEVER_YOU_WANT.pdf")
	w.Header().Set("Content-Type", "application/pdf")

	fmt.Println("error")
	params := mux.Vars(r)
	var id = params["id"]
	fmt.Println(id)
	kvs, err := db.Open()
	if err != nil {
		fmt.Println("error")
		return
	}

	var iWasStored = kvs.Get(id)
	_ = kvs.Close()

	// fp := path.Join("images", "foo.png")
	// http.ServeFile(w, r, fp)

	// json.NewEncoder(w).Encode(iWasStored)

	// data, err := ioutil.ReadFile("myfile")
	// if err != nil {
	// log.Fatal(err)
	// }

	http.ServeContent(w, r, "myfile", time.Now(), bytes.NewReader([]byte(iWasStored)))

}

func CreateStorage(w http.ResponseWriter, r *http.Request) {

	u, err := uuid.NewV4()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var value = string(b)
	fmt.Printf(string(b))

	kvs, err := db.Open()
	if err != nil {
		return
	}
	kvs.Put(u.String(), value)

	_ = kvs.Close()

	json.NewEncoder(w).Encode(u.String())

}

// func CreateStorage(w http.ResponseWriter, r *http.Request) {

// b, err := ioutil.ReadAll(r.Body)
// if err != nil {
// panic(err)
// }
// fmt.Printf(string(b))

// 	requestDump, err := httputil.DumpRequest(r, true)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(string(requestDump))

// }

func SaveFile(w http.ResponseWriter, r *http.Request) {

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
