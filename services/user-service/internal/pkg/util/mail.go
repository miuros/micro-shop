package util

import (
	gomail "gopkg.in/gomail.v2"
	"regexp"
)

var htmlMsg = `

<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8" http-equiv="Content-Type" content="text/html">
    <style>
    </style>
</head>
<body>
<p>dear {{username}},your order is closing to close</p>
</body>
</html>

`

var pattern, _ = regexp.Compile("{{username}}")

func SendMail(tomail, name string, from string, host, passwd string) error {
	msg := pattern.ReplaceAllString(htmlMsg, name)
	body := gomail.NewMessage()
	body.SetHeader("From", from)
	body.SetHeader("To", tomail)
	body.SetHeader("Subject", "order is outdating")
	body.SetBody("text/html", msg)
	dialer := gomail.NewDialer(host, 465, from, passwd)
	return dialer.DialAndSend(body)
}
