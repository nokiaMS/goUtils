package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
	"flag"
	"strings"
	"strconv"
	"math/rand"
	"time"
)

//向数据库中写入一条数据.
func ldbWrite(db *leveldb.DB, key string, info string) error {
	if db == nil {
		fmt.Println("Error: Db is nil.")
	}
	return db.Put([]byte(key), []byte(info),nil)
}

//生成随机字符串(不是绝对随机,时间相近的话值会相同.)
func GetRandomString(l int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//向数据库中批量写入数据.
func ldbWriteBatchBySingle(db *leveldb.DB, cnt int, l int) (int, error) {
	//构造数据.
	items := make(map[string]string)
	for i:= 0; i < cnt; i++ {
		k := KEYPREFIX + strconv.Itoa(i)
		items[k] = GetRandomString(l)
		//fmt.Printf("key: %s, value: %s\n", k, items[k])
	}

	startTime := time.Now().UnixNano()
	for i:=0; i < cnt; i++ {
		k := KEYPREFIX + strconv.Itoa(i)
		err := ldbWrite(db,k,items[k])
		if err != nil {
			return i, fmt.Errorf("Error when ldbWrite. key %s\n", k)
		}
	}
	endTime := time.Now().UnixNano()
	fmt.Printf("Total items: %d, time: %d ms\n", cnt, (endTime - startTime)/1000.0/1000.0)
	return cnt, nil
}

//从数据库中批量读取数据.
func ldbReadBatch(db * leveldb.DB, from int, to int) (int, error)  {
	n := 0
	startTime := time.Now().UnixNano()
	for i:= from; i < to; i++ {
		k := KEYPREFIX + strconv.Itoa(i)
		ldbRead(db,k)
		n += 1
	}
	endTime := time.Now().UnixNano()
	fmt.Printf("Total items: %d, time: %d ms\n", n, (endTime - startTime)/1000.0/1000.0)
	return n, nil
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
	KEYPREFIX = "TEST_KEY_"
)

func main() {
	path := flag.String("db", `F:\tmp\testldb`, "level db path.")
	mode := flag.String("mode", "batchRead", "test mode: batchWrite | batchRead | singleWrite | singleRead.")
	count := flag.Int("count", 100000, "Count of the values to write to or read from the db.")
	datalen := flag.Int("datalen", 100, "the length of per data to write into the database.")
	key := flag.String("key", "gx111", "key only for singleWrite or singleRead mode.")
	value := flag.String("value", "222222222222222222abfef", "value only for singleWrite mode.")
	from := flag.Int("from", 0,"Read batch from 'from'")
	to := flag.Int("to", 100000, "Read batch from 'to'")

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
		if *count == 0 || *datalen == 0 {
			fmt.Println("count or datalen error.")
			return
		}
		_, err := ldbWriteBatchBySingle(db, *count, *datalen)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Write batch done, count %d\n", *count)
		}
	} else if strings.Compare(*mode, BATCHREAD) == 0 {
		if *from < 0 || *to < 0 {
			fmt.Println("from or to parameter is wrong.")
			return
		}
		_, err := ldbReadBatch(db, *from, *to)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Read batch done, from %d, to %d\n", *from, *to)
		}
	}
}
