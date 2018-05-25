package tuili

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

var GMapTPLFile = make(map[string]*template.Template) //文件名->模板内容
var GMapTPLFuncs = make(map[string]interface{})       //函数名->函数地址
var GMapHTMFile = make(map[string]string)             //文件名->文件内容

func Init_TPLFiles() error {
	log.Println("初始化模板列表")
	err := filepath.Walk(filepath.Join(GExePath, "tuili/template"), func(path string, info os.FileInfo, err error) error {
		//模板
		if filepath.Ext(path) == ".tpl" {
			text, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("读取tpl文件错误", err)
				return err
			}
			tpl, err := template.New(info.Name()).Funcs(GMapTPLFuncs).Parse(string(text[:]))
			if err != nil {
				log.Println("解析模板文件失败", err)
				return err
			}
			GMapTPLFile[info.Name()] = tpl
		}
		//非模板
		if filepath.Ext(path) == ".htm" {
			text, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("读取hml文件错误", err)
				return err
			}
			GMapHTMFile[info.Name()] = string(text[:])
		}
		return nil
	})
	return err
}

var GMapHot = make(map[string]int32)     //存储前10的作家名
var GMapHotAuth = make(map[string]int32) //存储所有作家的下载情况
var GMuxHot = &sync.Mutex{}

//加载作者情况<启动>
func LoadHotAuth() {
	text, err := ioutil.ReadFile(filepath.Join(GExePath, "hotauth.txt"))
	if err != nil {
		log.Println("加载热门作家失败")
		return
	}
	//log.Println(string(text[:]), filepath.Join(GExePath, "hotauth.txt"))
	arrLines := strings.Split(string(text[:]), "\r\n")
	for _, line := range arrLines {
		//log.Println(line)
		if !strings.Contains(line, "||") {
			continue
		}
		arrItems := strings.Split(line, "||")
		if len(arrItems) != 2 {
			continue
		}
		dt, err := strconv.ParseInt(arrItems[1], 10, 32)
		if err != nil {
			log.Println("解析下载次数失败", err)
			continue
		}
		GMapHotAuth[arrItems[0]] = int32(dt)
	}
	//log.Println(GMapHotAuth)
}

//定期更新热门信息<定期执行>
func UpdateHotAuth() {
	//热门作家排行
	type STHotAuth struct {
		Name      string //作家名
		DownTimes int32  //下载次数
	}

	//转换成st数组
	arr := []*STHotAuth{}
	GMuxHot.Lock()
	defer GMuxHot.Unlock()

	for n, c := range GMapHotAuth {
		arr = append(arr, &STHotAuth{Name: n, DownTimes: c})
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].DownTimes > arr[j].DownTimes
	})

	//保存成map
	GMapHot = make(map[string]int32)
	for i, st := range arr {
		if i > 9 {
			break
		}
		GMapHot[st.Name] = st.DownTimes
	}
	//log.Println(GMapHot)

	//保存最新的map
	buf := bytes.NewBuffer(nil)
	for k, v := range GMapHotAuth {
		buf.WriteString(fmt.Sprintf("%v||%v\r\n", k, v))
	}
	ioutil.WriteFile(filepath.Join(GExePath, "hotauth.txt"), buf.Bytes(), 0x666)
}

//是否是热门
func IsInHot(name string) bool {
	GMuxHot.Lock()
	defer GMuxHot.Unlock()

	_, ok := GMapHot[name]
	return ok
}

//增加下载计数
func AddAuthCounts(name string) {
	GMuxHot.Lock()
	defer GMuxHot.Unlock()

	c, ok := GMapHotAuth[name]
	if ok {
		GMapHotAuth[name] = c + 1
	} else {
		GMapHotAuth[name] = 1
	}
}
