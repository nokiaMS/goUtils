package main

import (
	"diskChecker"
	"time"
)

func main() {
	checker := diskChecker.MakeDiskChecker("F:\\wanchain_with_istanbul",2)
	checker.Start()
	checker.Check()
	time.Sleep(time.Duration(3) * time.Minute)
	checker.Stop()
}
