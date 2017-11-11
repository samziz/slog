package main


// CopyError

type AuthError struct {
	error
	Body string
}

func (a AuthError) Error() string {
	return "AuthError: " + a.Body
}


// ConnectionError 

type ConnectionError struct {
	error
	Body string
}

func (c ConnectionError) Error() string {
	return "ConnectionError: " + c.Body
}


// CopyError

type CopyError struct {
	error
	Body string
}

func (c CopyError) Error() string {
	return "CopyError: " + c.Body
}


// OSError 

type OSError struct {
	error
	Body string
}

func (o OSError) Error() string {
	return "OSError: " + o.Body
}