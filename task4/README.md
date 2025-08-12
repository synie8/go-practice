# Go Blog API

一个使用 Go 语言开发的博客系统后端 API，支持用户注册登录、文章管理、评论功能等。

## 功能特性

- 用户注册和登录（JWT 认证）
- 文章的创建、查看、更新和删除
- 评论的创建和查看
- 基于角色的访问控制（只有作者可以修改自己的文章）
- 完整的日志记录系统
- SQLite 数据库存储

## 运行环境要求

- Go 1.16 或更高版本
- SQLite3 数据库（内置，无需额外安装）

## 依赖库

- Web框架: `github.com/gin-gonic/gin`
- ORM工具: `gorm.io/gorm` 、 `gorm.io/driver/sqlite`
- JWT支持: `github.com/dgrijalva/jwt-go`
- 密码加密: `golang.org/x/crypto/bcrypt`
- 日志记录: `github.com/sirupsen/logrus`

## 安装和设置

1. 克隆项目到本地：
   ```bash
   git clone <repository-url>
   ```

2. 进入项目目录：
   ```bash
   cd task4
   ```

3. 安装依赖包：
   ```bash
   go mod tidy
   ```

4. 运行项目：
   ```bash
   go run blog.go
   ```

服务将在 `localhost:8080` 上启动。

## API 接口文档

### 用户认证

#### 用户注册
```
POST /blog/api/v1/register
```

请求体：
```json
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```

#### 用户登录
```
POST /blog/api/v1/login
```

请求体：
```json
{
  "username": "testuser",
  "password": "password123"
}
```

响应：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 文章管理

需要在请求头中添加 `Authorization: <token>` 字段。

#### 创建文章
```
POST /blog/api/v1/post/createPost
```

请求体：
```json
{
  "title": "文章标题",
  "content": "文章内容",
  "user_id": 1
}
```

#### 获取所有文章
```
GET /blog/api/v1/post/getAllPost
```

#### 根据ID获取文章
```
GET /blog/api/v1/post/getPostById/:id
```

#### 更新文章
```
PUT /blog/api/v1/post/updatePost
```

请求体：
```json
{
  "ID": 1,
  "title": "更新后的标题",
  "content": "更新后的内容"
}
```

#### 删除文章
```
DELETE /blog/api/v1/post/deletePost/:id
```

### 评论管理

需要在请求头中添加 `Authorization: <token>` 字段。

#### 创建评论
```
POST /blog/api/v1/comment/createComment
```

请求体：
```json
{
  "content": "评论内容",
  "user_id": 1,
  "post_id": 1
}
```

#### 获取文章的所有评论
```
GET /blog/api/v1/comment/getCommentsByPostId/:id/comments
```

## 使用 Postman 进行测试

1. 安装 Postman
2. 创建一个新的请求
3. 选择相应的 HTTP 方法（GET、POST、PUT、DELETE）
4. 输入对应的 URL
5. 在 Headers 中添加 `Content-Type: application/json`
6. 对于需要认证的接口，在 Headers 中添加 `Authorization: <token>`
7. 在 Body 中选择 raw 和 JSON 格式，然后输入请求体
8. 点击 Send 发送请求

### 测试用例示例

1. **用户注册**
   - 方法: POST
   - URL: http://localhost:8080/blog/api/v1/register
   - Body:
     ```json
     {
       "username": "testuser",
       "password": "password123",
       "email": "test@example.com"
     }
     ```
   - 预期结果: 状态码 201，返回 "User registered successfully"

2. **用户登录**
   - 方法: POST
   - URL: http://localhost:8080/blog/api/v1/login
   - Body:
     ```json
     {
       "username": "testuser",
       "password": "password123"
     }
     ```
   - 预期结果: 状态码 200，返回 token

3. **创建文章**
   - 方法: POST
   - URL: http://localhost:8080/blog/api/v1/post/createPost
   - Headers: 
     - Authorization: <上一步获取的token>
   - Body:
     ```json
     {
       "title": "Test Post",
       "content": "This is a test post",
       "user_id": 1
     }
     ```
   - 预期结果: 状态码 200，返回 "Post created successfully"

4. **获取所有文章**
   - 方法: GET
   - URL: http://localhost:8080/blog/api/v1/post/getAllPost
   - 预期结果: 状态码 200，返回文章列表

5. **创建评论**
   - 方法: POST
   - URL: http://localhost:8080/blog/api/v1/comment/createComment
   - Headers:
     - Authorization: <token>
   - Body:
     ```json
     {
       "content": "This is a test comment",
       "user_id": 1,
       "post_id": 1
     }
     ```
   - 预期结果: 状态码 200，返回 "comment create successfully"

6. **获取文章评论**
   - 方法: GET
   - URL: http://localhost:8080/blog/api/v1/comment/getCommentsByPostId/1/comments
   - 预期结果: 状态码 200，返回评论列表

## 注意事项

1. token 在登录成功后返回，有效期为 24 小时
2. 数据库文件 `blog.db` 会在首次运行时自动创建
3. 用户密码在数据库中以加密形式存储
4. 文章和评论的创建、更新、删除操作会记录在日志中