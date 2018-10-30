package errors

import (
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	Handle(error)
}

type DefaultHandler struct {
	Logger      *log.Logger
}

func (h DefaultHandler) Handle(err error) {
	cerr, ok := err.(Error)
	if !ok {
		h.Logger.Errorf("unknown error: %s", cerr)
		return
	}
	h.Logger.Errorf("%s error: %s", cerr.Code, cerr.Msg)
}
