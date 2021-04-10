package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ashishkhuraishy/go_forum/utils"
	"github.com/stretchr/testify/require"
)

// This function will return a random user every time its called
func createRandomUser(t *testing.T) User {
	lengthOfName := 5
	testName := utils.RandomString(lengthOfName)

	result, err := testQueries.CreateUser(context.Background(), testName)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, testName, result.FullName)
	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreatedAt)

	return result
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)
	args := UpdateUserParams{
		ID:       user1.ID,
		FullName: utils.RandomString(5),
	}

	user2, err := testQueries.UpdateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, args.FullName, user2.FullName)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.ID)

	require.NoError(t, err)

	user, err := testQueries.GetUser(context.Background(), user1.ID)
	require.Empty(t, user)
	require.NotEmpty(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListUsers(t *testing.T) {
	// We can make sure that atleast 10 users
	// will be available on our db
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	// This should return 5 `Users` from db
	// after skipping the first 5
	args := ListUsersParams{
		Offset: 5,
		Limit:  5,
	}

	users, err := testQueries.ListUsers(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		// Check the available users are empty
		require.NotEmpty(t, user)
	}
}
