package main

import (
	"fmt"
	"os"

	"github.com/byuoitav/wso2services/wso2requests"
)

func main() {
	var punchResponse string

	key := os.Getenv("CLIENT_KEY")
	secret := os.Getenv("CLIENT_SECRET")

	fmt.Println("Timeclock Key", key)
	fmt.Println("Timeclock Secret", secret)
	body := ""
	byuID := "779147452"

	method := "POST"
	err, response, _ := wso2requests.MakeWSO2RequestWithHeadersReturnResponse(method, "https://api-sandbox.byu.edu/bdp/human_resources/worker_summary/v0?worker_id="+byuID, body, &punchResponse, map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})
	if err != nil {
		fmt.Println(punchResponse, response)
	}

	fmt.Println(punchResponse, response)
}
