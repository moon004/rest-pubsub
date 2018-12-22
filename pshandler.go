package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type PubBodyData struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type PullData struct {
	SubName string `json:"subname"`
}

func CheckError(w http.ResponseWriter, msg string, err error) {
	fmt.Printf("%s: %v \n", msg, err)
	if err != nil {
		fmt.Fprintf(w, "%s: %v", msg, err)
	}
}

func PubSubHandler() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Post("/create", HandleCreate)
	r.Options("/create", HandleCreate)
	r.Post("/publish", HandlePublish)
	r.Options("/publish", HandlePublish)
	r.Get("/pull", HandlePull)
	r.Options("/pull", HandlePull)

	return r
}

// CorsHandler Handle the Cross origin methods
func CorsHandler(w http.ResponseWriter) {
	// Handle CORS, RMB wildcard needs to change in production
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin")
}

// HandleCreate create new topic
func HandleCreate(w http.ResponseWriter, r *http.Request) {
	CorsHandler(w)
	ctx := context.Background()

	projectID := GetEnvVar("GOOGLE_CLOUD_PROJECT")
	client, err := pubsub.NewClient(ctx, projectID)
	CheckError(w, "Error Creating Client", err)
	t, err := client.CreateTopic(ctx, "top3")
	CheckError(w, "Error Creating topic", err)
	fmt.Fprintf(w, "Topic Created Successfully %s", t)
}

func HandlePublish(w http.ResponseWriter, r *http.Request) {
	CorsHandler(w)
	fmt.Println("Handle Publish")
	PubData := &PubBodyData{}
	ctx := context.Background()

	projectID := GetEnvVar("GOOGLE_CLOUD_PROJECT")

	client, err := pubsub.NewClient(ctx, projectID)
	CheckError(w, "Error creating pubsub client", err)

	body, err := ioutil.ReadAll(r.Body)
	CheckError(w, "Error parsing request body", err)

	err = json.Unmarshal(body, &PubData)
	CheckError(w, "Error Unmarshalling Data", err)

	t := client.Topic(PubData.Topic)

	ctx = context.Background()
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(PubData.Message),
	})
	id, err := result.Get(ctx)
	CheckError(w, "Error getting ID", err)

	fmt.Fprintf(w, "Published %s", id)
}

func HandlePull(w http.ResponseWriter, r *http.Request) {
	CorsHandler(w)
	fmt.Println("Handle Pull")
	Puller := &PullData{}

	projectID := GetEnvVar("GOOGLE_CLOUD_PROJECT")

	ctx := context.Background()
	cctx, cancel := context.WithCancel(ctx)

	client, err := pubsub.NewClient(ctx, projectID)
	CheckError(w, "Error creating pubsub client", err)

	body, err := ioutil.ReadAll(r.Body)
	CheckError(w, "Error parsing request body", err)

	err = json.Unmarshal(body, &Puller)
	CheckError(w, "Error Unmarshalling Data", err)

	sub := client.Subscription(Puller.SubName)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		fmt.Printf("Got message and acknowledged: %s\n", string(msg.Data))
		// Send back the Json format back to frontend
		w.Write(msg.Data)
		cancel()
	})
}

func GetEnvVar(str string) string {
	projectID := os.Getenv(str)
	if projectID == "" {
		fmt.Fprintf(os.Stderr, "%s environment variable must be set.\n", str)
		os.Exit(1)
	}
	return projectID
}

// func HandleTrigger(w http.ResponseWriter, r *http.Request) {
// 	CorsHandler(w)

// 	client := &http.Client{}
// 	body, err := ioutil.ReadAll(r.Body)
// 	CheckError("Error Reading Request Body in HandleTrigger", err)

// 	jsonData := `{"topic":"top3", "message":"` + string(body) + `"}`
// 	// Trigger publish Google Functions
// 	req, err := http.NewRequest("POST", ,strings.NewReader(jsonData))
// 	CheckError("Error POST request", err)
// 	_, err = client.Do(req)
// 	CheckError("Error on Client Do Post", err)
// }
