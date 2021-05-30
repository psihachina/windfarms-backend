package models

import (
	"errors"

	"github.com/jackc/pgx/pgtype"
)

type Windfarm struct {
	WindfarmID    string         `json:"windfarmId,omitempty" db:"windfarm_id"`
	WindfarmName  string         `json:"windfarmName" db:"windfarm_name" validate:"required"`
	PolygonDB     pgtype.Polygon `db:"polygon"`
	PolygonJSON   string         `json:"polygon"`
	NorthEast     string         `json:"northEast"`
	SouthWest     string         `json:"southWest"`
	Longitude     float64        `json:"windfarmLongitude" db:"longitude"`
	Latitude      float64        `json:"windfarmLatitude" db:"latitude"`
	Capacity      float64        `json:"windfarmCapacity" db:"capacity"`
	PolygonRadius float64        `json:"polygonRadius" db:"polygon_radius"`
	Description   string         `json:"windfarmDescription" db:"description"`
	Altitude      float64        `json:"altitude" db:"altitude"`
}

type UpdateWindfarmInput struct {
	Name        *string  `json:"windfarm_name"`
	Longitude   *float64 `json:"longitude"`
	Latitude    *float64 `json:"latitude"`
	Capacity    *float64 `json:"capacity"`
	Description *string  `json:"description"`
}

func (i UpdateWindfarmInput) Validate() error {
	if i.Name == nil && i.Longitude == nil &&
		i.Latitude == nil && i.Capacity == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
