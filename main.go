package main

import (
	E "errors"
	"flag"
	"fmt"
	"math"
	"net"
	"net/http"
	conv "strconv"
	"sync"
	"time"
)

var (
	addressFlag    *string
	requestNumFlag *int
	timeOutFlag    *int
	info           RequestsInfo
)

type Arguments map[string]string

type RequestsInfo struct {
	RequestsNum            int
	TotalTime              float64
	AverageTime            float64
	MaxTime                float64
	MinTime                float64
	NotAnsweredRequestsNum int
}

func (info *RequestsInfo) String() string {
	return fmt.Sprintf("\nRequests number: %d\nMin time: %f\nMax time: %f\nAverage time: %f\nTotal time: %f\nTime out requests: %d", info.RequestsNum,
		info.MinTime, info.MaxTime, info.AverageTime, info.TotalTime, info.NotAnsweredRequestsNum)
}

func (info *RequestsInfo) update(time float64) {
	info.TotalTime += time
	info.MaxTime = math.Max(info.MaxTime, time)
	info.MinTime = math.Min(info.MinTime, time)
	info.RequestsNum++
	info.calculateAvg()
}

func (info *RequestsInfo) calculateAvg() error {
	if info.RequestsNum > 0 {
		info.AverageTime = (float64)(info.TotalTime) / (float64)(info.RequestsNum)
		return nil
	}
	return E.New("Division by zero!!!")
}

func (info *RequestsInfo) init() {
	info.RequestsNum = 0
	info.TotalTime = 0
	info.AverageTime = 0
	info.MaxTime = -1
	info.MinTime = 1000
	info.NotAnsweredRequestsNum = 0
}

func init() {
	addressFlag = flag.String("address", "https://google.com", "Address of server.")
	requestNumFlag = flag.Int("num", 10, "Amount of requests.")
	timeOutFlag = flag.Int("timeOut", 1, "Time out for waiting response from server.")
	info.init()
}

func parseArgs() (args Arguments) {
	args = make(Arguments)
	args["address"] = *addressFlag
	args["num"] = conv.Itoa(*requestNumFlag)
	args["timeOut"] = conv.Itoa(*timeOutFlag)
	return
}

func sendRrequest(client http.Client, address string, mutex *sync.Mutex, wg *sync.WaitGroup) error {
	start := time.Now()
	_, err := client.Get(address)
	rTime := time.Since(start).Seconds()
	mutex.Lock()
	if err != nil {
		if isTimeoutError(err) {
			//		mutex.Lock()
			info.NotAnsweredRequestsNum++
			//		mutex.Unlock()
		}
		wg.Done()
		return err
	}
	info.update(rTime)
	wg.Done()
	mutex.Unlock()
	return err
}

func Perform(args Arguments) (err error) {
	address := args["address"]
	requestNum, err := conv.Atoi(args["num"])
	if err != nil {
		return err
	}
	timeOut, err := conv.Atoi(args["timeOut"])
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: time.Duration(timeOut * (int)(time.Second)),
	}
	var mutex sync.Mutex
	var wg sync.WaitGroup
	wg.Add(requestNum)
	for i := 0; i < requestNum; i++ {
		go sendRrequest(client, address, &mutex, &wg)
	}
	wg.Wait()
	fmt.Println(info.String())
	return err
}

func isTimeoutError(err error) bool {
	er, ok := err.(net.Error)
	return ok && er.Timeout()
}

func Success(time float64) {
	//If I use printf --- the program is not responding for ages!!! Why???
	fmt.Print("Operation was succesfylly complited in ")
	fmt.Print(time)
	fmt.Println(" seconds!!!")
}

func main() {
	start := time.Now()
	flag.Parse()
	err := Perform(parseArgs())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	rTime := time.Since(start).Seconds()
	Success(rTime)
}
