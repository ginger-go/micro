package micro

var errMap = make(map[interface{}]string)

func RegisterError(uuid string, message string) {
	errMap[uuid] = message
}

type Error interface {
	Code() string
	Error() string
}

func NewError(code string) Error {
	return &errorImp{
		code:    code,
		message: errMap[code],
	}
}

type errorImp struct {
	code    string
	message string
}

func (e *errorImp) Code() string {
	return e.code
}

func (e *errorImp) Error() string {
	return e.message
}
