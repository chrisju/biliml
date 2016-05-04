package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"golang.org/x/net/html"
)

func main() {
	sr := os.Args[1]
	sn := os.Args[2]
	r, _ := strconv.Atoi(sr)
	n, _ := strconv.Atoi(sn)

	chr := make(chan string)   //任务
	csvch := make(chan string) //结果 作为csv的一行
	quitch := make(chan int)
	go Deal(chr, csvch)
	go Save(csvch, quitch, n)
	for i := 0; i < n; i++ {
		id := r - i
		url := fmt.Sprintf("http://www.bilibili.tv/video/av%d/index.html", id)
		fmt.Println("url:", url)
		chr <- url
	}

	<-quitch
}

func Deal(ch chan string, ch2 chan string) {
	for i := 0; i < 1; i++ {
		go func() {
			for {
				url := <-ch
				fmt.Println("get task:" + url)
				s, _ := fetchBody(url)
				//fmt.Println(s)
				result := parseAv(s)
				fmt.Println(result)
				//sl, _ := getComments(url)
				//ch2 <- "save: " + strconv.Itoa(len(sl))
				ch2 <- "save: " + strconv.Itoa(len(result))
			}
		}()
	}
}

func Save(ch chan string, quitch chan int, n int) {
	for {
		line := <-ch
		fmt.Println(line)

		n--
		if n == 0 {
			quitch <- -1
		}
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
	var ss, ss1 [][]string

	// play,comm,coin,fav
	pattern = `(?s)<div\sclass="v-title-info">(.*?)<div\sclass="upinfo">`
	reg = regexp.MustCompile(pattern)
	ss = reg.FindAllStringSubmatch(s, -1)
	fmt.Printf("%q\n", ss[0][0])

	pattern = `<span\sid="dianji">(\d*)</span>`
	reg = regexp.MustCompile(pattern)
	ss1 = reg.FindAllStringSubmatch(ss[0][0], -1)
	fmt.Printf("%q\n", ss1)
	play := ss1[0][1]

	pattern = `<span\sid="dm_count">(\d*)</span>`
	reg = regexp.MustCompile(pattern)
	ss1 = reg.FindAllStringSubmatch(ss[0][0], -1)
	fmt.Printf("%q\n", ss1)
	comm := ss1[0][1]

	pattern = `<span\sid="v_ctimes">(\d*)</span>`
	reg = regexp.MustCompile(pattern)
	ss1 = reg.FindAllStringSubmatch(ss[0][0], -1)
	fmt.Printf("%q\n", ss1)
	coin := ss1[0][1]

	pattern = `<span\sid="stow_count">(\d*)</span>`
	reg = regexp.MustCompile(pattern)
	ss1 = reg.FindAllStringSubmatch(ss[0][0], -1)
	fmt.Printf("%q\n", ss1)
	fav := ss1[0][1]

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
	result["play"] = play
	result["comm"] = comm
	result["coin"] = coin
	result["fav"] = fav
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
