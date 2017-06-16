package data

import (
	"blog.ka1em.site/model"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type PageRepository struct {
	DB *sql.DB
}

func (r *PageRepository) GetByPageID(pageGUID string, p *model.Page) error {
	sql := "select page_title,page_content,page_date from pages where page_guid=?"
	err := r.DB.QueryRow(sql, pageGUID).Scan(&p.Title, &p.Content, &p.Date)
	if err != nil {
		return err
	}

	return nil
}
