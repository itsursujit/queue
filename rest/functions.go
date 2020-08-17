package rest

import (
	"fmt"
	"queue/workers"
)

var funcMapper = map[string]workers.JobFunc{
	"SendEmail": SendEmail,
	"SendSMS":   SendSMS,
}

func DoWork(message *workers.Msg) error {
	mp, _ := message.Map()
	function := mp["class"].(string)
	err := funcMapper[function](message)
	if err != nil {
		panic(err)
	}
	return nil
}
func SendEmail(message *workers.Msg) error {
	fmt.Println("I'm sending Email")
	return nil
}

func SendSMS(message *workers.Msg) error {
	fmt.Println("I'm sms sending")
	return nil
}
