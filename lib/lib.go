package lib

import (
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var unstableErr error

func init() {
	unstableErr = errors.New("包含非稳定版本")
}

// 判断版本a是否比版本b更新(必须为go稳定版本)
func GOVersionCompare(a string, b string) (bool, error) {
	splitA := strings.Split(a, ".")
	splitB := strings.Split(b, ".")
	lenSplitA := len(splitA)
	lenSplitB := len(splitB)
	maxLength := int(math.Max(float64(lenSplitA), float64(lenSplitB)))
	if lenSplitA < maxLength {
		for i := 0; i < (maxLength - lenSplitA); i++ {
			splitA = append(splitA, "0")
		}
	}
	if lenSplitB < maxLength {
		for i := 0; i < (maxLength - lenSplitB); i++ {
			splitB = append(splitB, "0")
		}
	}
	for i := 0; i < int(maxLength); i++ {
		versionA, e := strconv.Atoi(splitA[i])
		if e != nil {
			return false, unstableErr
		}
		versionB, e := strconv.Atoi(splitB[i])
		if e != nil {
			return false, unstableErr
		}
		if versionA == versionB {
			continue
		} else if versionA > versionB {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, nil
}

//noinspection ALL
func DownloadFromUrl(url string, filepath string) {
	out, err := os.Create(filepath)
	HandleErr(err)
	defer out.Close()
	resp, err := http.Get(url)
	HandleErr(err)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	HandleErr(err)
}

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
