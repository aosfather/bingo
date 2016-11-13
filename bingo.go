// bingo project bingo.go
package bingo

/**
200 OK 201 Created 204 No Content
304 Not Modified
400 Bad Request 401 Unauthorized 403 Forbidden 404 Not Found 409 Conflict
500 Internal Server Error


*/
const (
	Method_GET          = "GET"
	Method_POST         = "POST"
	Method_PUT          = "PUT"
	Method_DELETE       = "DELETE"
	Method_PATCH        = "PATCH"
	Code_OK             = 200
	Code_CREATED        = 201
	Code_EMPTY          = 204
	Code_NOT_MODIFIED   = 304
	Code_BAD            = 400
	Code_UNAUTHORIZED   = 401
	Code_FORBIDDEN      = 403
	Code_NOT_FOUND      = 404
	Code_CONFLICT       = 409
	Code_ERROR          = 500
	Code_NOT_IMPLEMENTS = 501
	Code_NOT_ALLOWED    = 405
)

type BingoError interface {
	error
	Code() int
}

type MethodError struct {
	code int
	msg  string
}

func (this MethodError) Code() int {
	return this.code
}

func (this MethodError) Error() string {
	return this.msg
}

func CreateError(c int, text string) MethodError {
	var err MethodError
	err.code = c
	err.msg = text
	return err
}
