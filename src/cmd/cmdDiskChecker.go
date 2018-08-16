package main

import (
	"diskChecker"
	"sync"
	"flag"
)

var wgs []*sync.WaitGroup

func registerSubUtil(wg *sync.WaitGroup)  {
	wgs = append(wgs, wg)
}

func main() {
	path := flag.String("foder","F:\\","Folder to monitor.")
	period := flag.Int("period", 5, "The period to check the folder size.")

	flag.Parse()

	//checker := diskChecker.MakeDiskChecker("F:\\wanchain_with_istanbul",2)
	checker := diskChecker.MakeDiskChecker(*path,*period)
	wg := checker.Start()
	registerSubUtil(wg)
	for _,wg := range wgs {
		wg.Wait()
	}
}
