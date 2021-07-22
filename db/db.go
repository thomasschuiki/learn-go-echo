package db

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/thomasschuiki/learn-go-echo/models"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

var userBucket = []byte("users")
var db *bolt.DB

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
	var u models.User
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucket)
		u.Name = name
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
		u.Password = string(hashedPassword)
		user, err := json.Marshal(u)
		if err != nil {
			return err
		}
		return b.Put([]byte(u.Name), user)
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func GetUser(name string) (models.User, error) {
	var user models.User
	err := db.View(func(t *bolt.Tx) error {
		b := t.Bucket(userBucket)
		v := b.Get([]byte(name))

		json.Unmarshal(v, &user)
		return nil
	})
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func AllUsers() ([]models.User, error) {
	var users []models.User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var u models.User
			json.Unmarshal(v, &u)
			users = append(users, u)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucket)
		return b.Delete([]byte(key))
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
