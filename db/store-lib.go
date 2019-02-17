package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type Store struct {
	db *bolt.DB
}

type MyStruct struct {
	Value string
	Type  string
}

var valueBucket = []byte("value-bucket")
var typeBucket = []byte("type-bucket")

func Open() (*Store, error) {
	opts := &bolt.Options{Timeout: 50 * time.Millisecond}
	if db, err := bolt.Open("db/my.db", 0640, opts); err != nil {
		return nil, err
	} else {
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(valueBucket)
			return err
		})
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(typeBucket)
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

func (store *Store) PutValue(key string, value string) {
	var encodedVal bytes.Buffer
	err := gob.NewEncoder(&encodedVal).Encode(value)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	err = store.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(valueBucket).Put([]byte(key), encodedVal.Bytes())
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
}
func (store *Store) PutType(key string, contentType string) {
	var encodedVal bytes.Buffer
	err := gob.NewEncoder(&encodedVal).Encode(contentType)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	err = store.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(typeBucket).Put([]byte(key), encodedVal.Bytes())
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func (store *Store) GetValue(key string) string {

	var value string

	err := store.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(valueBucket).Cursor()
		k, v := cursor.Seek([]byte(key))
		decodedValue, err := decode(k, v)
		value = decodedValue
		fmt.Printf("getting")
		fmt.Printf("%+v\n", v)
		fmt.Printf("%+v\n", decodedValue)
		return err
	})

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return value

}
func (store *Store) GetContentType(key string) string {

	var value string

	err := store.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(typeBucket).Cursor()
		k, v := cursor.Seek([]byte(key))
		decodedValue, err := decode(k, v)
		value = decodedValue
		fmt.Printf("getting")
		fmt.Printf("%+v\n", v)
		fmt.Printf("%+v\n", decodedValue)
		return err
	})

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return value

}

func (store *Store) List(numberOfValues int) []string {
	var storedValues []string
	_ = store.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(valueBucket).Cursor()
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

func (store *Store) ListAll() []string {
	var values []string
	_ = store.db.View(func(tx *bolt.Tx) error {
		_ = tx.Bucket(valueBucket).ForEach(func(k, v []byte) error {
			values = addToArray(k, v, values)
			return nil
		})
		return nil
	})
	return values
}

func (store *Store) Delete(key string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket(valueBucket).Cursor()
		if k, _ := c.Seek([]byte(key)); k != nil {
			return c.Delete()
		}
		return nil
	})
}

func addToArray(k []byte, v []byte, storedValues []string) []string {
	val, e := decode(k, v)
	if e == nil {
		storedValues = append(storedValues, val)
	}
	return storedValues
}

func decode(k []byte, v []byte) (string, error) {
	var value string
	if k != nil {
		d := gob.NewDecoder(bytes.NewReader(v))
		err := d.Decode(&value)
		return value, err
	}
	return value, errors.New("key not found")
}
