package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/joho/godotenv"
	"github.com/tebeka/selenium"
	_ "github.com/zserge/lorca"
)

func main() {

	/*
		// HTML로 UI를 생성. 수동 작업 할 동안 멈추는 용도
		ui, err := lorca.New("data:text/html,"+url.PathEscape(`
			<html>
				<head><title>Hello</title></head>
				<body>
					<h1>검색할 번호 입력</h1>

				</body>

			</html>
			`), "", 480, 320)
		if err != nil {
			log.Fatal(err)
		}
		defer ui.Close()

		// UI 가 닫힐 때까지 기다림
		<-ui.Done()
	*/

	fmt.Println("검색할 번호 입력")
	var pn string
	fmt.Scan(&pn)

	chromeDriver := webdriver.NewChromeDriver("./chromedriver.exe")
	err := chromeDriver.Start()
	if err != nil {
		log.Println(err)
	}
	desired := webdriver.Capabilities{"Platform": "Windows"}
	required := webdriver.Capabilities{}
	session, err := chromeDriver.NewSession(desired, required)
	if err != nil {
		log.Println(err)
	}
	err = session.Url("https://dwp.lotte.net/Group/LoginPage.bzr")
	if err != nil {
		log.Println(err)
	}

	nowUrl, _ := session.GetUrl()
	fmt.Println("세션 : ", nowUrl)

	time.Sleep(1 * time.Second)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println(os.Getenv("ID"), os.Getenv("PW"))

	//셀레니움 관련 제어 부분
	id, err := session.FindElement(selenium.ByCSSSelector, "#LoginPage_loginMain_tbxID")
	if err != nil {
		log.Println(err)
	}
	id.Click()
	id.SendKeys(os.Getenv("ID"))

	pw, _ := session.FindElement(selenium.ByCSSSelector, "#LoginPage_loginMain_tbxPwd")
	pw.SendKeys(os.Getenv("PW"))
	pw.SendKeys(selenium.EnterKey)

	nowUrl, _ = session.GetUrl()
	fmt.Println("세션 : ", nowUrl)

	time.Sleep(2 * time.Second)
	btn, _ := session.FindElement(selenium.ByCSSSelector, "li.e-sch")
	btn.Click()

	time.Sleep(1 * time.Second)
	combo, _ := session.FindElement(selenium.ByXPATH, "//*[@id='bzrForm']/div[1]/div[2]/div/div[1]/div[2]/div[2]/div[2]")
	combo, _ = combo.FindElement(selenium.ByXPATH, "//*[@id='bzrForm']/div[1]/div[2]/div/div[1]/div[2]/div[2]/div[2]/select/option[4]")
	combo.Click()

	inputPn, _ := session.FindElement(selenium.ByXPATH, "//*[@id='bzrForm']/div[1]/div[2]/div/div[1]/div[2]/div[2]/div[1]/input")
	inputPn.Click()
	fmt.Println(inputPn.Text())
	inputPn.SendKeys(pn)
	inputPn.SendKeys(selenium.EnterKey)

	fmt.Println("===================================================================================================")

	time.Sleep(5 * time.Second)
	session.Delete()
	chromeDriver.Stop()
}
