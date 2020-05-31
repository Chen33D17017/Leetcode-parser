package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

const url = "https://leetcode.com/problemset/all/"
const inputBox = `#question-app > div > div:nth-child(2) > div.question-list-base > div.question-filter-base > div.row.search-bar-control > div.col-sm-6.col-lg-7 > input`
const problemSelector = `#question-app > div > div:nth-child(2) > div.question-list-base > div.table-responsive.question-list-table > table > tbody.reactable-data > tr:nth-child(1) > `
const problemLevel = problemSelector + `td:nth-child(6)`
const problemName = problemSelector + `td:nth-child(3)`
const problemAddr = problemName + ` > div > a`

type config struct {
	StoreType string `json:"type"`
	ID        string `json:"id"`
	PW        string `json:"password"`
	DB        string `json:"database"`
}

var wg sync.WaitGroup
var tokens = make(chan int, 20)

func parseProblem(problem string, dbm dbManager) {
	var name, addr, level string
	problemID, err := strconv.Atoi(problem)
	tokens <- 1

	defer func() {
		wg.Done()
		<-tokens
	}()

	if dbm.checkExist(problemID) {
		fmt.Printf("Problem %s already exists\n", problem)
		return
	}
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	cmd := fmt.Sprintf(`document.querySelector("%s").href`, problemAddr)
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body > footer`),
		chromedp.SendKeys(inputBox, problem),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Text(problemName, &name),
		chromedp.Text(problemLevel, &level),
		chromedp.Evaluate(cmd, &addr),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = dbm.insertProblem([]string{problem, strings.TrimSpace(name), strings.ToUpper(level), addr})
	if err != nil {
		log.Fatalf("parseProblem: %s", err)
	}
	fmt.Printf("Search Result for id.%s: %s %s %s\n", problem, name, level, addr)

	chromedp.Cancel(ctx)
}

func main() {
	file, err := ioutil.ReadFile("config.json")

	if err != nil {
		log.Fatal("Fail to open config file")
	}

	configData := config{}
	err = json.Unmarshal([]byte(file), &configData)
	if err != nil {
		log.Fatalf("Read json err: %s", err.Error())
	}

	dbm, err := newDBManager(configData)
	if err != nil {
		log.Fatalf("main: fail to connect db: %s", err)
	}
	defer dbm.Close()

	proNum := 1459
	// use goroutine to parse data: 20 routines?
	for i := 1; i <= proNum; i++ {
		wg.Add(1)
		go parseProblem(fmt.Sprintf("%d", i), dbm)
	}

	wg.Wait()
}
