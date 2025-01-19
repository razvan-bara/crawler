package crawlerDb

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Article struct {
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
}

func InsertArticle(db *sql.DB, article Article) error {
	query := `
		INSERT INTO articles (title, authors)
		VALUES ($1, $2)
		RETURNING id;
	`

	var id int
	err := db.QueryRow(query, article.Title, pq.Array(article.Authors)).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to insert article: %w", err)
	}

	fmt.Printf("Article inserted with ID: %d\n", id)
	return nil
}

func BulkInsertArticles(db *sql.DB, articles []*Article) error {
	query := `
		INSERT INTO articles (title, authors)
		VALUES %s
	`

	values := []interface{}{}
	placeholders := ""
	for i, article := range articles {
		placeholder := fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
		placeholders += placeholder
		if i < len(articles)-1 {
			placeholders += ", "
		}
		values = append(values, article.Title, pq.Array(article.Authors))
	}

	finalQuery := fmt.Sprintf(query, placeholders)

	_, err := db.Exec(finalQuery, values...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert articles: %w", err)
	}

	return nil
}

// DeleteAllArticles deletes all articles from the database
func DeleteAllArticles(db *sql.DB) error {
	query := `DELETE FROM articles`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete all articles: %w", err)
	}

	return nil
}
