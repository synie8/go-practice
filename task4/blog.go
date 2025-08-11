package task4

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}
type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
	User    User
}
type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}

var db *gorm.DB
var secret_key = []byte("123456789AAAAA")
var logger *logrus.Logger

// 自定义错误类型
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

// 统一错误响应函数
func ErrorResponse(c *gin.Context, code int, message string) {
	// 记录错误日志
	logger.WithFields(logrus.Fields{
		"status": code,
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
	}).Error(message)

	c.JSON(code, gin.H{
		"error": message,
	})
}

// 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 如果是自定义API错误
			if apiError, ok := c.Errors.Last().Err.(APIError); ok {
				// 直接返回错误
				c.JSON(apiError.Code, gin.H{
					"error": apiError.Message,
				})
				return
			}

			// 默认服务器错误
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
	}
}

func CreateDb() {
	var err error
	db, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		// 记录致命错误日志
		logger.Fatal("failed to connect database: ", err)
	}

	// 自动迁移模型
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	// 记录数据库初始化成功日志
	logger.Info("Database connected and migrated successfully")
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// 记录注册尝试日志
	logger.WithFields(logrus.Fields{
		"username": user.Username,
		"email":    user.Email,
	}).Info("User registration attempt")

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user.Password = string(hashedPassword)
	//保存用户信息
	if err := db.Create(&user).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// 记录注册成功日志
	logger.WithField("username", user.Username).Info("User registered successfully")

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// 记录登录尝试日志
	logger.WithField("username", user.Username).Info("User login attempt")

	var storedUser User
	if err := db.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid username or password")
		} else {
			ErrorResponse(c, http.StatusInternalServerError, "Database error")
		}
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       storedUser.ID,
		"username": storedUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secret_key)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// 记录登录成功日志
	logger.WithField("username", user.Username).Info("User logged in successfully")

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// VerifyToken 是一个中间件函数，用于验证JWT Token
func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于登录路由，不需要验证token
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" {
			c.Next()
			return
		}

		// 从请求头部获取Token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			ErrorResponse(c, http.StatusUnauthorized, "Token is missing")
			c.Abort() // 阻止请求继续传递
			return
		}

		// 解析Token并验证其有效性
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret_key, nil // 假设SecretKey是task4包中定义的密钥
		})

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					ErrorResponse(c, http.StatusBadRequest, "Invalid token format")
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					ErrorResponse(c, http.StatusUnauthorized, "Token is expired or not active yet")
				} else {
					ErrorResponse(c, http.StatusUnauthorized, "Token is invalid")
				}
			} else {
				ErrorResponse(c, http.StatusUnauthorized, "Error parsing token")
			}
			c.Abort() // 阻止请求继续传递
			return
		}

		if !token.Valid {
			ErrorResponse(c, http.StatusUnauthorized, "Token is invalid")
			c.Abort() // 阻止请求继续传递
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["id"].(float64)) // JSON数字默认是float64
		username := claims["username"].(string)

		// 创建UserInfo实例并放入上下文
		userInfo := User{}
		userInfo.ID = userID
		userInfo.Username = username
		c.Set("user_info", userInfo)

		// 记录认证成功日志
		logger.WithFields(logrus.Fields{
			"user_id":  userID,
			"username": username,
			"path":     c.Request.URL.Path,
			"method":   c.Request.Method,
		}).Info("User authenticated successfully")

		// 继续处理请求
		c.Next()
	}
}

/**
** 实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容
 */
func CreatePost(c *gin.Context) {

	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	var userInfo User
	// 验证用户身份
	if err := db.Debug().Where("id = ?", post.UserID).First(&userInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid user")
		} else {
			ErrorResponse(c, http.StatusInternalServerError, "Database error")
		}
		return
	}
	if post.Title == "" {
		ErrorResponse(c, http.StatusBadRequest, "Title不能为空")
		return
	}
	if post.Content == "" {
		ErrorResponse(c, http.StatusBadRequest, "Content不能为空")
		return
	}
	if err := db.Create(&post).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// 记录文章创建成功日志
	logger.WithFields(logrus.Fields{
		"post_id":   post.ID,
		"title":     post.Title,
		"author_id": post.UserID,
	}).Info("Post created successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
}

/**
** 实现文章的读取功能
** 获取所有文章列表
 */
func GetAllPost(c *gin.Context) {
	var posts []Post
	if err := db.Preload("User").Find(&posts).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	// 记录获取文章列表日志
	logger.WithField("count", len(posts)).Info("Retrieved all posts")

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

/**
** 获取单个文章的详细信息
 */
func GetPostById(c *gin.Context) {
	id := c.Param("id")
	var result Post
	if err := db.Preload("User").Where("id = ?", id).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, http.StatusNotFound, "Post not found")
		} else {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve post")
		}
		return
	}

	// 记录获取文章详情日志
	logger.WithFields(logrus.Fields{
		"post_id": id,
		"title":   result.Title,
	}).Info("Post retrieved successfully")

	c.JSON(http.StatusOK, gin.H{"post": result})
}

