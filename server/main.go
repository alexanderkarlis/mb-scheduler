package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

var (
	// selenium opts
	opts = []selenium.ServiceOption{
		selenium.Output(os.Stderr),
	}
)

const (
	// SeleniumPort is the listening port for running the proxy
	SeleniumPort = 4444
	// UIPort for the running UI.
	UIPort                 = ":8888"
	seleniumPath           = "../adds/selenium-server-standalone-3.141.59.jar"
	geckoDriverPath        = "../adds/geckodriver"
	geckoDriverPathWindows = "../adds/geckodriver_windows.exe"
	baseURL                = "https://mindbody.io/"
	actualURL              = "https://mindbody.io/locations/elite-core-fitness"

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
)

// Frequency of signups
type Frequency int

const (
	every Frequency = iota
	single
)

// ServerStatus sends an ok to svelte ui on mount
type ServerStatus struct {
	Status string `json:"status"`
}

// User type for mindbody initialization
type User struct {
	FullName string   `json:"name"`
	UserName string   `json:"username"`
	Password string   `json:"password"`
	Schedule Schedule `json:"schedule"`
}

type Schedule struct {
	ClassTime string `json:"classtime"`
	Date      string `json:"date"`
	Frequency string `json:"frequency"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Content-Type", "application/json")
}

func handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	switch r.Method {
	case "GET":
		s := &ServerStatus{Status: "ok"}
		err := json.NewEncoder(w).Encode(&s)
		if err != nil {
			panic("error!")
		}
	case "POST":
		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		defer r.Body.Close()
		fmt.Printf("%+v\n", r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			panic("bad request")
		}
		_, weekday, err := parseDate(u.Schedule.Date)
		SignUp(weekday, u.Schedule.ClassTime, u.FullName, u.UserName, u.Password)
		err = json.NewEncoder(w).Encode(&u)
	}
}

func main() {
	// http.HandleFunc("/", handler)
	// log.Fatal(http.ListenAndServe(UIPort, nil))
	_, weekday, err := parseDate("11/02/2020")
	if err != nil {
		panic("panicked")
	}
	u := User{FullName: "Alexander Karlis", Password: "921921Zz?", UserName: "alexanderkarlis@gmail.com", Schedule: Schedule{ClassTime: "5:45pm", Date: "", Frequency: "1"}}
	// username: "alexanderkarlis@gmail.com",
	//         name: "Alexander Karlis",
	//         password: "921921Zz?",
	//         schedule: {
	//             classtime: "5:45pm",
	//             date: "11/02/2020",
	//             frequency: "1",
	//         },
	SignUp(weekday, u.Schedule.ClassTime, u.FullName, u.UserName, u.Password)
}

// SignUp signs up the user for a specified class time.
// ** @params weekday, date, classTime, fullName, userName, password string
func SignUp(weekday, classTime, fullName, userName, password string) bool {
	currentWorkingDirectory, err := os.Getwd()
	fmt.Println(currentWorkingDirectory)
	selenium.SetDebug(true)

	service, err := selenium.NewSeleniumService(seleniumPath, SeleniumPort, opts...)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "firefox", "headless": false}
	fmt.Println("TRYING TO START SELENIUM SERVER")
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", SeleniumPort))
	driver.ResizeWindow("", 1264, 1228)
	fmt.Println("AFTER \n")
	if err != nil {
		panic(err)
	}
	if err := driver.Get(baseURL); err != nil {
		panic(err)
	}
	defer driver.Quit()

	driver = login(userName, password, driver)

	driver.SetImplicitWaitTimeout(2000 * time.Millisecond)

	dismissButton, err := driver.FindElement(selenium.ByXPATH, dismissNotificationXPath)
	checkError(err)
	dismissButton.Click()
	var nameBox selenium.WebElement

	// some trouble finding the element, so used a recurse
	for {
		nameBox, err = driver.FindElement(selenium.ByXPATH, fmt.Sprintf(namedBoxXpath, fullName))
		if err != nil {
			fmt.Println(err)
			break
		}
		t, err := nameBox.Text()
		checkError(err)
		if t != "" {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}
	err = nameBox.Click()
	checkError(err)

	if err := driver.Get(actualURL); err != nil {
		panic(err)
	}

	dowDiv, err := driver.FindElement(selenium.ByXPATH, fmt.Sprintf(dayOfWeekXpath, weekday))
	dowDiv.Click()

	var classTimeBox selenium.WebElement
	for {
		classTimeBox, err = driver.FindElement(selenium.ByXPATH, fmt.Sprintf(timeOfDayXpath, classTime))
		cText, err := classTimeBox.Text()
		if cText != "" || err == nil {
			break
		}

		time.Sleep(1000 * time.Millisecond)
	}

	bookNowButton, err := driver.FindElement(selenium.ByXPATH, fmt.Sprintf(bookNowXpath, classTime))
	checkError(err)
	bookText, _ := bookNowButton.Text()
	fmt.Println(bookText)
	bookNowButton.Click()
	// another loop to make sure all times load before moving on

	for {
		orderAmt, err := driver.FindElement(selenium.ByXPATH, orderTotalXpath)
		if err == nil {
			orderTotalText, err := orderAmt.Text()
			fmt.Println(orderTotalText)

			if orderTotalText != "" || err == nil {
				fmt.Println(orderTotalText)
				break
			}
		}
		time.Sleep(2000 * time.Millisecond)
	}

	driverSnapshot(driver, "buy")
	for {
		time.Sleep(2000 * time.Millisecond)
		buyButton, err := driver.FindElement(selenium.ByXPATH, buyXpath)
		if err == nil {
			buyButton.Click()
			break
		}
		time.Sleep(2000 * time.Millisecond)
	}

	// Confirm a successful sign up
	driverSnapshot(driver, "")
	for {
		confirmationSpan, err := driver.FindElement(selenium.ByXPATH, orderConfimationXpath)
		if err == nil {
			confirmationSpanText, err := confirmationSpan.Text()

			if confirmationSpanText != "" || err == nil {
				fmt.Println(confirmationSpanText)
				break
			}
		}
		time.Sleep(2000 * time.Millisecond)
	}

	driverSnapshot(driver, "")
	return true
}

func parseDate(dt string) (string, string, error) {
	parsedTime, err := time.Parse("01/02/2006", dt)
	checkError(err)
	day := fmt.Sprintf("%d", parsedTime.Day())
	weekday := fmt.Sprintf("%s", parsedTime.Weekday())
	return day, shortDOW(weekday), err
}

func shortDOW(day string) string {
	var short string
	switch day {
	case "Monday":
		short = "Mon"
	case "Tuesday":
		short = "Tue"
	case "Wednesday":
		short = "Wed"
	case "Thursday":
		short = "Thu"
	case "Friday":
		short = "Fri"
	case "Saturday":
		short = "Sat"
	case "Sunday":
		short = "Sun"
	default:
		short = "err"
	}
	return short
}

// login function takes the username and password and signs into
// Mindbody.io
func login(username, password string, wd selenium.WebDriver) selenium.WebDriver {

	loginBtn, err := wd.FindElement("xpath", loginXpath)
	checkError(err)
	loginBtn.Click()

	emailInput, err := wd.FindElement("xpath", emailInputXpath)
	checkError(err)
	emailInput.SendKeys(username)

	passwordInput, err := wd.FindElement("xpath", passwordInputXpath)
	checkError(err)
	passwordInput.SendKeys(password)

	loginSubmitBtn, err := wd.FindElement("xpath", loginSubmitXpath)
	checkError(err)

	err = loginSubmitBtn.Click()
	checkError(err)
	return wd
}

func checkError(err error) {
	if err != nil {
		fmt.Println("#######")
		fmt.Println(err)
		fmt.Println("#######")
		panic(err)
	}
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
