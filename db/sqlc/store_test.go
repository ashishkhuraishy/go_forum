package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFeeds(t *testing.T) {
	store := NewStore(testDB)

	var postIDs []int32

	// This will create 10 posts with one like each
	// and add it to the end of the db. So that later
	// when we query we can confirm that we are fetching
	// the last posts from the db
	for i := 0; i < 10; i++ {
		vote := createRandomVote(t)
		postIDs = append(postIDs, vote.PostID)
	}

	// We will skip the first 5, then return the
	// next 5 from the list
	args := ListPostsParams{
		Offset: 5,
		Limit:  5,
	}

	feeds, err := store.GetFeeds(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, feeds)
	require.Len(t, feeds, 5)

	// Check the returned feeds are in decsending order
	// by checking their id
	require.Greater(t, feeds[0].Post.ID, feeds[1].Post.ID)

	for i, feed := range feeds {
		require.NotEmpty(t, feed)

		// We created the post with one like so the resulting
		// feed should retrun 1 as the vote count
		require.EqualValues(t, 1, feed.Votes)

		// Check the first element is the 5th element (arr[4])
		// and so on for each element in the feed
		require.Equal(t, postIDs[4-i], feed.Post.ID)

		// Checking the username is same as the user who created
		// the post
		//
		// Note : Tested it earlier so ignoring `error`
		user, _ := store.GetUser(context.Background(), feed.Post.UserID)
		require.Equal(t, user.Username, feed.UserName)
	}

}

func Benchmark(b *testing.B) {
	store := NewStore(testDB)

	for i := 0; i < b.N; i++ {
		args := ListPostsParams{
			Offset: 50,
			Limit:  50,
		}

		store.GetFeeds(context.Background(), args)
	}
}
