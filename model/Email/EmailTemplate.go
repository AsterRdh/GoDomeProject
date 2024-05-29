package Email

import "time"

var Template EMailTemplate = EMailTemplate{
	ResCheckEmail: "",
}

type EMailTemplate struct {
	ResCheckEmail string
}

type EMailKey struct {
	Email string
	Ts    time.Time
}
