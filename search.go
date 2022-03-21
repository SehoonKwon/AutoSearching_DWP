package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
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

	fmt.Println("이름 또는 휴대폰번호 입력")
	var searchKey string
	fmt.Scan(&searchKey)

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

	time.Sleep(2 * time.Second)
	btn, _ := session.FindElement(selenium.ByCSSSelector, "li.e-sch")
	btn.Click()

	time.Sleep(1 * time.Second)
	combo, _ := session.FindElement(selenium.ByXPATH, "//*[@id='bzrForm']/div[1]/div[2]/div/div[1]/div[2]/div[2]/div[2]")

	//입력받은 값이 이름인지, 휴대폰 번호인지 판별
	matched, _ := regexp.MatchString("[0-9]+", searchKey)
	if matched {
		combo, _ = combo.FindElement(selenium.ByXPATH, "//*[@id='bzrForm']/div[1]/div[2]/div/div[1]/div[2]/div[2]/div[2]/select/option[4]")
	} else {
		combo, _ = combo.FindElement(selenium.ByXPATH, "//*[@id='bzrForm']/div[1]/div[2]/div/div[1]/div[2]/div[2]/div[2]/select/option[1]")
	}
	combo.Click()

	inputPn, _ := session.FindElement(selenium.ByXPATH, "//*[@id='bzrForm']/div[1]/div[2]/div/div[1]/div[2]/div[2]/div[1]/input")
	inputPn.Click()
	fmt.Println(inputPn.Text())
	inputPn.SendKeys(searchKey)
	inputPn.SendKeys(selenium.EnterKey)

	time.Sleep(5 * time.Second)
	session.Delete()
	chromeDriver.Stop()
}
