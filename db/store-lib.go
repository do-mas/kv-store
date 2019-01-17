package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

type Store struct {
	db *bolt.DB
}
var bucketName = []byte("my-bucket")

func Open(path string) (*Store, error) {
	opts := &bolt.Options{Timeout: 50 * time.Millisecond,}
	if db, err := bolt.Open(path, 0640, opts); err != nil {
		return nil, err
	} else {
		err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucketName)
			return err
		})
		if err != nil {
			return nil, err
		} else {
			return &Store{db: db}, nil
		}
	}
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) Put(key string, value string) error {
	var eVal bytes.Buffer
	if err := gob.NewEncoder(&eVal).Encode(value); err != nil {
		return err
	}
	fmt.Println(eVal)
	fmt.Println(value)
	return store.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(key), eVal.Bytes())
	})
}


func (store *Store) Get(key string) string {
	var value string
	_ = store.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()
		if k, v := c.Seek([]byte(key)); k == nil {
			return errors.New("skv: bad value")
		} else {
			d := gob.NewDecoder(bytes.NewReader(v))
			_ = d.Decode(&value)
			return nil
		}
	})
	return value
}


func (store *Store) GetAllPairs() [] string {
	var value [] string
	_ = store.db.View(func(tx *bolt.Tx) error {
		_ = tx.Bucket(bucketName).ForEach(func(k, v []byte) error {
			value = append(value, decode(v))
			return nil
		})
		return nil
	});
	return value
}

func (store *Store) GetPairs(n int) [] string {
	var value [] string
	_ = store.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()
		k, v := c.First()
		for i := 1; i < n; i++ {
			k, v = c.Next()
			if k != nil {
				value = append(value, decode(v))
				fmt.Printf("key=%s, value=%s\n", k, v)
			}

		}
		return nil
	})
	return value
}

func (store *Store) Delete(key string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()
		if k, _ := c.Seek([]byte(key)); k != nil {
			fmt.Println("deleting")
			return c.Delete()
		}
		return nil
	})
}

func decode(eVal []byte) string {
	var dVal string
	d := gob.NewDecoder(bytes.NewReader(eVal))
	_ = d.Decode(&dVal)
	return dVal
}