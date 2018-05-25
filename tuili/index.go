package tuili

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var GExePath = GetExePath()

func GetExePath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

var (
	strOM    = ""
	strRB    = ""
	strZG    = ""
	strYC    = ""
	mapTList = make(map[string][]*STFileInfo) //作者文章列表
	GST      = struct {
		OM         string
		RB         string
		ZG         string
		YC         string
		TOTAL      int
		DownCounts int32
	}{}
	GFangWen    int32 = 0
	GTotalDowns int32 = 0
)

func HandleFrame(c *gin.Context) {
	filename := c.Param("name")
	if filename == "home" {
		HandleHome(c)
	} else if filename == "download" {
		HandleDown(c)
	} else if filename == "list" {
		HandleTXTList(c)
	} else if filename == "veriy" {
		HandleVeriy(c)
	} else if filename == "support" {
		HandleSupport(c)
	} else if filename == "upload" {
		HandleUpload(c)
	} else if filename == "uploadtxt" {
		HandleUploadTXT(c)
	} else if filename == "supportup" {
		HandleSupportUP(c)
	}
}

//格式化文本
func FormatInputText(input string) string {
	input = template.HTMLEscapeString(input)
	input = strings.Replace(input, "-", "_", -1)
	input = strings.Replace(input, " ", "&nbsp;", -1)
	input = strings.Replace(input, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)
	input = strings.Replace(input, "\r\n", "</br>", -1)
	return input
}

func HandleSupportUP(c *gin.Context) {
	txt := c.Request.FormValue("txt")
	if len(txt) > 0 {
		txt = FormatInputText(txt)
		t := time.Now().String()[:19]
		liuyan := fmt.Sprintf(`%v%v%v`, txt, SUPPORT_SEQ_SUB, t)
		GArrSupport = append([]string{liuyan}, GArrSupport...)
		GLiuyan = liuyan + SUPPORT_SEQ + GLiuyan
		ioutil.WriteFile(filepath.Join(GExePath, "support.txt"), []byte(GLiuyan), 0x666)
	}
	c.String(http.StatusOK, "")
}

func HandleUploadTXT(c *gin.Context) {

	fName := c.Request.FormValue("name")
	if c.Request.Method == "POST" && len(fName) > 0 {
		fileName := fmt.Sprintf(`%v_%s`, strings.Replace(time.Now().String()[:19], ":", "-", -1), fName)
		fName = filepath.Join(GetExePath(), "/tuili/upload/"+fileName)

		text, err := ioutil.ReadAll(c.Request.Body)
		if len(text) < 1024 || len(text) > 5*1024*1024 || err != nil {
			c.String(http.StatusOK, "文件大小不对.")
			return
		}

		if err := ioutil.WriteFile(fName, text, 0x666); err != nil {
			c.String(http.StatusOK, fmt.Sprintf("系统出错..%v", err.Error()))
			return
		}
		GMapVeriy[fileName] = &STFileInfo{FullName: fileName, Size: fmt.Sprintf("%vKB", len(text)/1024)}
		c.String(http.StatusOK, "")
		return
	}
	c.String(http.StatusOK, "参数错误.")
}

func HandleUpload(c *gin.Context) {
	buf := bytes.NewBufferString("")
	i := 0
	for _, v := range GMapVeriy {
		i++
		buf.WriteString(fmt.Sprintf(`<li>%v-<a download="" href="/upload/%v">%v</a><span>%v</span></li>`, i, v.FullName, v.FullName, v.Size))
	}
	tpl := GMapTPLFile["fm_upload.tpl"]
	if tpl != nil {
		s := buf.String()
		if err := tpl.Execute(c.Writer, &s); err != nil {
			c.String(http.StatusOK, "系统出错.")
		}
	}
}

func HandleSupport(c *gin.Context) {
	buf := bytes.NewBufferString("")
	i := 0
	for _, v := range GArrSupport {
		if len(v) == 0 {
			continue
		}
		i++
		if i > 20 {
			break
		}
		arr := strings.Split(v, SUPPORT_SEQ_SUB)
		buf.WriteString(fmt.Sprintf(`<div class="liuyan"><div>%v</div><div class="liuyanbt"><span>匿名网友</span><span>%v</span></div></div>`, arr[0], arr[1]))
	}
	tpl := GMapTPLFile["fm_support.tpl"]
	if tpl != nil {
		s := buf.String()
		if err := tpl.Execute(c.Writer, &s); err != nil {
			c.String(http.StatusOK, "系统出错.")
		}
	}
}

