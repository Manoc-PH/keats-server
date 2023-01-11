package models

import (
	"github.com/google/uuid"
)

type Measure_Unit struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Wt_Solid_Unit_Small   string    `json:"wt_solid_unit_small"`
	Wt_Solid_Desc_Small   string    `json:"wt_solid_desc_small"`
	Wt_Solid_Unit_Medium  string    `json:"wt_solid_unit_medium"`
	Wt_Solid_Desc_Medium  string    `json:"wt_solid_desc_medium"`
	Wt_Solid_Unit_Large   string    `json:"wt_solid_unit_large"`
	Wt_Solid_Desc_Large   string    `json:"wt_solid_desc_large"`
	Wt_Liquid_Unit_Small  string    `json:"wt_liquid_unit_small"`
	Wt_Liquid_Desc_Small  string    `json:"wt_liquid_desc_small"`
	Wt_Liquid_Unit_Medium string    `json:"wt_liquid_unit_medium"`
	Wt_Liquid_Desc_Medium string    `json:"wt_liquid_desc_medium"`
	Wt_Liquid_Unit_Large  string    `json:"wt_liquid_unit_large"`
	Wt_Liquid_Desc_Large  string    `json:"wt_liquid_desc_large"`
}
