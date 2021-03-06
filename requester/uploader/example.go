package uploader

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"github.com/iikira/BaiduPCS-Go/requester/rio"
)

// DoUpload 执行上传
func DoUpload(uploadURL string, readerlen64 rio.ReaderLen64, checkFunc CheckFunc) {
	u := NewUploader(uploadURL, readerlen64)
	u.SetCheckFunc(checkFunc)

	exit := make(chan struct{})

	u.OnExecute(func() {
		statusChan := u.GetStatusChan()
		for {
			select {
			case v, ok := <-statusChan:
				if !ok {
					return
				}

				fmt.Printf("\r ↑ %s/%s %s/s in %s ............",
					pcsutil.ConvertFileSize(v.Uploaded(), 2),
					pcsutil.ConvertFileSize(v.TotalSize(), 2),
					pcsutil.ConvertFileSize(v.SpeedsPerSecond(), 2),
					v.TimeElapsed(),
				)
			}
		}
	})

	u.OnFinish(func() {
		close(exit)
	})

	u.Execute()

	<-exit
	return
}
