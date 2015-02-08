package main

// A JSONResponse represents a JSON response.
type JSONResponse struct {
    meta Meta               `json:"meta"`
    data interface{}        `json:"data"`
}

// A Meta represents the meta field of a JSON response.
type Meta struct {
    Code int            `json:"code"`
    ErrorMessage string `json:"error_message"`  
}

// GenerateError creates a JSONResponse from information about a code.
func GenerateError(code int, errMsg string) (*JSONResponse) {
    return &JSONResponse{
        meta: Meta{
            Code: code,
            ErrorMessage: errMsg,
        },
    }
}

// GenerateError creates a JSONResponse from information about a code.
func GenerateSuccess(code int, successMsg string) (*JSONResponse) {
    return &JSONResponse{
        meta: Meta{
            Code: code,
        },
    }
}