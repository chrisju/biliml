package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/sjwhitworth/golearn/base"
	. "github.com/sjwhitworth/golearn/linear_models"
)

var mu sync.Mutex
var count int
var lr *LinearRegression

func main() {
	trainData1, err := base.ParseCSVToInstances("datatrain1.csv", true)
	errexit(err)
	lr = NewLinearRegression()
	err1 := lr.Fit(trainData1)
	errexit(err1)

	fmt.Println("start listening at 0.0.0.0:2233...")
	http.HandleFunc("/", handler2)
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe("0.0.0.0:2233", nil))
}

// handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// counter echoes the number of calls so far.
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

// handler echoes the HTTP request.
func handler2(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	//fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	//for k, v := range r.Header {
	//	fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	//}
	//fmt.Fprintf(w, "Host = %q\n", r.Host)
	//fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	title := "play,comments,danmu,favorites,coins\n"
	dealed := false
	aid := ""
	play := ""
	comm := ""
	danmu := ""
	fav := ""
	for k, v := range r.Form {
		switch k {
		case "aid":
			dealed = true
			aid = v[0]
			//fmt.Fprintf(w, "Form[%s] = %s\n", k, v)
		case "play":
			play = v[0]
			//fmt.Fprintf(w, "Form[%s] = %s\n", k, v)
		case "comm":
			comm = v[0]
			//fmt.Fprintf(w, "Form[%s] = %s\n", k, v)
		case "danmu":
			danmu = v[0]
			//fmt.Fprintf(w, "Form[%s] = %s\n", k, v)
		case "fav":
			fav = v[0]
			//fmt.Fprintf(w, "Form[%s] = %s\n", k, v)
		}
	}
	if aid != "" || (play != "" && comm != "" && danmu != "" && fav != "") {
		dealed = true
	}
	if !dealed {
		fmt.Fprintf(w, "wrong request!")
	} else {
		s := fmt.Sprintf("%s%s,%s,%s,%s,%s\n", title, play, comm, danmu, fav, "1")
		fname := "t/" + strconv.Itoa(count) + ".csv"
		Save(fname, s)

		fmt.Fprintf(w, s)

		mu.Lock()
		testData, err := base.ParseCSVToInstances(fname, true)
		errexit(err)
		predictions, err2 := lr.Predict(testData)
		errexit(err2)
		expectedValue, _ := strconv.ParseFloat(base.GetClass(predictions, 0), 64)
		mu.Unlock()
		fmt.Fprintf(w, fmt.Sprintf("expected value: %v", expectedValue))
	}
}

func Save(fname string, s string) {
	f, err := os.Create(fname)
	check(err)
	defer f.Close()
	_, err = f.Write([]byte(s))
	check(err)
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func errexit(err error) {
	if err != nil {
		panic(err)
	}
}
