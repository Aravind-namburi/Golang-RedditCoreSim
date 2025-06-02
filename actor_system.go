/*package engine

import "github.com/asynkron/protoactor-go/actor"

func SetupActorSystem(rootContext *actor.RootContext) (*actor.PID, *actor.PID, *actor.PID) {
	userProps := actor.PropsFromProducer(func() actor.Actor { return &UserActor{users: make(map[int]*User)} })
	subredditProps := actor.PropsFromProducer(func() actor.Actor { return &SubredditActor{subreddits: make(map[string]*Subreddit)} })
	postProps := actor.PropsFromProducer(func() actor.Actor {
		return &PostActor{posts: make(map[int]*Post), userActor: nil}
	})

	userActor := rootContext.Spawn(userProps)
	subredditActor := rootContext.Spawn(subredditProps)
	postActor := rootContext.Spawn(postProps)

	// Assign userActor PID to PostActor
	rootContext.Send(postActor, &AssignUserActor{UserActor: userActor})

	return userActor, subredditActor, postActor
}
*/

/*package engine

import (
	"github.com/asynkron/protoactor-go/actor"
)

type ActorSystem struct {
	RootContext    *actor.RootContext
	UserActor      *actor.PID
	SubredditActor *actor.PID
	PostActor      *actor.PID
}

func NewActorSystem() *ActorSystem {
	system := actor.NewActorSystem()
	rootContext := system.Root
	return &ActorSystem{RootContext: rootContext}
}

func (as *ActorSystem) SetupActors() {
	userProps := actor.PropsFromProducer(func() actor.Actor { return &UserActor{users: make(map[int]*User)} })
	subredditProps := actor.PropsFromProducer(func() actor.Actor { return &SubredditActor{subreddits: make(map[string]*Subreddit)} })
	postProps := actor.PropsFromProducer(func() actor.Actor { return &PostActor{posts: make(map[int]*Post)} })

	as.UserActor = as.RootContext.Spawn(userProps)
	as.SubredditActor = as.RootContext.Spawn(subredditProps)
	as.PostActor = as.RootContext.Spawn(postProps)

	// Link UserActor to PostActor
	as.RootContext.Send(as.PostActor, &AssignUserActor{UserActor: as.UserActor})
}

// Public methods to interact with actors
func (as *ActorSystem) RegisterUser(msg RegisterUser) {
	as.RootContext.Send(as.UserActor, &msg)
}

func (as *ActorSystem) CreateSubreddit(msg CreateSubreddit) {
	as.RootContext.Send(as.SubredditActor, &msg)
}

func (as *ActorSystem) CreatePost(msg PostMessage) {
	as.RootContext.Send(as.PostActor, &msg)
}

func (as *ActorSystem) AddComment(msg CommentMessage) {
	as.RootContext.Send(as.PostActor, &msg)
}

func (as *ActorSystem) VotePost(msg Vote) {
	as.RootContext.Send(as.PostActor, &msg)
}

func (as *ActorSystem) GetAllUsers() map[int]*User {
	response := make(map[int]*User)
	as.RootContext.Request(as.UserActor, &GetAllUsers{}).Result(&response)
	return response
}
*/

package engine

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type ActorSystem struct {
	RootContext    *actor.RootContext
	UserActor      *actor.PID
	SubredditActor *actor.PID
	PostActor      *actor.PID
}

func NewActorSystem() *ActorSystem {
	system := actor.NewActorSystem()
	rootContext := system.Root
	return &ActorSystem{RootContext: rootContext}
}

func (as *ActorSystem) SetupActors() {
	userProps := actor.PropsFromProducer(func() actor.Actor { return &UserActor{users: make(map[int]*User)} })
	subredditProps := actor.PropsFromProducer(func() actor.Actor { return &SubredditActor{subreddits: make(map[string]*Subreddit)} })
	postProps := actor.PropsFromProducer(func() actor.Actor { return &PostActor{posts: make(map[int]*Post)} })

	as.UserActor = as.RootContext.Spawn(userProps)
	as.SubredditActor = as.RootContext.Spawn(subredditProps)
	as.PostActor = as.RootContext.Spawn(postProps)

	// Link UserActor to PostActor
	as.RootContext.Send(as.PostActor, &AssignUserActor{UserActor: as.UserActor})
}

// Public methods for REST API interaction
func (as *ActorSystem) RegisterUser(msg RegisterUser) {
	as.RootContext.Send(as.UserActor, &msg)
}

func (as *ActorSystem) CreateSubreddit(msg CreateSubreddit) {
	as.RootContext.Send(as.SubredditActor, &msg)
}

func (as *ActorSystem) CreatePost(msg PostMessage) {
	as.RootContext.Send(as.PostActor, &msg)
}

func (as *ActorSystem) AddComment(msg CommentMessage) {
	as.RootContext.Send(as.PostActor, &msg)
}

func (as *ActorSystem) VotePost(msg Vote) {
	as.RootContext.Send(as.PostActor, &msg)
}

// Corrected Method for Fetching All Users
func (as *ActorSystem) GetAllUsers() map[int]*User {
	// Make a synchronous request
	future := as.RootContext.RequestFuture(as.UserActor, &GetAllUsers{}, 5*time.Second)
	result, err := future.Result()

	if err != nil {
		log.Printf("Error fetching users: %v\n", err)
		return nil
	}

	// Type-assert the result
	if users, ok := result.(map[int]*User); ok {
		return users
	}

	log.Println("Unexpected result type while fetching users")
	return nil
}
