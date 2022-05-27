package bttp

import (
	"encoding/json"
	"net/http"
)

// Response holds all information that is needed to translate into a http response.
type Response struct {
	// The http status code that is used in response.
	StatusCode int
	// The headers that are used in the response.
	Headers map[string]string
	// The data that should be responded with.
	Data any
}

// EmptyData is a placeholder that can be used inside the convenience function when
// no data should be responded.
var EmptyData any

// HandlerFunc is the custom function that can be translated into a http.HandlerFunc.
type HandlerFunc func(r *http.Request) Response

// Handle translates a HandlerFunc into a http.HandleFunc.
func Handle(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := f(r)
		if err := writeResponse(resp, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// DecodeBody takes a http.Request and parses it into a struct v.
// It returns a bool, which indicates if the parsing was successful and a Response.
// In case of a decoding error, the Response indicates a http.StatusBadRequest.
func DecodeBody(r *http.Request, v any) (bool, *Response) {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return false, &Response{
			StatusCode: http.StatusBadRequest,
			Headers:    nil,
			Data:       nil,
		}
	}
	return true, nil
}

// Ok is a convenience function to generate a Response with a http.StatusOK.
func Ok(data any) Response {
	return Response{
		StatusCode: http.StatusOK,
		Headers:    nil,
		Data:       data,
	}
}

// Created is a convenience function to generate a Response with a
// http.StatusCreated.
func Created(location string) Response {
	return Response{
		StatusCode: http.StatusCreated,
		Headers: map[string]string{
			"Location": location,
		},
		Data: nil,
	}
}

// BadRequest is a convenience function to generate a Response with a
// http.StatusBadRequest.
func BadRequest(data any) Response {
	return Response{
		StatusCode: http.StatusBadRequest,
		Headers:    nil,
		Data:       data,
	}
}

// InternalServerError is a convenience function to generate a Response with a
// http.StatusInternalServerError.
func InternalServerError(data any) Response {
	return Response{
		StatusCode: http.StatusInternalServerError,
		Headers:    nil,
		Data:       data,
	}
}

func writeResponse(r Response, w http.ResponseWriter) error {
	for k, v := range r.Headers {
		w.Header().Set(k, v)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(r.StatusCode)
	if r.Data != EmptyData {
		bytes, err := json.Marshal(r.Data)
		if err != nil {
			return err
		}
		_, err = w.Write(bytes)
		return err
	}
	return nil
}
