package mindreader

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

const (
	// SeleniumPort is the listening port for running the proxy
	SeleniumPort = 4444
	// serverPort for the running UI.
	seleniumPath    = "./adds/selenium-server-standalone-3.141.59.jar"
	geckoDriverPath = "./adds/geckodriver"
	baseURL         = "https://mindbody.io/"
	actualURL       = "https://www.mindbodyonline.com/explore/locations/elite-core-fitness"

	// some of these are full Xpaths, thus if the website under goes any cosmetic UI chnages,
	// it is likely these will break
	dismissNotificationXPath = "//span[@class=\"notification-dismiss\"]"
	loginXpath               = "//button[contains(text(), 'Login')]"
	emailInputXpath          = "/html/body/div[6]/div/div/div/form/div[1]/div[1]/div/input"
	passwordInputXpath       = "/html/body/div[6]/div/div/div/form/div[2]/div[1]/div/input"
	loginSubmitXpath         = "/html/body/div[6]/div/div/div/form/button[1]"
	loginSubmit2Xpath        = "/html/body/div[6]/div/div/div/form/button[1]"
	namedBoxXpath            = "//div[contains(text(), '%s')]"
	scheduleXpath            = "//a[contains(text(), 'Schedule')]"
	dayOfWeekXpath           = "//div[contains(text(), '%s')]"
	timeOfDayXpath           = "//h5[contains(text(), '%s')]"
	timeSlotClassName        = "columns is-vcentered ClassTimeScheduleItemDesktop_separator__1vvuL"
	bookNowXpath             = "//div/h5[contains(text(), \"%s\")]/../../div[@class=\"column\"]//button[contains(text(), 'Book Now')]"
	orderTotalXpath          = "//h5[contains(text(), '$0.00')]"
	buyXpath                 = "//button[contains(text(), 'Buy')]"
	orderConfimationXpath    = "//p[contains(text(), 'Order Confirmed')]"
	searchBoxXpath           = "//div[contains(text(), '%s')]/../div/div/input"
	sendSearchIconXpath      = "//button[@data-name='Search.Desktop.Button']"
	cardSearchSelection      = "//div[contains(@class, 'Card_card')]//div[contains(@class, 'CardContent')]"
	cardTitle                = "//h3[contains(@class, 'StudioDetails_title')]"
	cardAddress              = "//div[contains(@class, 'StudioDetails_address')]"
)

var (
	seleniumWdOpts = []selenium.ServiceOption{ // selenium opts
		selenium.Output(os.Stderr),
	}
)

// Run is the interafce to mindreader package
type Run interface {
	MindReader()
}

// MindReader , main function call
func MindReader() {
	fmt.Println("hello package from pkg")
}

// GetClassTimes gets all class times for selected location.
func GetClassTimes() bool {
	p := fmt.Println
	fmt.Println(os.Getwd())
	selenium.SetDebug(false)

	service, err := selenium.NewSeleniumService(seleniumPath, SeleniumPort)
	if err != nil {
		log.Println(err)
		return false
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "firefox", "headless": true}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://0.0.0.0:%d/wd/hub", SeleniumPort))
	if err != nil {
		log.Println(err)
		return false
	}

	driver.ResizeWindow("", 1800, 1550)

	if err := driver.Get(baseURL); err != nil {
		log.Println(err)
		return false
	}
	defer driver.Quit()
	driver = login("alexanderkarlis@gmail.com", "921921Zz?", driver)

	searchFor, err := driver.FindElement(selenium.ByXPATH, fmt.Sprintf(searchBoxXpath, "Search for anything"))
	if err != nil {
		panic(err)
	}
	searchFor.SendKeys("elite core fitness")
	time.Sleep(2 * time.Second)
	sendSearch, err := driver.FindElement(selenium.ByXPATH, sendSearchIconXpath)
	if err != nil {
		panic(err)
	}
	sendSearch.Click()

	time.Sleep(1 * time.Second)
	searchIcons, err := driver.FindElements("xpath", cardSearchSelection)

	if err != nil {
		panic(err)
	}
	p(searchIcons)
	for _, x := range searchIcons {
		gymText, err := x.Text()
		if err != nil {
			panic(err)
		}
		p(gymText)
		cardTitleObj, err := x.FindElement("xpath", cardTitle)
		if err != nil {
			panic(err)
		}
		_, err = cardTitleObj.Text()
		if err != nil {
			panic(err)
		}

		cardAddressObj, err := x.FindElement("xpath", cardAddress)
		if err != nil {
			panic(err)
		}
		_, err = cardAddressObj.Text()
		if err != nil {
			panic(err)
		}
		// fmt.Printf("%s - %s\n", gymTitleName, gymTitleAddress)
	}
	for {
	}
	return true
}

