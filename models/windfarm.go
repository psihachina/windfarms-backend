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

type Point struct {
	Longitude float64 `json:"lat" db:"latitude"`
	Latitude  float64 `json:"lng" db:"longitude"`
}

type UpdateWindfarmInput struct {
	Name            *string  `json:"windfarmName"`
	Longitude       *float64 `json:"windfarmLongitude"`
	Latitude        *float64 `json:"windfarmLatitude"`
	Capacity        *float64 `json:"windfarmCapacity"`
	RangeToCity     *float64 `json:"rangeToCity"`
	RangeToRoad     *float64 `json:"rangeToRoad"`
	RangeToCityLine *float64 `json:"rangeToCityLine"`
	CityLatitude    *float64 `json:"cityLongitude"`
	CityLongitude   *float64 `json:"cityLatitude"`
	Description     *string  `json:"windfarmDescription"`
}

func (i UpdateWindfarmInput) Validate() error {
	if i.Name == nil && i.Longitude == nil &&
		i.Latitude == nil && i.Capacity == nil &&
		i.RangeToCity == nil && i.RangeToRoad == nil &&
		i.RangeToCityLine == nil && i.CityLatitude == nil &&
		i.CityLongitude == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
