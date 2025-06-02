package engine

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
)

// UserActor
type UserActor struct {
	users map[int]*User
	mu    sync.Mutex
}

func (u *UserActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *RegisterUser:
		u.mu.Lock()
		id := len(u.users) + 1
		u.users[id] = &User{ID: id, Username: msg.Username, Password: msg.Password, Karma: 0, PostKarma: 0, CommentKarma: 0}
		fmt.Printf("User %s registered with ID %d\n", msg.Username, id)
		u.mu.Unlock()

	case *UpdateKarma:
		u.mu.Lock()
		user, exists := u.users[msg.UserID]
		if exists {
			user.Karma += msg.KarmaChange
			fmt.Printf("User %d's karma updated to %d\n", msg.UserID, user.Karma)
		} else {
			fmt.Printf("User ID %d does not exist\n", msg.UserID)
		}
		u.mu.Unlock()
	case *GetAllUsers:
		u.mu.Lock()
		fmt.Println("Current Karma for all users:")
		for id, user := range u.users {
			fmt.Printf("User %d (%s): Karma: %d\n", id, user.Username, user.Karma)
		}
		u.mu.Unlock()
	}

}

// SubredditActor
type SubredditActor struct {
	subreddits map[string]*Subreddit
	mu         sync.Mutex
}

func (s *SubredditActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CreateSubreddit:
		s.mu.Lock()
		if _, exists := s.subreddits[msg.Name]; exists {
			fmt.Printf("Subreddit %s already exists\n", msg.Name)
		} else {
			id := len(s.subreddits) + 1
			s.subreddits[msg.Name] = &Subreddit{
				ID:      id,
				Name:    msg.Name,
				Members: make(map[int]bool),
				Posts:   []int{},
			}
			fmt.Printf("Subreddit %s created\n", msg.Name)
		}
		s.mu.Unlock()
	}
}

// PostActor
type PostActor struct {
	posts     map[int]*Post
	userActor *actor.PID
	mu        sync.Mutex
}

func (p *PostActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *AssignUserActor:
		p.userActor = msg.UserActor
		fmt.Println("UserActor assigned to PostActor")

	case *PostMessage:
		p.mu.Lock()
		id := len(p.posts) + 1
		p.posts[id] = &Post{
			ID:        id,
			UserID:    msg.UserID,
			Subreddit: msg.Subreddit,
			Content:   msg.Content,
		}
		fmt.Printf("Post %d created in subreddit %s by user %d\n", id, msg.Subreddit, msg.UserID)
		p.mu.Unlock()

	case *CommentMessage:
		p.mu.Lock()
		post, exists := p.posts[msg.PostID]
		if !exists {
			fmt.Printf("Post ID %d does not exist\n", msg.PostID)
			p.mu.Unlock()
			return
		}
		commentID := len(post.Comments) + 1
		comment := &Comment{
			ID:       commentID,
			PostID:   msg.PostID,
			ParentID: msg.ParentID,
			UserID:   msg.UserID,
			Content:  msg.Content,
		}
		post.Comments = append(post.Comments, comment)
		fmt.Printf("Comment added to post %d by user %d\n", msg.PostID, msg.UserID)
		p.mu.Unlock()

	case *Vote:
		p.mu.Lock()
		if msg.Target == "post" {
			post, exists := p.posts[msg.ID]
			if !exists {
				fmt.Printf("Post ID %d does not exist\n", msg.ID)
				p.mu.Unlock()
				return
			}
			if msg.Type == "upvote" {
				post.Upvotes++
				ctx.Send(p.userActor, &UpdateKarma{UserID: post.UserID, KarmaChange: 1})
			} else if msg.Type == "downvote" {
				post.Downvotes++
				ctx.Send(p.userActor, &UpdateKarma{UserID: post.UserID, KarmaChange: -1})
			}
			fmt.Printf("Post %d %sd by user %d\n", msg.ID, msg.Type, msg.UserID)
		}
		p.mu.Unlock()
	}
}
