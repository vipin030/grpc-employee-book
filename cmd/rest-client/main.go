package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	// get configuration
	address := flag.String("server", "http://localhost:8080", "HTTP gateway url, e.g. http://localhost:8080")
	flag.Parse()
	var body string

	//Call CreateEmployee
	resp, err := http.Post(*address+"/v1/employee/create", "application/json", strings.NewReader(fmt.Sprintf(`
		{
				"name":"Veena",
				"city":"Stockholm",
				"salary":55000
		}
	`)))
	if err != nil {
		log.Fatalf("failed to call CreateEmployee method: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read CreateEmployee response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Create response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call ReadEmployee
	resp, err = http.Get(fmt.Sprintf("%s%s/%s", *address, "/v1/employee", "1"))
	fmt.Println(fmt.Sprintf("%s%s/%s", *address, "/v1/employee", "1"))
	if err != nil {
		log.Fatalf("failed to call Read method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read Read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Read response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call UpdateEmployee
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", *address, "/v1/employee", "1"),
		strings.NewReader(fmt.Sprintf(`
		{
			"name":"Fanny",
			"city":"Stockholm",
			"salary":50000
		}
	`)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call UpdateEmployee method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read UpdateEmployee response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Update response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", *address, "/v1/employee", "9"), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call DeleteEmployee method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read DeleteEmployee response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}
