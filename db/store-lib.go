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

type MyStruct struct {
	Val string
}

var bucketName = []byte("my-bucket")

func Open() (*Store, error) {
	opts := &bolt.Options{Timeout: 50 * time.Millisecond,}
	if db, err := bolt.Open("db/my.db", 0640, opts); err != nil {
		return nil, err
	} else {
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucketName)
			return err
		})
		if err != nil {
			return nil, err
		}
		return &Store{db: db}, nil

	}
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) Put(key string, value MyStruct) {
	var encodedVal bytes.Buffer
	err := gob.NewEncoder(&encodedVal).Encode(value)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	err = store.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(key), encodedVal.Bytes())
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func (store *Store) Get(key string) *MyStruct {

	var value MyStruct

	err := store.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(bucketName).Cursor()
		k, v := cursor.Seek([]byte(key))
		decodedValue, err := decode(k, v)
		value = decodedValue
		return err
	})

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &value

}

func (store *Store) List(numberOfValues int) []MyStruct {
	var storedValues [] MyStruct
	_ = store.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()
		k, v := c.First()
		storedValues = addToArray(k, v, storedValues)
		for i := 1; i < numberOfValues; i++ {
			if k, v = c.Next(); k != nil {
				storedValues = addToArray(k, v, storedValues)
			}
		}
		return nil
	})
	return storedValues
}

func (store *Store) ListAll() [] MyStruct {
	var values [] MyStruct
	_ = store.db.View(func(tx *bolt.Tx) error {
		_ = tx.Bucket(bucketName).ForEach(func(k, v []byte) error {
			values = addToArray(k, v, values)
			return nil
		})
		return nil
	})
	return values
}

func (store *Store) Delete(key string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()
		if k, _ := c.Seek([]byte(key)); k != nil {
			return c.Delete()
		}
		return nil
	})
}

func addToArray(k []byte, v []byte, storedValues []MyStruct) []MyStruct {
	val, e := decode(k, v)
	if e == nil {
		storedValues = append(storedValues, val)
	}
	return storedValues
}

func decode(k []byte, v []byte) (MyStruct, error) {
	var value MyStruct
	if k != nil {
		d := gob.NewDecoder(bytes.NewReader(v))
		err := d.Decode(&value)
		return value, err
	}
	return value, errors.New("key not found")
}
