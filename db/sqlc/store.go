package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	DB *sql.DB
}

// NewStore returns a new store struct that
// can be accessed outside the `db` package
// and used to query the db
func NewStore(db *sql.DB) *Store {
	return &Store{
		DB:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)
	err = fn(queries)
	if err != nil {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			return fmt.Errorf("transaction error : %v \n rollback error : %v", err, rollBackErr)
		}

		return err
	}

	return tx.Commit()
}

type Feed struct {
	UserName string `json:"username"`
	Votes    int64  `json:"likes"`
	Post     Post   `json:"post"`
}

func (s *Store) GetFeeds(ctx context.Context, params ListPostsParams) ([]Feed, error) {
	var feeds []Feed

	err := s.execTx(ctx, func(q *Queries) error {

		posts, err := s.Queries.ListPosts(context.Background(), params)
		if err != nil {
			return err
		}

		for _, post := range posts {
			var feed Feed
			feed.Post = post

			user, err := s.Queries.GetUser(context.Background(), post.UserID)
			if err != nil {
				return err
			}

			votes, err := s.Queries.CountPostLikes(context.Background(), post.ID)
			if err != nil {
				return nil
			}

			feed.UserName = user.FullName
			feed.Votes = votes

			feeds = append(feeds, feed)
		}

		return nil
	})

	return feeds, err
}
