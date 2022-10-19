package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Park struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	InPark       string `json:"inPark"`
	Manufacturer string `json:"manufacturer"`
}

var parks []Park

// Function that returns the instance from the list of parks by searching by its name
func parks_search_from_name(name string) (park *Park, index int) {
	if len(parks) == 0 {
		return nil, -1
	}
	for i, v := range parks {
		if v.Name == name {
			return &v, i
		}
	}
	return nil, -1
}

// Function that returns the instance from the list of parks by searching by its id
func parks_search_from_id(id int64) (park *Park, index int) {
	if len(parks) == 0 {
		return nil, -1
	}
	for i, v := range parks {
		if v.Id == id {
			return &v, i
		}
	}
	return nil, -1
}

// Function that finds the largest id and return this id+1
func parks_get_next_id() int64 {
	var res int64 = 0

	var tmp int64
	for _, park := range parks {
		if park.Id > tmp {
			tmp = park.Id
			res = tmp
		}
	}

	return res + 1
}

// Function that removes element by combining two slices made of all elements to the one we remove and
// Another slice made of all elements starting from the one we remove to the end of the original slice
func parks_remove_from_index(index uint) error {
	if index > uint(len(parks)) {
		errors.New("Out of range index")
	}
	if index == uint(len(parks))-1 {
		// Remove last element of the slice
		parks = append(parks[:index])
	} else {
		parks = append(parks[:index], parks[index+1:]...)
	}

	return nil
}

// Function to handle bad request
func send_bad_req(reason string, w http.ResponseWriter) {
	log.Print("Bad request from user: " + reason)
	w.WriteHeader(http.StatusBadRequest)
	// Return empty response
	w.Write([]byte("{\"error\": \"" + reason + "\"}"))
}

// Function that handle not found pages
func send_not_found(w http.ResponseWriter) {
	// Send a 404
	w.WriteHeader(http.StatusNotFound)
	// Return empty response
	w.Write([]byte("{}"))
}

// Function to handle good requests
func send_ok(w http.ResponseWriter) {
	// Just send a 200 OK
	w.WriteHeader(http.StatusOK)
	// Return empty response
	w.Write([]byte("{}"))
}

// Function to send a park instance to the client
func send_park(park *Park, w http.ResponseWriter) {
	w.WriteHeader(http.StatusFound)
	json_encoder := json.NewEncoder(w)
	json_encoder.Encode(park)
}

func parse_request_json(w http.ResponseWriter, r *http.Request) (*Park, error) {
	var search Park
	// Set default id to an invalid one
	search.Id = -1
	// Decode json body
	json_decoder := json.NewDecoder(r.Body)
	json_decoder.DisallowUnknownFields()
	err := json_decoder.Decode(&search)
	if err != nil {
		// could not understand the request's body
		return nil, err
	}

	return &search, nil
}

// Function to handle the GET requests
func handle_get(w http.ResponseWriter, r *http.Request) {
	log.Print("Got a GET request from ", r.RemoteAddr)

	search, err := parse_request_json(w, r)

	if err != nil {
		send_bad_req("Got error when parsing JSON", w)
		return
	}
	if search.Id < -1 {
		// Invalid id to search
		send_bad_req("'id' cannot be negative", w)
		return
	}
	if search.Id == -1 && search.InPark == "" && search.Manufacturer == "" && search.Name == "" {
		// Send all parks if nothing was specified
		w.WriteHeader(http.StatusFound)
		value, _ := json.Marshal(parks)
		w.Write(value)
		return
	}

	var res *Park
	if search.Id != -1 && search.Name == "" {
		// Seach by Id
		res, _ = parks_search_from_id(search.Id)
	} else if search.Id == -1 && search.Name != "" {
		// Seach by Name
		res, _ = parks_search_from_name(search.Name)
	} else {
		send_bad_req("Cannot filter with 'id' and 'name' at the same time", w)
		return
	}
	if res == nil {
		send_not_found(w)
		return
	}
	send_park(res, w)
}

