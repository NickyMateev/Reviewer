package job

import (
	"fmt"
	"github.com/NickyMateev/Reviewer/models"
	"testing"
)

func TestContainsWhenUserIsMissing(t *testing.T) {
	user := generateDummyUser(4)

	users := []*models.User{
		generateDummyUser(1),
		generateDummyUser(2),
		generateDummyUser(3),
	}

	if contains(users, user) != false {
		t.Errorf("user should not be part of user slice, but it is")
	}
}

func TestContainsWhenUserIsInSlice(t *testing.T) {
	user := generateDummyUser(4)

	users := []*models.User{
		generateDummyUser(1),
		generateDummyUser(2),
		generateDummyUser(3),
		user,
	}

	if contains(users, user) != true {
		t.Errorf("user should be part of user slice, but it is not")
	}
}

func TestFindIdlersWhenThereAreNoIdlers(t *testing.T) {
	user1 := generateDummyUser(1)
	user2 := generateDummyUser(2)
	user3 := generateDummyUser(3)

	requestedReviewers := []*models.User{
		user1,
		user2,
		user3,
	}

	actualReviewers := []*models.User{
		user1,
		user2,
		user3,
	}

	idlers := findIdlers(requestedReviewers, actualReviewers)

	if len(idlers) != len(requestedReviewers)-len(actualReviewers) {
		t.Errorf("(%d) idlers were found when there were suppossed to be none", len(idlers))
	}
}

func TestFindIdlersWhenThereAreIdlers(t *testing.T) {
	user1 := generateDummyUser(1)
	user2 := generateDummyUser(2)
	user3 := generateDummyUser(3)

	requestedReviewers := []*models.User{
		user1,
		user2,
		user3,
	}

	actualReviewers := []*models.User{
		user1,
	}

	idlers := findIdlers(requestedReviewers, actualReviewers)

	if len(idlers) != len(requestedReviewers)-len(actualReviewers) {
		t.Errorf("(%d) idlers were found when there were suppossed to be none", len(idlers))
	}
}

func generateDummyUser(id int64) *models.User {
	return &models.User{
		Username: fmt.Sprint("user", id),
		GithubID: id,
		Metadata: []byte{},
	}
}
