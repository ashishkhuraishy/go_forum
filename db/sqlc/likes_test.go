package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	DUPLICATE_LIKE_ERROR = "pq: duplicate key value violates unique constraint \"post_user_unique\""
)

func createRandomLike(t *testing.T) Like {
	post := createRandomPost(t, 0)
	liked := true
	args := AddLikeParams{
		UserID: post.UserID,
		PostID: post.ID,
		Liked:  liked,
	}

	like, err := testQueries.AddLike(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, like)

	require.NotZero(t, like.ID)
	require.NotZero(t, like.LikedAt)

	require.Equal(t, post.ID, like.PostID)
	require.Equal(t, post.UserID, like.UserID)
	require.Equal(t, args.Liked, like.Liked)

	return like
}

func TestAddLike(t *testing.T) {
	like := createRandomLike(t)

	// Should throw error if a like with the same userID
	// and postID exists
	args := AddLikeParams{
		UserID: like.UserID,
		PostID: like.PostID,
		Liked:  true,
	}

	like1, err := testQueries.AddLike(context.Background(), args)

	require.Empty(t, like1)
	require.Error(t, err)

	require.EqualErrorf(t, err, DUPLICATE_LIKE_ERROR, "formatted")

}

func TestGetLike(t *testing.T) {
	like := createRandomLike(t)

	// CASE 1
	// Check for a post with a userId which actually liked
	// the post, which should return valid data
	args := GetLikeInfoParams{
		UserID: like.UserID,
		PostID: like.PostID,
	}

	like1, err := testQueries.GetLikeInfo(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, like1)

	require.Equal(t, like.ID, like1.ID)
	require.Equal(t, like.Liked, like1.Liked)
	require.Equal(t, like.LikedAt, like1.LikedAt)
	require.Equal(t, like.PostID, like1.PostID)
	require.Equal(t, like.UserID, like1.UserID)

	// CASE 2
	// Check with a user who hasnt liked the post yet
	user := createRandomUser(t)
	arg2 := GetLikeInfoParams{
		UserID: user.ID,
		PostID: like.PostID,
	}

	// If a userId who is not liked is queried, then
	// the db should return no rows error
	like2, err := testQueries.GetLikeInfo(context.Background(), arg2)
	require.Empty(t, like2)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

}

func TestUpdateLike(t *testing.T) {
	like := createRandomLike(t)

	args := UpdateLikeParams{
		PostID: like.PostID,
		UserID: like.UserID,
		Liked:  !like.Liked,
	}

	like1, err := testQueries.UpdateLike(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, like1)

	require.Equal(t, like.ID, like1.ID)
	require.Equal(t, like.LikedAt, like1.LikedAt)
	require.Equal(t, like.PostID, like1.PostID)
	require.Equal(t, like.UserID, like1.UserID)

	require.Equal(t, args.Liked, like1.Liked)
}

func TestDeleteLike(t *testing.T) {
	like := createRandomLike(t)

	args := DeleteLikeParams{
		PostID: like.PostID,
		UserID: like.UserID,
	}

	err := testQueries.DeleteLike(context.Background(), args)
	require.NoError(t, err)

	arg1 := GetLikeInfoParams{
		UserID: like.UserID,
		PostID: like.PostID,
	}

	like1, err := testQueries.GetLikeInfo(context.Background(), arg1)

	require.Error(t, err)
	require.Empty(t, like1)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListLikesOfUser(t *testing.T) {
	// Create a user and add likes of that user
	// to random posts
	user := createRandomUser(t)

	// Creates 10 random posts and likes 5 posts
	// then dislikes 5 posts.
	for i := 0; i < 10; i++ {
		post := createRandomPost(t, 0)
		args := AddLikeParams{
			UserID: user.ID,
			PostID: post.ID,
			Liked:  i%2 == 0,
		}

		// Beacuse we already tested this fn [TestAddLike], we will skip
		// the return value of these
		testQueries.AddLike(context.Background(), args)
	}

	likes, err := testQueries.ListLikesOfUser(context.Background(), user.ID)

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
		require.Equal(t, true, like.Liked)
	}

}

func TestCountPostLikes(t *testing.T) {
	// Create a random post and lets generated
	// random users like this post
	post := createRandomPost(t, 0)

	// Will create 10 random users and let
	// 5 users like the post and other 5
	// dislike the post
	for i := 0; i < 10; i++ {
		user := createRandomUser(t)

		args := AddLikeParams{
			UserID: user.ID,
			PostID: post.ID,
			Liked:  i%2 == 0,
		}

		// Already tested this fn... dont mind the return
		testQueries.AddLike(context.Background(), args)
	}

	count, err := testQueries.CountPostLikes(context.Background(), post.ID)

	// We liked the post by 5 users and disliked by 5
	// Our count luke should return a diffeence between
	// the sum of liked and disliked (like the upvote system
	// in reddit). Hence here the sum should be 0 (5 - 5)
	require.NoError(t, err)
	require.EqualValues(t, 0, count)

}
