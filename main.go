package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"net"
	"net/http"
)

var (
	location string
)

func main() {

	http.HandleFunc("/set", Transactions)
	http.HandleFunc("/statistics", Statistics)
	http.HandleFunc("/reset", Delete)

	er := http.ListenAndServe(net.JoinHostPort("localhost", "9090"), nil)
	if er != nil {
		fmt.Println("Error initializing Http rest Api" + er.Error())
	}

	fmt.Println("Client initialized successfully")

}

func Transactions(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Method not allowed, For setting location use POST"))
		fmt.Println("error in request method ", request.Method)
		return
	}
	//log.Println("Body in post method for transaction is ", request.Body)
	SubmitTransaction(request.Body, writer)
	return
}

func SubmitTransaction(request io.ReadCloser, writer http.ResponseWriter) {

	if request==nil{
		fmt.Println("body is nil ")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Please fill the required fields in body"))
		return
	}
	bytes, err := ioutil.ReadAll(request)
	defer request.Close()
	if err != nil {
		fmt.Println("erron in reading ", err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Invalid Json"))
		return
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(bytes, &jsonMap)
	if err != nil {
		fmt.Println("error while marshaling json ")
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte("fields are not parsable: marshaling failed"))
		return
	}

	if len(jsonMap)>1{
		fmt.Println("Only one parameter city is allowed")
		writer.WriteHeader(http.StatusRequestHeaderFieldsTooLarge)
		writer.Write([]byte("only one field is allowed city"))
		return
	}
	city:=jsonMap["city"].(string)
	location = city
	fmt.Println("Data after Marshaling ", location)

	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write([]byte("Location Set successfully"))

}

func Statistics(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = writer.Write([]byte("Method not allowed, USE GET for statistics"))
		log.Println("Error in Request method type, for Statistics use GET method", request.Method)
		return
	}
	//Todo code can be written to get the current location. For now hard coding it.
	currentLocation :="Bangalore"

	if location =="" || currentLocation==location{

		var resp = make(map[string]string)
		resp["statistic"] = "statistics api is accessible"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("error while converting map to response json format")
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("Status Internal Server Error"))
			return
		}
		//	fmt.Println("jsonResponse is ", jsonResp)
		writer.Header().Add("content-type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(jsonResp)
		return
	}
	fmt.Println("Accesible from all places")
	writer.WriteHeader(http.StatusUnauthorized)
	writer.Write([]byte("UnAuthorized Access"))
	return
}

func Delete(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "DELETE" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Method not allowed, Use DELETE for reset"))
		fmt.Println("error in request menthod ", request.Method)
		return
	}
	location=""

	writer.WriteHeader(http.StatusNoContent)
}

