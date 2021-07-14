package db

import (
	"encoding/binary"
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

var userBucket = []byte("users")
var db *bolt.DB

type User struct {
	id       int
	name     string
	password string
}

func Init(dbpath string) error {
	var err error
	db, err = bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(userBucket)
		return err
	})
}

func CreateUser(name, password string) (int, error) {
	var id int
	var u User
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		u.id = id
		u.name = name
		u.password = password
		user, err := json.Marshal(u)
		if err != nil {
			return err
		}
		return b.Put(key, user)
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
