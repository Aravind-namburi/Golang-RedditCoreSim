# Golang-RedditCoreSim

**A Distributed Reddit-like Engine with REST API and Simulation**

## Overview

Golang-RedditCoreSim is a Reddit-inspired distributed application built with the actor model using Go and the [protoactor-go](https://github.com/asynkron/protoactor-go) framework. It supports core Reddit functionalities such as user registration, subreddit creation, posting, commenting, voting, and karma tracking. It simulates user interactions and also exposes a REST API for client integration.

This project was developed in two phases:
- **Part 1:** Core engine using actor-based concurrency with client simulation
- **Part 2:** RESTful API integration to expose engine functionality to external clients

---

## Features

- ✅ Register accounts
- ✅ Create/join/leave subreddits
- ✅ Text-based posts and comments (including hierarchical commenting)
- ✅ Upvotes/downvotes and karma computation
- ✅ Direct messaging between users
- ✅ Zipf-distributed subreddit popularity
- ✅ REST API with endpoints for user, post, comment, voting
- ✅ Interactive simulator to simulate thousands of users
- ✅ Digital signature verification for post authenticity (Bonus)

---

## API Endpoints

| Method | Endpoint               | Description                 |
|--------|------------------------|-----------------------------|
| POST   | `/api/users`           | Register a user             |
| POST   | `/api/subreddits`      | Create a subreddit          |
| POST   | `/api/posts`           | Create a post               |
| POST   | `/api/comments`        | Add a comment               |
| POST   | `/api/votes`           | Upvote or downvote a post   |
| GET    | `/api/users/karma`     | Get all users with karma    |

---

## How to Run the Project

### Prerequisites

- [Go](https://golang.org/dl/) 1.18 or higher installed
- `go mod` support enabled

### Steps

1. **Clone the Repository**

```bash
git clone https://github.com/Aravind-namburi/Golang-RedditCoreSim.git
cd Golang-RedditCoreSim
```

2. **Install Dependencies**

```bash
go mod tidy
```

3. **Run the Server and Simulation**

```bash
go run main.go
```

This starts:
- The REST API server at `http://localhost:8080`
- The user activity simulator in the background

---

## Team Members

- Aravind Namburi

---

## Performance Metrics

- Simulated 1000+ users in stress tests
- Subreddits follow Zipf distribution for realistic traffic
- Handled multiple concurrent API requests and simulated clients

---

## What’s Working

- Full Reddit-like behavior via actors
- REST API with Gorilla Mux
- REST + simulation integration
- Digital signature support (via crypto libraries)

---

## Folder Structure

```
.
├── engine/              # Actor definitions and core logic
├── simulator/           # Simulator to test interactions
├── models.go            # Models for users, posts, etc.
├── main.go              # API server + simulator launcher
├── go.mod / go.sum      # Dependency management
```

---

## License

This project is licensed under the MIT License.
