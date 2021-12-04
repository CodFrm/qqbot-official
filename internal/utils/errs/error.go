package errs

import "net/http"

type ReplyError struct {
	msg string
}

func NewReplyError(msg string) error {
	return &ReplyError{msg: msg}
}

func (r *ReplyError) Error() string {
	return r.msg
}

type JsonRespondError struct {
	Status int    `json:"-"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

func NewError(status, code int, msg string) error {
	return &JsonRespondError{
		Status: status,
		Code:   code,
		Msg:    msg,
	}
}

func (j *JsonRespondError) Error() string {
	return j.Msg
}

func NewBadRequestError(code int, err string) error {
	return &JsonRespondError{
		Status: http.StatusBadRequest,
		Code:   code,
		Msg:    err,
	}
}
