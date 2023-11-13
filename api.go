package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	store      Store
	listenAddr string
}

func NewApiServer(store Store, listenAddr string) *ApiServer {
	return &ApiServer{store, listenAddr}
}

func (api *ApiServer) Start() *gin.Engine {
	r := gin.Default()

	// setting up the db and the routes
	err := api.store.createTables()
	// Panic Id: 3
	if err != nil {
		panic(fmt.Sprintf("Error creating tables (Panic Id: 3): %s", err))
	}

	// setting up the static file server
	r.Static("/static", "./public")
	r.LoadHTMLGlob("templates/*")

	// protecting the routes
	protected := r.Group("/").Use(api.handleAuthMiddleware())

	// small health check
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "✨ OK")
	})

	// setting up the routes
	r.GET("/", api.HomePage)

	// the auth stuff
	r.GET("/auth", api.AuthPage)
	r.POST("/signup", api.handleSignup)
	r.POST("/login", api.handleLogin)
	r.GET("/get_auth", api.handleAuth)

	// normal routes
	r.GET("/user/:user_id", api.handleGetUser)
	r.GET("/post/:post_id", api.handleGetFullPost)

	// the post stuff
	protected.POST("/create_post", api.handlePostCreation)
	protected.POST("/update_user", api.handleUserUpdate)
	protected.POST("/like/:post_id", api.handleLike)
	protected.POST("/unlike/:post_id", api.handleUnlike)
	protected.POST("/get_likes", api.handleGetLikes)
	protected.POST("/get_posts", api.handleUserPosts)
	protected.POST("/get_like_status/:post_id", api.handleGetLikeStatus)

	return r
}

func (api *ApiServer) handleAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := c.Cookie("session_id")

		if err != nil {
			c.Redirect(302, "/auth")
			return
		}

		userID, err := api.store.GetSession(sessionId)

		if err != nil {
			c.Redirect(302, "/auth")
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}

func (api *ApiServer) HomePage(c *gin.Context) {
	_, err := c.Cookie("session_id")
	if err != nil {
		c.Redirect(302, "/auth")
		return
	}

	c.HTML(200, "index.html", nil)
}

func (api *ApiServer) AuthPage(c *gin.Context) {

	//check if the user is already logged in
	_, err := c.Cookie("session_id")
	if err == nil {
		c.Redirect(302, "/")
		return
	}

	c.HTML(200, "auth.html", nil)
}

func (api *ApiServer) handleSignup(c *gin.Context) {
	var req CreateUserRequest

	email := c.PostForm("email")
	password := c.PostForm("password")

	req.Email = email
	req.Password = password

	userId, err := api.store.CreateUser(req)

	if err != nil {
		c.String(500, "Error creating user")
		return
	}

	//create the session
	sessionId, err := api.store.CreateSession(userId)

	if err != nil {
		fmt.Println(err)
		c.String(500, "Error creating session")
		return
	}

	//set the cookie
	c.SetCookie("session_id", sessionId, 3600, "/", "localhost", false, true)
	c.String(200, "User created")
}

func (api *ApiServer) handleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.String(400, "Email or password missing")
		return
	}

	//check if the user exists
	user, err := api.store.GetUserByEmail(email)

	if err != nil {
		fmt.Println(err)
		c.String(500, "Error getting user")
		return
	}

	//check if the password is correct
	if !MatchPasswords(password, user.Password) {
		c.String(400, "Incorrect password")
		return
	}

	//create the session
	sessionId, err := api.store.CreateSession(user.ID)

	if err != nil {
		fmt.Println(err)
		c.String(500, "Error creating session")
		return
	}

	//set the cookie
	c.SetCookie("session_id", sessionId, 3600, "/", "localhost", false, false)

	c.String(200, "User logged in")
}

func (api *ApiServer) handleAuth(c *gin.Context) {
	sessionId, err := c.Cookie("session_id")

	fmt.Println(sessionId)

	if err != nil {
		c.String(400, "Session cookie not found")
		return
	}

	userId, err := api.store.GetSession(sessionId)

	if err != nil {
		c.String(500, "Error getting user")
		return
	}

	c.String(200, userId)
}

func (api *ApiServer) handlePostCreation(c *gin.Context) {

	userId, _ := c.Get("user_id")
	content := c.PostForm("content")

	if content == "" {
		c.String(400, "Content cannot be empty")
		return
	}

	postId, err := api.store.CreatePost(CreatePostRequest{
		UserID:  userId.(string),
		Content: content,
	})

	if err != nil {
		c.String(500, "Error creating post")
		return
	}

	c.String(200, postId)
}

func (api *ApiServer) handleUserUpdate(c *gin.Context) {
	userID, _ := c.Get("user_id")

	name := c.PostForm("name")
	bio := c.PostForm("bio")

	err := api.store.UpdateUser(userID.(string), name, bio)

	if err != nil {
		c.String(500, "Error updating user")
		return
	}

	c.String(200, "User updated")

}

func (api *ApiServer) handleLike(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Param("post_id")

	fmt.Println(userID, postID)

	err := api.store.LikePost(userID.(string), postID)

	if err != nil {
		fmt.Println(err)
		c.String(500, "Error liking post")
		return
	}

	c.String(200, "liked")
}

func (api *ApiServer) handleUnlike(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Param("post_id")

	err := api.store.UnlikePost(userID.(string), postID)

	if err != nil {
		c.String(500, "Error unliking post")
		return
	}

	c.String(200, "unliked")
}

func (api *ApiServer) handleGetLikes(c *gin.Context) {
	postID := c.PostForm("post_id")

	likes, err := api.store.PostLikes(postID)

	if err != nil {
		c.String(500, "Error getting likes")
		return
	}

	c.JSON(200, likes)
}

func (api *ApiServer) handleGetLikeStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Param("post_id")

	isLiked, err := api.store.IsLiked(userID.(string), postID)

	if err != nil {
		c.String(200, "false")
		return
	}

	if isLiked {
		c.String(200, "true")
		return
	}

	c.String(200, "false")
}

func (api *ApiServer) handleUserPosts(c *gin.Context) {
	userID, _ := c.Get("user_id")

	posts, err := api.store.GetPostsByUser(userID.(string))

	if err != nil {
		c.String(500, "Error getting posts")
		return
	}

	c.HTML(200, "post.html", gin.H{
		"Posts": posts,
	})
}

func (api *ApiServer) handleGetUser(c *gin.Context) {
	userID := c.Param("user_id")

	user, err := api.store.GetUser(userID)

	if err != nil {
		fmt.Println(err)
		c.String(500, "Error getting user")
		return
	}

	c.JSON(200, user)
}

func (api *ApiServer) handleGetFullPost(c *gin.Context) {
	postID := c.Param("post_id")

	post, err := api.store.GetFullPost(postID)
	if err != nil {
		fmt.Println(err)
		c.String(500, "Error getting post")
		return
	}

	c.JSON(200, post)
}
