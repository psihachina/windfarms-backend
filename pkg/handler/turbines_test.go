package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/service"
	service_mocks "github.com/psihachina/windfarms-backend/pkg/service/mocks"
	"github.com/psihachina/windfarms-backend/utils"
)

func TestHandler_create(t *testing.T) {
	// Init test table
	type mockBehavior func(r *service_mocks.MockTurbines, UserID string, turbine models.Turbine, outputs models.Outputs)

	turbine := randomTurbine(false)
	turbineID := utils.RandomString(18)

	tests := []struct {
		name                 string
		inputBody            string
		inputTurbine         models.Turbine
		inputOutputs         models.Outputs
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			inputBody: fmt.Sprintf(`{"userID":"%s","turbine_name":"%s",
			"maximum_power":%v,"max_wind_speed":%v,"min_wind_speed":%v,
			"number_blades":%v,"tower_height":%v,"rotor_diameter":%v,
			"annual_turbine_maintenance":%v,"outputs": [{"production":%v}]}`,
				turbine.UserID, turbine.TurbineName, turbine.MaximumPower,
				turbine.MaxWindSpeed, turbine.MinWindSpeed,
				turbine.Blades, turbine.TowerHeight, turbine.RotorDiameter,
				turbine.AnnualTurbineMaintenance, turbine.Outputs[0].Production,
			),
			inputTurbine: turbine,
			inputOutputs: models.Outputs{
				Outputs: turbine.Outputs,
			},
			mockBehavior: func(r *service_mocks.MockTurbines, userID string, turbine models.Turbine, outputs models.Outputs) {
				r.EXPECT().Create(turbine.UserID, turbine, outputs).Return(turbineID, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: fmt.Sprintf(`{"id":"%s"}`, turbineID),
		},
		{
			name: "Wrong Input",
			inputBody: fmt.Sprintf(`{"userID":"%s","turbine_name":"%s"}`,
				turbine.UserID, turbine.TurbineName,
			),
			inputTurbine: models.Turbine{},
			inputOutputs: models.Outputs{
				Outputs: turbine.Outputs,
			},
			mockBehavior: func(r *service_mocks.MockTurbines, userID string, turbine models.Turbine, outputs models.Outputs) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name: "Service Error",
			inputBody: fmt.Sprintf(`{"userID":"%s","turbine_name":"%s",
			"maximum_power":%v,"max_wind_speed":%v,"min_wind_speed":%v,
			"number_blades":%v,"tower_height":%v,"rotor_diameter":%v,
			"annual_turbine_maintenance":%v,"outputs": [{"production":%v}]}`,
				turbine.UserID, turbine.TurbineName, turbine.MaximumPower,
				turbine.MaxWindSpeed, turbine.MinWindSpeed,
				turbine.Blades, turbine.TowerHeight, turbine.RotorDiameter,
				turbine.AnnualTurbineMaintenance, turbine.Outputs[0].Production,
			),
			inputTurbine: turbine,
			inputOutputs: models.Outputs{
				Outputs: turbine.Outputs,
			},
			mockBehavior: func(r *service_mocks.MockTurbines, userID string, turbine models.Turbine, outputs models.Outputs) {
				r.EXPECT().Create(turbine.UserID, turbine, outputs).Return("", errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Init Depedencies
			c := gomock.NewController(t)
			defer c.Finish()

			turbine := service_mocks.NewMockTurbines(c)
			test.mockBehavior(turbine, test.inputTurbine.UserID, test.inputTurbine, test.inputOutputs)

			services := &service.Service{Turbines: turbine}
			handler := NewHandler(services)

			//Init Endpoint
			r := gin.New()
			r.POST("/api/turbine/", func(c *gin.Context) { c.Set(userCtx, test.inputTurbine.UserID) }, handler.createTurbine)

			//Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/turbine/",
				bytes.NewBufferString(test.inputBody))

			//Make Request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func randomTurbine(emptyField bool) models.Turbine {
	var turbineName string
	if emptyField {
		turbineName = ""
	} else {
		turbineName = utils.RandomString(6)
	}
	return models.Turbine{
		UserID:                   utils.RandomString(18),
		TurbineName:              turbineName,
		MaximumPower:             float64(utils.RandomInt(0, 3000)),
		MaxWindSpeed:             float64(utils.RandomInt(0, 3000)),
		MinWindSpeed:             float64(utils.RandomInt(0, 3000)),
		Blades:                   3,
		TowerHeight:              float64(utils.RandomInt(18, 100)),
		RotorDiameter:            float64(utils.RandomInt(18, 100)),
		AnnualTurbineMaintenance: float64(utils.RandomInt(18, 100)),
		Outputs: []models.Output{
			{
				Production: float64(utils.RandomInt(1, 1000)),
			},
		},
	}
}
