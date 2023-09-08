package infa_auth

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

func testAuth() {
	url := "https://pod.ics.dev:444/session-service/api/v1/session/Agent"
	sessionToken := "64Vjmeewe81iwbIPfgUmqu"
	validateToken(url, sessionToken)
}

func validateToken2(url string, sessionToken string) bool {
	// Create an HTTP client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Create a client with the custom transport
	client := &http.Client{Transport: tr}

	//client := &http.Client{}

	// Create a GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}
	req.Header.Add("IDS-AGENT-SESSION-ID", sessionToken)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return false
	}
	if resp.StatusCode != 200 {
		fmt.Println("http status is not 200:")
		return false
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}

	// Print the response
	fmt.Println(string(body))
	return true

}
