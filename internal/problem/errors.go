package problem

import "errors"

// Error definitions for the problem package
var (
	// ErrProblemNotFound is returned when a problem ID is not found
	ErrProblemNotFound = errors.New("problem not found")
	
	// ErrInvalidProblemData is returned when problem data is invalid
	ErrInvalidProblemData = errors.New("invalid problem data")
)