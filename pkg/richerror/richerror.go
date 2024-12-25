package richerror

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
	KindBadRequest
)

type Operation string

type RichError struct {
	operation    Operation
	wrappedError error
	message      string
	kind         Kind
	meta         map[string]any
}

func New(op Operation) RichError {
	return RichError{operation: op}
}

func (r RichError) WithErr(err error) RichError {
	r.wrappedError = err
	return r
}

func (r RichError) WithMessage(msg string) RichError {
	r.message = msg
	return r
}

func (r RichError) WithKind(kind Kind) RichError {
	r.kind = kind
	return r
}

func (r RichError) WithMeta(meta map[string]any) RichError {
	r.meta = meta
	return r
}

func (r RichError) Error() string {
	return r.message
}

func (r RichError) Kind() Kind {
	if r.kind != 0 {
		return r.kind
	}

	re, ok := r.wrappedError.(RichError)
	if !ok {
		return 0
	}

	return re.Kind()
}

func (r RichError) Message() string {
	if r.message != "" {
		return r.message
	}

	re, ok := r.wrappedError.(RichError)
	if !ok {
		return r.wrappedError.Error()
	}

	return re.Message()
}
