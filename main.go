package kv_store

import (
	"kv-store/main/store"
	"fmt"
)

func main() {
	store, err := store.Open("my.db")
	var info = "rtyryrytrytr"

	err = store.Put("shbhjbjmmm35", info)
	if err != nil {
		return
	}

	var val = store.Get("ses34d135")
	fmt.Println("value " + val);
	var arr = store.GetPairs(1)
	store.GetPairs1(5)
	fmt.Println("array value" + arr[0])
	fmt.Println("array size", len(arr))
	fmt.Println("2d: ", arr)
	store.Delete("ses34d135")
}