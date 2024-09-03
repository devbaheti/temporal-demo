package utils

import (
	"errors"
	"fmt"
	"temporal-demo/models"
)

func FindWorkflowByID(workflows []models.Workflow, id string) (*models.Workflow, error) {
	for _, workflow := range workflows {
		if workflow.Name == id {
			return &workflow, nil
		}
	}
	return nil, errors.New("workflow not found")
}

// ConvertInterfaceToString attempts to convert an interface{} to a string
func ConvertInterfaceToString(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}
