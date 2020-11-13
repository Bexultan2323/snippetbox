package postgres

import (
	"aitu.com/snippetbox/pkg/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	insertSql                 = "INSERT INTO snippets (company,content,created,update,profits,founder,location,employees,capital) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id"
	getSnippetById            = "SELECT id, company, content, created, update,profits,founder,location,employees,capital FROM snippets where id=$1 AND update > now()"
	getLastTenCreatedSnippets = "SELECT id, company, content, created, update,profits,founder,location,employees,capital FROM snippets WHERE update > now() ORDER BY created DESC LIMIT 10"
	delete                    = "Delete from snippets where id = $1"
)

type SnippetModel struct {
	Pool *pgxpool.Pool
}

func (m *SnippetModel) Insert(company, content, created, update, profits, founder, location, employees, capital string) (int, error) {
	var id uint64
	row := m.Pool.QueryRow(context.Background(), insertSql, company, content, created, update, profits, founder, location, employees, capital)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}
	err := m.Pool.QueryRow(context.Background(), getSnippetById, id).
		Scan(&s.ID, &s.Company, &s.Content, &s.Created, &s.Update, &s.Profits, &s.Founder, &s.Location, &s.Employees, &s.Capital)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	snippets := []*models.Snippet{}
	rows, err := m.Pool.Query(context.Background(), getLastTenCreatedSnippets)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Company, &s.Content, &s.Created, &s.Update, &s.Profits, &s.Founder, &s.Location, &s.Employees, &s.Capital)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

func (m *SnippetModel) Delete(id int) error {
	err := m.Pool.QueryRow(context.Background(), delete, id).
		Scan()
	if err != nil {
		if err.Error() == "no rows in result set" {
			return models.ErrNoRecord
		} else {
			return err
		}
	}
	return nil
}
