package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/repository"
)

type WindsService struct {
	repo             repository.Winds
	windfarmsService Windfarms
}

func NewWindsService(repo repository.Winds, windfarmsService Windfarms) *WindsService {
	return &WindsService{
		repo:             repo,
		windfarmsService: windfarmsService,
	}
}

// Create - функция создания данных о ветрах на територии ветрянной электростанции
func (s *WindsService) Create(userID string, windfarmID string) error {
	var winds []models.Wind

	windfarm, err := s.windfarmsService.GetByID(userID, windfarmID)
	if err != nil {
		return err
	}

	if err := exec.Command("./scripts/grib",
		fmt.Sprintf("%f", windfarm.Longitude)+":1:0.0001",
		fmt.Sprintf("%f", windfarm.Latitude)+":1:0.0001").Run(); err != nil {
		return err
	}

	if err := exec.Command("./scripts/csv").Run(); err != nil {
		return err
	}

	csvfile, err := os.Open(os.Getenv("HOME") + "/weather/filterSpeed.csv")
	if err != nil {
		return err
	}

	csvfile2, err := os.Open(os.Getenv("HOME") + "/weather/filterDirection.csv")
	if err != nil {
		return err
	}

	r1 := csv.NewReader(csvfile)
	r2 := csv.NewReader(csvfile2)

	for {
		record1, err := r1.Read()
		if err == io.EOF {
			break
		}
		record2, err := r2.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var weather models.Wind
		weather.Date = strings.Split(record1[0], " ")[0]
		weather.Time = strings.Split(record1[0], " ")[1]
		weather.WindDirection, _ = strconv.ParseFloat(record2[6], 64)
		weather.WindfarmID = windfarmID
		weather.Altitude, _ = strconv.ParseFloat(strings.Split(record1[3], " ")[0], 64)
		weather.WindSpeed, _ = strconv.ParseFloat(record1[6], 64)

		winds = append(winds, weather)
	}

	windRange := make(map[string]map[float64]map[string]models.Wind)

	for _, item := range winds {
		if windRange[item.Date] == nil {
			windRange[item.Date] = map[float64]map[string]models.Wind{}
		}
		if windRange[item.Date][item.Altitude] == nil {
			windRange[item.Date][item.Altitude] = map[string]models.Wind{}
		}
		windRange[item.Date][item.Altitude][item.Time] = item
	}

	dateKeys := reflect.ValueOf(windRange).MapKeys()

	for dateIndex, date := range dateKeys {
		for heightKey, height := range windRange[date.String()] {
			timeKeys := reflect.ValueOf(height).MapKeys()

			for timeIndex, timeRange := range timeKeys {
				var current float64
				var next float64

				t, err := time.Parse("15:04:05", timeRange.String())
				if err != nil {
					return err
				}
				if timeIndex != 3 {
					current = height[timeRange.String()].WindSpeed
					next = height[timeKeys[timeIndex+1].String()].WindSpeed
				} else if len(dateKeys) != dateIndex+1 {
					current = height[timeRange.String()].WindSpeed
					next = windRange[dateKeys[dateIndex+1].String()][heightKey][timeKeys[1].String()].WindSpeed
				} else {
					continue
				}
				step := (next - current) / 6

				for i := 1; i < 6; i++ {
					newItem := height[timeRange.String()]
					t = t.Add(time.Hour)
					newItem.Time = t.Format("15:04:05")
					newItem.WindSpeed += float64(i) * step
					newItem.WindSpeed = math.Round(newItem.WindSpeed*100000) / 100000
					winds = append(winds, newItem)
				}

			}

		}
	}

	return s.repo.Create(userID, windfarmID, winds)
}

func (s *WindsService) GetAll(userID, windfarmID string) ([]models.Wind, error) {
	return s.repo.GetAll(userID, windfarmID)
}

func (s *WindsService) GetWindForChart(userID, windfarmID, from, to string, height int) (models.ChartData, error) {

	winds, err := s.repo.GetWindForChart(userID, windfarmID, from, to, height)
	if err != nil {
		return models.ChartData{}, err
	}

	//Experemental distribution
	var expDist [31]int

	//Distribution density
	var distDensity [31]int

	//Wind energy
	var windEnergy [31]float64

	for i := 0; i < 30; i++ {
		for _, wind := range winds {
			if wind.WindSpeed > float64(i) {
				expDist[i] += 1
			}
		}
	}

	for i := 0; i < 30; i++ {
		distDensity[i] = int(math.Abs(float64(expDist[i]) - float64(expDist[i+1])))

		var buf []models.Wind

		for _, wind := range winds {
			if wind.WindSpeed > float64(i) {
				buf = append(buf, wind)
			}
		}

		var avg float64

		if expDist[i] != 0 {
			for _, b := range buf {
				avg += b.WindSpeed
			}
			avg = avg / float64(expDist[i])
		} else {
			avg = 0
		}

		windEnergy[i] = math.Round(((0.5*(1.2*(math.Abs(float64(expDist[i])-float64(expDist[i+1]))*60))*math.Pow(avg, 3))/3600000)*100) / 100
	}

	return models.ChartData{ExperementalDistribution: expDist, DistributionDensity: distDensity, WindEnergy: windEnergy}, err
}

func (s *WindsService) GetWindForTable(userID, windfarmID string) (models.TableData, error) {
	winds, err := s.repo.GetAll(userID, windfarmID)
	if err != nil {
		return models.TableData{}, err
	}

	windRange := make(map[float64]map[string][]models.Wind)

	for _, wind := range winds {
		if windRange[wind.Altitude] == nil {
			windRange[wind.Altitude] = map[string][]models.Wind{}
			windRange[wind.Altitude][wind.Time] = append(windRange[wind.Altitude][wind.Time], wind)
		} else {
			windRange[wind.Altitude][wind.Time] = append(windRange[wind.Altitude][wind.Time], wind)
		}
	}

	var tableData models.TableData

	tableData.Avg = map[string]map[int]float64{}
	tableData.Dispersion = map[string]map[int]float64{}
	tableData.StandardDeviation = map[string]map[int]float64{}

	for heightKey, height := range windRange {
		for timeKey, timeRange := range height {
			if tableData.Avg[timeKey] == nil {
				tableData.Avg[timeKey] = map[int]float64{}
				tableData.StandardDeviation[timeKey] = map[int]float64{}
				tableData.Dispersion[timeKey] = map[int]float64{}
			}
			var avg float64
			var dispersion float64

			for _, wind := range timeRange {
				avg += wind.WindSpeed
			}
			avg = math.Round((avg/float64(len(timeRange)))*100) / 100

			for _, wind := range timeRange {
				dispersion += math.Pow(wind.WindSpeed-avg, 2)
			}

			dispersion = dispersion / float64(len(timeRange))
			dispersion = math.Round(dispersion*100) / 100

			tableData.Avg[timeKey][int(heightKey)] = avg
			tableData.StandardDeviation[timeKey][int(heightKey)] = math.Round(math.Sqrt(dispersion)*100) / 100
			tableData.Dispersion[timeKey][int(heightKey)] = dispersion
		}
	}

	return tableData, err
}
