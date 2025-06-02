/*package main

import (
	"reddit_clone2/engine"
	"reddit_clone2/simulator"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func main() {
	// Set up the actor system
	actorSystem := actor.NewActorSystem()                                        // Initialize the actor system
	rootContext := actorSystem.Root                                              // Get the root context
	userActor, subredditActor, postActor := engine.SetupActorSystem(rootContext) // Create actors

	// Simulate user activity
	simulator.SimulateUsers(rootContext, userActor, subredditActor, postActor)

	// Allow the system to process all messages before exiting
	time.Sleep(5 * time.Second)
}
*/

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reddit_clone2/engine"
	"reddit_clone2/simulator"

	"github.com/gorilla/mux"
)

var (
	actorSystem *engine.ActorSystem
)

func main() {
	// Initialize the Actor System
	actorSystem = engine.NewActorSystem()
	actorSystem.SetupActors()

	// Setup the router
	r := mux.NewRouter()

	// API Endpoints
	r.HandleFunc("/api/users", RegisterUser).Methods("POST")
	r.HandleFunc("/api/subreddits", CreateSubreddit).Methods("POST")
	r.HandleFunc("/api/posts", CreatePost).Methods("POST")
	r.HandleFunc("/api/comments", AddComment).Methods("POST")
	r.HandleFunc("/api/votes", VotePost).Methods("POST")
	r.HandleFunc("/api/users/karma", GetAllUsers).Methods("GET")

	// Start REST API Server
	go func() {
		log.Println("Starting REST API server on :8080")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	// Start the interactive simulator
	simulator.SimulateUsers(actorSystem)
}

// REST API Handlers
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user engine.RegisterUser
	json.NewDecoder(r.Body).Decode(&user)
	actorSystem.RegisterUser(user)
	w.WriteHeader(http.StatusCreated)
}

func CreateSubreddit(w http.ResponseWriter, r *http.Request) {
	var subreddit engine.CreateSubreddit
	json.NewDecoder(r.Body).Decode(&subreddit)
	actorSystem.CreateSubreddit(subreddit)
	w.WriteHeader(http.StatusCreated)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post engine.PostMessage
	json.NewDecoder(r.Body).Decode(&post)
	actorSystem.CreatePost(post)
	w.WriteHeader(http.StatusCreated)
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	var comment engine.CommentMessage
	json.NewDecoder(r.Body).Decode(&comment)
	actorSystem.AddComment(comment)
	w.WriteHeader(http.StatusCreated)
}

func VotePost(w http.ResponseWriter, r *http.Request) {
	var vote engine.Vote
	json.NewDecoder(r.Body).Decode(&vote)
	actorSystem.VotePost(vote)
	w.WriteHeader(http.StatusOK)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := actorSystem.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}
