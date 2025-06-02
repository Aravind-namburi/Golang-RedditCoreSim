/*
package simulator

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"reddit_clone2/engine"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

// To see all the action sperformed while simulating many users (option 7)
type Metrics struct {
	TotalActions  int64         // Total actions performed
	TotalDuration time.Duration // Total time taken for all actions
	ActionCounts  map[string]int
	mu            sync.Mutex // Protect concurrent writes to metrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		ActionCounts: make(map[string]int),
	}
}

func (m *Metrics) RecordAction(actionType string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalActions++
	m.TotalDuration += duration
	m.ActionCounts[actionType]++
}

func (m *Metrics) Report() {
	m.mu.Lock()
	defer m.mu.Unlock()
	fmt.Println("\n--- Performance Metrics ---")
	fmt.Printf("Total Actions: %d\n", m.TotalActions)
	fmt.Printf("Average Latency per Action: %v\n", m.TotalDuration/time.Duration(m.TotalActions))
	fmt.Println("Action Counts:")
	for action, count := range m.ActionCounts {
		fmt.Printf("  %s: %d\n", action, count)
	}
	fmt.Println("---------------------------\n")
}

// To see the resource cinsumption while performing simulating many users (option 7)
func ReportResources() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Println("\n--- Resource Utilization ---")
	fmt.Printf("CPU Count: %d\n", runtime.NumCPU())
	fmt.Printf("Threads count: %d\n", runtime.NumGoroutine())
	fmt.Printf("Allocated Memory: %d KB\n", memStats.Alloc/1024)
	fmt.Printf("Total Memory Allocated: %d KB\n", memStats.TotalAlloc/1024)
	fmt.Printf("System Memory: %d KB\n", memStats.Sys/1024)
	fmt.Println("----------------------------\n")
}

// Function to register users
func RegisterUsers(rootContext *actor.RootContext, userActor *actor.PID) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1]

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	rootContext.Send(userActor, &engine.RegisterUser{Username: username, Password: password})
	time.Sleep(100 * time.Millisecond)
}

// Function to create subreddits
func CreateSubreddits(rootContext *actor.RootContext, subredditActor *actor.PID) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter subreddit name: ")
	subredditName, _ := reader.ReadString('\n')
	subredditName = subredditName[:len(subredditName)-1]

	rootContext.Send(subredditActor, &engine.CreateSubreddit{Name: subredditName})
	time.Sleep(100 * time.Millisecond)
}

// Function to create posts
func CreatePosts(rootContext *actor.RootContext, postActor *actor.PID) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userID, _ := strconv.Atoi(userIDStr[:len(userIDStr)-1])

	fmt.Print("Enter subreddit name: ")
	subredditName, _ := reader.ReadString('\n')
	subredditName = subredditName[:len(subredditName)-1]

	fmt.Print("Enter post content: ")
	content, _ := reader.ReadString('\n')
	content = content[:len(content)-1]

	rootContext.Send(postActor, &engine.PostMessage{UserID: userID, Subreddit: subredditName, Content: content})
	time.Sleep(200 * time.Millisecond)
}

// Function to add comments
func AddComments(rootContext *actor.RootContext, postActor *actor.PID) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userID, _ := strconv.Atoi(userIDStr[:len(userIDStr)-1])

	fmt.Print("Enter post ID: ")
	postIDStr, _ := reader.ReadString('\n')
	postID, _ := strconv.Atoi(postIDStr[:len(postIDStr)-1])

	fmt.Print("Enter comment content: ")
	content, _ := reader.ReadString('\n')
	content = content[:len(content)-1]

	rootContext.Send(postActor, &engine.CommentMessage{UserID: userID, PostID: postID, Content: content})
	time.Sleep(150 * time.Millisecond)
}

// Function to display user karma
func DisplayKarma(rootContext *actor.RootContext, userActor *actor.PID) {
	rootContext.Request(userActor, &engine.GetAllUsers{})
	time.Sleep(100 * time.Millisecond)
}

// Function to upvote or Downvote a post/ comment
func UpvoteOrDownvote(rootContext *actor.RootContext, postActor *actor.PID) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userID, _ := strconv.Atoi(userIDStr[:len(userIDStr)-1])

	fmt.Print("Enter 'post' or 'comment' to vote on: ")
	target, _ := reader.ReadString('\n')
	target = target[:len(target)-1]

	fmt.Print("Enter the ID of the post/comment to vote on: ")
	targetIDStr, _ := reader.ReadString('\n')
	targetID, _ := strconv.Atoi(targetIDStr[:len(targetIDStr)-1])

	fmt.Print("Enter 'upvote' or 'downvote': ")
	voteType, _ := reader.ReadString('\n')
	voteType = voteType[:len(voteType)-1]

	rootContext.Send(postActor, &engine.Vote{
		UserID: userID,
		Target: target,
		ID:     targetID,
		Type:   voteType,
	})
	time.Sleep(100 * time.Millisecond)
}


func SimulateManyUsers(rootContext *actor.RootContext, userActor, subredditActor, postActor *actor.PID) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // Local random generator
	metrics := NewMetrics()                                // Create a metrics tracker
	numUsers := 1000                                       // Simulate 1000 users
	maxActionsPerUser := 10                                // Limit to 10 actions per user
	var wg sync.WaitGroup                                  // WaitGroup to track user completion

	fmt.Println("Starting simulation...")
	ReportResources() // Report initial resources

	for i := 1; i <= numUsers; i++ {
		wg.Add(1) // Increment WaitGroup counter for each user
		username := fmt.Sprintf("user%d", i)
		password := "password123"

		// Register the user
		rootContext.Send(userActor, &engine.RegisterUser{Username: username, Password: password})

		// Assign random actions for the user
		go func(userID int) {
			defer wg.Done() // Decrement WaitGroup counter when the goroutine finishes
			for actionCount := 0; actionCount < maxActionsPerUser; actionCount++ {
				startTime := time.Now() // Record start time of the action
				action := rng.Intn(3)   // Randomly choose an action: 0 = post, 1 = comment, 2 = vote
				switch action {
				case 0: // Post
					rootContext.Send(postActor, &engine.PostMessage{
						UserID:    userID,
						Subreddit: fmt.Sprintf("random%d", actionCount),
						Content:   fmt.Sprintf("This is post by user %d", userID),
					})
					metrics.RecordAction("Post", time.Since(startTime))
				case 1: // Comment
					rootContext.Send(postActor, &engine.CommentMessage{
						UserID:  userID,
						PostID:  rng.Intn(100), // Random post ID
						Content: fmt.Sprintf("This is a comment by user %d", userID),
					})
					metrics.RecordAction("Comment", time.Since(startTime))
				case 2: // Vote
					rootContext.Send(postActor, &engine.Vote{
						UserID: userID,
						Target: "post",        // Randomly target posts
						ID:     rng.Intn(100), // Random post ID
						Type:   "upvote",
					})
					metrics.RecordAction("Vote", time.Since(startTime))
				}
				time.Sleep(time.Millisecond * time.Duration(rng.Intn(1000))) // Random delay between actions
			}
			fmt.Printf("User %d completed %d actions\n", userID, maxActionsPerUser)
		}(i)
	}

	wg.Wait()         // Wait for all user goroutines to complete
	metrics.Report()  // Print the metrics report
	ReportResources() // Report resources after simulation
	fmt.Println("Simulating many users completed.")
}

// Menu-driven simulation

func SimulateUsers(rootContext *actor.RootContext, userActor, subredditActor, postActor *actor.PID) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Register a new user")
		fmt.Println("2. Create a new subreddit")
		fmt.Println("3. Create a new post")
		fmt.Println("4. Add a comment to a post")
		fmt.Println("5. Display user karma")
		fmt.Println("6. Upvote/Downvote a post or comment")
		fmt.Println("7. Simulate many users")
		fmt.Println("8. Exit")

		choice, _ := reader.ReadString('\n')
		choice = choice[:len(choice)-1]

		switch choice {
		case "1":
			RegisterUsers(rootContext, userActor)
		case "2":
			CreateSubreddits(rootContext, subredditActor)
		case "3":
			CreatePosts(rootContext, postActor)
		case "4":
			AddComments(rootContext, postActor)
		case "5":
			DisplayKarma(rootContext, userActor)
		case "6":
			UpvoteOrDownvote(rootContext, postActor)
		case "7":
			fmt.Println("Simulating many users...")
			SimulateManyUsers(rootContext, userActor, subredditActor, postActor)
			fmt.Println("Simulation completed. Returning to the menu...")
		case "8":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option, try again.")
		}
	}
}
*/

