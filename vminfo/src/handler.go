package main

import (
	"github.com/cloudevents/sdk-go"
	"fmt"
	"context"
	"time"
	"net/http"
	"bytes"
	"encoding/json"
	"errors"
	"log"
)

// Sample cloud event Meta data and body
type YugenEvent struct { 
    CloudEventsVersion string    `json:"cloudEventsVersion"`
	EventType          string    `json:"eventType"`
	EventTypeVersion   string    `json:"eventTypeVersion"`
	Source             string    `json:"source"`
	EventID            string    `json:"eventID"`
	EventTime          time.Time `json:"eventTime"`
	EventData    	   Data `json:"data"`	 
}


type Data struct {
	SlackHook          string `json:"slackHook"`
	VCenterURL         string `json:"vcenterURL"`
	SessionID          string `json:"sessionId"`
}

 type Context struct {
	
   }


// Replace the function body with the implementation
func Event(ctx context.Context, event cloudevents.Event) error {
	data := ""
	if err := event.DataAs(&data); err != nil {
		fmt.Printf("Got Data Error: %s\n", err.Error())
		return err
	}

	newData := YugenEvent{}
	log.Println("Got data ++ ", data)
	err := json.Unmarshal([]byte(data), &newData)
	if err != nil {
		fmt.Println("error:", err)
	}


	fmt.Printf("Got Data: %+v\n", newData)
	fmt.Printf("Got Transport Context: %+v\n", cloudevents.HTTPTransportContextFrom(ctx))

	vmInfo, _ := GetVMinfo(newData.EventData.VCenterURL, newData.EventData.SessionID)

	SendSlackNotification(newData.EventData.SlackHook, vmInfo)

	return nil
}


func GetVMinfo(url, sessionId string) (string,error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Recieved error", err.Error())
		return "",err
	}
	req.Header.Set("Vmware-Api-Session-Id", sessionId)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Received error", err.Error())
		return "",err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    vmInfo := ""
	if resp !=nil {
		vmInfo = buf.String()
	}
	return vmInfo, nil
}

func SendSlackNotification(webhookUrl string, msg string) error {

    slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
    req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    if buf.String() != "ok" {
        return errors.New("Non-ok response returned from Slack")
    }
    return nil
}


type SlackRequestBody struct {
    Text string `json:"text"`
}