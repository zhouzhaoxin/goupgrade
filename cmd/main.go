package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/k0kubun/pp"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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

func init() {
	ARCH = make(map[string]string)
	ARCH["amd64"] = "x86-64"
}

func main() {
	url := "https://studygolang.com/dl"
	versions := buildGOVersion(url)
	latestVersion := versions[0]
	_, _ = pp.Println(latestVersion)
}

// 获取golang最新版本
func buildGOVersion(url string) (versions []GOVersion) {
	resp, err := http.Get(url)
	handleErr(err)
	document, err := goquery.NewDocumentFromReader(resp.Body)
	handleErr(err)
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
