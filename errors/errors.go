package errors

type ErrorType int

const (
	UnknownCommand            ErrorType = iota
	UnavailableCommand        ErrorType = iota
	WrongState                ErrorType = iota
	MultipleStatesFromCommand ErrorType = iota
	SendMessageError          ErrorType = iota
)

type StateMachineError struct {
	Code    ErrorType
	Message string
}

func (err StateMachineError) Error() string {
	return err.Message
}

type CommandHandlerError struct {
	Code    ErrorType
	Message string
}

func (err CommandHandlerError) Error() string {
	return err.Message
}
