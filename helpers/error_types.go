package helpers

type NotFoundError struct { //404
	Message    string
	MessageDev string
}

type BadRequestError struct { //400
	Message    string
	MessageDev string
}

type InternalServerError struct { //500
	Message    string
	MessageDev string
}

type UnauthorizedError struct { //401
	Message    string
	MessageDev string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}
