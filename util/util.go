package util

import (
	"fmt"
	"net/http"
	"strconv"
)

// helper function to check status of given api
func CheckAPIStatus(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer resp.Body.Close()
	return strconv.Itoa(resp.StatusCode)
}
