package database

import (
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/dylanmccormick/go-redis/internal/resp"
	// "github.com/dylanmccormick/go-redis/internal/resp"
)

type RedisObject struct {
	Data any
	Typ  string
	TTL  int64
}
type Database struct {
	mu   sync.Mutex
	name string
	data map[string]*RedisObject
}

func InitializeDB() *Database {
	return &Database{
		mu:   sync.Mutex{},
		data: map[string]*RedisObject{},
	}
}

func (db *Database) SetWithOptions(key, value, options string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	ro := RedisObject{
		Data: value,
		Typ:  reflect.TypeOf(value).String(),
		TTL:  -1,
	}
	db.data[key] = &ro

	return nil
}

func (db *Database) Set(key, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	ro := RedisObject{
		Data: value,
		Typ:  reflect.TypeOf(value).String(),
		TTL:  -1,
	}
	db.data[key] = &ro

	return nil
}

func (db *Database) Get(key string) (any, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if val, ok := db.data[key]; ok {
		return val.Data, nil
	}

	return nil, fmt.Errorf("key not in database")
}

func (db *Database) RPush(key, value string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.data[key]; !ok {
		ro := RedisObject{
			Data: []string{},
			Typ:  reflect.TypeOf([]string{}).String(),
			TTL:  -1,
		}
		db.data[key] = &ro
	}
	if db.data[key].Typ != "[]string" {
		return "", fmt.Errorf("incorrect type for key %s", key)
	}
	db.data[key].Data = append(db.data[key].Data.([]string), value)

	return fmt.Sprintf("(integer) + %d", len(db.data[key].Data.([]string))), nil
}

func (db *Database) LPush(key, value string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.data[key]; !ok {
		ro := RedisObject{
			Data: []string{},
			Typ:  reflect.TypeOf([]string{}).String(),
			TTL:  -1,
		}
		db.data[key] = &ro
	}
	val, _ := db.data[key]
	if val.Typ != "[]string" {
		return "", fmt.Errorf("incorrect type for key %s", key)
	}
	d, _ := val.Data.([]string)
	db.data[key].Data = append(db.data[key].Data.([]string), value)
	db.data[key].Data = append([]string{value}, d...)

	return fmt.Sprintf("(integer) + %d", len(db.data[key].Data.([]string))), nil
}

func (db *Database) RPop(key string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.data[key]; !ok {
		return "", fmt.Errorf("Key not found in db")
	}
	val, _ := db.data[key]
	if val.Typ != "[]string" {
		return "", fmt.Errorf("incorrect type for key %s", key)
	}
	d, _ := val.Data.([]string)
	lastIdx := len(d) - 1
	db.data[key].Data = d[:lastIdx]

	return d[lastIdx], nil
}

func (db *Database) LPop(key string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.data[key]; !ok {
		return "", fmt.Errorf("Key not found in db")
	}
	val, _ := db.data[key]
	if val.Typ != "[]string" {
		return "", fmt.Errorf("incorrect type for key %s", key)
	}
	d, _ := val.Data.([]string)
	db.data[key].Data = d[1:]

	return d[0], nil
}

func (db *Database) LRange(key string, start, stop int) (string, error) {
	stringBuilder := ""
	val, ok := db.data[key]
	if !ok {
		return "", fmt.Errorf("Key does not exist in database")
	}
	d, ok := val.Data.([]string)
	if !ok {
		return "", fmt.Errorf("Data is not of type array")
	}

	if stop > len(d) {
		stop = len(d)
	}

	for i, v := range d[start:stop] {
		stringBuilder += fmt.Sprintf("%d) %s\n", i+1, v)
	}

	return stringBuilder, nil
}

func (db *Database) Save() {
	var f *os.File
	var err error

	f, err = os.Create("data.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	for k, v := range db.data {
		fmt.Printf("Key: %#v, Val: %#v\n", k, v)
		strVal, err := resp.Serialize(v.Data)
		if err != nil {
			panic(err)
		}
		strCat := fmt.Sprintf("%#v, %#v\n", k, strVal)
		line := []byte(strCat)
		n, err := f.Write(line)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Wrote %d bytes\n", n)

	}
}
