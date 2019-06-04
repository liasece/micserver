package mail

import (
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
	"sync"
	"time"
)

type MailNotifyBuffer struct {
	WContent       string
	wlastappendsec uint64
	wlatertime     uint64

	EContent       string
	elastappendsec uint64
	elatertime     uint64

	wmailbuffermutex sync.Mutex
	emailbuffermutex sync.Mutex
}

var mailnotifybuffer_s *MailNotifyBuffer

func init() {
	mailnotifybuffer_s = &MailNotifyBuffer{}
}

func GetMailNotifyBuffer() *MailNotifyBuffer {
	return mailnotifybuffer_s
}

func (this *MailNotifyBuffer) StartMailBuffer() {
	go this.syncSendMail()
}

func (this *MailNotifyBuffer) AppendWarning(content string, latersec uint64) {
	this.wmailbuffermutex.Lock()
	defer this.wmailbuffermutex.Unlock()

	this.WContent += content
	this.wlastappendsec = uint64(time.Now().Unix())
	this.wlatertime = latersec
}

func (this *MailNotifyBuffer) AppendError(content string, latersec uint64) {
	this.emailbuffermutex.Lock()
	defer this.emailbuffermutex.Unlock()

	this.EContent += content
	this.elastappendsec = uint64(time.Now().Unix())
	this.elatertime = latersec
}

func (this *MailNotifyBuffer) syncSendMail() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[syncSendMail] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	for {
		this.wmailbuffermutex.Lock()
		if this.WContent != "" &&
			this.wlastappendsec+this.wlatertime < uint64(time.Now().Unix()) {
			GetMailManager().SendMailServerWarning(this.WContent)
			log.Debug("[SendMail] [%s]", this.WContent)
			this.WContent = ""
		}
		this.wmailbuffermutex.Unlock()

		this.emailbuffermutex.Lock()
		if this.EContent != "" &&
			this.elastappendsec+this.elatertime < uint64(time.Now().Unix()) {
			GetMailManager().SendMailServerError(this.EContent)
			log.Debug("[SendMail] [%s]", this.EContent)
			this.EContent = ""
		}
		this.emailbuffermutex.Unlock()

		time.Sleep(300 * time.Millisecond)
	}
}
