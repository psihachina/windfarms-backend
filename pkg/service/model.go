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

func (s *ModelsService) CreateModel(userID string, windfarmID string, model models.Model) (string, error) {
	return s.repo.CreateModel(userID, windfarmID, model)
}

// GenerateModel - ...
func (s *ModelsService) GenerateModel(userID, windfarmID, modelID string, model models.Model) (string, error) {

	turbines, err := s.turbinesService.GetMap(userID)
	if err != nil {
		log.Fatal(err)
	}

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	winds, err := s.windsService.GetAll(userID, windfarmID)
	if err != nil {
		log.Fatal(err)
	}

	windMap := make(map[string]map[string]map[string]models.Wind)

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
		turbine := turbines[model.Turbines[i].TurbineName]

		elevations, err := getElevation(c, model.Turbines[i].Latitude, model.Turbines[i].Longitude)
		if err != nil {
			log.Fatal(err)
		}

		height := model.WindfarmAltitude - elevations[0].Elevation + turbine.TowerHeight

		wg.Add(1)

		go func(i int, turbine models.Turbine) {
			defer wg.Done()

			var shadingTurbines []shadingTurbine

			for j := 0; j < len(model.Turbines); j++ {
				if i != j && checkInclude(model.Turbines[i], model.Turbines[j], 3*turbines[model.Turbines[j].TurbineName].RotorDiameter) {
					shadingTurbines = append(shadingTurbines, shadingTurbine{
						distance: getDistance(
							maps.LatLng{
								Lat: model.Turbines[i].Latitude,
								Lng: model.Turbines[i].Longitude,
							},
							maps.LatLng{
								Lat: model.Turbines[j].Latitude,
								Lng: model.Turbines[j].Longitude,
							}),

						radius:  3 * turbines[model.Turbines[j].TurbineName].RotorDiameter,
						turbine: model.Turbines[j],
					})
				}
			}

			for date := range windMap {
				for time := range windMap[date] {

					windSpeed, err := getWindSpeed(height, windMap[date][time])
					if err != nil {
						log.Fatal(err)
					}

					windDirection, err := getWindDirection(height, windMap[date][time])
					if err != nil {
						log.Fatal(err)
					}

					var shading float64

					for j := 0; j < len(shadingTurbines); j++ {
						alpha := windDirection - math.Atan2(
							model.Turbines[i].Latitude-shadingTurbines[j].turbine.Latitude,
							model.Turbines[i].Longitude-shadingTurbines[j].turbine.Longitude,
						)
						shading = shadingFactor(alpha, shadingTurbines[j].radius, shadingTurbines[j].distance)
					}

					if len(shadingTurbines) == 0 {
						shading = 1
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

					speedWithShading := math.Round(windSpeed*shading*100) / 100

					minSpeed, maxSpeed, err := getClosestPair(speed, speedWithShading)
					if err != nil {
						log.Fatal(err)
					}

					output := linearApproximation(speedWithShading, minSpeed, maxSpeed, outputMap[fmt.Sprint(maxSpeed)].Production, outputMap[fmt.Sprint(minSpeed)].Production)

					if math.IsNaN(output) {
						output = 0
					}

					production := models.Production{
						Time:             time,
						Date:             date,
						Value:            output,
						ICUF:             output / turbine.MaximumPower,
						WindSpeed:        windSpeed,
						Altitude:         height,
						TurbineModelID:   model.Turbines[i].TurbineModelID,
						WindDirection:    windDirection,
						Shading:          shading,
						SpeedWithShading: speedWithShading,
					}

					model.Turbines[i].Productions = append(model.Turbines[i].Productions, &production)
				}
			}
		}(i, turbine)
	}
	wg.Wait()
	return s.repo.GenerateModel(userID, windfarmID, modelID, model)
}

// GetAll - ...
func (s *ModelsService) GetAll(userID, windfarmID string) ([]models.Model, error) {
	return s.repo.GetAll(userID, windfarmID)
}

// GetMapData - ...
func (s *ModelsService) GetMapData(userID, windfarmID string, modelID string) (models.ModelMap, error) {
	model, err := s.repo.GetByIDMap(userID, windfarmID, modelID)

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
	return s.repo.Delete(userID, windfarmID, modelID)
}

// Delete - ...
func (s *ModelsService) DeleteTurbine(modelID, modelTrubineID string) error {
	return s.repo.DeleteTurbine(modelID, modelTrubineID)
}

//Update - ...
func (s *ModelsService) Update(userID, windfarmID, modelID string, input models.UpdateModelInput) error {
	return s.repo.Update(userID, windfarmID, modelID, input)
}

func getElevation(c *maps.Client, lat, lng float64) ([]maps.ElevationResult, error) {
	r := &maps.ElevationRequest{
		Locations: []maps.LatLng{
			{Lat: lat, Lng: lng},
		},
	}

	return c.Elevation(context.Background(), r)
}

func getWindSpeed(height float64, winds map[string]models.Wind) (float64, error) {
	altitudes := []float64{15, 30, 50, 75, 100, 150, 200}

	minH, maxH, err := getClosestPair(altitudes, height)
	if err != nil {
		return 0, err
	}

	ws := linearApproximation(height, minH, maxH, winds[fmt.Sprint(maxH)].WindSpeed, winds[fmt.Sprint(minH)].WindSpeed)
	if minH == maxH {
		ws = winds[fmt.Sprint(maxH)].WindSpeed
	}

	return ws, nil
}

func getWindDirection(height float64, winds map[string]models.Wind) (float64, error) {
	altitudes := []float64{15, 30, 50, 75, 100, 150, 200}

	minH, maxH, err := getClosestPair(altitudes, height)
	if err != nil {
		return 0, err
	}

	wd := linearApproximation(height, minH, maxH, winds[fmt.Sprint(maxH)].WindDirection, winds[fmt.Sprint(minH)].WindDirection)
	if minH == maxH {
		wd = winds[fmt.Sprint(maxH)].WindDirection
	}

	return wd, nil
}

func checkInclude(t1, t2 models.TurbineModel, radius float64) bool {

	distance := getDistance(maps.LatLng{
		Lat: t1.Latitude,
		Lng: t1.Longitude,
	},
		maps.LatLng{
			Lat: t2.Latitude,
			Lng: t2.Longitude,
		})

	return distance < radius
}

func getDistance(p1, p2 maps.LatLng) float64 {
	R := float64(6378137)
	dLat := rad(p2.Lat - p1.Lat)
	dLong := rad(p2.Lng - p1.Lng)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rad(p1.Lat))*math.Cos(rad(p2.Lat))*
			math.Sin(dLong/2)*math.Sin(dLong/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c
	return d
}

func rad(x float64) float64 {
	return x * math.Pi / 180
}

func shadingFactor(alpha, radius, distance float64) float64 {
	if math.Cos(alpha) < 0 {
		return 1
	}
	return ((1 - math.Cos(alpha)) * (1 - (distance / radius)))
}

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

type shadingTurbine struct {
	distance float64
	radius   float64
	turbine  models.TurbineModel
}
