////////////////////////////////////
// RPC Proxy implementation in Go //
////////////////////////////////////
// Test cURL call: curl -H "Content-Type: application/json" -X POST -d '{"method":"methodName","parameters":{"param1":"value1"}}' http://localhost:9898/rpc

package main

import (
	"encoding/json"
	// "errors"
	// "github.com/7Linternational/libs/producer"
	// "html"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"net/http"
	// "reflect"
)

type api_request struct {
	Parameters map[string]string
}

type api_response struct {
	result string
}

/**
 * [checkErr description]
 * @param  {[type]} err error         [description]
 * @return {[type]}     [description]
 */
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/**
 * [getUserData description]
 * @param  {[type]} userId int           [description]
 * @return {[type]}        [description]
 */
// func getUserData(userId int) {
// 	log.Fatal("getUserData")
// }

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}

	return string(bytes)
}

/**
 * [getData description]
 * @param  {[type]} rw  http.ResponseWriter [description]
 * @param  {[type]} req *http.Request       [description]
 * @return {[type]}     [description]
 */
func getUserData(rw http.ResponseWriter, req *http.Request) {
	////////////////////////////
	// Get request parameters //
	////////////////////////////
	decoder := json.NewDecoder(req.Body)
	var r api_request
	err := decoder.Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(r.Parameters)

	///////////////////////////
	// Push request to Queue //
	///////////////////////////
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	q, _ := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	msgs, _ := ch.Consume(
		q.Name, // queue
		"",     // consume
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	corrId := randomString(32)
	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte("RPC SEND"),
		})
	checkErr(err)
	for d := range msgs {
		if corrId == d.CorrelationId {
			////////////////////////
			// Send Response back //
			////////////////////////
			rw.WriteHeader(http.StatusOK)
			res := []byte(string(d.Body))
			rw.Write(res)
			break
		}
	}

	////////////////////////
	// Send Response back //
	////////////////////////
	// rw.WriteHeader(http.StatusOK)
	// res := []byte("API Response")
	// rw.Write(res)
}

/**
 * [Call description]
 * @param {[type]} m      map[string]interface{} [description]
 * @param {[type]} name   string                   [description]
 * @param {[type]} params ...interface{})        (result       []reflect.Value, err error [description]
 */
// func Call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
// 	f := reflect.ValueOf(m[name])
// 	if len(params) != f.Type().NumIn() {
// 		err = errors.New("The number of params is not adapted.")
// 		return
// 	}
// 	in := make([]reflect.Value, len(params))
// 	for k, param := range params {
// 		in[k] = reflect.ValueOf(param)
// 	}
// 	result = f.Call(in)
// 	return
// }

/**
 * [RPCHandler description]
 * @param {[type]} rw  http.ResponseWriter [description]
 * @param {[type]} req *http.Request       [description]
 */
// func RPCHandler(rw http.ResponseWriter, req *http.Request) {
// 	if req.Method == "POST" {
// 		log.Printf("Hello, %q", html.EscapeString(req.URL.Path))
// 		decoder := json.NewDecoder(req.Body)
// 		var r api_request
// 		err := decoder.Decode(&r)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		log.Println(r.Method)
// 		for key, value := range r.Parameters {
// 			log.Println(key + " = " + value)
// 		}
// 	}
// }

/**
 * [main description]
 * @return {[type]} [description]
 */
func main() {
	// funcs := map[string]interface{}{
	// 	"getUserData": getUserData,
	// }

	// _, err := Call(funcs, "getUserData", 1)
	// checkErr(err)

	// http.HandleFunc("/rpc", RPCHandler)
	http.HandleFunc("/getUserData", getUserData)
	http.ListenAndServe(":9898", nil)
}
