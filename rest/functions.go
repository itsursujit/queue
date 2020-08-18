package rest

import (
	"fmt"
	"queue/workers"
)

func init() {
	workers.FuncMapper.Map = map[string]workers.JobFunc{
		"SendEmail": SendEmail,
		"SendSMS":   SendSMS,
	}
}

func SendEmail(message *workers.Msg) error {
	fmt.Println("I'm sending Email")
	return nil
}

func SendSMS(message *workers.Msg) error {
	fmt.Println("I'm sms sending")
	return nil
}