func HandleVeriy(c *gin.Context) {
	buf := bytes.NewBufferString("")
	i := 0
	for _, v := range GMapVeriy {
		i++
		buf.WriteString(fmt.Sprintf(`<li>%v-<a download="" href="/upload/%v">%v</a><span>%v</span></li>`, i, v.FullName, v.FullName, v.Size))
	}
	tpl := GMapTPLFile["fm_txtlist.tpl"]
	if tpl != nil {
		s := buf.String()
		if err := tpl.Execute(c.Writer, &s); err != nil {
			c.String(http.StatusOK, "系统出错.")
		}
	}
}

//作家小说列表
func HandleTXTList(c *gin.Context) {
	key := c.Request.FormValue("id")
	if len(key) < 4 {
		return
	}
	t := key[:2]
	buf := bytes.NewBufferString("")
	if arr, ok := mapTList[key]; ok {
		for i, v := range arr {
			if strings.Contains(v.FullName, ".txt") || strings.Contains(v.FullName, ".TXT") || strings.Contains(v.FullName, ".Txt") {
				buf.WriteString(fmt.Sprintf(`<li>%v-%v<span>%v</span> <a download="" href="/%v/%v">下载</a> <a href="/%v/%v?type=1">阅读</a></li>`, i+1, v.FullName, v.Size, t, v.FullName, t, v.FullName))
			} else {
				buf.WriteString(fmt.Sprintf(`<li>%v-%v<span>%v</span> <a download="" href="/%v/%v">下载</a></li>`, i+1, v.FullName, v.Size, t, v.FullName))
			}
		}
		tpl := GMapTPLFile["fm_txtlist.tpl"]
		if tpl != nil {
			s := buf.String()
			if err := tpl.Execute(c.Writer, &s); err != nil {
				c.String(http.StatusOK, "系统出错.")
			}
		}
	}
}

func HandleHome(c *gin.Context) {
	atomic.AddInt32(&GFangWen, 1)
	fw := atomic.LoadInt32(&GFangWen)

	tpl := GMapTPLFile["fm_index.tpl"]
	if tpl != nil {
		if err := tpl.Execute(c.Writer, &fw); err != nil {
			c.String(http.StatusOK, "系统出错.")
		}
	}
}

//作家列表
func HandleDown(c *gin.Context) {
	GMuxZone.Lock()
	defer GMuxZone.Unlock()

	log.Println("收到访问", c.ClientIP())

	tpl := GMapTPLFile["fm_download.tpl"]
	if tpl != nil {
		if err := tpl.Execute(c.Writer, &GST); err != nil {
			c.String(http.StatusOK, "系统出错.")
		}
	}
}

func GetUrul(src string) string {
	if strings.Index(src, "om") != -1 {
		return "/tuili/books/oumei/"
	}
	if strings.Index(src, "rb") != -1 {
		return "/tuili/books/riben/"
	}
	if strings.Index(src, "zg") != -1 {
		return "/tuili/books/zhguo/"
	}
	if strings.Index(src, "yc") != -1 {
		return "/tuili/books/chuang/"
	}
	return ""
}

