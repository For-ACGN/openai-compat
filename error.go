package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIError represents an error returned by the API.
type APIError struct {
	StatusCode int    // HTTP status code
	APICode    string // Business error code from API response
	Message    string // Human-readable error message
}

func (e *APIError) Error() string {
	if e.StatusCode != 0 {
		return e.Message
	}
	return fmt.Sprintf("code: %s, %s", e.APICode, e.Message)
}

type rawError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
	} `json:"error"`
}

func handleError(resp *http.Response) error {
	data, _ := io.ReadAll(resp.Body)
	if strings.HasPrefix(string(data), "<html>") {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "unexpected HTML response",
		}
	}
	// return the standard api error
	re := new(rawError)
	err := json.Unmarshal(data, re)
	if err == nil {
		return &APIError{
			APICode: re.Error.Code,
			Message: re.Error.Message,
		}
	}
	// fallback to simple error
	apiErr := &APIError{
		StatusCode: 1,
	}
	switch resp.StatusCode {
	case http.StatusBadRequest:
		apiErr.Message = "Bad request"
	case http.StatusUnauthorized:
		apiErr.Message = "Invalid authentication credentials"
	case http.StatusPaymentRequired:
		apiErr.Message = "Insufficient account balance"
	case http.StatusTooManyRequests:
		apiErr.Message = "Rate limit exceeded"
	case http.StatusNotFound:
		apiErr.Message = "Requested resource not found"
	case http.StatusInternalServerError:
		apiErr.Message = "Internal server error"
	default:
		apiErr.Message = fmt.Sprintf("Unexpected API response (HTTP %d)", resp.StatusCode)
	}
	return apiErr
}
