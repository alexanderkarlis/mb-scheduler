package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/tebeka/selenium"
)

var (
	// selenium opts
	opts = []selenium.ServiceOption{
		selenium.Output(os.Stderr),
	}
	dbport     = 5432
	dbuser     = "postgres"
	dbpassword = "postgres"
	dbname     = "mb_scheduler_db"
)

var db *sql.DB
var err error
var psqlInfo string

// opts here to combat whether or not env var is set
func init() {
	log.Println("in init")
	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		dbhost = "0.0.0.0"
	}
	if err != nil {
		panic(err)
	}
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)
	db, err = sql.Open("postgres", psqlInfo)
}

const (
	// SeleniumPort is the listening port for running the proxy
	SeleniumPort = 4444
	// serverPort for the running UI.
	serverPort             = ":8888"
	seleniumPath           = "./adds/selenium-server-standalone-3.141.59.jar"
	geckoDriverPath        = "./adds/geckodriver"
	geckoDriverPathWindows = "./adds/geckodriver_windows.exe"
	baseURL                = "https://mindbody.io/"
	actualURL              = "https://www.mindbodyonline.com/explore/locations/elite-core-fitness"

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
// postgres
type ScheduleDatum struct {
	// TimeToExecute is in epoch time ms of the class
	// `datetime - 1 day + 10min`
	TimeToExecute int64  `json:"runtime" db:"runtime"`
	FullName      string `json:"fullname" db:"fullname"`
	UserName      string `json:"username" db:"username"`
	Password      string `json:"password" db:"password"`
	ClassTime     string `json:"classtime" db:"classtime"`
	DayOfWeek     string `json:"weekday" db:"weekday"`
	Date          string `json:"date" db:"date"`
	Status        string `json:"status" db:"status"`
}

// ScheduleData - Slice of ScheduleDatum
type ScheduleData []ScheduleDatum

func statusUpdate() string { return "" }

func main() {
	log.Println("Starting services...")
	log.Printf("Using port -> %+s\n", serverPort)

	http.HandleFunc("/", newSignupHandler)
	http.HandleFunc("/status", serverStatusHandler)
	http.HandleFunc("/all_times", getAllSchedulesHandler)
	http.HandleFunc("/delete_schedule", deleteScheduledDateHandler)
	http.HandleFunc("/get_run_history", getRunHistoryHandler)

	go func() {
		log.Fatal(http.ListenAndServe(serverPort, nil))
	}()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	var nextRowData *ScheduleData
	var firstRow ScheduleDatum
	for {
		select {
		case <-ticker.C:
			nextRowData = getAllSchedules(1)
			log.Println("Checking rows")
			if len(*nextRowData) > 0 {
				firstRow = (*nextRowData)[0]
				log.Printf("firstRow: %+v\n", firstRow)
				if firstRow.TimeToExecute < time.Now().Unix() && firstRow.TimeToExecute != 0 {
					success := SignUp(shortDOW(firstRow.DayOfWeek), firstRow.ClassTime, firstRow.FullName, firstRow.UserName, firstRow.Password)
					log.Printf("%+v\n", success)
					firstRow.setRowHistory(success)
					firstRow.TimeToExecute = 0
					if success {
						log.Println("SUCCESS!")
					} else {
						log.Println("FAILED!")
					}
				} else {
					log.Printf("NOW %d\n", time.Now().Unix())
					log.Printf("TESTING %d\n", firstRow.TimeToExecute)
					log.Println("NOT executed")
				}
			}
		}
	}

}

// Pretty prints ScheduleData struct.
func (d *ScheduleData) pp() {
	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		log.Println("error:", err)
	}
	log.Println(string(b))
}

func (u *User) pp() {
	b, err := json.MarshalIndent(u, "", "    ")
	if err != nil {
		log.Println("error:", err)
	}
	log.Println(string(b))
}

