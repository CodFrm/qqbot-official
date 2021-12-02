package db

import (
	"encoding/binary"
	"encoding/json"

	"github.com/xujiajun/nutsdb"
)

var db *nutsdb.DB

const MainBucket = "main-bucket"

func InitDatabase() error {
	var err error
	opt := nutsdb.DefaultOptions
	opt.Dir = "./db/main"
	db, err = nutsdb.Open(opt)
	return err
}

func Get(key string) ([]byte, error) {
	var val []byte
	if err := db.View(func(tx *nutsdb.Tx) error {
		e, err := tx.Get(MainBucket, []byte(key))
		if err != nil {
			return err
		}
		val = e.Value
		return nil
	}); err != nil {
		return nil, err
	}
	return val, nil
}

func Put(key string, val []byte, ttl uint32) error {
	if err := db.Update(func(tx *nutsdb.Tx) error {
		err := tx.Put(MainBucket, []byte(key), val, ttl)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func Incr(key string, num int64, ttl uint32) (int64, error) {
	var val int64
	if err := db.Update(func(tx *nutsdb.Tx) error {
		e, err := tx.Get(MainBucket, []byte(key))
		if err != nil {
			val = 0
		} else {
			val, _ = binary.Varint(e.Value)
		}
		b := make([]byte, binary.MaxVarintLen64)
		val = val + num
		binary.PutVarint(b, val)
		return tx.Put(MainBucket, []byte(key), b, ttl)
	}); err != nil {
		return 0, err
	}
	return val, nil
}

func GetOrSet(key string, get interface{}, set func() (interface{}, error), ttl uint32) error {
	b, err := Get(key)
	if err != nil {
		val, err := set()
		if err != nil {
			return err
		}
		b, err := json.Marshal(val)
		if err != nil {
			return err
		}
		if err := Put(key, b, ttl); err != nil {
			return err
		}
	} else {
		return json.Unmarshal(b, get)
	}
	return nil
}
