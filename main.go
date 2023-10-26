package main

import (
	"fmt"

	"github.com/jakepeery/O-Auth-Testing/wso2services/wso2requests"
	//"github.com/byuoitav/wso2services/wso2requests"
)

func main() {
	var punchResponse string

	body := ""
	//byuID := "779147452"
	url := "https://api-sandbox.byu.edu/bdp/human_resources/worker_summary/v0" //?worker_id=" + byuID

	method := "GET"

	fmt.Println(method, url, body, &punchResponse)
	err, response, _ := wso2requests.MakeWSO2RequestWithHeadersReturnResponse(method, url, body, &punchResponse, map[string]string{
		"Host": "api-sandbox.byu.edu",
	})
	if err != nil {
		fmt.Println(err, punchResponse, response)
	}

	fmt.Println(punchResponse, response)
}
