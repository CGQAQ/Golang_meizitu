/**
	always report target time out issue
 */


// Command text is a chromedp example demonstrating how to extract text from a
// specific element.
package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}


	// run task list
	var res string
	err = c.Run(ctxt, text(&res))
	if err != nil {
		log.Fatal(err)
	}
	c.Wait()

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("overview: %s", res)
}

func text(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://baidu.com`),
		chromedp.Text(`#pkg-overview`, res, chromedp.NodeVisible, chromedp.ByID),
	}
}