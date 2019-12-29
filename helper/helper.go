package helper

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func ExtractRecordingUrl(r *http.Request) string {
	return r.FormValue("RecordingUrl")
}

func RecordingUrlExist(r *http.Request) bool {
	if ExtractRecordingUrl(r) == "" {
		return false
	}
	return true
}

func getEmergencyContacts() string {
	return os.Getenv("EmergencyContacts")
}

func getTwilioNumber() string {
	return os.Getenv("TwilioNumber")
}

func getOutboundHandlerUrl() string {
	return os.Getenv("OutboundHandlerUrl")
}

func CallFolks() {
	toNumbers := strings.Split(getEmergencyContacts(), " ")

	fromNumber := getTwilioNumber()

	for _, toNumber := range toNumbers {
		go func(toNumber string) {
			callFolks(fromNumber, toNumber, getOutboundHandlerUrl())
		}(toNumber)
	}
}

func callFolks(fromNumber string, toNumber string, outboundHandlerUrl string) {
	data := url.Values{}
	data.Set("To", toNumber)
	data.Set("From", fromNumber)
	data.Set("Url", outboundHandlerUrl)

	accountSid := os.Getenv("AccountSid")
	authToken := os.Getenv("AuthToken")

	request, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/"+accountSid+"/Calls.json", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("unable to create a new reqeust, error is %+v\n", err)
	}
	request.SetBasicAuth(accountSid, authToken)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("unable to execute the request, error is %+v\n", err)
	}
	log.Printf("response status is %v\n", response.Status)
}

func LogRequestBody(r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}
	for k, v := range r.PostForm {
		log.Printf("%s: %v", k, v)
	}
}

//TODO: Have some persistent data store for storing recording url instead of global variables
//TODO: This should be a map[string]string -> AccId : LatestRecordingUrl
var recordingUrl string

func SetRecordingUrl(rurl string) {
	recordingUrl = rurl
}
func GetRecordingUrl() string {
	return recordingUrl
}
