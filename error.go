package go2bosh

/* Descriptive error struct to inform the user for all the error issues */
type handlerError struct {
	Error   error
	Message string
	Code    int
}
