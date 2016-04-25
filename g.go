package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	sr := os.Args[1]
	sn := os.Args[2]
	r, _ := strconv.Atoi(sr)
	n, _ := strconv.Atoi(sn)
	for i := 0; i < n; i++ {
		id := r - i
		url := fmt.Sprintf("http://www.bilibili.com/video/%d/index.html", id)
		fmt.Println("url:", url)
	}
}