// Function to handle the POST requests
func handle_post(w http.ResponseWriter, r *http.Request) {
	// Here POST will update or create the park if it is not in the list
	log.Print("Got a POST request from ", r.RemoteAddr)

	search, err := parse_request_json(w, r)

	if err != nil {
		send_bad_req("Got error when parsing JSON", w)
		return
	}
	if search.Id < -1 {
		// Invalid id to search
		send_bad_req("'id' cannot be negative", w)
		return
	}
	if search.Id == -1 {
		// Tried to update a park without its id, should have made a PUT request instead
		send_bad_req("You need to specify 'id'", w)
		return
	}
	_, index := parks_search_from_id(search.Id)

	if index == -1 {
		// Value does not exist
		send_not_found(w)
	} else {
		// Value exists, updating it, note that some parts might have bee ommited in the input,
		// So we are updating only things that was specified
		if search.Name != "" {
			parks[index].Name = search.Name
		}
		if search.InPark != "" {
			parks[index].InPark = search.InPark
		}
		if search.Manufacturer != "" {
			parks[index].Manufacturer = search.Manufacturer
		}
		send_park(&parks[index], w)
	}
}

// Function to handle the DELETE requests
func handle_delete(w http.ResponseWriter, r *http.Request) {
	// Delete a park from its id
	log.Print("Got a DELETE request from ", r.RemoteAddr)

	search, err := parse_request_json(w, r)

	if err != nil {
		send_bad_req("Got error when parsing JSON", w)
		return
	}
	if search.Id <= -1 {
		// Invalid id to search
		send_bad_req("'id' cannot be negative", w)
		return
	}
	var index int
	_, index = parks_search_from_id(search.Id)

	if index == -1 {
		// Value not found
		send_not_found(w)
		return
	}

	err = parks_remove_from_index(uint(index))
	if err != nil {
		// An error occured while removing
		send_not_found(w)
		return
	}
	send_ok(w)
}

// Function to handle the PUT requests
func handle_put(w http.ResponseWriter, r *http.Request) {
	// Create a new park
	log.Print("Got a PUT request from ", r.RemoteAddr)
	search, err := parse_request_json(w, r)
	if err != nil {
		send_bad_req("Got error when parsing JSON", w)
		return
	}
	if search.Name == "" || search.InPark == "" || search.Manufacturer == "" {
		// If missing some elements
		send_bad_req("'name', 'inPark' and 'manufacturer' have to be set when PUTing a new element.", w)
		return
	}
	// Generate a new id and append to the list of parks
	search.Id = parks_get_next_id()
	parks = append(parks, *search)
	send_park(search, w)
}

// Function to handle all requests
func handle_requests(w http.ResponseWriter, r *http.Request) {
	// This function will call the corresponding function for the http method used
	switch r.Method {
	case http.MethodGet:
		handle_get(w, r)
	case http.MethodPost:
		handle_post(w, r)
	case http.MethodDelete:
		handle_delete(w, r)
	case http.MethodPut:
		handle_put(w, r)
	default:
		// Other http method than RESTful one
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		// Send empty response
		w.Write([]byte("{}"))
	}
}

func main() {
	// Intialize list of parks
	parks = make([]Park, 0, 5)
	fmt.Println("Welcome to your Rest API\nTo add a new data, you can try : curl -X PUT 127.0.0.1:8080/endpoint --data '{\"name\":\"Labyrinthe\", \"inPark\":\"Paris\", \"manufacturer\":\"Vortex\"}'")
	fmt.Println("To retrieve a data : curl -X GET 127.0.0.1:8080/endpoint --data '{\"id\":1}'; you can get with id or name ")
	fmt.Println("To modify a data :curl -X POST 127.0.0.1:8080/endpoint --data '{\"id\":1, \"inPark\":\"NewYork\"}'")
	fmt.Println("To delete a data : curl -X DELETE 127.0.0.1:8080/endpoint --data '{\"id\":1}'")
	http.HandleFunc("/endpoint", handle_requests)
	log.Print("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
