package main

import (
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"os/exec"
)

// Mailer struct
type Mailer struct {
	gmail, pass, to string
}

// BuildAndSendEmail send alert email
func (m *Mailer) BuildAndSendEmail(unreachable []error) {
	if m.to != "" {
		var (
			count   = len(unreachable)
			subject string
			body    string
		)
		if count == 1 {
			urlErr, ok := unreachable[0].(*url.Error)
			if ok {
				subject = fmt.Sprintf("'%s' is unreachable", urlErr.URL)
			} else {
				subject = "A site is unreachable"
			}
		} else {
			subject = fmt.Sprintf("%d sites are unreachable", count)
		}
		for _, error := range unreachable {
			body += fmt.Sprintf("* %s\n\n", error.Error())
		}
		err := m.SendMail(subject, body)
		if err != nil {
			log.Println(err)
		}
	}
}

// SendMail send a email
func (m *Mailer) SendMail(subject, body string) error {
	if m.gmail == "" {
		return m.sendEmailWithMail(subject, body)
	}
	return m.sendEmailWithGmail(subject, body)
}

func (m *Mailer) sendEmailWithGmail(subject, body string) error {
	msg := "From: Bulldog <bulldog>\n" +
		"To: " + m.to + "\n" +
		"Subject: [Bulldog] " + subject + "\n\n" +
		body
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", m.gmail, m.pass, "smtp.gmail.com"),
		m.gmail, []string{m.to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (m *Mailer) sendEmailWithMail(subject, body string) error {
	c1 := exec.Command("echo", body)
	c2 := exec.Command("mail", "-s [Bulldog] "+subject, "-r Bulldog <bulldog>", m.to)
	c2.Stdin, _ = c1.StdoutPipe()
	err1 := c1.Start()
	err2 := c2.Run()
	err3 := c1.Wait()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}
	return nil
}

// func errorToEmoji(err error) string {
// 	if strings.Contains(err.Error(), "Client.Timeout exceeded") {
// 		return "⌛️"
// 	}
// 	if strings.Contains(err.Error(), "unsupported protocol scheme") {
// 		return "😓"
// 	}
// 	if strings.Contains(err.Error(), "no such host") {
// 		return "🛑"
// 	}
// 	if strings.Contains(err.Error(), "status code is 403") {
// 		return "🔒"
// 	}
// 	if strings.Contains(err.Error(), "status code is 404") {
// 		return "🔎"
// 	}
// 	if strings.Contains(err.Error(), "status code is 500") {
// 		return "🆘"
// 	}
// 	return "⚠️"
// }
