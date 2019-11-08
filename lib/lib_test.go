package lib

import (
	"fmt"
	"log"
	"testing"
)

func TestGOVersionCompare(t *testing.T) {
	a := "1.13.4"
	b := "1.13rcl"
	c, err := GOVersionCompare(a, b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
}

func TestDownloadFromUrl(t *testing.T) {
	url := "https://studygolang.com/dl/golang/go1.13.4.linux-amd64.tar.gz"
	path := "/home/apple/go/src/goupgrade/lib/a.tar.gz"
	DownloadFromUrl(url, path)
}

func TestSnipper(t *testing.T) {

}
