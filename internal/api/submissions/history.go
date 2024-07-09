package submissions

import (
	"github.com/edulinq/autograder/internal/api/core"
	"github.com/edulinq/autograder/internal/db"
	"github.com/edulinq/autograder/internal/model"
)

type HistoryRequest struct {
	core.APIRequestAssignmentContext
	core.MinCourseRoleStudent

	TargetUser core.TargetCourseUserSelfOrGrader `json:"target-email"`
}

type HistoryResponse struct {
	FoundUser bool                           `json:"found-user"`
	History   []*model.SubmissionHistoryItem `json:"history"`
}

func HandleHistory(request *HistoryRequest) (*HistoryResponse, *core.APIError) {
	response := HistoryResponse{
		FoundUser: false,
		History:   make([]*model.SubmissionHistoryItem, 0),
	}

	if !request.TargetUser.Found {
		return &response, nil
	}

	response.FoundUser = true

	history, err := db.GetSubmissionHistory(request.Assignment, request.TargetUser.Email)
	if err != nil {
		return nil, core.NewInternalError("-603", &request.APIRequestCourseUserContext, "Failed to get submission history.").
			Err(err).Assignment(request.Assignment.GetID()).
			Add("target-user", request.TargetUser.Email)
	}

	response.History = history

	return &response, nil
}
