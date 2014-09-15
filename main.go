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
	wait   sync.WaitGroup
	c      = flag.Int("c", 100, "concurrency")
	n      = flag.Int("n", 10, "number")
	url    = flag.String("h", "", "url")
	file   = flag.String("f", "", "params file")
	params []string
)

func main() {
	runtime.GOMAXPROCS(4)

	flag.Parse()

	fi, err := os.Open(*file)
	if err != nil {
		fmt.Println(*file)
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

	start := time.Now()
	for i := 0; i < *c; i++ {
		wait.Add(1)

		go func() {
			sendHttp(*url, *n)
			wait.Done()
		}()
	}
	wait.Wait()

	cost := time.Now().Sub(start).Nanoseconds() / int64(time.Millisecond)
	num := *c * *n
	average := float64(cost) / float64(num)
	fmt.Println("num : ", num, "cost : ", cost, "averange : ", average)
}

func sendHttp(url string, n int) {
	var data map[string]interface{}
	var d string
	if strings.Contains(url, "?") {
		d = "&"
	} else {
		d = "?"
	}
	client := &http.Client{}

	for i := 0; i < n; i++ {
		param := params[rand.Intn(len(params))]
		req, err := http.NewRequest("GET", url+d+param, nil)
		req.Header.Set("User-Agent", "go benchmark v 0.1")
		res, err := client.Do(req)
		//res, err := http.Get(url + d + param)
		if err != nil {
			fmt.Println("http get error, ", err, res)
			continue
		}
		result, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println("http get error, ", err)
			continue
		}
		json.Unmarshal(result, &data)
	}
}
