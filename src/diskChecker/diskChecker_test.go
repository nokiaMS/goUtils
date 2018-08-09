package diskChecker

import (
	"testing"
	"strings"
	"github.com/wanchain/go-wanchain/log"
)

const absolutePath = "F:\\wanchain_with_istanbul"
const relativePath = ".\\wanchain_with_istanbul"
const targetRelativePath = "C:\\Users\\asa\\AppData\\Local\\Temp\\.\\wanchain_with_istanbul"

func Test_check_absolutePath(t *testing.T)  {
	checker := MakeDiskChecker(absolutePath, 10)
	checker.check()
	if strings.Compare(checker.fullPath, absolutePath) == 0 {
		t.Log("PASS.")
	} else {
		log.Error("FAILED.")
	}
}

func Test_check_relativePath(t *testing.T)  {
	checker := MakeDiskChecker(relativePath, 10)
	checker.check()
	if strings.Compare(checker.fullPath, targetRelativePath) == 0 {
		t.Log("PASS.")
	} else {
		t.Error("FAILED.")
	}
}