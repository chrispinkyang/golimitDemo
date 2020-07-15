package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	"./golimit"
)

const (
	routineCountTotal = 5 //限制线程数
)

func main() {
	var numberTasks = [5]string{"13456755448", " 13419385751", "13419317885", " 13434343439", "13438522395"}

	g := golimit.NewG(routineCountTotal)
	wg := &sync.WaitGroup{}
	client = &http.Client{}
	beg := time.Now()
	for i := 0; i < len(numberTasks); i++ {
		wg.Add(1)
		task := numberTasks[i]
		g.Run(func() {
			respBody, err := NumberQueryRequest(task)
			if err != nil {
				fmt.Printf("error occurred in NumberQueryRequest: %s\n", task)
			} else {
				fmt.Printf("response data: %s\n", string(respBody))
			}
			wg.Done()
		})
	}
	wg.Wait()
	fmt.Printf("time consumed: %fs", time.Now().Sub(beg).Seconds())
}
var client *http.Client

func NumberQueryRequest(keyword string) (body []byte, err error) {
	url := fmt.Sprintf("https://api.binstd.com/shouji/query?appkey=df2720f76a0991fa&shouji=%s", keyword)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("response status code is not OK, response code is %d, body:%s", resp.StatusCode, string(data))
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
