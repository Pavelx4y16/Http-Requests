package main

import (
	E "errors"
	"flag"
	"fmt"
	"math"
	"net"
	"net/http"
	conv "strconv"
	"time"
)

var (
	addressFlag    *string
	requestNumFlag *int
	timeOutFlag    *int
	info           InfoServer
)

type Arguments map[string]string

type InfoServer struct {
	RequestsNum            int
	TotalTime              float64
	AverageTime            float64
	MaxTime                float64
	MinTime                float64
	NotAnsweredRequestsNum int
}

func (info *InfoServer) String() string {
	return fmt.Sprintf("\nRequests number: %d\nMin time: %f\nMax time: %f\nAverage time: %f\nTotal time: %f\nTime out requests: %d", info.RequestsNum,
		info.MinTime, info.MaxTime, info.AverageTime, info.TotalTime, info.NotAnsweredRequestsNum)
}

func (info *InfoServer) update(time float64) {
	info.TotalTime += time
	info.MaxTime = math.Max(info.MaxTime, time)
	info.MinTime = math.Min(info.MinTime, time)
	info.RequestsNum++
	info.calculateAvg()
}

func (info *InfoServer) calculateAvg() error {
	if info.RequestsNum > 0 {
		info.AverageTime = (float64)(info.TotalTime) / (float64)(info.RequestsNum)
		return nil
	}
	return E.New("Division by zero!!!")
}

func (info *InfoServer) init() {
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
	timeOutFlag = flag.Int("timeOut", 3, "Time out for waiting response from server.")
	info.init()
}

func FlagInfo() {
	fmt.Printf("Server: %s\nRequestNum: %d\nTimeOut: %d\n", *addressFlag, *requestNumFlag, *timeOutFlag)
}

func parseArgs() (args Arguments) {
	args = make(Arguments)
	args["address"] = *addressFlag
	args["num"] = conv.Itoa(*requestNumFlag)
	args["timeOut"] = conv.Itoa(*timeOutFlag)
	return
}

func sendRrequest(client http.Client, address string) error {
	start := time.Now()
	resp, err := client.Get(address)
	if err != nil {
		if checkTimeoutError(err) {
			info.NotAnsweredRequestsNum++
		}
		return err
	}
	defer resp.Body.Close()
	rTime := time.Since(start).Seconds()
	info.update(rTime)
	//fmt.Println("\n\n\n-------------------------------------\n\n\n")
	//_, err = io.Copy(os.Stdout, resp.Body)
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
	for i := 0; i < requestNum; i++ {
		sendRrequest(client, address)
	}
	fmt.Println(info.String())
	return err
}

func checkTimeoutError(err error) bool {
	er, ok := err.(net.Error)
	return ok && er.Timeout()
}

func Success() {
	fmt.Println("Operation was succesfylly complited!!!")
}

func main() {
	flag.Parse()
	FlagInfo()
	err := Perform(parseArgs())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	Success()
}
