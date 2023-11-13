package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store interface {
	createTables() error
	CreateUser(req CreateUserRequest) (string, error)
	GetUser(ID string) (User, error)
	GetUserByEmail(email string) (PrivateUser, error)
	UpdateUser(userID, name string, bio string) error
	CreateSession(userID string) (string, error)
	GetSession(sessionID string) (string, error)
	DeleteSession(sessionID string) error
	CreatePost(req CreatePostRequest) (string, error)
	GetPost(ID string) (Post, error)
	GetPostsByUser(userID string) ([]Post, error)
	LikePost(userID, postID string) error
	UnlikePost(userID, postID string) error
	PostLikes(postID string) ([]Like, error)
	IsLiked(userID, postID string) (bool, error)
	GetFullPost(postID string) (FullPost, error)
}

type DbInstance struct {
	db *sql.DB
}

func NewStore(env *Env) (*DbInstance, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", env.User, env.Password, env.Host, env.DB)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &DbInstance{db}, nil
}

func (pq *DbInstance) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		name TEXT,
		bio TEXT
	);

	CREATE TABLE IF NOT EXISTS posts (
		id TEXT PRIMARY KEY,
		user_id TEXT REFERENCES users(id),
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS comments (
		id TEXT PRIMARY KEY,
		user_id TEXT REFERENCES users(id),
		post_id TEXT REFERENCES posts(id),
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		update_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS likes (
		id TEXT PRIMARY KEY,
		user_id TEXT REFERENCES users(id),
		post_id TEXT REFERENCES posts(id),
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS follows (
		id TEXT PRIMARY KEY,
		follower_id TEXT REFERENCES users(id),
		followee_id TEXT REFERENCES users(id),
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS sessions (
		session_id TEXT PRIMARY KEY,
		user_id TEXT REFERENCES users(id),
		created_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err := pq.db.Exec(query)
	return err
}

func (pq *DbInstance) CreateUser(req CreateUserRequest) (string, error) {
	query := `
	INSERT INTO users (id, email, password) VALUES ($1, $2, $3)
	`

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	userID := RanHash(8)

	_, err = pq.db.Exec(query, userID, req.Email, hashedPassword)

	if err != nil {
		return "", err
	}

	return userID, nil
}

func (pq *DbInstance) GetUser(ID string) (User, error) {
	query := `
	SELECT * from users WHERE id = $1
	`

	var user User
	var temp string

	err := pq.db.QueryRow(query, ID).Scan(&user.ID, &user.Email, &temp, &user.CreatedAt, &user.Name, &user.Bio)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (pq *DbInstance) GetUserByEmail(email string) (PrivateUser, error) {
	query := `
	SELECT * FROM users WHERE email = $1
	`

	var user PrivateUser

	err := pq.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.Name, &user.Bio)
	if err != nil {
		return PrivateUser{}, err
	}

	return user, nil
}

func (pq *DbInstance) CreateSession(userID string) (string, error) {
	query := `
	INSERT INTO sessions (session_id, user_id) VALUES ($1, $2)
	`

	sessionID := RanHash(16)

	_, err := pq.db.Exec(query, sessionID, userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (pq *DbInstance) GetSession(sessionID string) (string, error) {
	query := `
	SELECT user_id FROM sessions WHERE session_id = $1
	`

	var userID string

	err := pq.db.QueryRow(query, sessionID).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (pq *DbInstance) DeleteSession(sessionID string) error {
	query := `
	DELETE FROM sessions WHERE session_id = $1
	`

	_, err := pq.db.Exec(query, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (pq *DbInstance) CreatePost(req CreatePostRequest) (string, error) {
	query := `
	INSERT INTO posts (id, user_id, content) VALUES ($1, $2, $3)
	`

	postID := RanHash(8)

	_, err := pq.db.Exec(query, postID, req.UserID, req.Content)
	if err != nil {
		return "", err
	}

	return postID, nil
}

func (pq *DbInstance) GetPost(ID string) (Post, error) {
	query := `
	SELECT * FROM posts WHERE id = $1
	`

	var post Post

	err := pq.db.QueryRow(query, ID).Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return Post{}, err
	}

	return post, nil
}

func (pq *DbInstance) GetPostsByUser(userID string) ([]Post, error) {
	query := `
	SELECT * FROM posts WHERE user_id = $1
	`

	rows, err := pq.db.Query(query, userID)
	if err != nil {
		return []Post{}, err
	}

	var posts []Post

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return []Post{}, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (pq *DbInstance) UpdateUser(userID, name string, bio string) error {
	query :=
		`
	UPDATE users SET name = $1, bio = $2 WHERE id = $3
	`

	_, err := pq.db.Exec(query, name, bio, userID)
	return err
}

func (pq *DbInstance) LikePost(userID, postID string) error {
	query :=
		`
	INSERT INTO likes (id, user_id, post_id) VALUES ($1, $2, $3)
	`

	likeID := RanHash(8)

	_, err := pq.db.Exec(query, likeID, userID, postID)
	return err
}

func (pq *DbInstance) UnlikePost(userID, postID string) error {
	query :=
		`
	DELETE FROM likes WHERE user_id = $1 AND post_id = $2
	`

	_, err := pq.db.Exec(query, userID, postID)
	return err
}

func (pq *DbInstance) PostLikes(postID string) ([]Like, error) {
	query :=
		`
	SELECT * FROM likes WHERE post_id = $1
	`

	rows, err := pq.db.Query(query, postID)
	if err != nil {
		return []Like{}, err
	}

	var likes []Like

	for rows.Next() {
		var like Like
		err := rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.CreatedAt)
		if err != nil {
			return []Like{}, err
		}
		likes = append(likes, like)
	}

	return likes, nil
}

func (pq *DbInstance) IsLiked(userID, postID string) (bool, error) {
	query :=
		`
	SELECT * FROM likes WHERE user_id = $1 AND post_id = $2
	`

	rows, err := pq.db.Query(query, userID, postID)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (pq *DbInstance) GetFullPost(postID string) (FullPost, error) {
	query := `
	SELECT * FROM posts INNER JOIN users ON posts.user_id = users.id WHERE posts.id = $1
	`

	rows, err := pq.db.Query(query, postID)
	if err != nil {
		return FullPost{}, err
	}

	var fullPost FullPost

	for rows.Next() {
		var post Post
		var user User
		var temp string

		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt, &user.ID, &user.Email, &temp, &user.CreatedAt, &user.Name, &user.Bio)
		if err != nil {
			return FullPost{}, err
		}

		likes, err := pq.PostLikes(postID)
		if err != nil {
			return FullPost{}, err
		}

		fullPost = FullPost{post, user, likes}
	}

	return fullPost, nil
}
