package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"strconv"
)

func commandHandler(cmd string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving command")
		w.Header().Set("Content-Type", "application/x-64")
		w.Write([]byte(cmd))
	}
}

func responseHandler(fin chan<- bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			log.Println("Response received")
			scanner := bufio.NewScanner(base64.NewDecoder(base64.StdEncoding, r.Body))
			for scanner.Scan() {
				log.Println("> ", scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				log.Println("Error reading response: ", err)
			}
			log.Println("Response finished")
			fin <- true
		}
	}
}

func main() {
	command := flag.String("command", "hostname", "a command to execute on remote service")
	port := flag.Uint("port", 8080, "port to listen on")
	flag.Parse()

	log.Println("Running command server on port", *port, "with the command:", *command)

	finished := make(chan bool)

	http.HandleFunc("/command", commandHandler(base64.StdEncoding.EncodeToString([]byte(*command))))
	http.HandleFunc("/response", responseHandler(finished))

	go func() {
		log.Fatal(http.ListenAndServe(":"+strconv.FormatUint(uint64(*port), 10), nil))
	}()

	log.Println("Waiting")
	<-finished
	log.Println("Finished")
}
