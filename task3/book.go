package task3

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Book struct {
	ID     uint64  `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

// 自定义表名
func (Book) TableName() string {
	return "book"
}

func (b *Book) Create() {
	gormDb := ConnectDb()
	gormDb.AutoMigrate(&Book{})
	var books = []Book{{Title: "何为人父", Author: "a1", Price: 30.89}, {Title: "人性的弱点", Author: "a2", Price: 50.66}, {Title: "斗破苍穹", Author: "a3", Price: 60.88}}

	gormDb.Debug().Create(&books)
}

func (b *Book) GetExpensiveBooks50(db *sqlx.DB) ([]Book, error) {
	var books []Book
	query := `SELECT id, title, author, price FROM book WHERE price > ? ORDER BY price DESC`

	err := db.Select(&books, query, 50)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %w", err)
	}
	fmt.Println("价格高于50元的书籍:")
	for _, book := range books {
		fmt.Printf("%d: %s (%s) - ￥%.2f\n",
			book.ID, book.Title, book.Author, book.Price)
	}
	return books, nil
}
