package errors

import "fmt"

type DuplicateError struct{
	Duplicate string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("Value %s already exists", e.Duplicate)
}

func NewDuplicateError(duplicate string) *DuplicateError {
	return &DuplicateError{
		Duplicate: duplicate,
	}
}

type WrongDataError struct{
	Data string
}

func (e *WrongDataError) Error() string {
	return fmt.Sprintf("Couldn't find or accept data %s", e.Data)
}

func NewWrongDataError(login string) *WrongDataError {
	return &WrongDataError{
		Data: login,
	}
}
