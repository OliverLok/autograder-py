package admin

import (
	"testing"

	"github.com/edulinq/autograder/internal/api/core"
	"github.com/edulinq/autograder/internal/db"
	"github.com/edulinq/autograder/internal/model"
)

func TestUpdateCourse(test *testing.T) {
	// Remove a user and ensure the user is back after a reload.

	// Leave the course in a good state after the test.
	db.ResetForTesting()
	defer db.ResetForTesting()

	course := db.MustGetTestCourse()

	count := countUsers(test, course)
	if count != 5 {
		test.Fatalf("Unexpected pre-remove user count. Expected 5, found %d.", count)
	}

	exists, enrolled, err := db.RemoveUserFromCourse(course, "course-student@test.edulinq.org")
	if err != nil {
		test.Fatalf("Error when removing the user: '%v'.", err)
	}

	if !exists {
		test.Fatalf("User does not exist.")
	}

	if !enrolled {
		test.Fatalf("User was not enrolled in course.")
	}

	count = countUsers(test, course)
	if count != 4 {
		test.Fatalf("Unexpected post-remove user count. Expected 4, found %d.", count)
	}

	reloadRequest(test)

	count = countUsers(test, course)
	if count != 5 {
		test.Fatalf("Unexpected post-reload user count. Expected 5, found %d.", count)
	}
}

func reloadRequest(test *testing.T) {
	response := core.SendTestAPIRequest(test, core.NewEndpoint(`courses/admin/update`), nil)
	if !response.Success {
		test.Errorf("Response is not a success when it should be: '%v'.", response)
	}
}

func countUsers(test *testing.T, course *model.Course) int {
	users, err := db.GetCourseUsers(course)
	if err != nil {
		test.Fatalf("Failed to get users: '%v'.", err)
	}

	return len(users)
}
