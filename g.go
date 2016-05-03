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
				url2 := parseCommentUrl(s)
				fmt.Println(url2)
				//sl, _ := getComments(url)
				//ch2 <- "save: " + strconv.Itoa(len(sl))
				ch2 <- "save: " + strconv.Itoa(len(url2))
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

func parseCommentUrl(s string) string {
	//fmt.Println(s)
	//pattern := `<div(.*?)>`
	pattern := `(?s)<div\sclass="scontent"\sid="bofqi">(.*?)</div>`
	reg := regexp.MustCompile(pattern)
	fmt.Printf("%q\n", reg.FindAllString(s, -1))
	return s[:20]
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
