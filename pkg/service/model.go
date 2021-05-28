package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/repository"
	"googlemaps.github.io/maps"
)

type ModelsService struct {
	repo            repository.Models
	turbinesService Turbines
	windsService    Winds
}

func NewModelsService(repo repository.Models, turbinesService Turbines, windsService Winds) *ModelsService {
	return &ModelsService{repo: repo, turbinesService: turbinesService, windsService: windsService}
}

// Create - ...
func (s *ModelsService) Create(userID string, windfarmID string, model models.Model) (string, error) {

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	winds, err := s.windsService.GetAll(userID, windfarmID)
	if err != nil {
		log.Fatal(err)
	}

	windMap := make(map[string]map[string]map[string]models.Wind)

	altitudes := []float64{15, 30, 50, 75, 100, 150, 200}

	for _, item := range winds {
		if windMap[item.Date] == nil {
			windMap[item.Date] = map[string]map[string]models.Wind{}
		}
		if windMap[item.Date][item.Time] == nil {
			windMap[item.Date][item.Time] = map[string]models.Wind{}
		}
		windMap[item.Date][item.Time][fmt.Sprint(item.Altitude)] = item
	}

	workerCount := len(model.Turbines)

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		turbine, err := s.turbinesService.GetByID(userID, model.Turbines[i].TurbineName)
		if err != nil {
			log.Fatal(err)
		}

		r := &maps.ElevationRequest{
			Locations: []maps.LatLng{
				{Lat: model.Turbines[i].Longitude, Lng: model.Turbines[i].Latitude},
			},
		}

		elvsRes, err := c.Elevation(context.Background(), r)
		if err != nil {
			log.Fatal(err)
		}

		height := model.WindfarmAltitude - elvsRes[0].Elevation + turbine.TowerHeight
		fmt.Println(model.WindfarmAltitude, elvsRes[0].Elevation, height)

		wg.Add(1)

		go func(i int, turbine models.Turbine) {
			defer wg.Done()

			for date := range windMap {
				for time := range windMap[date] {
					minAlt, maxAlt, err := getClosestPair(altitudes, height)
					if err != nil {
						log.Fatal(err)
					}
					windSpeed := linearApproximation(height, minAlt, maxAlt, windMap[date][time][fmt.Sprint(maxAlt)].WindSpeed, windMap[date][time][fmt.Sprint(minAlt)].WindSpeed)
					if minAlt == maxAlt {
						windSpeed = windMap[date][time][fmt.Sprint(maxAlt)].WindSpeed
					}

					outputMap := make(map[string]models.Output)

					for _, item := range turbine.Outputs {
						outputMap[fmt.Sprint(item.Speed)] = item
					}

					var speed []float64

					for k := range outputMap {
						key, err := strconv.ParseFloat(k, 64)
						if err != nil {
							log.Fatal(err)
						}
						speed = append(speed, key)
					}

					sort.Slice(speed, func(i, j int) bool {
						return speed[i] < speed[j]
					})

					minSpeed, maxSpeed, err := getClosestPair(speed, windSpeed)
					if err != nil {
						log.Fatal(err)
					}

					output := linearApproximation(windSpeed, minSpeed, maxSpeed, outputMap[fmt.Sprint(maxSpeed)].Production, outputMap[fmt.Sprint(minSpeed)].Production)

					if math.IsNaN(output) {
						output = 0
					}

					production := models.Production{
						Time:           time,
						Date:           date,
						Value:          output,
						ICUF:           output / turbine.MaximumPower,
						WindSpeed:      windSpeed,
						Altitude:       height,
						TurbineModelID: model.Turbines[i].TurbineModelID,
					}

					model.Turbines[i].Productions = append(model.Turbines[i].Productions, &production)
				}
			}
		}(i, turbine)
	}
	wg.Wait()
	return s.repo.Create(userID, windfarmID, model)
}

// GetAll - ...
func (s *ModelsService) GetAll(userID, windfarmID string) ([]models.Model, error) {
	return s.repo.GetAll(userID, windfarmID)
}

// GetMapData - ...
func (s *ModelsService) GetMapData(userID, windfarmID string, modelID string) (models.ModelMap, error) {
	model, err := s.repo.GetByID(userID, windfarmID, modelID)

	if err != nil {
		return models.ModelMap{}, err
	}

	modelMap := models.ModelMap{}

	for _, turbine := range model.Turbines {
		for _, production := range turbine.Productions {

			if modelMap.Production == nil {
				modelMap.Production = map[string]map[string]map[string]models.Production{}
			}
			if modelMap.Production[production.Date] == nil {
				modelMap.Production[production.Date] = map[string]map[string]models.Production{}
			}
			if modelMap.Production[production.Date][production.Time] == nil {
				modelMap.Production[production.Date][production.Time] = map[string]models.Production{}
			}
			modelMap.Production[production.Date][production.Time][turbine.TurbineModelID] = *production
		}
	}

	return modelMap, err
}

// GetByID - ...
func (s *ModelsService) GetByID(userID, windfarmID, modelID string) (models.Model, error) {
	return s.repo.GetByID(userID, windfarmID, modelID)
}

// Delete - ...
func (s *ModelsService) Delete(userID, windfarmID, modelID string) error {
	return s.repo.Delete(userID, windfarmID)
}

// Update - ...
// func (s *ModelsService) Update(userID, windfarmID, modelID string, input models.UpdateModelInput) error {
// 	return s.repo.Update(userID, windfarmID, modelID, input)
// }

func getClosestPair(arr []float64, number float64) (float64, float64, error) {
	var res_l, res_r int

	l := 0
	r := len(arr) - 1

	if number < 2 {
		return 2, 2, nil
	}

	for {
		if l > len(arr)-1 {
			res_l = len(arr) - 1
			break
		} else {
			if arr[l] > number {
				res_l = l
				break
			}
			l++
		}
	}

	for {
		if r < 0 {
			res_r = 0
			break
		} else {
			if arr[r] < number {
				res_r = r
				break
			}
			r--
		}
	}

	return arr[res_l], arr[res_r], nil
}

func linearApproximation(x, xmin, xmax, ymax, ymin float64) float64 {
	return ymin + (((x - xmin) / (xmax - xmin)) * (ymax - ymin))
}
