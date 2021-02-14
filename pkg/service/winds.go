package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

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

	if err := exec.Command("wgrib2",
		"./assets/wfilter.gbr",
		"-csv",
		"./assets/wfilter.csv").Run(); err != nil {
		return err
	}

	csvfile, err := os.Open("./assets/wfilter.csv")
	if err != nil {
		return err
	}

	r := csv.NewReader(csvfile)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var weather models.Wind
		weather.Date = strings.Split(record[0], " ")[0]
		weather.Time = strings.Split(record[0], " ")[1]
		weather.Humidity = 1
		weather.WindDirection = "WIND"
		weather.WindfarmID = windfarmID
		weather.Altitude, _ = strconv.ParseFloat(strings.Split(record[3], " ")[0], 64)
		weather.WindSpeed, _ = strconv.ParseFloat(record[6], 64)

		winds = append(winds, weather)
	}

	return s.repo.Create(userID, windfarmID, winds)
}

func (s *WindsService) GetAll(userID, windfarmID string) ([]models.Wind, error) {
	return s.repo.GetAll(userID, windfarmID)
}
