package sendto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type MailRequest struct {
	ToEmail     string `json:"toEmail"`
	MessageBody string `json:"messageBody"`
	Subject     string `json:"subject"`
	Attachment  string `json:"attachment"` // Optional attachment field
}

func SendEmailToApi(otp string, email string, purpose string) error {
	// URL api
	postUrl := "http://localhost:8080/api/v1/send-email"
	mailRequest := MailRequest{
		ToEmail:     email,
		MessageBody: "OTP is::" + otp,
		Subject:     "OTP for " + purpose,
		Attachment:  "path/to/email", // Optional, can be set if needed
	}

	// convert the request to JSON

	requestBody, err := json.Marshal(mailRequest)
	if err != nil {
		return err
	}

	// send the request
	req, err := http.NewRequest(postUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	// put header
	req.Header.Set("Content-Type", "application/json")
	// excute the request
	client := &http.Client{}
	resq, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resq.Body.Close()
	fmt.Sprintln("Response Status:", resq.Status)
	return nil
}
