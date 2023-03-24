package clubhouse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type logger interface {
	Debug(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
}

func NewApiClient(baseUrl string, httpClient *http.Client, logger logger) *ApiClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if httpClient.Jar == nil {
		jar, _ := cookiejar.New(nil)
		httpClient.Jar = jar
	}

	return &ApiClient{
		baseUrl:    baseUrl,
		httpClient: httpClient,
		logger:     logger,
	}
}

func (c *ApiClient) Login(username string, password string) (*User, error) {
	viewstate, err := c.getViewState(fmt.Sprintf("%s/login.aspx", c.baseUrl))
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	loginUrl := fmt.Sprintf("%s/login.aspx", c.baseUrl)
	params := url.Values{}
	params.Add("p$lt$page_content$pageplaceholder$p$lt$zoneLeft$CHOLogin$LoginControl$ctl00$Login1$UserName", username)
	params.Add("p$lt$page_content$pageplaceholder$p$lt$zoneLeft$CHOLogin$LoginControl$ctl00$Login1$Password", password)
	params.Add("p$lt$page_content$pageplaceholder$p$lt$zoneLeft$CHOLogin$LoginControl$ctl00$Login1$LoginButton", "Login")
	params.Add("__VIEWSTATE", viewstate)

	resp, err := c.httpClient.PostForm(loginUrl, params)
	if err != nil {
		c.logger.Error("failed to perform http request", "err", err)
		return nil, fmt.Errorf("Login: %w", err)
	}
	defer resp.Body.Close()

	return c.GetCurrentUser()
}

func (c *ApiClient) GetCurrentUser() (*User, error) {
	url := fmt.Sprintf("%s/api/v1/GetCurrentUser", c.baseUrl)

	var userResponse UserResponse
	if err := c.get(url, &userResponse); err != nil {
		return nil, fmt.Errorf("GetCurrentUser: %w", err)
	}

	return &userResponse.User, nil
}

func (c *ApiClient) GetGolfCourses() ([]GolfCourse, error) {
	url := fmt.Sprintf("%s/api/v1/teetimes/GetGolfCourses", c.baseUrl)

	var golfCourseResponse GolfCourseResponse
	if err := c.get(url, &golfCourseResponse); err != nil {
		return nil, fmt.Errorf("GetGolfCourses: %w", err)
	}

	return golfCourseResponse.GolfCourses, nil
}

func (c *ApiClient) GetGolfCourse(courseName string) (GolfCourse, error) {
	courses, err := c.GetGolfCourses()
	if err != nil {
		return GolfCourse{}, fmt.Errorf("GetGolfCourse: %w", err)
	}

	for _, course := range courses {
		if course.Name == courseName {
			return course, nil
		}
	}

	return GolfCourse{}, fmt.Errorf("course \"%s\" not found", courseName)
}

func (c *ApiClient) GetTeeTimes(date time.Time, courses ...GolfCourse) ([]TeeTime, error) {
	var ids []string
	for _, course := range courses {
		ids = append(ids, strconv.Itoa(course.Id))
	}

	url := fmt.Sprintf("%s/api/v1/teetimes/GetAvailableTeeTimes/%s/%s/0/null/0", c.baseUrl, date.Format("20060102"), strings.Join(ids, ";"))

	var teeTimesResponse TeeTimesResponse
	if err := c.get(url, &teeTimesResponse); err != nil {
		return nil, fmt.Errorf("GetTeeTimes: %w", err)
	}

	return teeTimesResponse.Data.TeeSheet, nil
}

func (c *ApiClient) LookupMembers() ([]Member, error) {
	url := fmt.Sprintf("%s/api/v1/participantlookup/LookupMembers/1", c.baseUrl)

	var lookupMembersResponse LookupMembersResponse
	if err := c.get(url, &lookupMembersResponse); err != nil {
		return nil, fmt.Errorf("LookupMemebers: %w", err)
	}

	return lookupMembersResponse.Members, nil
}

func (c *ApiClient) BookTeeTime(teeTime TeeTime, bookingRequest *BookingRequest) (BookingResponse, error) {
	c.proceedBooking(teeTime)

	return c.commitBooking(bookingRequest)
}

func (c *ApiClient) proceedBooking(teeTime TeeTime) error {
	url := fmt.Sprintf("%s/api/v1/teetimes/ProceedBooking/%d", c.baseUrl, teeTime.TeeSheetTimeID)

	if err := c.get(url, nil); err != nil {
		return fmt.Errorf("proceedBooking: %w", err)
	}

	return nil
}

func (c *ApiClient) commitBooking(bookingRequest *BookingRequest) (BookingResponse, error) {
	url := fmt.Sprintf("%s/api/v1/teetimes/CommitBooking/0", c.baseUrl)

	var bookingResponse BookingResponse
	if err := c.post(url, bookingRequest, bookingResponse); err != nil {
		return BookingResponse{}, fmt.Errorf("commitBooking: %w", err)
	}

	return bookingResponse, nil
}

func (c *ApiClient) get(url string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("failed to create http request", "error", err)
		return fmt.Errorf("get: %w", err)
	}

	return c.do(req, result)
}

func (c *ApiClient) post(url string, data interface{}, result interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		c.logger.Error("failed to marshal data", "error", err)
		return fmt.Errorf("post: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		c.logger.Error("failed to create http request", "error", err)
		return fmt.Errorf("post: %w", err)
	}

	return c.do(req, result)
}

func (c *ApiClient) do(req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("failed to make http request", "error", err)
		return fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("failed to read http response body", "error", err)
		return fmt.Errorf("get: %w", err)
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		var err ErrorResponse
		if err := json.Unmarshal(body, &err); err != nil {
			c.logger.Error("failed to unmarshal http response object", "error", err)
			return fmt.Errorf("get: %w", err)
		}

		return err
	}

	if result != nil {
		if err := json.Unmarshal(body, &result); err != nil {
			c.logger.Error("failed to unmarshal http response object", "error", err)
			return fmt.Errorf("get: %w", err)
		}
	}

	return nil
}

func (c *ApiClient) isLoggedIn() bool {
	hasDotNetSessionCookie := false
	aspxFormsAuthCookieIsValid := false
	siteUrl, _ := url.Parse(c.baseUrl)

	for _, cookie := range c.httpClient.Jar.Cookies(siteUrl) {
		if cookie.Name == "ASP.NET_SessionId" {
			hasDotNetSessionCookie = true
		} else if cookie.Name == ".ASPXFORMSAUTH" && cookie.Expires.Local().Before(time.Now()) {
			aspxFormsAuthCookieIsValid = true
		}
	}

	return hasDotNetSessionCookie && aspxFormsAuthCookieIsValid
}

func (c *ApiClient) getViewState(url string) (string, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		c.logger.Error("failed to make http request", "error", err)
		return "", fmt.Errorf("getLoginPayload: %w", err)
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)

	for tt := tokenizer.Next(); tt != html.ErrorToken; tt = tokenizer.Next() {
		if tt == html.SelfClosingTagToken {
			token := tokenizer.Token()

			if token.Data == "input" {
				var viewstate string
				var viewstateInput bool

				for _, attr := range token.Attr {
					if attr.Key == "id" {
						viewstateInput = attr.Val == "__VIEWSTATE"
					}
					if attr.Key == "value" {
						viewstate = attr.Val
					}
				}

				if viewstateInput {
					return viewstate, nil
				}
			}
		}
	}

	c.logger.Error("failed to find viewstate")
	return "", fmt.Errorf("failed to find viewstate")
}
