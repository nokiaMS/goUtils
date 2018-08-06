package diskChecker

type DiskChecker struct {
	path [] string	//path to check.
	subResultCh chan uint64 //channel to send the sub result to.
	quitCh chan struct{} //quit information.
}

//make a new diskChecker object and return it.
func MakeDiskChecker() *DiskChecker {
	diskChecker := &DiskChecker{
		subResultCh: make(chan uint64),
	}
	return diskChecker
}

//set path
func (dc *DiskChecker)SetPath(paths []string)  {
	dc.path = paths
}

func (dc *DiskChecker)sumResults()  {
	var totalSize uint64 = 0

	select {
	case result := <-dc.subResultCh:
		totalSize = totalSize + result
	case <-dc.quitCh:
		return
	}
}

//start
func (dc *DiskChecker) start()  {
	
}