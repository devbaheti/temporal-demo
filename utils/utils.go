package utils

import (
	"errors"
	"temporal-demo/models"
)

func FindWorkflowByID(workflows []models.Workflow, id string) (*models.Workflow, error) {
	for _, workflow := range workflows {
		if workflow.ID == id {
			return &workflow, nil
		}
	}
	return nil, errors.New("workflow not found")
}