// CalculateSignUpTimes take the desired Class day of the
// week and the class time, and calculates the next occurence
// of that desired class.
func (u *User) CalculateSignUpTimes() *ScheduleData {
	reqDOW := days[u.Schedule.DayOfWeek]
	todayDOW := days[time.Now().Weekday().String()]
	// how many days from today
	var daysPastFromToday int
	var parsedClassTime time.Time
	parsedClassTime, err := time.Parse("03:04pm EST", fmt.Sprintf("%011s", u.Schedule.ClassTime+" EST"))
	if err != nil {
		panic("error parsing date")
	}

	// has the classTime (on the same day) already happened??
	if reqDOW == todayDOW {
		nh, nm, ns := time.Now().Local().Clock()
		ph, pm, ps := parsedClassTime.Hour(), parsedClassTime.Minute(), parsedClassTime.Second()

		// check and see if classtime has already passed for today
		if nh*3600+nm*60+ns < ph*3600+pm*60+ps {
			daysPastFromToday = 0
		} else {
			daysPastFromToday = 7
		}
	} else if reqDOW > todayDOW {
		daysPastFromToday = reqDOW - todayDOW
	} else {
		// get to the end of the week, then
		// count back up to the reqTime.
		daysPastFromToday = (7 - todayDOW) + (reqDOW)
	}

	// fcrt - first class calculated run time
	fcrt := time.Now().Local().AddDate(0, 0, daysPastFromToday)
	// add the calculated class date to the incoming object
	m := fcrt.AddDate(0, 0, -1)

	runT := time.Date(
		m.Year(),
		m.Month(),
		m.Day(),
		parsedClassTime.Hour(),
		parsedClassTime.Minute()+10, // add ten minutes to runtime
		parsedClassTime.Second(),
		0,
		parsedClassTime.Location(),
	)
	freq, err := strconv.Atoi(u.Schedule.Frequency)
	var d ScheduleData

	for i := 0; i < freq; i++ {
		runDate := runT.AddDate(0, 0, 7*i).Add(5 * time.Hour)
		classDate := fcrt.AddDate(0, 0, 7*i)
		d = append(
			d,
			ScheduleDatum{
				TimeToExecute: runDate.Unix(),
				FullName:      u.FullName,
				UserName:      u.UserName,
				Password:      u.Password,
				ClassTime:     u.Schedule.ClassTime,
				DayOfWeek:     u.Schedule.DayOfWeek,
				Date:          classDate.Format("01/02/2006"),
				Status:        "scheduled",
			},
		)
	}
	return &d
}

func replaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func getAllSchedules(limit int) *ScheduleData {
	var s ScheduleDatum
	var id string
	var data ScheduleData
	var orderedStmt string
	if limit > 0 {
		orderedStmt = fmt.Sprintf(`SELECT * FROM schedule_rt WHERE status = 'scheduled' ORDER BY runtime ASC LIMIT %d`, limit)
	} else {
		orderedStmt = `SELECT * FROM schedule_rt WHERE status = 'scheduled' ORDER BY runtime ASC`
	}

	rows, err := db.Query(orderedStmt)
	if err != nil {
		log.Println(rows)
		log.Println("error on getting rows")
		panic(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&id, &s.TimeToExecute, &s.FullName, &s.UserName, &s.Password, &s.ClassTime, &s.DayOfWeek, &s.Date, &s.Status)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, s)
	}

	return &data
}

// deletes row based on runtime key(?), runs query to remove dupes
// returns total dataset
func deleteScheduledDate(runtime int64) *ScheduleData {
	var s ScheduleDatum
	var id string
	var data ScheduleData

	deleteQuery := `DELETE FROM schedule_rt
					WHERE runtime = $1`

	orderedStmt := `SELECT * FROM schedule_rt ORDER BY runtime ASC`
	delStmt, err := db.Prepare(deleteQuery)
	if err != nil {
		log.Println(err)
	}
	resp, err := delStmt.Exec(runtime)
	log.Println("execute sql query response", resp)
	if err != nil {
		log.Println(err)
	}

	rows, err := db.Query(orderedStmt)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &s.TimeToExecute, &s.FullName, &s.UserName, &s.Password, &s.ClassTime, &s.DayOfWeek, &s.Date, &s.Status)

		if err != nil {
			log.Fatal(err)
		}
		log.Println(s)
		data = append(data, s)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return &data
}

// deletes row based on runtime key(?), runs query to remove dupes
// returns total dataset
func getRunHistory() *ScheduleData {
	var s ScheduleDatum
	var id string
	var data ScheduleData

	runHistoryQuery := `SELECT * 
					FROM schedule_rt
					WHERE status != 'scheduled'
					ORDER BY runtime ASC;`

	rows, err := db.Query(runHistoryQuery)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &s.TimeToExecute, &s.FullName, &s.UserName, &s.Password, &s.ClassTime, &s.DayOfWeek, &s.Date, &s.Status)

		if err != nil {
			log.Fatal(err)
		}
		data = append(data, s)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return &data
}

