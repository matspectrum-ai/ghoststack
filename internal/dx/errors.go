package dx

import (
	"errors"
	"fmt"
)

type UserError struct {
	Code    string
	Message string
	Hint    string
}

func (e *UserError) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("%s: %s (dica: %s)", e.Code, e.Message, e.Hint)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func Wrap(err error, code, message string) *UserError {
	if err == nil {
		return &UserError{Code: code, Message: message}
	}
	return &UserError{Code: code, Message: message, Hint: err.Error()}
}

func CommandNotFound(name string) error {
	return &UserError{
		Code:    "COMMAND_NOT_FOUND",
		Message: "comando nao encontrado",
		Hint:    "execute 'ghost --help' para ver comandos disponiveis",
	}
}

func IsUserError(err error) bool {
	var ue *UserError
	return errors.As(err, &ue)
}
