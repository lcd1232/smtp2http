package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/alash3al/go-smtpsrv"
	"github.com/go-resty/resty"
	"github.com/pkg/errors"
	"github.com/zaccone/spf"
)

func handler(req *smtpsrv.Request) (err error) {
	defer func() {
		if err != nil {
			logger.WithError(err).Error("got error")
		}
	}()
	// validate the from data
	if *flagStrictValidation {
		if req.SPFResult != spf.Pass {
			err = errors.New("your host isn't configured correctly or you are a spammer -_-")
			return
		} else if !req.Mailable {
			err = errors.New("your mail isn't valid because it cannot receive emails -_-")
			return
		}
	}

	msg, err := parsemail.Parse(req.Message)
	if err != nil {
		err = errors.New("cannot read your message, it may be because of it exceeded the limits")
		return
	}

	logger.WithField("mail-from", req.From).
		WithField("mail-to", strings.Join(extractEmails(msg.To), ",")).
		WithField("mail-cc", strings.Join(extractEmails(msg.Cc), ",")).
		WithField("mail-bcc", strings.Join(extractEmails(msg.Bcc), ",")).
		Info("got message")

	rq := resty.R()

	// set the url-encoded-data
	rq.SetFormData(map[string]string{
		"id":              msg.Header.Get("Message-ID"),
		"subject":         msg.Subject,
		"body[text]":      string(msg.TextBody),
		"body[html]":      string(msg.HTMLBody),
		"addresses[from]": req.From,
		"addresses[to]":   strings.Join(extractEmails(msg.To), ","),
		"addresses[cc]":   strings.Join(extractEmails(msg.Cc), ","),
		"addresses[bcc]":  strings.Join(extractEmails(msg.Bcc), ","),
	})

	// set the files "attachments"
	for i, file := range msg.Attachments {
		is := strconv.Itoa(i)
		rq.SetFileReader("file["+is+"]", file.Filename, file.Data)
	}

	// submit the form
	resp, err := rq.Post(*flagWebhook)
	if err != nil {
		err = errors.New("cannot accept your message due to internal error, please report that to our engineers, '" + (err.Error()) + "'")
		return
	} else if resp.StatusCode() != http.StatusOK {
		err = errors.New("backend status code: " + resp.Status())
		return
	}

	return nil
}
