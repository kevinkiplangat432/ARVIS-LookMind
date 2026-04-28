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

		// create a new request for the destination server
		proxyUrl := "https://api.openai.com" + r.URL.Path
		newReq, err := http.NewRequest(r.method, proxyUrl, bytes.NewReader(bodyBytes))
		if err != nill {
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}
		//copy headers
		newReq.Header = r.Header.clone()

		client := &http.Client{}
		resp, err := client.Do(newReq)
		if err != nil{
			http.Error(w, "Failed to reach destination", http.StatusBadGateway)
			return
		}


		// close the body when done 
		defer r.Body.Close()
		// use the data from the body
		log.Println(r.Method, r.URL.Path, "Body:", string(bodyBytes))

		// restore the Body  put a fresh one in there 
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	

		switch r.Method{
		case http.MethodGet:
			w.Write([]byte("Received!"))
		case http.MethodPost:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not supported\n"))
		}
		
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}


