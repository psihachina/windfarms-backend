package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/repository"
	"googlemaps.github.io/maps"
)

// WindfarmsService - ...
type WindfarmsService struct {
	repo repository.Windfarms
}

// NewWindfarmsService - ...
func NewWindfarmsService(repo repository.Windfarms) *WindfarmsService {

	return &WindfarmsService{repo: repo}
}

// Create - ...
func (s *WindfarmsService) Create(userID string, windfarm models.Windfarm) (string, error) {

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		return "", err
	}

	r := maps.ElevationRequest{}

	fmt.Println(windfarm)

	northEastLng, err := strconv.ParseFloat(strings.Split(windfarm.NorthEast, ",")[1], 64)
	if err != nil {
		return "", err
	}
	northEastLat, err := strconv.ParseFloat(strings.Split(windfarm.NorthEast, ",")[0], 64)
	if err != nil {
		return "", err
	}
	southWestLng, err := strconv.ParseFloat(strings.Split(windfarm.SouthWest, ",")[1], 64)
	if err != nil {
		return "", err
	}
	southWestLat, err := strconv.ParseFloat(strings.Split(windfarm.SouthWest, ",")[0], 64)
	if err != nil {
		return "", err
	}

	stepLng := (northEastLng - southWestLng) / 20
	fmt.Println("stepLng", stepLng)
	stepLat := (northEastLat - southWestLat) / 20
	fmt.Println("stepLat", stepLat)

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			r.Locations = append(r.Locations, maps.LatLng{
				Lat: southWestLat + float64(i)*stepLat,
				Lng: southWestLng + float64(j)*stepLng,
			})
		}
	}
	fmt.Println("OPS")
	elevationsResult, err := c.Elevation(context.Background(), &r)
	if err != nil {
		return "", err
	}

	avgElv := 0.0

	for i, element := range elevationsResult {
		fmt.Println(i, element.Elevation)
		avgElv += element.Elevation
	}

	avgElv = avgElv / 400

	fmt.Println(avgElv)

	windfarm.Altitude = avgElv

	return s.repo.Create(userID, windfarm)
}

// GetAll - ...
func (s *WindfarmsService) GetAll(userID string) ([]models.Windfarm, error) {
	return s.repo.GetAll(userID)
}

// GetByID - ...
func (s *WindfarmsService) GetByID(userID, windfarmID string) (models.Windfarm, error) {
	return s.repo.GetByID(userID, windfarmID)
}

// Delete - ...
func (s *WindfarmsService) Delete(userID, windfarmID string) error {
	return s.repo.Delete(userID, windfarmID)
}

// Update - ...
func (s *WindfarmsService) Update(userID string, itemID string, input models.UpdateWindfarmInput) error {
	return s.repo.Update(userID, itemID, input)
}