/*

package simulator

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reddit_clone2/engine"
	"strconv"
)

func SimulateUsers(actorSystem *engine.ActorSystem) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Register a new user")
		fmt.Println("2. Create a new subreddit")
		fmt.Println("3. Create a new post")
		fmt.Println("4. Add a comment to a post")
		fmt.Println("5. Display user karma")
		fmt.Println("6. Upvote/Downvote a post or comment")
		fmt.Println("7. Test API Endpoints")
		fmt.Println("8. Exit")

		choice, _ := reader.ReadString('\n')
		choice = choice[:len(choice)-1]

		switch choice {
		case "1":
			registerUserCLI(actorSystem, reader)
		case "2":
			createSubredditCLI(actorSystem, reader)
		case "3":
			createPostCLI(actorSystem, reader)
		case "4":
			addCommentCLI(actorSystem, reader)
		case "5":
			displayKarmaCLI(actorSystem)
		case "6":
			upvoteOrDownvoteCLI(actorSystem, reader)
		case "7":
			runAPITests()
		case "8":
			fmt.Println("Exiting simulation...")
			return
		default:
			fmt.Println("Invalid option. Try again.")
		}
	}
}

// Command Line Interaction Functions
func registerUserCLI(actorSystem *engine.ActorSystem, reader *bufio.Reader) {
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1]

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	actorSystem.RegisterUser(engine.RegisterUser{Username: username, Password: password})
}

func createSubredditCLI(actorSystem *engine.ActorSystem, reader *bufio.Reader) {
	fmt.Print("Enter subreddit name: ")
	subredditName, _ := reader.ReadString('\n')
	subredditName = subredditName[:len(subredditName)-1]

	actorSystem.CreateSubreddit(engine.CreateSubreddit{Name: subredditName})
}

func createPostCLI(actorSystem *engine.ActorSystem, reader *bufio.Reader) {
	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userID, _ := strconv.Atoi(userIDStr[:len(userIDStr)-1])

	fmt.Print("Enter subreddit name: ")
	subredditName, _ := reader.ReadString('\n')
	subredditName = subredditName[:len(subredditName)-1]

	fmt.Print("Enter post content: ")
	content, _ := reader.ReadString('\n')
	content = content[:len(content)-1]

	actorSystem.CreatePost(engine.PostMessage{
		UserID:    userID,
		Subreddit: subredditName,
		Content:   content,
	})
}

func addCommentCLI(actorSystem *engine.ActorSystem, reader *bufio.Reader) {
	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userID, _ := strconv.Atoi(userIDStr[:len(userIDStr)-1])

	fmt.Print("Enter post ID: ")
	postIDStr, _ := reader.ReadString('\n')
	postID, _ := strconv.Atoi(postIDStr[:len(postIDStr)-1])

	fmt.Print("Enter comment content: ")
	content, _ := reader.ReadString('\n')
	content = content[:len(content)-1]

	actorSystem.AddComment(engine.CommentMessage{
		UserID:  userID,
		PostID:  postID,
		Content: content,
	})
}

func displayKarmaCLI(actorSystem *engine.ActorSystem) {
	users := actorSystem.GetAllUsers()
	fmt.Println("Current User Karma:")
	for id, user := range users {
		fmt.Printf("User ID %d (%s): Karma: %d\n", id, user.Username, user.Karma)
	}
}

func upvoteOrDownvoteCLI(actorSystem *engine.ActorSystem, reader *bufio.Reader) {
	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userID, _ := strconv.Atoi(userIDStr[:len(userIDStr)-1])

	fmt.Print("Enter 'post' or 'comment' to vote on: ")
	target, _ := reader.ReadString('\n')
	target = target[:len(target)-1]

	fmt.Print("Enter the ID of the post/comment to vote on: ")
	targetIDStr, _ := reader.ReadString('\n')
	targetID, _ := strconv.Atoi(targetIDStr[:len(targetIDStr)-1])

	fmt.Print("Enter 'upvote' or 'downvote': ")
	voteType, _ := reader.ReadString('\n')
	voteType = voteType[:len(voteType)-1]

	actorSystem.VotePost(engine.Vote{
		UserID: userID,
		Target: target,
		ID:     targetID,
		Type:   voteType,
	})
}

// --- API Test Functions ---

func runAPITests() {
	fmt.Println("\n--- Running API Endpoint Tests ---")

	// Register a User
	registerPayload := map[string]string{"Username": "testuser", "Password": "password"}
	sendPostRequest("http://localhost:8080/api/users", registerPayload)

	// Create a Subreddit
	subredditPayload := map[string]string{"Name": "testsub"}
	sendPostRequest("http://localhost:8080/api/subreddits", subredditPayload)

	// Create a Post
	postPayload := map[string]interface{}{
		"UserID":    1,
		"Subreddit": "testsub",
		"Content":   "This is a post from API test",
	}
	sendPostRequest("http://localhost:8080/api/posts", postPayload)

	// View All Users’ Karma
	sendGetRequest("http://localhost:8080/api/users/karma")

	// Add a Comment
	commentPayload := map[string]interface{}{
		"UserID":  1,
		"PostID":  1,
		"Content": "Great post!",
	}
	sendPostRequest("http://localhost:8080/api/comments", commentPayload)

	// Upvote a Post
	votePayload := map[string]interface{}{
		"UserID": 1,
		"Target": "post",
		"ID":     1,
		"Type":   "upvote",
	}
	sendPostRequest("http://localhost:8080/api/votes", votePayload)
}

func sendPostRequest(url string, payload interface{}) {
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		fmt.Printf("Failed POST to %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("POST to %s - Status: %s\n", url, resp.Status)
}

func sendGetRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed GET from %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("GET from %s - Status: %s\n", url, resp.Status)
}
*/

