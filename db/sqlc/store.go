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
	// Voted will indicate if a post is voted by the logged
	// in user or not. It will have 3 possible value
	// `0`  : If user havent responded to the post
	// `1`  : If user have upvoted the post
	// `-1` : If user has downvoted the post
	Voted int8 `josn:"voted"`
	Post  Post `json:"post"`
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

			// userchan := make(chan User)
			// votesChan := make(chan Vote)

			user, err := s.Queries.GetUser(context.Background(), post.UserID)
			if err != nil {
				return err
			}

			votes, err := s.Queries.CountVotesOfPost(context.Background(), post.ID)
			if err != nil {
				return err
			}

			// args := GetVoteInfoParams{
			// 	UserID: post.UserID,
			// 	PostID: post.ID,
			// }

			// vote, err := s.Queries.GetVoteInfo(context.Background(), args)
			// if err != nil {
			// 	if err != sql.ErrNoRows {
			// 		return err
			// 	}
			// 	feed.Voted = 0
			// }

			// if vote.Voted {
			// 	feed.Voted = 1
			// } else {
			// 	feed.Voted = -1
			// }

			feed.UserName = user.Username
			feed.Votes = votes

			feeds = append(feeds, feed)
		}

		return nil
	})

	return feeds, err
}
