package utils

import (
	"bytes"
	"fmt"
	"text/template"
)

type MailData struct {
	Name string
}

func SendEmail(tempPath string, mailData map[string]int) error {

	var body bytes.Buffer

	temp, err := template.ParseFiles(tempPath)

	if err != nil {
		return err
	}

	err = temp.Execute(&body, mailData)

	if err != nil {
		return err
	}

	fmt.Println(body)


	return nil

}



