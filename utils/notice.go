package utils

import "gopkg.in/gomail.v2"

func SendInformation(toEmail string) string {
	//V8Example()
	m := gomail.NewMessage()
	m.SetHeader("From", "204292960@qq.com")
	//m.SetHeader("To", "javaandroidxml@163.com")
	m.SetHeader("To", toEmail)
	//m.SetAddressHeader("Cc", "204292960@qq.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	//m.Attach("git.pdf")	// 发送附件

	d := gomail.NewDialer("smtp.qq.com", 25, "204292960@qq.com", "opakqsscbwbpbjfc")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
		return ""
	}
	return "邮件发送成功"
}
