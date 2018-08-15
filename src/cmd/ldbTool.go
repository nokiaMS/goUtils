package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
	"flag"
	"strings"
)

//向数据库中写入一条数据.
func ldbWrite(db *leveldb.DB, key string, info string) error {
	if db == nil {
		fmt.Println("Error: Db is nil.")
	}
	return db.Put([]byte(key), []byte(info),nil)
}

//向数据库中批量写入数据.
func ldbWriteBatch(db *leveldb.DB, cnt int) (int, error) {
	return 0, nil
}

func ldbReadBatch(db * leveldb.DB, from int, to int) (int, error)  {
	return 0, nil
}

//从数据库中读取一条数据.
func ldbRead(db *leveldb.DB, key string)([]byte, error)  {
	if db == nil {
		fmt.Println("Error: Db is nil.")
	}
	return db.Get([]byte(key), nil)
}

const (
	BATCHWRITE = "batchWrite"
	BATCHREAD = "batchRead"
	SINGLEWRITE = "singleWrite"
	SINGLEREAD = "singleRead"
)

func main() {
	path := flag.String("db", `F:\tmp\testldb`, "level db path.")
	mode := flag.String("mode", "singleRead", "test mode: batchWrite | batchRead | singleWrite | singleRead.")
	//count := flag.Int("count", 0, "Count of the values to write to or read from the db.")
	//datalen := flag.Int("datalen", 0, "the length of per data to write into the database.")
	key := flag.String("key", "gx111", "key only for singleWrite or singleRead mode.")
	value := flag.String("value", "222222222222222222abfef", "value only for singleWrite mode.")

	//创建或者打开ldb.
	db, err := leveldb.OpenFile(*path,nil)	// ``代表是原始字符串,不需要转义.
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	//single write.
	if strings.Compare(*mode,SINGLEWRITE) == 0 {
		//check parameter validation.
		if strings.Compare(*key, "") == 0 {
			fmt.Println("Key should not be null.")
			return
		}
		err := ldbWrite(db, *key, *value)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Write data done [ key: %s, value: %s ]\n", *key, *value)
		}
	} else if strings.Compare(*mode, SINGLEREAD) == 0 {
		//check parameter validation.
		if strings.Compare(*key, "") == 0 {
			fmt.Println("Key should not be null.")
			return
		}
		ret, err := ldbRead(db, *key)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Read data done [ key: %s, value: %s ]\n", *key, string(ret))
		}
	} else if strings.Compare(*mode, BATCHWRITE) == 0 {

	} else if strings.Compare(*mode, BATCHREAD) == 0 {

	}
}