package simulator

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reddit_clone2/engine"
	"strconv"
)

func SimulateUsers(actorSystem *engine.ActorSystem) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Register a new user")
		fmt.Println("2. Create a new subreddit")
		fmt.Println("3. Create a new post")
		fmt.Println("4. Add a comment to a post")
		fmt.Println("5. Display user karma")
		fmt.Println("6. Upvote/Downvote a post or comment")
		fmt.Println("7. Test API Endpoints")
		fmt.Println("8. Exit")

		choice, _ := reader.ReadString('\n')
		choice = choice[:len(choice)-1]

		switch choice {
		case "1":
			registerUserCLI(actorSystem, reader)
		case "2":
			createSubredditCLI(actorSystem, reader)
		case "3":
			createPostCLI(actorSystem, reader)
		case "4":
			addCommentCLI(actorSystem, reader)
		case "5":
			displayKarmaCLI(actorSystem)
		case "6":
			upvoteOrDownvoteCLI(actorSystem, reader)
		case "7":
			runAPITests(reader)
		case "8":
			fmt.Println("Exiting simulation...")
			return
		default:
			fmt.Println("Invalid option. Try again.")
		}
	}
}

// --- CLI Interaction Functions ---
func registerUserCLI(actorSystem *engine.ActorSystem, reader *bufio.Reader) {
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1]

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	actorSystem.RegisterUser(engine.RegisterUser{Username: username, Password: password})
}