// SignUp signs up the user for a specified class time.
// ** @params weekday, date, classTime, fullName, userName, password string
func SignUp(weekday, classTime, fullName, userName, password string) bool {
	selenium.SetDebug(true)

	service, err := selenium.NewSeleniumService(seleniumPath, SeleniumPort, seleniumWdOpts...)
	if err != nil {
		log.Println(err)
		return false
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "firefox", "headless": true}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://0.0.0.0:%d/wd/hub", SeleniumPort))
	if err != nil {
		log.Println(err)
		return false
	}

	driver.ResizeWindow("", 1264, 1228)

	if err := driver.Get(baseURL); err != nil {
		log.Println(err)
		return false
	}
	defer driver.Quit()

	driver = login(userName, password, driver)

	driver.SetImplicitWaitTimeout(2000 * time.Millisecond)

	dismissButton, err := driver.FindElement(selenium.ByXPATH, dismissNotificationXPath)
	if err != nil {
		log.Println(err)
		return false
	}
	dismissButton.Click()
	var nameBox selenium.WebElement

	// some trouble finding the element, so used a recurse
	loopBackoff := 0
	for {
		nameBox, err = driver.FindElement(selenium.ByXPATH, fmt.Sprintf(namedBoxXpath, fullName))
		if err != nil {
			log.Println(err)
			break
		}
		t, err := nameBox.Text()
		checkError(err)
		if t != "" {
			break
		}
		if loopBackoff > 20 {
			log.Println("CRASH LOOP, BACKOFF -- orderTotalText")
			return false
		}
		time.Sleep(2000 * time.Millisecond)
		log.Printf("loopBackoff at %d\n", loopBackoff)
		loopBackoff++
	}
	err = nameBox.Click()
	if checkError(err) {
		return false
	}

	if err := driver.Get(actualURL); err != nil {
		return false
	}

	time.Sleep(1000 * time.Millisecond)
	var dowDiv selenium.WebElement
	for {
		dowDiv, err = driver.FindElement(selenium.ByXPATH, fmt.Sprintf(dayOfWeekXpath, weekday))
		if err == nil {
			log.Println("found DOW button")
			break
		}
		if loopBackoff > 20 {
			log.Println("CRASH LOOP, BACKOFF -- orderTotalText")
			return false
		}
		time.Sleep(2000 * time.Millisecond)
		log.Printf("loopBackoff at %d\n", loopBackoff)
		loopBackoff++
	}

	dowDiv.Click()

	var classTimeBox selenium.WebElement
	loopBackoff = 0
	for {
		classTimeBox, err = driver.FindElement(selenium.ByXPATH, fmt.Sprintf(timeOfDayXpath, classTime))
		cText, err := classTimeBox.Text()
		if cText != "" || err == nil {
			break
		}

		if loopBackoff > 20 {
			log.Println("CRASH LOOP, BACKOFF -- orderTotalText")
			return false
		}
		time.Sleep(2000 * time.Millisecond)
		log.Printf("loopBackoff at %d\n", loopBackoff)
		loopBackoff++
	}

	bookNowButton, err := driver.FindElement(selenium.ByXPATH, fmt.Sprintf(bookNowXpath, classTime))
	if checkError(err) {
		return false
	}
	bookText, _ := bookNowButton.Text()
	log.Println(bookText)
	bookNowButton.Click()

	time.Sleep(2000 * time.Millisecond)
	loopBackoff = 0
	for {
		orderAmt, err := driver.FindElement(selenium.ByXPATH, orderTotalXpath)
		if err == nil {
			orderTotalText, err := orderAmt.Text()
			log.Println(orderTotalText)

			if orderTotalText != "" || err == nil {
				log.Println(orderTotalText)
				break
			}
		}
		if loopBackoff > 20 {
			log.Println("CRASH LOOP, BACKOFF -- orderTotalText")
			return false
		}
		time.Sleep(2000 * time.Millisecond)
		log.Printf("loopBackoff at %d\n", loopBackoff)
		loopBackoff++
	}

	// driverSnapshot(driver, "buy")
	loopBackoff = 0
	for {
		time.Sleep(2000 * time.Millisecond)
		buyButton, err := driver.FindElement(selenium.ByXPATH, buyXpath)
		if err == nil {
			buyButton.Click()
			break
		}
		if loopBackoff > 20 {
			log.Println("CRASH LOOP, BACKOFF -- orderTotalText")
			return false
		}
		time.Sleep(2000 * time.Millisecond)
		log.Printf("loopBackoff at %d\n", loopBackoff)
		loopBackoff++
	}

	// Confirm a successful sign up
	// driverSnapshot(driver, "")
	loopBackoff = 0
	for {
		confirmationSpan, err := driver.FindElement(selenium.ByXPATH, orderConfimationXpath)
		if err == nil {
			confirmationSpanText, err := confirmationSpan.Text()

			if confirmationSpanText != "" || err == nil {
				log.Println(confirmationSpanText)
				break
			}
		}
		if loopBackoff > 20 {
			log.Println("CRASH LOOP, BACKOFF -- orderTotalText")
			return false
		}
		time.Sleep(2000 * time.Millisecond)
		log.Printf("loopBackoff at %d\n", loopBackoff)
		loopBackoff++
	}

	// driverSnapshot(driver, "")
	return true
}

// login function takes the username and password and signs into Mindbody.io
func login(username, password string, wd selenium.WebDriver) selenium.WebDriver {
	loginBtn, err := wd.FindElement("xpath", loginXpath)
	checkError(err)
	loginBtn.Click()

	emailInput, err := wd.FindElement("xpath", emailInputXpath)
	checkError(err)
	emailInput.SendKeys(username)

	passwordInput, err := wd.FindElement("xpath", passwordInputXpath)
	// wd.WaitWithTimeout
	checkError(err)
	passwordInput.SendKeys(password)

	loginSubmitBtn, err := wd.FindElement("xpath", loginSubmitXpath)
	checkError(err)

	err = loginSubmitBtn.Click()
	checkError(err)
	return wd
}

func checkError(err error) bool {
	if err != nil {
		log.Println("#######")
		log.Println(err)
		log.Println("#######")
		log.Println("Error: \n", err)
		return true
	}
	return false
}

func driverSnapshot(webdriver selenium.WebDriver, fileName string) bool {
	if fileName == "" {
		fileName = "sc"
	}
	imgBytes, err := webdriver.Screenshot()

	checkError(err)

	err = ioutil.WriteFile(fileName, imgBytes, 0777)
	if err != nil {
		return false
	}
	return true
}
