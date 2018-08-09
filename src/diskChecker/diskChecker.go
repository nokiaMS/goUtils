package diskChecker

import (
	"fmt"
	"path/filepath"
	"os"
	"runtime"
	"strings"
	"errors"
)

type sizeResult struct {
	fileSize int64		//file size.
	showFlag bool		//where show flag to user.
}

type DiskChecker struct {
	Path string					//path to check.
	fullPath string				//full path for checking.
	SubResultCh chan *sizeResult 	//channel to send the sub result to.
	quitCh chan struct{} 		//quit information.
	totalSize int64			//total size of the target size.
	TimerPeriod int				//Period to show the result.(second)
}

var slash string

func init() {	//init函数在main前执行.
	if strings.Compare(runtime.GOOS,"windows") == 0 {
		slash = "\\"
	} else {
		slash = "/"
	}
}

//make a new diskChecker object and return it.
func MakeDiskChecker(path string, period int) *DiskChecker {
	diskChecker := &DiskChecker{
		Path: path,
		TimerPeriod: period,
		SubResultCh: make(chan *sizeResult),
		quitCh: make(chan struct{}),
	}
	return diskChecker
}

//show result to user.
func (dc *DiskChecker)showResult()  {
	fmt.Printf("%s size: %u", dc.totalSize)
}

func (dc *DiskChecker) Start()  {
	go dc.sumResults()
}

func (dc *DiskChecker) Check()  {
	dc.check()
}
//get the result
func (dc *DiskChecker)sumResults()  {
	//timer := time.NewTimer(time.Duration(dc.TimerPeriod) * time.Second)

	for {
		select {
		case result := <-dc.SubResultCh:
			dc.totalSize += result.fileSize
			if result.showFlag == true { //show result to user.
				dc.showResult()
			}
		//case <-timer.C:
		//	fmt.Println("Timer start.")
		//	dc.check()
		case <-dc.quitCh:
			return
		}
	}
}

func (dc *DiskChecker) Stop()  {
	dc.quitCh <- struct{}{} //创建一个空struct{}实例.
}

//遍历文件夹,统计结果.
func (dc *DiskChecker)walkPath() {
	filepath.Walk(dc.fullPath, dc.fileSizeCheck)
}

//发送消息给dc,使其展示结果给用户.
func (dc *DiskChecker)resultToUser()  {
	msg := &sizeResult{fileSize:0,showFlag:true}
	dc.SubResultCh <- msg
}

func (dc *DiskChecker)fileSizeCheck(path string, info os.FileInfo, err error) error {
	fmt.Println(path)

	if err != nil {
		return nil
	}
	//判断file是否需要返回大小.
	f, err := os.Stat(path)
	if err != nil {
		return nil
	}
	//判断文件是否为一个常规文件.
	if !f.Mode().IsRegular() {
		return nil
	}
	//对于常规文件返回文件大小.
	size := f.Size()
	result := &sizeResult{fileSize:size,showFlag:false}
	dc.SubResultCh <- result
	return nil
}

//start check the folder.
func (dc *DiskChecker) check() error {
	//检查path是相对路径还是绝对路径,如果是相对路径那么获得绝对路径.
	if dc.Path[0] == '.' {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return errors.New("Can not solve the full path.")
		}
		dc.fullPath = dir + slash + dc.Path
	} else {
		dc.fullPath = dc.Path
	}

	dc.walkPath()

	return nil
}
