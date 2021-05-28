package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

const (
	chromeDriverPath = "chromedriver"
	port             = 8080
)

func main() {
	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr),
	}
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer wd.Quit()

	if err := wd.Get("http://play.golang.org"); err != nil {
		log.Fatalln(err.Error())
	}

	elem, err := wd.FindElement(selenium.ByCSSSelector, "#code")
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err := elem.Clear(); err != nil {
		log.Fatalln(err.Error())
	}

	code := `package main
import "fmt"
func main() { fmt.Println("Hello World!") }
`

	if err = elem.SendKeys(code); err != nil {
		log.Fatalln(err.Error())
	}

	btn, err := wd.FindElement(selenium.ByCSSSelector, "#run")
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err := btn.Click(); err != nil {
		log.Fatalln(err.Error())
	}

	outputDiv, err := wd.FindElement(selenium.ByCSSSelector, "#output")
	if err != nil {
		log.Fatalln(err.Error())
	}

	var output string
	for {
		output, err = outputDiv.Text()
		if err != nil {
			log.Fatalln(err.Error())
		}
		if output != "Waiting for remote server..." {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Printf("%s", strings.Replace(output, "\n\n", "\n", -1))

	time.Sleep(5 * time.Second)
}
