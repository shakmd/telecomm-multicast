# telecomm-multicast

### What does this app do?
The main purpose of this application is to immediately multi-cast a single message(via a phone call, sms etc) across different channels(or people) via phone call, sms etc

### How does this app works?
This application highly relies on [Twilio, a cloud-communication-platform as a service firm](https://www.twilio.com/) </br>

We can visualise this app as made up of two components: </br>

#### Twilio setup
A toll-free number is set up with Twilio with certain configuration like which web end-point to call when an incoming call is received </br>
Whenever a call is made to this toll-free number, it invokes a POST request to configured end point. The response body(TwiML, xml for Twilio) from the POST request dictates how to respond to the incoming call</br>
 
Currently, the response body dictates following:
- Ask the caller to record their message
- Record the message for 10 seconds
- Disconnect the call
- Perform a POST request to a web end-point when the recording completes

When Twilio invokes a POST request after the completion of the recording, the end point multi-cast the recorded message(via phone call) to all the contacts that user has specified in the configuration file(or env variable). The phone number used to make the outbound call is the same toll-free number that we setup with Twilio  

#### Web Serive setup
Three enpoints are exposed by this web service: </br>
- "/"           -> Returns Twiml. Invoked by Twilio whenver an incoming call is received
- "/multiplex"  -> Returns Twiml. Invoked by Twilio when the incoming caller has finished recording the message
- "/outbound"   -> Returns Twiml. Invoked by Twilio when it needs to make an outgoing call


### How this app can be used in real world
This application can be used in number of real world problems, some of it are:
- To send a message to multiple people in an emergency or time-constrained situation
- To contact a number of people when out of call/sms/data balance 

