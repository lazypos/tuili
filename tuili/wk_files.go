package tuili

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type STFileInfo struct {
	Zone     string //区域
	Auth     string //作者
	Name     string //书名
	FullName string //全名
	Size     string //大小
}

const (
	SUPPORT_SEQ     = "|'|tlw|'|"
	SUPPORT_SEQ_SUB = `|"|tlw|"|`
)

//地区-》作家-》信息
var GMapTXTs = make(map[string]map[string]*STFileInfo)
var GTotalFiles int = 0
var GMapVeriy = make(map[string]*STFileInfo)
var GArrSupport = []string{}
var GLiuyan = ""

func Init_Files() error {
	log.Println("初始化小说列表")
	loadTXTInfo("/tuili/books/oumei", "om")
	loadTXTInfo("/tuili/books/riben", "rb")
	loadTXTInfo("/tuili/books/zhguo", "zg")
	loadTXTInfo("/tuili/books/chuang", "yc")
	loadVeriyFile()
	loadSupport("support.txt")
	return nil
}

//加载留言
func loadSupport(fpath string) error {
	text, err := ioutil.ReadFile(filepath.Join(GExePath, fpath))
	if err != nil {
		log.Println("加载留言失败", err)
		return err
	}

	GLiuyan = string(text[:])
	GArrSupport = strings.Split(GLiuyan, SUPPORT_SEQ)
	return nil
}

//加载过审列表
func loadTXTInfo(fpath, zone string) error {
	m := make(map[string]*STFileInfo)
	log.Println(filepath.Join(GExePath, fpath))
	err := filepath.Walk(filepath.Join(GExePath, fpath), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || info.Size() < 1024 || info.Size() > 10*1024*1024 {
			return nil
		}
		//log.Println(filepath.Base(info.Name()))
		arrInfo := strings.Split(filepath.Base(info.Name()), "-")
		if len(arrInfo) != 2 {
			return nil
		}
		size := fmt.Sprintf("%vKB", info.Size()/1024)
		st := &STFileInfo{Zone: zone, Auth: arrInfo[0], Name: arrInfo[1], FullName: info.Name(), Size: size}
		m[info.Name()] = st
		GTotalFiles++
		return nil
	})
	log.Println("加载小说完毕", zone, err)
	GMapTXTs[zone] = m
	return err
}

//加载审核列表
func loadVeriyFile() {
	log.Println("加载审核小说列表")
	filepath.Walk(filepath.Join(GExePath, "tuili/upload"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || info.Size() < 1024 || info.Size() > 10*1024*1024 {
			return nil
		}
		size := fmt.Sprintf("%vKB", info.Size()/1024)
		st := &STFileInfo{Zone: "", Auth: "", Name: "", FullName: info.Name(), Size: size}
		GMapVeriy[info.Name()] = st
		return nil
	})
}

func GetAuthList(zone string) []string {
	arrAuth := []string{}
	m, ok := GMapTXTs[zone]
	if ok && m != nil {
		mAuth := make(map[string]bool)
		for _, st := range m {
			if _, is := mAuth[st.Auth]; !is {
				mAuth[st.Auth] = true
				arrAuth = append(arrAuth, st.Auth)
			}
		}
	}
	sort.Strings(arrAuth)
	return arrAuth
}

//根据作家名字得到作品列表
func GetTXTList(name string) []*STFileInfo {
	arr := []*STFileInfo{}
	for _, m := range GMapTXTs {
		for _, st := range m {
			if st.Auth == name {
				arr = append(arr, st)
			}
		}
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].FullName < arr[j].FullName
	})
	return arr
}
