package mail

import (
	"base"
	"base/logger"
	// "errors"
	"base/util"
	"fmt"
	"gopkg.in/gomail.v2"
	"net/smtp"
	// "os"
	"strings"
)

type loginAuth struct {
	username, password string
}

// loginAuth returns an Auth that implements the LOGIN authentication
// mechanism as defined in RFC 4616.
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	command := string(fromServer)
	command = strings.TrimSpace(command)
	command = strings.TrimSuffix(command, ":")
	command = strings.ToLower(command)

	if more {
		if command == "username" {
			return []byte(a.username), nil
		} else if command == "password" {
			return []byte(a.password), nil
		} else {
			// We've already sent everything.
			return nil, fmt.Errorf("unexpected server challenge: %s", command)
		}
	}
	return nil, nil
}

type MailContent struct {
	ToAddr  []string
	Title   string
	Content string

	mailmanager *MailManager
}

type MailManager struct {
	SenderAddr        string
	SenderKeyCode     string
	SenderNickName    string
	ServerAddr        string
	ServerPort        int
	GlobalToAddr      []string
	DefultContentType string
}

var mailmanager_s *MailManager

func init() {
	mailmanager_s = &MailManager{}
	mailmanager_s.GlobalToAddr = make([]string, 0)
	mailmanager_s.DefultContentType = "Content-Type: text/plain; charset=UTF-8"
}

func GetMailManager() *MailManager {
	return mailmanager_s
}

func (this *MailManager) InitMailManagerByConfig() {
	this.SenderAddr = base.GetGBServerConfigM().GetProp("mailsenderaddr")
	this.SenderKeyCode = base.GetGBServerConfigM().GetProp("mainsenderkeycode")
	this.SenderNickName = base.GetGBServerConfigM().GetProp("mailsendernickname")
	this.ServerAddr = base.GetGBServerConfigM().GetProp("mailserveraddr")
	this.ServerPort = int(base.GetGBServerConfigM().GetPropInt("mailserverport"))
	times := 1
	for {
		tmpstr := fmt.Sprintf("mailglobaltoaddr%d", times)
		times++
		toaddr := base.GetGBServerConfigM().GetProp(tmpstr)
		if toaddr == "" {
			break
		}
		this.GlobalToAddr = append(this.GlobalToAddr, toaddr)
	}
}

func (this *MailManager) SendMailServerError(content string) {
	mail := MailContent{}
	mail.InitMail(this.GlobalToAddr, this.SenderNickName+"-Server Error", content, this)
	mail.SendMail()
}

func (this *MailManager) SendMailServerWarning(content string) {
	mail := MailContent{}
	mail.InitMail(this.GlobalToAddr, this.SenderNickName+"-Server Warning", content, this)
	go mail.SendMail()
}

func (this *MailContent) InitMail(toaddr []string, title string,
	content string, mailmanager *MailManager) {
	this.ToAddr = toaddr
	this.Title = title
	this.Content = content
	this.mailmanager = mailmanager
}

func (this *MailContent) SendMail() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			logger.Error("[SendMail] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	sender := this.mailmanager.SenderAddr
	serveraddr := this.mailmanager.ServerAddr
	serverport := this.mailmanager.ServerPort
	keycode := this.mailmanager.SenderKeyCode
	// contenttype := this.mailmanager.DefultContentType
	// auth := smtp.PlainAuth("", sender, keycode, serveraddr)

	// msg := []byte("To: " + strings.Join(this.ToAddr, ",") + "\r\n" +
	// 	"From: " + this.mailmanager.SenderNickName + "<" + sender + ">\r\n" +
	// 	"Subject: " + this.Title + "\r\n" +
	// 	contenttype + "\r\n\r\n" +
	// 	this.Content)
	// url := fmt.Sprintf("%s:%d", serveraddr, serverport)
	// err := smtp.SendMail(url, auth, sender, this.ToAddr, msg)
	// if err != nil {
	// 	logger.Error("send mail error: %v", err)
	// }

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", this.ToAddr...)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", this.Title)
	m.SetBody("text/html", this.Content)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(serveraddr, serverport, sender, keycode)
	// d.SSL = true
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		logger.Error("send mail error: %v", err)
	}
}
