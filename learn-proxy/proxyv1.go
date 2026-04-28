package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func main(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		// read the body 
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil{
			http.Error(w, "can't read Body", http.StatusInternalServerError)
			return
		}
		// close the body when done 
		defer r.Body.Close()
		log.Println(r.Method, r.URL.Path, "Body:", string(bodyBytes))


		// create a new request for the destination server
		proxyUrl := "https://openai.com" + r.URL.Path
		newReq, err := http.NewRequest(r.Method, proxyUrl, bytes.NewReader(bodyBytes))
		if err != nil {
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}
		

		//copy headers
		newReq.Header = r.Header.Clone()
		newReq.Host = "localhost:8080"

		// read the response4
		client := &http.Client{}
		resp, err := client.Do(newReq)
		if err != nil{
			http.Error(w, "Failed to reach destination", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}


