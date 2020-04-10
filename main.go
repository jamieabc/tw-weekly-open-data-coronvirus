package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

const (
	JSONURL    = "https://od.cdc.gov.tw/eic/Weekly_Age_County_Gender_19CoV.json"
	OSExitCode = -1
)

type covid struct {
	Year           string `json:"診斷年份"`
	Week           int    `json:"診斷週別,string"`
	County         string `json:"縣市"`
	Gender         string `json:"性別"`
	Foreign        string `json:"是否為境外移入"`
	Age            string `json:"年齡層"`
	ConfirmedCount int    `json:"確定病例數,string"`
}

func (c covid) String() string {
	return fmt.Sprintf("%d week %d, county: %s, foreign: %s", c.Year, c.Week, c.County, c.Foreign)
}

func main() {
	resp, err := http.Get(JSONURL)
	if nil != err {
		fmt.Printf("get %s with error: %s\n", JSONURL, err)
		os.Exit(OSExitCode)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		fmt.Println("read response with error: ", err)
		os.Exit(OSExitCode)
	}

	var arr []covid
	_ = json.Unmarshal([]byte(string(data)), &arr)

	weekCount, countyCount, total := aggregate(arr)
	var weeks []int
	for k := range weekCount {
		weeks = append(weeks, k)
	}

	fmt.Println()
	fmt.Println("aggregate by week")
	fmt.Println()

	sort.Ints(weeks)
	for _, i := range weeks {
		fmt.Printf("week %d, count: %d\n", i, weekCount[i])
	}

	fmt.Println("aggregate by county")
	fmt.Println()
	for k, v := range countyCount {
		fmt.Printf("county %s, count: %d\n", k, v)
	}

	fmt.Println()
	fmt.Println("total: ", total)
}

func aggregate(data []covid) (map[int]int, map[string]int, int) {
	weekCount := make(map[int]int)
	countyCount := make(map[string]int)
	total := 0

	for _, d := range data {
		if _, ok := weekCount[d.Week]; !ok {
			weekCount[d.Week] = 1
		} else {
			weekCount[d.Week] += d.ConfirmedCount
		}
		if _, ok := countyCount[d.County]; !ok {
			countyCount[d.County] = 1
		} else {
			countyCount[d.County] += d.ConfirmedCount
		}

		total += d.ConfirmedCount
	}
	return weekCount, countyCount, total
}
