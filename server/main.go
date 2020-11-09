package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"strconv"
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
	// serverPort for the running UI.
	serverPort             = ":8888"
	seleniumPath           = "./adds/selenium-server-standalone-3.141.59.jar"
	geckoDriverPath        = "./adds/geckodriver"
	geckoDriverPathWindows = "./adds/geckodriver_windows.exe"
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

const (
	every Frequency = iota
	single
)

// Days - slice of day
type Days []string

var days = map[string]int{
	"Monday":    0,
	"Tuesday":   1,
	"Wednesday": 2,
	"Thursday":  3,
	"Friday":    4,
	"Saturday":  5,
	"Sunday":    6,
}

// Frequency of signups
type Frequency int

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

// Schedule type for mindbody initialization
type Schedule struct {
	ClassTime string `json:"classtime"`
	DayOfWeek string `json:"weekday"`
	Date      string `json:"date"`
	Frequency string `json:"frequency"`
}

// ScheduleDatum is the data to be prepared and inserted into
// mongodb
type ScheduleDatum struct {
	// TimeToExecute is in epoch time ms of the class
	// `datetime - 1 day + 10min`
	TimeToExecute int64 `json:"time"`
	User          User  `json:"user"`
}

// ScheduleData - Slice of ScheduleDatum
type ScheduleData []ScheduleDatum

func main() {
	// p := fmt.Println
	// fmt.Println("Starting services...")
	// fmt.Printf("Using port -> %+s\n", serverPort)

	s := Schedule{ClassTime: "5:45pm", Date: "", DayOfWeek: "Saturday", Frequency: "10"}
	u := User{FullName: "Alexander Karlis", UserName: "alexanderkarlis@gmail.com", Password: "password", Schedule: s}
	u.calculateSignUpTimes()
	// http.HandleFunc("/", newSignupHandler)
	// go func() {
	// 	log.Fatal(http.ListenAndServe(serverPort, nil))
	// }()

	// _, weekday, err := parseDate("11/02/2020")
	// if err != nil {
	// 	panic("panicked")
	// }
	// u := User{FullName: "Alexander Karlis", Password: "921921Zz?", UserName: "alexanderkarlis@gmail.com", Schedule: Schedule{ClassTime: "5:45pm", Date: "", Frequency: "1"}}
	// SignUp(weekday, u.Schedule.ClassTime, u.FullName, u.UserName, u.Password)
}

// returns the number of days until the next weeekday, int
func (d *ScheduleData) pp() {
	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}

func (u *User) calculateSignUpTimes() *ScheduleDatum {
	p := fmt.Println
	reqTime := days[u.Schedule.DayOfWeek]
	todayTime := days[time.Now().Weekday().String()]

	// how many days from today
	var daysPastFromToday int
	var parsedClassTime time.Time
	parsedClassTime, err := time.Parse("03:04pm", fmt.Sprintf("%07s", u.Schedule.ClassTime))
	if reqTime == todayTime {
		p("today -> ", time.Now().Weekday())
		fmt.Printf("todayDate: %d\n", todayTime)
		fmt.Printf("reqDate: %d\n", reqTime)
		p(u.Schedule.ClassTime)

		if err != nil {
			panic("error parsing date")
		}

		// has the classTime (on the same day) already happened??
		nh, nm, ns := time.Now().Local().Clock()
		ph, pm, ps := parsedClassTime.Hour(), parsedClassTime.Minute(), parsedClassTime.Second()

		// check and see if classtime has already passed for today
		if nh*3600+nm*60+ns < ph*3600+pm*60+ps {
			daysPastFromToday = 0
			p("class has NOT happened yet.")
		} else {
			daysPastFromToday = 7
			p("class has ALREADY happened. starting with next instance of day")
		}

	} else if reqTime > todayTime {
		daysPastFromToday = reqTime - todayTime
	} else {
		// get to the end of the week, then
		// count back up to the reqTime.
		daysPastFromToday = (6 - reqTime - 1) + (todayTime)
	}

	// fcrt first class calculated run time
	fcrt := time.Now().AddDate(0, 0, daysPastFromToday)
	m := fcrt.AddDate(0, 0, -1)
	fmt.Printf("days from today: %d\n", daysPastFromToday)

	runT := time.Date(
		m.Year(),
		m.Month(),
		m.Day(),
		parsedClassTime.Hour(),
		parsedClassTime.Minute()+10,
		parsedClassTime.Second(),
		0,
		parsedClassTime.Location(),
	)

	var d ScheduleData
	freq, err := strconv.Atoi(u.Schedule.Frequency)
	for i := 0; i < freq; i++ {
		runDate := runT.AddDate(0, 0, 7*i)
		u.Schedule.Date = runDate.Format("01/02/2006")
		d = append(d, ScheduleDatum{TimeToExecute: runDate.Unix(), User: *u})
	}

	d.pp()

	return &ScheduleDatum{TimeToExecute: runT.Unix(), User: *u}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Content-Type", "application/json")
}

func newSignupHandler(w http.ResponseWriter, r *http.Request) {
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
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			panic("bad request")
		}
		_, weekday, err := parseDate(u.Schedule.Date, u.Schedule.DayOfWeek)
		SignUp(weekday, u.Schedule.ClassTime, u.FullName, u.UserName, u.Password)
		err = json.NewEncoder(w).Encode(&u)
	}
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

func parseDate(dt, weekday string) (string, string, error) {
	parsedTime, err := time.Parse("01/02/2006", dt)
	checkError(err)
	day := fmt.Sprintf("%d", parsedTime.Day())
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
