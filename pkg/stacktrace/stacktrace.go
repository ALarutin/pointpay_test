package stacktrace

import (
	stError "github.com/go-errors/errors"
)

type Stacker interface {
	Stack() []byte
}

const Key = "error_stack_trace"

func Get(err error) string {
	switch e := err.(type) {
	case Stacker:
		return string(e.Stack())
	default:
		// К сожалению, в данной ситуации, найти место возникновения ошибки невозможно,
		// поэтому берётся за основу stacktrac-а место вызова функции Get
		return string(stError.Wrap(err, 1).Stack())
	}
}
