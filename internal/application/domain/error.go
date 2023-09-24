package domain

import "fmt"

type ErrorCategory int

const (
	InvalidInput ErrorCategory = iota
	Unauthorized
	InternalError
)

type Error interface {
	error
	Category() ErrorCategory
	Message() string
	Details() string
}

type err struct {
	category ErrorCategory
	message  string
	details  string
}

func (e err) Category() ErrorCategory {
	return e.category
}

func (e err) Error() string {
	return fmt.Sprintf("error: %v", e.message)
}

func (e err) Message() string {
	return e.message
}

func (e err) Details() string {
	return e.details
}

var (
	ErrInvalidProductName Error = err{
		category: InvalidInput,
		message:  "Invalid Name",
		details:  `Product name cannot be an empty string.`,
	}
	ErrInvalidProductPrice Error = err{
		category: InvalidInput,
		message:  "Invalid Price",
		details:  `Product price cannot be 0.`,
	}
)