//作家列表
func HandleFile(c *gin.Context) {
	tp := c.Request.FormValue("type")
	name := c.Param("name")
	//下载计数
	if !strings.Contains(name, "-") {
		return
	}
	auth := strings.Split(name, "-")[0]
	AddAuthCounts(auth)
	//log.Println(auth)

	//在线阅读
	if tp == "1" {
		//fmt.Println(c.Request.RequestURI)
		c.Writer.Header().Set("Content-Type", "text/txt;charset=gb2312")
		c.Writer.WriteHeader(http.StatusOK)
		log.Println("阅读小说", name, c.ClientIP())
	} else {
		//下载
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename=%v`, name))
		c.Writer.WriteHeader(http.StatusOK)
		log.Println("下载小说", name, c.ClientIP())
	}

	handle, err := os.Open(filepath.Join(GExePath, GetUrul(c.Request.RequestURI)+name))
	if err != nil {
		c.String(http.StatusOK, "系统出错 %v", err.Error())
		return
	}
	defer handle.Close()
	if _, err = io.Copy(c.Writer, handle); err != nil {
		c.String(http.StatusOK, "系统出错")
		return
	}
	atomic.AddInt32(&GTotalDowns, 1)
	GST.DownCounts = atomic.LoadInt32(&GTotalDowns)
}

//主框架
func HandleIndex(c *gin.Context) {
	tpl := GMapTPLFile["index.tpl"]
	if tpl != nil {
		if err := tpl.Execute(c.Writer, nil); err != nil {
			c.String(http.StatusOK, "系统出错.")
		}
	}
}

func getZoneString(zone string) string {
	buf := bytes.NewBufferString("")
	arr := GetAuthList(zone)
	for i, v := range arr {
		num := fmt.Sprintf(`%v_%v`, zone, i)
		//是否热门
		if IsInHot(v) {
			buf.WriteString(fmt.Sprintf(`<a style="color: red;" href="javascript:;" onclick="loadPage1('/frame/list?id=%v')" >%v</a>`, num, v))
		} else {
			buf.WriteString(fmt.Sprintf(`<a style="color: blue;" href="javascript:;" onclick="loadPage1('/frame/list?id=%v')" >%v</a>`, num, v))
		}
		mapTList[num] = GetTXTList(v)
	}
	return buf.String()
}

var GMuxZone sync.Mutex

func Init_Index() {
	GMuxZone.Lock()
	defer GMuxZone.Unlock()

	strOM = getZoneString("om")
	strRB = getZoneString("rb")
	strZG = getZoneString("zg")
	strYC = getZoneString("yc")
	GST = struct {
		OM         string
		RB         string
		ZG         string
		YC         string
		TOTAL      int
		DownCounts int32
	}{OM: strOM, RB: strRB, ZG: strZG, YC: strYC, TOTAL: GTotalFiles, DownCounts: GTotalDowns}
}

func SaveFangwen() {
	LoadHotAuth()

	//访问次数
	line, err := ioutil.ReadFile(filepath.Join(GExePath, "sum.txt"))
	if err != nil {
		log.Println("加载访问数失败.", err)
		return
	}
	arr := strings.Split(string(line[:]), "||")
	fw, err := strconv.ParseInt(arr[0], 10, 32)
	if err != nil {
		log.Println("转换访问数失败.", err)
		return
	}
	GFangWen = int32(fw)

	dw, err := strconv.ParseInt(arr[1], 10, 32)
	if err != nil {
		log.Println("转换下载数失败.", err)
		return
	}
	GTotalDowns = int32(dw)
	log.Println("访问下载次数", fw, dw)

	for {
		s := fmt.Sprintf("%v||%v", atomic.LoadInt32(&GFangWen), atomic.LoadInt32(&GTotalDowns))
		err = ioutil.WriteFile(filepath.Join(GExePath, "sum.txt"), []byte(s), 0x666)
		if err != nil {
			log.Println("保存访问数失败.", err)
		}
		//更新热门
		UpdateHotAuth()
		Init_Index()

		time.Sleep(time.Minute)
	}
}

//下载的次数
func HandleUpdateFile(c *gin.Context) {
	name := c.Param("name")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename=%v`, name))
	c.Writer.WriteHeader(http.StatusOK)
	handle, err := os.Open(filepath.Join(GExePath, "/tuili/upload/"+name))
	if err != nil {
		c.String(http.StatusOK, "系统出错 %v", err.Error())
		return
	}
	defer handle.Close()
	if _, err = io.Copy(c.Writer, handle); err != nil {
		c.String(http.StatusOK, "系统出错")
		return
	}
	atomic.AddInt32(&GTotalDowns, 1)
	GST.DownCounts = atomic.LoadInt32(&GTotalDowns)
}

func WebInit(eg *gin.Engine) {
	log.Println("程序启动")
	go SaveFangwen()
	Init_TPLFiles()
	Init_Files()
	Init_Index()

	eg.StaticFS("/layui", http.Dir(GExePath+"/layui"))
	eg.StaticFS("/image", http.Dir(GExePath+"/tuili/images"))

	eg.POST("/frame/:name", HandleFrame)
	eg.GET("/frame/:name", HandleFrame)
	eg.GET("/om/:name", HandleFile)
	eg.GET("/rb/:name", HandleFile)
	eg.GET("/zg/:name", HandleFile)
	eg.GET("/yc/:name", HandleFile)
	eg.GET("/upload/:name", HandleUpdateFile)
	eg.GET("/", HandleIndex)
}
