package engine

type User struct {
	ID           int
	Username     string
	Password     string
	Karma        int
	PostKarma    int
	CommentKarma int
}

type Subreddit struct {
	ID      int
	Name    string
	Members map[int]bool
	Posts   []int
}

type Post struct {
	ID        int
	UserID    int
	Subreddit string
	Content   string
	Upvotes   int
	Downvotes int
	Comments  []*Comment
}

type Comment struct {
	ID        int
	PostID    int
	ParentID  int
	UserID    int
	Content   string
	Upvotes   int
	Downvotes int
	Replies   []*Comment
}
