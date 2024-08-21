package apierrors

import (
	"errors"
	"net/http"
	"testing"
)

type apiErrorTestCase struct {
	name     string
	function func() ApiError
	expected apiErr
}

func TestApiErrorFunctions(t *testing.T) {
	// Define the test cases
	testCases := []apiErrorTestCase{
		{
			name: "NewNotFoundApiError",
			function: func() ApiError {
				return NewNotFoundApiError("Not found")
			},
			expected: apiErr{
				ErrorMessage: "Not found",
				ErrorCode:    "not_found",
				ErrorStatus:  http.StatusNotFound,
				ErrorCause:   CauseList{},
			},
		},
		{
			name: "NewBadRequestApiError",
			function: func() ApiError {
				return NewBadRequestApiError("Bad request")
			},
			expected: apiErr{
				ErrorMessage: "Bad request",
				ErrorCode:    "bad_request",
				ErrorStatus:  http.StatusBadRequest,
				ErrorCause:   CauseList{},
			},
		},
		{
			name: "NewInternalServerApiError",
			function: func() ApiError {
				return NewInternalServerApiError("Internal server error", nil)
			},
			expected: apiErr{
				ErrorMessage: "Internal server error",
				ErrorCode:    "internal_server_error",
				ErrorStatus:  http.StatusInternalServerError,
				ErrorCause:   CauseList{},
			},
		},
		{
			name: "NewUnauthorizedApiError",
			function: func() ApiError {
				return NewUnauthorizedApiError()
			},
			expected: apiErr{
				ErrorMessage: "You're not authorized to use this resource",
				ErrorCode:    "unauthorized",
				ErrorStatus:  http.StatusUnauthorized,
				ErrorCause:   CauseList{},
			},
		},
	}

	// Execute the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.function()
			if errors.Is(result, tc.expected) {
				t.Errorf("Expected %+v, but got %+v", tc.expected, result)
			}
		})
	}
}
