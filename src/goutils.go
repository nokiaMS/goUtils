package main

import (
	"diskChecker"
	"sync"
)

var wgs []*sync.WaitGroup

func registerSubUtil(wg *sync.WaitGroup)  {
	wgs = append(wgs, wg)
}

func main() {
	checker := diskChecker.MakeDiskChecker("F:\\wanchain_with_istanbul",2)
	wg := checker.Start()
	registerSubUtil(wg)
	for _,wg := range wgs {
		wg.Wait()
	}
}
