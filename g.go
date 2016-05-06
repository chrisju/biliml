package main

import (
	j "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	. "github.com/WhiteBlue/bilibili-service/lib"
)

func main() {
	sr := os.Args[1]
	sn := os.Args[2]
	r, _ := strconv.Atoi(sr)
	n, _ := strconv.Atoi(sn)

	chr := make(chan int)      //任务
	csvch := make(chan string) //结果 作为csv的一行
	quitch := make(chan int)
	go Deal(chr, csvch)
	go Save(csvch, quitch, n+1)
	for i := 0; i < n; i++ {
		aid := r - i
		chr <- aid
	}

	<-quitch
}

func Deal(ch chan int, ch2 chan string) {
	client := NewBiliClient()
	var b []byte
	ch2 <- "play,comments,danmu,favorites,coins"
	for i := 0; i < 1; i++ {
		go func() {
			for {
				aid := <-ch
				p := map[string]string{}
				json, err := client.GetVideoInfo2(strconv.Itoa(aid))
				if err != nil {
					fmt.Println(err)
					ch2 <- ""
					continue
				}
				fmt.Println(json)
				b, _ = j.Marshal(json.Get("coins"))
				p["coin"] = strings.Trim(string(b), "\"")
				b, _ = j.Marshal(json.Get("favorites"))
				p["fav"] = strings.Trim(string(b), "\"")
				b, _ = j.Marshal(json.Get("play"))
				p["play"] = strings.Trim(string(b), "\"")
				b, _ = j.Marshal(json.Get("review"))
				p["comm"] = strings.Trim(string(b), "\"")
				b, _ = j.Marshal(json.Get("video_review"))
				p["danmu"] = strings.Trim(string(b), "\"")
				//url := fmt.Sprintf("http://www.bilibili.tv/video/av%d/index.html", aid)
				//fmt.Println("get task:" + url)
				//s, _ := fetchBody(url)
				////fmt.Println(s)
				//result := parseAv(s)
				//fmt.Println(result)
				//sl, _ := getComments(url)
				//ch2 <- "save: " + strconv.Itoa(len(sl))
				fmt.Println(p)
				line := p["play"] + "," + p["comm"] + "," + p["danmu"] + "," + p["fav"] + "," + p["coin"]
				ch2 <- line
			}
		}()
	}
}

func Save(ch chan string, quitch chan int, n int) {
	f, err := os.Create("data.csv")
	check(err)
	defer f.Close()
	for {
		line := <-ch
		if line != "" {
			fmt.Println(line)
			_, err := f.Write([]byte(line + "\n"))
			check(err)
		}

		n--
		if n == 0 {
			quitch <- -1
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fetchBody(url string) (string, error) {
	fmt.Println("fetch: " + url)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func parseAv(s string) map[string]string {
	var pattern string
	var reg *regexp.Regexp
	var ss [][]string

	// play,comm,coin,fav
	//pattern = `(?s)<div\sclass="v-title-info">(.*?)<div\sclass="upinfo">`
	//reg = regexp.MustCompile(pattern)
	//ss = reg.FindAllStringSubmatch(s, -1)
	//fmt.Printf("%q\n", ss[0][0])

	//pattern = `<span\sid="dianji">(\d*)</span>`
	//reg = regexp.MustCompile(pattern)
	//ss1 = reg.FindAllStringSubmatch(ss[0][0], -1)
	//fmt.Printf("%q\n", ss1)
	//play := ss1[0][1]

	//pattern = `<span\sid="dm_count">(\d*)</span>`
	//reg = regexp.MustCompile(pattern)
	//ss1 = reg.FindAllStringSubmatch(ss[0][0], -1)
	//fmt.Printf("%q\n", ss1)
	//comm := ss1[0][1]

	//pattern = `<span\sid="v_ctimes">(\d*)</span>`
	//reg = regexp.MustCompile(pattern)
	//ss1 = reg.FindAllStringSubmatch(ss[0][0], -1)
	//fmt.Printf("%q\n", ss1)
	//coin := ss1[0][1]

	//pattern = `<span\sid="stow_count">(\d*)</span>`
	//reg = regexp.MustCompile(pattern)
	//ss1 = reg.FindAllStringSubmatch(ss[0][0], -1)
	//fmt.Printf("%q\n", ss1)
	//fav := ss1[0][1]

	// comment url
	pattern = `(?s)<div\sclass="scontent"\sid="bofqi">(.*?)</div>`
	reg = regexp.MustCompile(pattern)
	ss = reg.FindAllStringSubmatch(s, -1)
	fmt.Printf("%q\n", ss[0][0])
	pattern = `[^a]id=(\d+)`
	reg = regexp.MustCompile(pattern)
	ss = reg.FindAllStringSubmatch(ss[0][0], -1)
	fmt.Printf("%q\n", ss)
	id := ss[0][1]
	commurl := fmt.Sprintf("http://comment.bilibili.tv/%s.xml", id)

	var result = map[string]string{}
	result["url"] = commurl
	//result["play"] = play
	//result["comm"] = comm
	//result["coin"] = coin
	//result["fav"] = fav
	return result
}

func getComments(url string) ([]string, error) {
	fmt.Println("fetch: " + url)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var comments []string
	visitNode := func(n *html.Node) {
		fmt.Println(n)
		if n.Type == html.ElementNode && n.Data == "d" {
			fmt.Println(n.Namespace)
			comments = append(comments, n.Namespace)
		}
	}

	//fmt.Println(doc)
	forEachNode(doc, visitNode, nil)
	fmt.Println(comments)
	return comments, nil
}

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
