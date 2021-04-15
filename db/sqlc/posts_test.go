package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ashishkhuraishy/go_forum/utils"
	"github.com/stretchr/testify/require"
)

func createRandomPost(t *testing.T, userID int32) Post {
	if userID == 0 {
		user := createRandomUser(t)
		userID = user.ID
	}
	args := CreatePostParams{
		UserID:  userID,
		Title:   utils.RandomString(5),
		Content: utils.RandomString(30),
	}

	post, err := testQueries.CreatePost(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, post)

	require.NotZero(t, post.ID)
	require.NotZero(t, post.CreatedAt)

	require.Equal(t, userID, post.UserID)
	require.Equal(t, args.Title, post.Title)
	require.Equal(t, args.Content, post.Content)

	return post
}

func TestCreatePost(t *testing.T) {
	createRandomPost(t, 0)
}

func TestGetPost(t *testing.T) {
	post1 := createRandomPost(t, 0)
	post2, err := testQueries.GetPost(context.Background(), post1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, post1.ID, post2.ID)
	require.Equal(t, post1.Title, post2.Title)
	require.Equal(t, post1.Content, post2.Content)
	require.Equal(t, post1.UserID, post2.UserID)
	require.Equal(t, post1.CreatedAt, post2.CreatedAt)
}

func TestListPosts(t *testing.T) {
	// We can make sure that atleast 10 posts
	// will be available on our db
	for i := 0; i < 10; i++ {
		createRandomPost(t, 0)
	}

	// This should return 5 `Posts` from db
	// after skipping the first 5
	args := ListPostsParams{
		Offset: 5,
		Limit:  5,
	}

	posts, err := testQueries.ListPosts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, posts, 5)

	// Check the id of first element is greater than the second one
	// This is to make sure that the returned posts is sorted in
	// descending order
	require.Greater(t, posts[0].ID, posts[1].ID)

	for _, post := range posts {
		// Check the available users are empty
		require.NotEmpty(t, post)
	}
}

func TestListPostsFromUser(t *testing.T) {
	// Create a random User
	user := createRandomUser(t)

	// We can make sure that atleast 10 posts
	// with the user we created will be available
	// on our db
	for i := 0; i < 10; i++ {
		createRandomPost(t, user.ID)
	}

	// This should return 5 `Posts` from db
	// with the user we specifies after
	// skipping the first 5
	args := ListPostsFromUserParams{
		UserID: user.ID,
		Offset: 5,
		Limit:  5,
	}

	posts, err := testQueries.ListPostsFromUser(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, posts, 5)

	// Check the id of first element is greater than the second one
	// This is to make sure that the returned posts is sorted in
	// descending order
	require.Greater(t, posts[0].ID, posts[1].ID)

	for _, post := range posts {
		// Check the available users are empty
		require.NotEmpty(t, post)
		require.Equal(t, user.ID, post.UserID)
	}
}

func TestUpdatePost(t *testing.T) {
	post := createRandomPost(t, 0)
	args := UpdatePostParams{
		ID:      post.ID,
		Title:   utils.RandomString(5),
		Content: utils.RandomString(100),
	}

	post2, err := testQueries.UpdatePost(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, post.ID, post2.ID)
	require.Equal(t, args.Title, post2.Title)
	require.Equal(t, args.Content, post2.Content)
	require.Equal(t, post.UserID, post2.UserID)
	require.Equal(t, post.CreatedAt, post2.CreatedAt)
}

func TestDeletePost(t *testing.T) {
	post := createRandomPost(t, 0)

	err := testQueries.DeletePost(context.Background(), post.ID)
	require.NoError(t, err)

	post1, err := testQueries.GetPost(context.Background(), post.ID)

	require.Error(t, err)
	require.Empty(t, post1)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
