package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"telecomm-multicast/helper"
	"time"
)

const startupMessage = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Say>Hello there, please record your emergency message</Say>
	<Record timeout="3" recordingStatusCallback="/multiplex"/>
</Response>`

const echoMessage = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
	<Play>%s</Play>
</Response>`

const outboundMessage = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
	<Say voice="alice">This is an emergency message from Shakim. Please act upon after hearing it</Say>
	<Play>%s</Play>
</Response>`

//This is the main handler. Invoked when a)first inbound call to twilio b)after completion of the recording in the first inbound call
func indexHandler(w http.ResponseWriter, r *http.Request) {
	helper.LogRequestBody(r)

	w.Header().Set("Content-Type", "text/xml")

	if helper.RecordingUrlExist(r) {
		//TODO: Ideally should hang up the call
		fmt.Fprintf(w, echoMessage, helper.ExtractRecordingUrl(r))
	} else {
		fmt.Fprint(w, startupMessage)
	}
}

//This is callback method invoked when the recording completes
func multiplexHandler(w http.ResponseWriter, r *http.Request) {
	helper.LogRequestBody(r)

	if !helper.RecordingUrlExist(r) {
		return
	}

	helper.SetRecordingUrl(helper.ExtractRecordingUrl(r))

	time.Sleep(2000 * time.Millisecond)

	helper.CallFolks()
}

//This is the method invoked when an outbound call goes from twilio. It returns Twiml which the receiver hears
func outboundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprintf(w, outboundMessage, helper.GetRecordingUrl())
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/multiplex", multiplexHandler)
	http.HandleFunc("/outbound", outboundHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
