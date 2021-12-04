package db

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/CodFrm/qqbot/utils"
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
			if err == nutsdb.ErrKeyNotFound {
				return nil
			}
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

func Del(key string) error {
	if err := db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(MainBucket, []byte(key))
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
			val = utils.StringToInt64(string(e.Value))
		}
		val = val + num
		return tx.Put(MainBucket, []byte(key), []byte(strconv.FormatInt(val, 10)), ttl)
	}); err != nil {
		return 0, err
	}
	return val, nil
}

func GetOrSet(key string, get interface{}, set func() (interface{}, error), ttl uint32) error {
	b, err := Get(key)
	if err != nil || b == nil {
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
		copyInterface(get, val)
	} else {
		return json.Unmarshal(b, get)
	}
	return nil
}

func copyInterface(dst interface{}, src interface{}) {
	dstof := reflect.ValueOf(dst)
	if dstof.Kind() == reflect.Ptr {
		el := dstof.Elem()
		srcof := reflect.ValueOf(src)
		if srcof.Kind() == reflect.Ptr {
			el.Set(srcof.Elem())
		} else {
			el.Set(srcof)
		}
	}
}
