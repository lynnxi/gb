package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	//"strconv"
	"strings"
	"sync"
	"time"
)

var (
	wait      sync.WaitGroup
	c         = flag.Int("c", 400, "concurrency")
	n         = flag.Int("n", 5000, "number")
	paramFile = flag.String("pf", "", "params file")
	urlFile   = flag.String("uf", "", "url file")
	params    []string
	urls      []string
)

func parseParamFile(file string) {
	fi, err := os.Open(file)
	if err != nil {
		fmt.Println(file)
		panic(err)
	}
	defer fi.Close()
	reader := bufio.NewReader(fi)

	for {
		line, err := reader.ReadSlice('\n')
		if err == io.EOF {
			break
		}
		i := len(line) - 2
		if i < 0 {
			fmt.Println("bad line terminator:" + string(line))
		}
		line = line[:i]
		params = append(params, string(line))
	}
}

func parseUrlFile(file string) {
	fi, err := os.Open(file)
	if err != nil {
		fmt.Println(file)
		panic(err)
	}
	defer fi.Close()
	reader := bufio.NewReader(fi)

	for {
		line, err := reader.ReadSlice('\n')
		if err == io.EOF {
			break
		}
		i := len(line) - 2
		if i < 0 {
			fmt.Println("bad line terminator:" + string(line))
		}
		line = line[:i]
		urls = append(urls, string(line))
	}
}

func main() {
	runtime.GOMAXPROCS(8)

	flag.Parse()

	parseParamFile(*paramFile)
	parseUrlFile(*urlFile)

	start := time.Now()
	for i := 0; i < *c; i++ {
		wait.Add(1)

		go func() {
			startHttpTest(*n)
			wait.Done()
		}()
	}
	wait.Wait()

	cost := time.Now().Sub(start).Nanoseconds() / int64(time.Millisecond)
	num := *c * *n
	average := float64(cost) / float64(num)
	fmt.Println("num : ", num, "cost : ", cost, "averange : ", average)
}

func startHttpTest(n int) {
	for i := 0; i < n; i++ {
		param := params[rand.Intn(len(params))]
		for _, url := range urls {
			sendHttp(url, param)
		}

	}
}

func sendHttp(url string, param string) {
	var data map[string]interface{}
	var d string
	if strings.Contains(url, "?") {
		d = "&"
	} else {
		d = "?"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url+d+param, nil)
	req.Header.Set("User-Agent", "go benchmark v 0.1")
	res, err := client.Do(req)
	//res, err := http.Get(url + d + param)
	if err != nil {
		fmt.Println("http get error, ", err, res)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println("http get error, ", err)
		return
	}
	json.Unmarshal(result, &data)
}
