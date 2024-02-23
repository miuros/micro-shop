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
<p>{{code}}</p>
</body>
</html>

`

var pattern, _ = regexp.Compile("{{code}}")

func SendMail(tomail, content string, from string, host, passwd string) error {
	msg := pattern.ReplaceAllString(htmlMsg, content)
	body := gomail.NewMessage()
	body.SetHeader("From", from)
	body.SetHeader("To", tomail)
	body.SetHeader("Subject", "verify code")
	body.SetBody("text/html", msg)
	dialer := gomail.NewDialer(host, 465, from, passwd)
	return dialer.DialAndSend(body)
}