func createSubredditCLI(actorSystem *engine.ActorSystem, reader *bufio.Reader) {
	fmt.Print("Enter subreddit name: ")
	subredditName, _ := reader.ReadString('\n')
	subredditName = subredditName[:len(subredditName)-1]

	actorSystem.CreateSubreddit(engine.CreateSubreddit{Name: subredditName})
}

func createPostCLI(actorSystem *engine.ActorSystem, reader *bufio.Reader) {
	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userID, _ := strconv.Atoi(userIDStr[:len(userIDStr)-1])

	fmt.Print("Enter subreddit name: ")
	subredditName, _ := reader.ReadString('\n')
	subredditName = subredditName[:len(subredditName)-1]

	fmt.Print("Enter post content: ")
	content, _ := reader.ReadString('\n')
	content = content[:len(content)-1]

	actorSystem.CreatePost(engine.PostMessage{
		UserID:    userID,
		Subreddit: subredditName,
		Content:   content,
	})
}

func addCommentCLI(actorSystem *engine.ActorSystem, reader *bufio.Reader) {
	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userID, _ := strconv.Atoi(userIDStr[:len(userIDStr)-1])

	fmt.Print("Enter post ID: ")
	postIDStr, _ := reader.ReadString('\n')
	postID, _ := strconv.Atoi(postIDStr[:len(postIDStr)-1])

	fmt.Print("Enter comment content: ")
	content, _ := reader.ReadString('\n')
	content = content[:len(content)-1]

	actorSystem.AddComment(engine.CommentMessage{
		UserID:  userID,
		PostID:  postID,
		Content: content,
	})
}

