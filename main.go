package main

import (
	"fmt"
	"kv-store/main/db"
)

func main() {

	// open connection
	kvs, err := db.Open()
	if err != nil {
		return
	}

	// save
	var key0 = "im-the-key-0"
	var key1 = "im-the-key-1"
	var key2 = "im-the-key-2"
	var storeMe0 = db.MyStruct{Val: "dummy-stuff-0"}
	var storeMe1 = db.MyStruct{Val: "dummy-stuff-1"}
	var storeMe2 = db.MyStruct{Val: "dummy-stuff-2"}
	kvs.Put(key0, storeMe0)
	kvs.Put(key1, storeMe1)
	kvs.Put(key2, storeMe2)

	// get by key
	var iWasStored = kvs.Get(key0)
	fmt.Printf("%+v\n", iWasStored)
	// get by key
	var iWasStored1 = kvs.Get(key1)
	fmt.Printf("%+v\n", iWasStored1)

	// get number of values
	var storedValues = kvs.List(2)
	fmt.Println(len(storedValues))
	fmt.Printf("%+v\n", storedValues)

	// get all stored values
	var allStoredValues = kvs.ListAll()
	fmt.Println(len(allStoredValues))
	fmt.Printf("%+v\n", allStoredValues)


	// delete by key
	_ = kvs.Delete(key0)

	allStoredValues = kvs.ListAll()
	fmt.Println(len(allStoredValues))
	fmt.Printf("%+v\n", allStoredValues)

	// close connection
	_ = kvs.Close()

}