// setRowHistory sets the current row data to success or
// failure based on whether the SignUp function is successful
func (s *ScheduleDatum) setRowHistory(success bool) {
	if success {
		s.Status = "success"
	} else {
		s.Status = "failed"

	}
	updateRowQuery := `UPDATE schedule_rt SET status = $1 WHERE runtime = $2;`
	preparedQuery, err := db.Prepare(updateRowQuery)
	if err != nil {
		log.Printf("Prepare query: %+v\n", s)
	}
	updateResult, err := preparedQuery.Exec(s.Status, s.TimeToExecute)
	if err != nil {
		log.Printf("Could not update row: %+v\n", s)
	}
	log.Printf("update result: %+v\n", updateResult)
}

// select groupby -> calculate next run time -> return value to ui
func getGroupedSchedules() {
	var username, weekday, classtime string
	groupStmt := `SELECT username, weekday, classtime 
					FROM schedule_rt 
					GROUP BY weekday, classtime, username`

	rows, err := db.Query(groupStmt)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&username, &weekday, &classtime)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(username, weekday, classtime)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Println(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Content-Type", "application/json")
}

// PrepareQuery function, puts data into Postgres db
func (d *ScheduleData) PrepareQuery() (string, []interface{}) {
	insertQuery := "INSERT INTO schedule_rt(runtime, fullname, username, password, classtime, weekday, date, status) VALUES "
	vals := []interface{}{}
	for _, row := range *d {
		insertQuery += "(?, ?, ?, ?, ?, ?, ?, ?), "
		vals = append(vals, row.TimeToExecute, row.FullName, row.UserName, row.Password, row.ClassTime, row.DayOfWeek, row.Date, row.Status)
	}

	insertQuery = strings.TrimSuffix(insertQuery, ", ")
	insertQuery = replaceSQL(insertQuery, "?")

	log.Println(insertQuery)
	return insertQuery, vals
}

func getRunHistoryHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	switch r.Method {
	case "GET":
		data := getRunHistory()

		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Println(err)
		}
	case "POST":
		http.Error(w, "Bad request verb", http.StatusBadRequest)
	}
}

// TODO: FINISH DELETE. DO DB QUERY AND RETURN ALL RESULTS
func deleteScheduledDateHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	switch r.Method {
	case "GET":
		http.Error(w, "Bad request verb", http.StatusBadRequest)
	case "POST":
		var rt struct {
			Runtime int64 `json:"runtime"`
		}
		log.Println("in delete")
		err := json.NewDecoder(r.Body).Decode(&rt)
		defer r.Body.Close()

		data := deleteScheduledDate(rt.Runtime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("bad request")
		}

		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Println(err)
		}
	}
}

func getAllSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	a := getAllSchedules(-1)
	err := json.NewEncoder(w).Encode(&a)
	if err != nil {
		panic("error!")
	}
}

func serverStatusHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	switch r.Method {
	case "GET":
		s := &ServerStatus{Status: "ok"}
		err := json.NewEncoder(w).Encode(&s)
		if err != nil {
			panic("error!")
		}
	default:
		http.Error(w, "Bad request verb", http.StatusBadRequest)
	}
}

// TOOD: move POST function to separate function
func newSignupHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	switch r.Method {
	case "GET":
		http.Error(w, "Bad request verb", http.StatusBadRequest)
	case "POST":
		var data ScheduleData

		var u User
		err = json.NewDecoder(r.Body).Decode(&u)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("bad request")
		}
		data = *u.CalculateSignUpTimes()
		pstmt, values := data.PrepareQuery()
		stmt, err := db.Prepare(pstmt)
		if err != nil {
			log.Println(err)
		}
		resp, err := stmt.Exec(values...)
		log.Println("execute sql query response", resp)
		if err != nil {
			log.Println(err)
		}

		deleteDupesQuery, _ := ioutil.ReadFile("../scripts/delete_dupes.sql")

		removeDupesStmt, err := db.Prepare(string(deleteDupesQuery))
		if err != nil {
			log.Println(err)
		}
		removeDupesResp, err := removeDupesStmt.Exec()
		log.Println("removeDupesStmt sql query response", removeDupesResp)
		if err != nil {
			log.Println(err)
		}

		err = json.NewEncoder(w).Encode(&u)
		if err != nil {
			log.Println(err)
		}
	}
}

// SignUp signs up the user for a specified class time.
// ** @params weekday, date, classTime, fullName, userName, password string
func SignUp(weekday, classTime, fullName, userName, password string) bool {
	currentWorkingDirectory, err := os.Getwd()
	log.Println(currentWorkingDirectory)
	selenium.SetDebug(true)

	service, err := selenium.NewSeleniumService(seleniumPath, SeleniumPort, opts...)
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
