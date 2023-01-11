package constants

import (
	"server/models"

	"github.com/google/uuid"
)

var id, _ = uuid.FromBytes([]byte("fc2618e8-8b9e-4da7-a352-e2dddae595be"))
var Default_Measure_Unit = models.Measure_Unit{
	ID:   id,
	Name: "metric",
}
