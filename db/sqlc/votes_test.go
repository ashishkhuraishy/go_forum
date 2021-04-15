package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	DUPLICATE_LIKE_ERROR = "pq: duplicate key value violates unique constraint \"user_post_index\""
)

func createRandomVote(t *testing.T) Vote {
	post := createRandomPost(t, 0)
	voted := true
	args := AddVoteParams{
		UserID: post.UserID,
		PostID: post.ID,
		Voted:  voted,
	}

	like, err := testQueries.AddVote(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, like)

	require.NotZero(t, like.ID)

	require.Equal(t, post.ID, like.PostID)
	require.Equal(t, post.UserID, like.UserID)
	require.Equal(t, args.Voted, like.Voted)

	return like
}

func TestAddVote(t *testing.T) {
	like := createRandomVote(t)

	// Should throw error if a like with the same userID
	// and postID exists
	args := AddVoteParams{
		UserID: like.UserID,
		PostID: like.PostID,
		Voted:  true,
	}

	like1, err := testQueries.AddVote(context.Background(), args)

	require.Empty(t, like1)
	require.Error(t, err)

	require.EqualErrorf(t, err, DUPLICATE_LIKE_ERROR, "formatted")

}

func TestGetVote(t *testing.T) {
	like := createRandomVote(t)

	// CASE 1
	// Check for a post with a userId which actually liked
	// the post, which should return valid data
	args := GetVoteInfoParams{
		UserID: like.UserID,
		PostID: like.PostID,
	}

	like1, err := testQueries.GetVoteInfo(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, like1)

	require.Equal(t, like.ID, like1.ID)
	require.Equal(t, like.Voted, like1.Voted)
	require.Equal(t, like.PostID, like1.PostID)
	require.Equal(t, like.UserID, like1.UserID)

	// CASE 2
	// Check with a user who hasnt liked the post yet
	user := createRandomUser(t)
	arg2 := GetVoteInfoParams{
		UserID: user.ID,
		PostID: like.PostID,
	}

	// If a userId who is not liked is queried, then
	// the db should return no rows error
	like2, err := testQueries.GetVoteInfo(context.Background(), arg2)
	require.Empty(t, like2)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

}

func TestUpdateVote(t *testing.T) {
	like := createRandomVote(t)

	args := UpdateVoteParams{
		PostID: like.PostID,
		UserID: like.UserID,
		Voted:  !like.Voted,
	}

	like1, err := testQueries.UpdateVote(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, like1)

	require.Equal(t, like.ID, like1.ID)
	require.Equal(t, like.PostID, like1.PostID)
	require.Equal(t, like.UserID, like1.UserID)

	require.Equal(t, args.Voted, like1.Voted)
}

func TestDeleteVote(t *testing.T) {
	like := createRandomVote(t)

	args := DeleteVoteParams{
		PostID: like.PostID,
		UserID: like.UserID,
	}

	err := testQueries.DeleteVote(context.Background(), args)
	require.NoError(t, err)

	arg1 := GetVoteInfoParams{
		UserID: like.UserID,
		PostID: like.PostID,
	}

	like1, err := testQueries.GetVoteInfo(context.Background(), arg1)

	require.Error(t, err)
	require.Empty(t, like1)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListVotesOfUser(t *testing.T) {
	// Create a user and add likes of that user
	// to random posts
	user := createRandomUser(t)

	// Creates 10 random posts and likes 5 posts
	// then dislikes 5 posts.
	for i := 0; i < 10; i++ {
		post := createRandomPost(t, 0)
		args := AddVoteParams{
			UserID: user.ID,
			PostID: post.ID,
			Voted:  i%2 == 0,
		}

		// Beacuse we already tested this fn [TestAddVote], we will skip
		// the return value of these
		testQueries.AddVote(context.Background(), args)
	}

	likes, err := testQueries.ListVotesOfUser(context.Background(), user.ID)

	// We created 10 likes from the user and 5 of
	// them have a status of `liked = true`. our fn
	// should only return those
	require.NoError(t, err)
	require.Len(t, likes, 5)

	for _, like := range likes {
		// Loop through all the likes and make sure
		// all the like we recived are from the same
		// are user and `liked` is `true` for every
		// one of them
		require.NotEmpty(t, like)
		require.Equal(t, user.ID, like.UserID)
		require.Equal(t, true, like.Voted)
	}

}

func TestCountPostVotes(t *testing.T) {
	// Create a random post and lets generated
	// random users like this post
	post := createRandomPost(t, 0)

	// Will create 10 random users and let
	// 5 users like the post and other 5
	// dislike the post
	for i := 0; i < 10; i++ {
		user := createRandomUser(t)

		args := AddVoteParams{
			UserID: user.ID,
			PostID: post.ID,
			Voted:  i%2 == 0,
		}

		// Already tested this fn... dont mind the return
		testQueries.AddVote(context.Background(), args)
	}

	count, err := testQueries.CountVotesOfPost(context.Background(), post.ID)

	// We liked the post by 5 users and disliked by 5
	// Our count luke should return a diffeence between
	// the sum of liked and disliked (like the upvote system
	// in reddit). Hence here the sum should be 0 (5 - 5)
	require.NoError(t, err)
	require.EqualValues(t, 0, count)

}
