package main

// A JSONResponse represents a JSON response.
type JSONResponse struct {
    Meta M                  `json:"meta"`
    Data interface{}        `json:"data,omitempty"`
}

// A M represents the meta field of a JSON response.
type M struct {
    Code int            `json:"code"`
    ErrorMessage string `json:"error_message,omitempty"`  
}

// GenerateError creates an error JSONResponse.
func GenerateError(code int, errMsg string) (*JSONResponse) {
    return &JSONResponse{
        Meta: M {
            Code: code,
            ErrorMessage: errMsg,
        },
    }
}

// GenerateSuccess creates a success JSONResponse.
func GenerateSuccess(code int, successMsg string) (*JSONResponse) {
    return &JSONResponse{
        Meta: M {
            Code: code,
        },
    }
}