package main

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/front-ck996/csy/browserSpider"
)

func main() {
	handle := browserSpider.New(browserSpider.BrowserHandleInit{})

	err := chromedp.Run(handle.Ctx, handle.Navigate("https://www.baidu.com", ""))
	fmt.Println(err)
	//navigate :=
	//fmt.Print(navigate)
}
