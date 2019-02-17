package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kv-store/main/db"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/nu7hatch/gouuid"
)

func GetStorage(w http.ResponseWriter, r *http.Request) {

	kvs, err := db.Open()
	if err != nil {
		fmt.Println("error")
		return
	}

	id := mux.Vars(r)["id"]
	contentType := kvs.GetContentType(id)
	content := kvs.GetValue(id)

	fmt.Println(contentType)
	fmt.Println(content)
	// if contentType == "application/octet-stream" {
	w.Header().Set("Content-Disposition", "attachment;")
	// w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Type", contentType)
	http.ServeContent(w, r, "myfile", time.Now(), bytes.NewReader([]byte(content)))
	// } else {
	// json.NewEncoder(w).Encode(content)
	// }

	_ = kvs.Close()

}

func Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	fmt.Println(handler.Header)
	fmt.Println(handler.Filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func CreateStorage(w http.ResponseWriter, r *http.Request) {

	printRequest(r)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var value = string(b)

	kvs, err := db.Open()
	if err != nil {
		return
	}

	u, err := uuid.NewV4()
	kvs.PutValue(u.String(), value)
	kvs.PutType(u.String(), getContentType(r))

	_ = kvs.Close()
	// ---

	json.NewEncoder(w).Encode(u.String())

}

func getContentType(r *http.Request) string {
	contentType := r.Header.Get("Content-type")
	fmt.Println(contentType)
	return contentType

}

func printRequest(r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

}

// func GetFile(w http.ResponseWriter, r *http.Request) {

// w.Header().Set("Content-Disposition", "attachment; filename=WHATEVER_YOU_WANT.pdf")

// fmt.Println("error")
// params := mux.Vars(r)
// var id = params["id"]
// fmt.Println(id)
// kvs, err := db.Open()
// if err != nil {
// 	fmt.Println("error")
// 	return
// }

// var iWasStored = kvs.Get(id)
// _ = kvs.Close()

// fp := path.Join("images", "foo.png")
// http.ServeFile(w, r, fp)

// json.NewEncoder(w).Encode(iWasStored)

// data, err := ioutil.ReadFile("myfile")
// if err != nil {
// log.Fatal(err)
// }

// http.ServeContent(w, r, "myfile", time.Now(), bytes.NewReader([]byte(iWasStored)))

// }

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

// func SaveFile(w http.ResponseWriter, r *http.Request) {

// 	u, err := uuid.NewV4()
// 	var value string
// 	_ = json.NewDecoder(r.Body).Decode(&value)

// 	fmt.Printf("%+v\n", r.Body)
// 	fmt.Printf("%+v\n", value)

// 	kvs, err := db.Open()
// 	if err != nil {
// 		return
// 	}
// 	kvs.Put(u.String(), value)

// 	_ = kvs.Close()

// 	json.NewEncoder(w).Encode(u.String())

// }
