package main

import (
	"fmt"
	"kv-store/main/db"
)

func main() {

	var my_map = make(map [string] string)
	my_map["1"] = "123"
	//fmt.Printf(my_map["1"])

	kvs, err := db.Open("db/my.db")
	if err != nil {
		return
	}

	var content = "content"
	err = kvs.Put("key", content)
	err = kvs.Put("key1", "n")
	err = kvs.Put("key2", "uhi")
	err = kvs.Put("key3", "uhiuhiu")
	err = kvs.Put("key4", "uhiuhiu")

	my_map = kvs.GetPairs(8)

	fmt.Println(len(my_map))
	fmt.Println(my_map["key"])
	fmt.Println(my_map["key1"])
	fmt.Println(my_map["key2"])
	fmt.Println(my_map["key3"])
	fmt.Println(my_map["key4"])

	//var val = kvs.Get("key")
	//fmt.Println("value: " + val)
	//
	var arr = kvs.GetAllPairs()
	fmt.Println("array content:", arr)
	//
	//kvs.GetPairs(1)
	//
	//_ = kvs.Delete("key")
}
