////////////////////////////////////
// RPC Proxy implementation in Go //
////////////////////////////////////
// Test cURL call:
// curl -H "Content-Type: application/json" -X POST -d '{"userId":"1", "options":[1,2,3,4]}' http://localhost:9898/getUserData

package main

import (
	"encoding/json"
	"github.com/7Linternational/libs/mq"
	"log"
	"net/http"
)

/**
 * [checkErr description]
 * @param  {[type]} err error         [description]
 * @return {[type]}     [description]
 */
func checkErr(err error) {
	if err != nil {
		// log.Fatal(err)
		log.Println(err)
	}
}

/**
 * [getData description]
 * @param  {[type]} rw  http.ResponseWriter [description]
 * @param  {[type]} req *http.Request       [description]
 * @return {[type]}     [description]
 */
func getUserData(rw http.ResponseWriter, req *http.Request) {
	///////////////////////////////////
	// Declare method request struct //
	///////////////////////////////////
	type request struct {
		UserId  int   `json:"userId"`
		Options []int `json:"options"`
	}
	////////////////////////////
	// Get request parameters //
	////////////////////////////
	decoder := json.NewDecoder(req.Body)
	var r request
	err := decoder.Decode(&r)
	checkErr(err)
	log.Println(r)

	///////////////////////////
	// Push request to Queue //
	///////////////////////////
	body, err := json.Marshal(r)
	checkErr(err)
	res := mq.Consume(string(body))

	////////////////////////
	// Send Response back //
	////////////////////////
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(res))
}

/**
 * [main description]
 * @return {[type]} [description]
 */
func main() {
	http.HandleFunc("/getUserData", getUserData)
	http.ListenAndServe(":9898", nil)
}
