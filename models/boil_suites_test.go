// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Projects", testProjects)
	t.Run("PullRequests", testPullRequests)
	t.Run("Users", testUsers)
}

func TestDelete(t *testing.T) {
	t.Run("Projects", testProjectsDelete)
	t.Run("PullRequests", testPullRequestsDelete)
	t.Run("Users", testUsersDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Projects", testProjectsQueryDeleteAll)
	t.Run("PullRequests", testPullRequestsQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Projects", testProjectsSliceDeleteAll)
	t.Run("PullRequests", testPullRequestsSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Projects", testProjectsExists)
	t.Run("PullRequests", testPullRequestsExists)
	t.Run("Users", testUsersExists)
}

func TestFind(t *testing.T) {
	t.Run("Projects", testProjectsFind)
	t.Run("PullRequests", testPullRequestsFind)
	t.Run("Users", testUsersFind)
}

func TestBind(t *testing.T) {
	t.Run("Projects", testProjectsBind)
	t.Run("PullRequests", testPullRequestsBind)
	t.Run("Users", testUsersBind)
}

func TestOne(t *testing.T) {
	t.Run("Projects", testProjectsOne)
	t.Run("PullRequests", testPullRequestsOne)
	t.Run("Users", testUsersOne)
}

func TestAll(t *testing.T) {
	t.Run("Projects", testProjectsAll)
	t.Run("PullRequests", testPullRequestsAll)
	t.Run("Users", testUsersAll)
}

func TestCount(t *testing.T) {
	t.Run("Projects", testProjectsCount)
	t.Run("PullRequests", testPullRequestsCount)
	t.Run("Users", testUsersCount)
}

func TestHooks(t *testing.T) {
	t.Run("Projects", testProjectsHooks)
	t.Run("PullRequests", testPullRequestsHooks)
	t.Run("Users", testUsersHooks)
}

func TestInsert(t *testing.T) {
	t.Run("Projects", testProjectsInsert)
	t.Run("Projects", testProjectsInsertWhitelist)
	t.Run("PullRequests", testPullRequestsInsert)
	t.Run("PullRequests", testPullRequestsInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("PullRequestToUserUsingAuthor", testPullRequestToOneUserUsingAuthor)
	t.Run("PullRequestToProjectUsingProject", testPullRequestToOneProjectUsingProject)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("ProjectToContributors", testProjectToManyContributors)
	t.Run("ProjectToPullRequests", testProjectToManyPullRequests)
	t.Run("PullRequestToApprovers", testPullRequestToManyApprovers)
	t.Run("PullRequestToCommenters", testPullRequestToManyCommenters)
	t.Run("PullRequestToIdlers", testPullRequestToManyIdlers)
	t.Run("PullRequestToReviewers", testPullRequestToManyReviewers)
	t.Run("UserToApprovedPullRequests", testUserToManyApprovedPullRequests)
	t.Run("UserToCommentedPullRequests", testUserToManyCommentedPullRequests)
	t.Run("UserToIdledPullRequests", testUserToManyIdledPullRequests)
	t.Run("UserToProjects", testUserToManyProjects)
	t.Run("UserToAuthoredPullRequests", testUserToManyAuthoredPullRequests)
	t.Run("UserToRequestedReviews", testUserToManyRequestedReviews)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("PullRequestToUserUsingAuthoredPullRequests", testPullRequestToOneSetOpUserUsingAuthor)
	t.Run("PullRequestToProjectUsingPullRequests", testPullRequestToOneSetOpProjectUsingProject)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("ProjectToContributors", testProjectToManyAddOpContributors)
	t.Run("ProjectToPullRequests", testProjectToManyAddOpPullRequests)
	t.Run("PullRequestToApprovers", testPullRequestToManyAddOpApprovers)
	t.Run("PullRequestToCommenters", testPullRequestToManyAddOpCommenters)
	t.Run("PullRequestToIdlers", testPullRequestToManyAddOpIdlers)
	t.Run("PullRequestToReviewers", testPullRequestToManyAddOpReviewers)
	t.Run("UserToApprovedPullRequests", testUserToManyAddOpApprovedPullRequests)
	t.Run("UserToCommentedPullRequests", testUserToManyAddOpCommentedPullRequests)
	t.Run("UserToIdledPullRequests", testUserToManyAddOpIdledPullRequests)
	t.Run("UserToProjects", testUserToManyAddOpProjects)
	t.Run("UserToAuthoredPullRequests", testUserToManyAddOpAuthoredPullRequests)
	t.Run("UserToRequestedReviews", testUserToManyAddOpRequestedReviews)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("ProjectToContributors", testProjectToManySetOpContributors)
	t.Run("PullRequestToApprovers", testPullRequestToManySetOpApprovers)
	t.Run("PullRequestToCommenters", testPullRequestToManySetOpCommenters)
	t.Run("PullRequestToIdlers", testPullRequestToManySetOpIdlers)
	t.Run("PullRequestToReviewers", testPullRequestToManySetOpReviewers)
	t.Run("UserToApprovedPullRequests", testUserToManySetOpApprovedPullRequests)
	t.Run("UserToCommentedPullRequests", testUserToManySetOpCommentedPullRequests)
	t.Run("UserToIdledPullRequests", testUserToManySetOpIdledPullRequests)
	t.Run("UserToProjects", testUserToManySetOpProjects)
	t.Run("UserToRequestedReviews", testUserToManySetOpRequestedReviews)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("ProjectToContributors", testProjectToManyRemoveOpContributors)
	t.Run("PullRequestToApprovers", testPullRequestToManyRemoveOpApprovers)
	t.Run("PullRequestToCommenters", testPullRequestToManyRemoveOpCommenters)
	t.Run("PullRequestToIdlers", testPullRequestToManyRemoveOpIdlers)
	t.Run("PullRequestToReviewers", testPullRequestToManyRemoveOpReviewers)
	t.Run("UserToApprovedPullRequests", testUserToManyRemoveOpApprovedPullRequests)
	t.Run("UserToCommentedPullRequests", testUserToManyRemoveOpCommentedPullRequests)
	t.Run("UserToIdledPullRequests", testUserToManyRemoveOpIdledPullRequests)
	t.Run("UserToProjects", testUserToManyRemoveOpProjects)
	t.Run("UserToRequestedReviews", testUserToManyRemoveOpRequestedReviews)
}

func TestReload(t *testing.T) {
	t.Run("Projects", testProjectsReload)
	t.Run("PullRequests", testPullRequestsReload)
	t.Run("Users", testUsersReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Projects", testProjectsReloadAll)
	t.Run("PullRequests", testPullRequestsReloadAll)
	t.Run("Users", testUsersReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Projects", testProjectsSelect)
	t.Run("PullRequests", testPullRequestsSelect)
	t.Run("Users", testUsersSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Projects", testProjectsUpdate)
	t.Run("PullRequests", testPullRequestsUpdate)
	t.Run("Users", testUsersUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Projects", testProjectsSliceUpdateAll)
	t.Run("PullRequests", testPullRequestsSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
}
