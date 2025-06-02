/*package engine

import "github.com/asynkron/protoactor-go/actor"

type RegisterUser struct {
	Username string
	Password string
}

type CreateSubreddit struct {
	Name string
}

type PostMessage struct {
	UserID    int
	Subreddit string
	Content   string
}

type CommentMessage struct {
	UserID   int
	PostID   int
	ParentID int
	Content  string
}

type Vote struct {
	UserID int
	Target string // "post" or "comment"
	ID     int
	Type   string // "upvote" or "downvote"
}

type GetFeed struct {
	UserID int
}

type UpdateKarma struct {
	UserID      int
	KarmaChange int
}

type AssignUserActor struct {
	UserActor *actor.PID
}

type GetAllUsers struct{}
*/


/*package engine

type RegisterUser struct {
	Username string
	Password string
}

type CreateSubreddit struct {
	Name string
}

type PostMessage struct {
	UserID    int
	Subreddit string
	Content   string
}

type CommentMessage struct {
	UserID   int
	PostID   int
	ParentID int
	Content  string
}

type Vote struct {
	UserID int
	Target string // "post" or "comment"
	ID     int
	Type   string // "upvote" or "downvote"
}

type GetAllUsers struct{}
*/


package engine

import "github.com/asynkron/protoactor-go/actor"

// User Registration
type RegisterUser struct {
	Username string
	Password string
}

// Subreddit Management
type CreateSubreddit struct {
	Name string
}

// Post Management
type PostMessage struct {
	UserID    int
	Subreddit string
	Content   string
}

// Comment Management
type CommentMessage struct {
	UserID   int
	PostID   int
	ParentID int // For hierarchical comments
	Content  string
}

// Voting System
type Vote struct {
	UserID int
	Target string // "post" or "comment"
	ID     int    // ID of the target being voted on
	Type   string // "upvote" or "downvote"
}

// Karma Management
type UpdateKarma struct {
	UserID      int
	KarmaChange int
}

// User-Post Linking
type AssignUserActor struct {
	UserActor *actor.PID
}

// Retrieve All Users
type GetAllUsers struct{}
