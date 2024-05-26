package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

/*
	Notes for myself:
	1. Naming convetion
	 - Camel case but capital first letter if public

	2. Installing dependencies (testify)
	 - https://pkg.go.dev/github.com/stretchr/testify#section-readme
	 - To use go get you need to have a go.mod (module) file
	 - use command: go mod init <module name>

	3. Send HTTP requests

	Send request to add item:
	curl --header "Content-Type: application/json" \
	--request POST \
	--data '{"title":"Learn go","description":"complete intern project"}' \
	http://localhost:8090/addItem

	Send request to get all items:
	curl --header "Accept: application/json" \
	--request GET \
	http://localhost:8090/getAllItems

	Send request to mark item as completed
	curl --request PUT \
	http://localhost:8090/completeItem?id=0

	4. Datadog
	 - The client (Datadog StatsD client) is stored in the global client variable.

*/

// Note that the field names in the struct should start with an uppercase letter
// to be exported and accessible by the json package

// `json:"id"`, in the request the field is lowercase id
type Item struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:completed`
}

var items []Item
var itemId = 0 //used to define id (no persistant datastore)
var client *statsd.Client

func initStatsdClient() {
	var err error
	client, err = statsd.New("127.0.0.1:8125", statsd.WithTags([]string{"env:local", "service:todo-app", "version:1.0.0"}))
	if err != nil {
		log.Fatalf("Failed to create StatsD client: %v", err)
	}
}

func initLogFile() {
	// Opens a log file in write mode
	// The logs will be written to a file named app.log in the current directory
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file %v", err)
		return
	}
	log.SetOutput(logFile)
	log.Println("Logging to file started")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			When a new HTTP request comes in, the loggingMiddleware function is called.
		*/
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.RequestURI)

		// responseRecorder is created using &responseRecorder{w, http.StatusOK}.
		// This instance captures the original http.ResponseWriter (w) and initializes the status code to http.StatusOK.
		// Note how you instantiate a struct
		rr := &responseRecorder{w, http.StatusOK}

		// ServeHTTP method of the next HTTP handler (next) is called with
		// the responseRecorder as the http.ResponseWriter argument.
		next.ServeHTTP(rr, r)

		duration := time.Since(start)
		log.Printf("Completed %s %s %d %s", r.Method, r.RequestURI, rr.statusCode, duration)
	})
}

// Supposedly its easier to get statuscode by encapsulating ResponseWriter
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	/*
		WriteHeader defined. This method overrides the standard WriteHeader method of the http.ResponseWriter.
		When WriteHeader is called on a responseRecorder instance, it captures the status code passed to it and
		also calls WriteHeader on the embedded http.ResponseWriter, ensuring that the response's status code is
		set correctly.
	*/
	rr.ResponseWriter.WriteHeader(code)
}

func getAllItems(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items) //NewEncoder writes to w

	// Sends the timing data to datadog
	duration := time.Since(start)
	/*
	 func (c *Client) Timing(name string, value time.Duration, tags []string, rate float64) error
	 tags ([]string): An array of tags to associate with the metric.
	 Tags are used for filtering and aggregating metrics in Datadog.
	 If there are no specific tags to add, you can pass nil.
	 rate (float64): The sample rate for the metric. It is a value between 0 and 1 that indicates the probability of the metric being sent.
	 A rate of 1 means that all data points will be sent, while a rate of 0.5 means that, on average,
	 only half of the data points will be sent.
	 This is useful for controlling the volume of data sent to Datadog.
	*/
	client.Timing("todo_app.get_all_items.duration", duration, nil, 1)
	client.Incr("todo_app.get_all_items.count", nil, 1)

	log.Printf("getAllItems function: duration=%s, items_count=%d", duration, len(items))
}

func addItem(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	var newItem Item
	err := json.NewDecoder(req.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		client.Incr("todo_app.add_item.errors", nil, 1)
		log.Printf("addItem function: failed to add item")
		return
	}
	newItem.Id = itemId
	itemId++
	items = append(items, newItem)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Item successfully added")

	duration := time.Since(start)
	client.Timing("todo_app.get_all_items.duration", duration, nil, 1)
	client.Incr("todo_app.add_item.count", nil, 1)

	log.Printf("addItem function: new item added %v", newItem.Title)
}

func completeItem(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	idString := req.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Please enter the item Id", http.StatusBadRequest)
		client.Incr("todo_app.complete_item.errors", nil, 1)
		return
	}

	id, err := strconv.Atoi(idString) //converts string to integer
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		client.Incr("todo_app.complete_item.errors", nil, 1)
		return
	}
	for i := range items {
		if items[i].Id == id {
			items[i].Completed = true
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Item is completed")
			duration := time.Since(start)
			client.Timing("todo_app.complete_item.duration", duration, nil, 1)
			log.Printf("completeItem function: item completed %v", items[i].Title)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
	client.Incr("todo_app.complete_items.errors", nil, 1)
}

func main() {

	initStatsdClient()
	/*
		The line defer client.Close() ensures that the client.Close() method is called when the
		main function completes, regardless of how it completes (e.g., whether it exits normally or
		because of an error).
	*/
	defer client.Close()

	initLogFile()

	mux := http.NewServeMux() // multiplexer (or router) for handling HTTP requests,
	// allows us to provide middleware without affecting global handler
	mux.HandleFunc("/getAllItems", getAllItems)
	mux.HandleFunc("/addItem", addItem)
	mux.HandleFunc("/completeItem", completeItem)

	loggedMux := loggingMiddleware(mux)

	log.Println("Server is starting on :8090")
	http.ListenAndServe(":8090", loggedMux)
}