func displayKarmaCLI(actorSystem *engine.ActorSystem) {
	users := actorSystem.GetAllUsers()
	fmt.Println("Current User Karma:")
	for id, user := range users {
		fmt.Printf("User ID %d (%s): Karma: %d\n", id, user.Username, user.Karma)
	}
}

func upvoteOrDownvoteCLI(actorSystem *engine.ActorSystem, reader *bufio.Reader) {
	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userID, _ := strconv.Atoi(userIDStr[:len(userIDStr)-1])

	fmt.Print("Enter 'post' or 'comment' to vote on: ")
	target, _ := reader.ReadString('\n')
	target = target[:len(target)-1]

	fmt.Print("Enter the ID of the post/comment to vote on: ")
	targetIDStr, _ := reader.ReadString('\n')
	targetID, _ := strconv.Atoi(targetIDStr[:len(targetIDStr)-1])

	fmt.Print("Enter 'upvote' or 'downvote': ")
	voteType, _ := reader.ReadString('\n')
	voteType = voteType[:len(voteType)-1]

	actorSystem.VotePost(engine.Vote{
		UserID: userID,
		Target: target,
		ID:     targetID,
		Type:   voteType,
	})
}

// --- API Test Functions ---
func runAPITests(reader *bufio.Reader) {
	fmt.Println("\n--- Running API Endpoint Tests ---")

	// Register a User
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1]

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]
	registerPayload := map[string]string{"Username": username, "Password": password}
	sendPostRequest("http://localhost:8080/api/users", registerPayload)

	// Create a Subreddit
	fmt.Print("Enter subreddit name: ")
	subredditName, _ := reader.ReadString('\n')
	subredditName = subredditName[:len(subredditName)-1]
	subredditPayload := map[string]string{"Name": subredditName}
	sendPostRequest("http://localhost:8080/api/subreddits", subredditPayload)

	// Create a Post
	fmt.Print("Enter post content: ")
	postContent, _ := reader.ReadString('\n')
	postContent = postContent[:len(postContent)-1]
	postPayload := map[string]interface{}{
		"UserID":    1, // Can make dynamic
		"Subreddit": subredditName,
		"Content":   postContent,
	}
	sendPostRequest("http://localhost:8080/api/posts", postPayload)

	// Fetch All Users’ Karma
	sendGetRequest("http://localhost:8080/api/users/karma")

	// Add a Comment
	fmt.Print("Enter comment content: ")
	commentContent, _ := reader.ReadString('\n')
	commentContent = commentContent[:len(commentContent)-1]
	commentPayload := map[string]interface{}{
		"UserID":  1, // Example
		"PostID":  1, // Example Post ID
		"Content": commentContent,
	}
	sendPostRequest("http://localhost:8080/api/comments", commentPayload)

	// Upvote a Post
	votePayload := map[string]interface{}{
		"UserID": 1, // Example
		"Target": "post",
		"ID":     1, // Example Post ID
		"Type":   "upvote",
	}
	sendPostRequest("http://localhost:8080/api/votes", votePayload)
}

func sendPostRequest(url string, payload interface{}) {
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		fmt.Printf("Failed POST to %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("POST to %s - Status: %s\n", url, resp.Status)
}

func sendGetRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed GET from %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("GET from %s - Status: %s\n", url, resp.Status)
}