/*
*
实现文章的更新功能，只有文章的作者才能更新自己的文章
*/
func updatePost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	// 校验输入数据
	if post.Title == "" || post.Content == "" {
		ErrorResponse(c, http.StatusBadRequest, "Title and content are required")
		return
	}
	var dbPost Post
	if err := db.Where("id = ?", post.ID).First(&dbPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, http.StatusNotFound, "Post not found")
		} else {
			ErrorResponse(c, http.StatusInternalServerError, "Database query failed")
		}
		return
	}

	//校验是否是作者
	userInfo := c.MustGet("user_info").(User)
	if dbPost.UserID != userInfo.ID {
		ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// 执行更新操作
	if err := db.Model(&dbPost).Updates(post).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to update post")
		return
	}

	// 记录文章更新成功日志
	logger.WithFields(logrus.Fields{
		"post_id":   post.ID,
		"title":     post.Title,
		"author_id": userInfo.ID,
	}).Info("Post updated successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})

}

/**
** 实现文章的删除功能，只有文章的作者才能删除自己的文章
 */
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var dbPost Post
	if err := db.Where("id = ?", id).First(&dbPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, http.StatusNotFound, "Post not found")
		} else {
			ErrorResponse(c, http.StatusInternalServerError, "Database query failed")
		}
		return
	}

	//校验是否是作者
	userInfo := c.MustGet("user_info").(User)
	if dbPost.UserID != userInfo.ID {
		ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// 执行删除操作
	if err := db.Delete(&dbPost).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	// 记录文章删除成功日志
	logger.WithFields(logrus.Fields{
		"post_id":   id,
		"author_id": userInfo.ID,
	}).Info("Post deleted successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Post delete successfully"})

}

/**
*实现评论的创建功能，已认证的用户可以对文章发表评论
 */
func CreateComment(c *gin.Context) {
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request comment: "+err.Error())
		return
	}

	// 验证用户是否存在
	var user User
	if err := db.Where("id = ?", comment.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, http.StatusBadRequest, "User not found")
		} else {
			ErrorResponse(c, http.StatusInternalServerError, "Database error")
		}
		return
	}

	// 验证文章是否存在
	var post Post
	if err := db.Where("id = ?", comment.PostID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, http.StatusBadRequest, "Post not found")
		} else {
			ErrorResponse(c, http.StatusInternalServerError, "Database error")
		}
		return
	}

	if err := db.Create(&comment).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	// 记录评论创建成功日志
	logger.WithFields(logrus.Fields{
		"comment_id": comment.ID,
		"post_id":    comment.PostID,
		"user_id":    comment.UserID,
	}).Info("Comment created successfully")

	c.JSON(http.StatusOK, gin.H{"message": "comment create successfully"})

}

/**
*实现评论的读取功能，支持获取某篇文章的所有评论列表。
**/
func GetCommentsByPostId(c *gin.Context) {
	var comments []Comment
	postId := c.Param("id")
	if err := db.Preload("User").Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, http.StatusNotFound, "No comments found for this post")
		} else {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve comments")
		}
		return
	}

	// 记录获取评论列表日志
	logger.WithFields(logrus.Fields{
		"post_id": postId,
		"count":   len(comments),
	}).Info("Comments retrieved successfully")

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func main() {
	// 初始化日志记录器
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// 记录服务启动日志
	logger.Info("Starting blog service")

	r := gin.Default()

	// 注册全局错误处理中间件
	r.Use(ErrorHandler())

	// 注册全局中间件
	r.Use(VerifyToken())

	// 设置基础URL路径
	api := r.Group("/blog/api/v1")
	{
		// 登录路由,不需要Token验证
		api.POST("/login", Login)
		// 用户注册
		api.POST("/register", Register)

		// 受保护的路由
		api.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "This is a protected route"})
		})

		// 文章相关路由
		api.POST("/post/createPost", CreatePost)
		api.GET("/post/getAllPost", GetAllPost)
		api.GET("/post/getPostById/:id", GetPostById)
		api.PUT("/post/updatePost", updatePost)
		api.DELETE("/post/deletePost/:id", DeletePost)

		// 评论相关路由
		api.POST("/comment/createComment", CreateComment)
		api.GET("/comment/getCommentsByPostId/:id/comments", GetCommentsByPostId)
	}

	// 记录服务运行日志
	logger.Info("Blog service is running on :8080")

	r.Run(":8080")
}
