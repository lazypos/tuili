package main

import (
	"./tuili"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	hd, err := os.OpenFile(filepath.Join(tuili.GExePath, "tuili.log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0x666)
	if err != nil {
		log.Println("打开日志失败", err)
		return
	}
	defer hd.Close()
	//log.SetOutput(hd)

	//启动
	eg := gin.Default()
	tuili.WebInit(eg)
	log.Println(eg.Run(":80"))
}
