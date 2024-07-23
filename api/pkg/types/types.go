package types

type Response struct {
	ID        string
	Body      string
	ErrCode   int
	ErrString string
}

type Request struct {
	ID     string
	Body   string
	Method string
}
