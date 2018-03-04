package mvc
type BingoError interface {
	error
	Code() int
}

