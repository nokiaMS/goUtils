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

const batchCnt = 10000

//向buffer中填充数据.
func fillData(from int, cnt int, l int) *map[string]string{
	items := make(map[string]string)
	for j := 0; j < cnt; j++ {
		k := KEYPREFIX + strconv.Itoa(j+cnt * from)
		items[k] = GetRandomString(l)
		//fmt.Printf("key: %s, value: %s\n", k, items[k])
	}
	return &items
}

func showTimeResult(start int64, end int64, total int, curRound int)  {
	timems := float32(end-start)/float32(1000)/float32(1000)
	fmt.Printf("Total items: %d, time: %v ms, rate %v items per sec.\n", total, timems, float32(curRound)/timems * 1000)
}

func writeToDB(items *map[string]string, db *leveldb.DB, round int) (int, error)  {
	for k, v := range *items {
		err := ldbWrite(db, k, v)
		if err != nil {
			index,_ := strconv.Atoi(k)
			return 0, fmt.Errorf("Error when ldbWrite %d. key %s\n", batchCnt * round + index, k)
		}
	}
	return 0, nil
}

//向数据库中批量写入数据.
func ldbWriteBatchBySingle(db *leveldb.DB, cnt int, l int) (int, error) {
	//构造数据.
	partCnt := cnt/batchCnt
	remainItemCnt := cnt - batchCnt * partCnt

	for i := 0; i < partCnt; i++ {		//strconv.itoa很耗时.
		items := fillData(batchCnt * (i-1), batchCnt, l)
		startTime := time.Now().UnixNano()
		writeToDB(items, db, i)
		endTime := time.Now().UnixNano()
		showTimeResult(startTime, endTime, batchCnt*(i+1), batchCnt)
	}

	//处理最后不足batchCnt的部分.
	items := fillData(batchCnt * partCnt, remainItemCnt, l)
	writeToDB(items, db, partCnt)
	return cnt, nil
}

//从数据库中批量读取数据.
func ldbReadBatchBySingle(db * leveldb.DB, from int, to int, size int, checkLen bool) (int, error)  {
	n := 0
	keys := []string{}	//变量不能作为数组长度,数组长度必须是在编译的时候能够确定的,这点和C语言是一样的.
	//构造keys.
	for i := from; i < to; i++ {
		k := KEYPREFIX + strconv.Itoa(i)
		keys = append(keys, k)
	}

	startTime := time.Now().UnixNano()
	for _, v := range keys {
		ret,_ := ldbRead(db, v)
		if checkLen {
			if len(ret) != size {
				fmt.Println("Size of data does not match with the data length in db.")
			}
		}
		n += 1
	}
	endTime := time.Now().UnixNano()
	showTimeResult(startTime, endTime, n, n)
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

var modeDesc = "test mode: " + BATCHWRITE + " " + BATCHREAD + " " + SINGLEWRITE + " " + SINGLEREAD + "."

func main() {
	path := flag.String("db", `F:\tmp\testldb`, "level db path.")
	mode := flag.String("mode", BATCHREAD, modeDesc)
	count := flag.Int("count", 1000000, "Count of the values to write to db.")
	datalen := flag.Int("datalen", 10, "the length of per data to write into the database.")
	key := flag.String("key", "gx111", "key only for singleWrite or singleRead mode.")
	value := flag.String("value", "abcdefghigklmnopqrstuvwxyz", "value only for singleWrite mode.")
	from := flag.Int("from", 0,"Read batch from 'from'")
	to := flag.Int("to", 10000, "Read batch from 'to'")

	flag.Parse()	//如果不添加这句话程序运行不会出错,但是 program --help的话不会打印出帮助信息.

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
		_, err := ldbReadBatchBySingle(db, *from, *to, *datalen, true)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Read batch done, from %d, to %d\n", *from, *to)
		}
	}
}
