package main

import (
	"fmt"
	"kv-store/main/db"
)

func main() {
	kvs, err := db.Open("db/my.db")
	if err != nil {
		return
	}

	var content = "content"
	err = kvs.Put("key", content)

	var val = kvs.Get("key")
	fmt.Println("value: " + val)

	var arr = kvs.GetAllPairs()
	fmt.Println("array content:", arr)

	kvs.GetPairs(1)

	_ = kvs.Delete("key")
}
