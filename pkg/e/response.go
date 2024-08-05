package e

import (
	"net/http"
	"time"
)

// 運行成功
func StatusSuccess(message string, data interface{}) (int, *Response) {
	return http.StatusOK, newResponse(http.StatusOK, SUCCESS, message, time.Now(), data)
}

// 建立成功
func StatusCreated(message string, data interface{}) (int, *Response) {
	return http.StatusCreated, newResponse(http.StatusCreated, CREATED, message, time.Now(), data)
}

func StatusNoContent(message string) (int, *Response) {
	return http.StatusNoContent, newResponse(http.StatusNoContent, NO_CONTENT, message, time.Now(), nil)
}

func StatusAccept(message string, data interface{}) (int, *Response) {
	return http.StatusAccepted, newResponse(http.StatusAccepted, ACCEPT, message, time.Now(), data)
}
