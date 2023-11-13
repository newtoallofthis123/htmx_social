package main

type Env struct {
	User     string
	Password string
	DB       string
	Host     string
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Name      any    `json:"name"`
	Bio       any    `json:"bio"`
}

type PrivateUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Password  string `json:"password"`
	Name      any    `json:"name"`
	Bio       any    `json:"bio"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Post struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updated_at"`
}

type CreatePostRequest struct {
	Content string `json:"content"`
	UserID  string `json:"user_id"`
}

type Like struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	PostID    string `json:"post_id"`
	CreatedAt string `json:"created_at"`
}

type FullPost struct {
	post  Post
	user  User
	likes []Like
}
