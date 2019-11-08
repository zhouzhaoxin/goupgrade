package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/briandowns/spinner"
	"goupgrade/lib"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type GOVersionDetail struct {
	FileName string `json:"文件名"`
	Type     string `json:"类型"`
	OS       string `json:"操作系统(OS)"`
	Arch     string `json:"架构(Arch)"`
	Size     string `json:"大小"`
	SHA256   string `json:"SHA256 Checksum"`
	Href     string `json:"下载地址"`
}

type GOVersion struct {
	Version string
	Detail  []GOVersionDetail
}

var ARCH map[string]string
var needUpgrade string

const host = "https://studygolang.com"

func init() {
	ARCH = make(map[string]string)
	ARCH["amd64"] = "x86-64"
}

func main() {
	url := host + "/dl"
	versions := buildGOVersion(url)
	latestVersion := versions[0]
	if !checkIfCanUpgrade(latestVersion) {
		fmt.Println("已取消升级")
		return
	}
	fmt.Println("开始升级")
	filepath := download(latestVersion)
	fmt.Println(filepath)
}

// 下载最新的安装包
func download(latestVersion GOVersion) (filepath string) {
	dirname := "upgrade"
	detail := latestVersion.Detail[0]
	downloadURL := host + detail.Href
	dir, err := os.Getwd()
	err = os.Mkdir(path.Join(dir, dirname), 0777)
	lib.HandleErr(err)
	filepath = path.Join(dir, dirname, detail.FileName)
	lib.HandleErr(err)
	s := spinner.New(spinner.CharSets[36], 1000*time.Millisecond) // Build our new spinner
	s.Start()
	lib.DownloadFromUrl(downloadURL, filepath)
	s.Stop()
	return
}

// 检查是否需要升级
func checkIfCanUpgrade(latestVersion GOVersion) bool {
	latestVersionNumber := latestVersion.Version[2:]
	currentVersionNumber := runtime.Version()[2:]
	canUpgrade, err := lib.GOVersionCompare(latestVersionNumber, currentVersionNumber)
	lib.HandleErr(err)
	if !canUpgrade {
		log.Fatalf("当前安装版本为%s, 目前最新版本为%s --无法升级\n", currentVersionNumber, latestVersionNumber)
	}
	fmt.Printf("当前版本为%s, 可升级到%s\n", currentVersionNumber, latestVersionNumber)
	for {
		fmt.Printf("是否进行升级yes/no")
		_, err = fmt.Scan(&needUpgrade)
		lib.HandleErr(err)
		if strings.ToLower(needUpgrade) == "yes" {
			return true
		} else if strings.ToLower(needUpgrade) == "no" {
			return false
		} else {
			continue
		}
	}
}

// 获取golang最新版本
func buildGOVersion(url string) (versions []GOVersion) {
	resp, err := http.Get(url)
	lib.HandleErr(err)
	document, err := goquery.NewDocumentFromReader(resp.Body)
	lib.HandleErr(err)
	document.Find("#stable").Each(func(i int, selection *goquery.Selection) {
		selection.NextAll().EachWithBreak(func(i int, selection *goquery.Selection) bool {
			version, exists := selection.Attr("id")
			if exists && version == "unstable" {
				return false
			}
			versionObj := GOVersion{}
			versionObj.Version = version
			selection.Find("table[class=codetable]").Each(func(i int, selection *goquery.Selection) {
				versionDetail := GOVersionDetail{}
				versionDetails := make([]GOVersionDetail, 0)
				selection.Find("tbody>tr").EachWithBreak(func(i int, selection *goquery.Selection) bool {

					selection.Find("td").EachWithBreak(func(i int, selection *goquery.Selection) bool {
						if i == 0 {
							href, exists := selection.Find("a").Attr("href")
							if !exists {
								return false
							}
							versionDetail.Href = href
							versionDetail.FileName = selection.Text()
						} else if i == 1 {
							versionDetail.Type = selection.Text()
						} else if i == 2 {
							versionDetail.OS = selection.Text()
						} else if i == 3 {
							versionDetail.Arch = selection.Text()
						} else if i == 4 {
							versionDetail.Size = selection.Text()
						} else if i == 5 {
							versionDetail.SHA256 = selection.Text()
						}
						return true
					})
					// todo Archive|Installer|Source check
					if strings.ToLower(versionDetail.OS) == runtime.GOOS &&
						versionDetail.Arch == ARCH[runtime.GOARCH] {
						versionDetails = append(versionDetails, versionDetail)
					}
					return true
				})
				versionObj.Detail = versionDetails
			})
			versions = append(versions, versionObj)
			return true
		})
	})
	return
}
