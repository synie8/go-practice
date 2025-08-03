package task3

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	ID      uint64
	Name    string
	Posts   []Post
	PostNum int
}
type Post struct {
	ID         uint64
	Text       string
	UserID     uint64
	Comments   []Comment
	CommentNum int
}
type Comment struct {
	ID     uint64
	Text   string
	UserID uint64 //评论人
	PostID uint64
}

func (u *User) Create() {
	db := ConnectDb()
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	var users []User = []User{{Name: "a1"}, {Name: "a2"}, {Name: "a3"}}
	db.Debug().Create(&users)

	var posts []Post = []Post{{Text: "a1", UserID: 1}, {Text: "a2", UserID: 1}, {Text: "a3", UserID: 2}}
	db.Debug().Create(&posts)

	var comments []Comment = []Comment{{Text: "ss", PostID: 1, UserID: 2}, {Text: "22", PostID: 2, UserID: 2}, {Text: "33", PostID: 3, UserID: 2}, {Text: "66", PostID: 3, UserID: 2}}
	db.Debug().Create(&comments)
}

func (u *User) Query(userID uint64) {
	db := ConnectDb()
	err := db.Preload("Posts.Comments").First(&u, userID).Error
	if err != nil {
		fmt.Errorf("查询失败: %v", err)
	}
	fmt.Println(u)
}
func (u *User) Query2() {
	db := ConnectDb()
	var post Post
	query := db.Table("comments").Select("COUNT(*) as ct,comments.post_id as post_id ").Group("comments.post_id")
	db.Debug().Table("posts").Model(&Post{}).Joins("join (?) q on posts.id = q.post_id", query).Order("q.ct desc").First(&post)
	fmt.Println(post)
}
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 在文章创建时自动更新用户的文章数量统计字段
	tx.Debug().Table("users as u").Where("id = ?", p.UserID).Update("post_num", tx.Table("posts as p").Select("COUNT(*) as ct").Group("p.user_id").Where("p.user_id = u.id"))
	return nil
}
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	// 在评论删除时检查文章的评论数量,更新评论量
	tx.Debug().Table("posts as p").Where("id = ?", c.PostID).Update("comment_num", tx.Table("comments as c").Select("COUNT(*) as ct").Group("c.post_id").Where("c.post_id = ?", c.PostID))
	return nil
}
