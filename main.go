package main

import (
	"context"
	"fmt"
	"log"
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

var wg sync.WaitGroup
var tokens = make(chan int, 5)

func parseProblem(problem string) {
	var name, addr, level string
	tokens <- 1
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := fmt.Sprintf(`document.querySelector("%s").href`, problemAddr)
	err := chromedp.Run(ctx,
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

	fmt.Printf("Search Result for id.%s: %s %s %s\n", problem, name, level, addr)
	wg.Done()
	chromedp.Cancel(ctx)
	<-tokens
}

func main() {
	proNum := 200
	// use goroutine to parse data: 20 routines?
	for i := 1; i <= proNum; i++ {
		wg.Add(1)
		go parseProblem(fmt.Sprintf("%d", i))
	}

	wg.Wait()
}
