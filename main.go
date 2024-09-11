package main

import (
	"fmt"
	"gin-rush-template/tools"
	"github.com/jordan-wright/email"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"regexp"
	"sync"
)

// TODO: 如果要设置多个邮箱的话，可以考虑在 Bushi 里面加一个 SendNum 的字段，然后每次SendNum 到达一百的时候就进行下一个邮箱，同时把邮箱的设定放到一个数组里面，不同的邮箱对应不同的密码

const SMTPAddr = ""
const SMTPPassword = ""
const SMTPUsername = ""
const SMTPHost = ""
const GistID = ""

var bushi *Bushi

type Bushi struct {
	*email.Email
}

func main() {
	Init(
		"hdu <connect@hduhelp.com>",
		"Welcome to hduhelp!",
		"Text Body is, of course, supported!",
		"<h1>Today</h1>",
	)
}

func Init(from, subject, text, html string) {
	e := email.NewEmail()
	e.From = from
	e.Subject = subject
	e.Text = []byte(text)
	e.HTML = []byte(html)
	bushi = &Bushi{e}
}

func getBushi() *Bushi {
	return bushi
}

func (b *Bushi) getSubject() string {
	return b.Subject
}

func (b *Bushi) getText() string {
	return string(b.Text)
}

func (b *Bushi) getHTML() string {
	return string(b.HTML)
}

func (b *Bushi) getSender() string {
	return b.From
}

func (b *Bushi) getTo() {
	resp, err := http.Get(GistID)
	tools.PanicOnErr(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	tools.PanicOnErr(err)

	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	matches := re.FindAll(body, -1)

	var res []string
	res = make([]string, 0, len(matches))
	for i, m := range matches {
		fmt.Println("第", i, "位", string(m))
		res = append(res, string(m))
	}
	fmt.Println(len(res))
	b.To = res
}

func (b *Bushi) SendEmail() {
	var ch = make(chan int, 10) // 限制并发数
	var wg sync.WaitGroup

	for i, person := range b.To {
		wg.Add(1)
		ch <- i
		go func(per string) {
			defer wg.Done()
			bushi.To = []string{per} // 正式版
			err := bushi.Send(SMTPAddr, smtp.PlainAuth("hduhelp", SMTPUsername, SMTPPassword, SMTPHost))
			if err != nil {
				fmt.Println("send email to", person, "failed---", "err:", err)
			}
			<-ch
		}(person)
	}

	wg.Wait()
	fmt.Println("send email success")
}
