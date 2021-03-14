package tests

import (
	"errors"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/psihachina/windfarms-backend/models"
	mockdb "github.com/psihachina/windfarms-backend/pkg/repository/mock"
	"github.com/psihachina/windfarms-backend/utils"
	"github.com/stretchr/testify/assert"
)

func TestTurbinePostgres_Create(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockdb.NewMockTurbines(ctrl)

	type args struct {
		UserID  string
		turbine models.Turbine
		outputs models.Outputs
	}

	tests := []struct {
		name    string
		mock    func(args args)
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			input: args{
				turbine: randomTurbine(false),
				outputs: randomOutputs(),
				UserID:  strconv.Itoa(utils.RandomInt(1, 1000)),
			},
			mock: func(args args) {
				r.EXPECT().
					Create(args.UserID, gomock.Eq(args.turbine), gomock.Eq(args.outputs)).
					Times(1).
					Return(args.turbine.TurbineID, nil)
			},
		},
		{
			name: "Empty field",
			input: args{
				turbine: randomTurbine(true),
				outputs: randomOutputs(),
				UserID:  strconv.Itoa(utils.RandomInt(1, 1000)),
			},
			mock: func(args args) {
				r.EXPECT().
					Create(args.UserID, gomock.Eq(args.turbine), gomock.Eq(args.outputs)).
					Times(1).
					Return("", errors.New("no rows in result set"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.Create(tt.input.UserID, tt.input.turbine, tt.input.outputs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func randomTurbine(emptyName bool) models.Turbine {
	var name string
	if emptyName {
		name = ""
	} else {
		name = utils.RandomString(6)
	}
	return models.Turbine{
		TurbineName:              name,
		MaximumPower:             float64(utils.RandomInt(1, 1000)),
		MaxWindSpeed:             float64(utils.RandomInt(1, 1000)),
		MinWindSpeed:             float64(utils.RandomInt(1, 1000)),
		Blades:                   utils.RandomInt(1, 1000),
		TowerHeight:              float64(utils.RandomInt(1, 1000)),
		RotorDiameter:            float64(utils.RandomInt(1, 1000)),
		AnnualTurbineMaintenance: float64(utils.RandomInt(1, 1000)),
	}
}

func randomOutputs() models.Outputs {
	return models.Outputs{
		Outputs: []models.Output{
			models.Output{
				Production: float64(utils.RandomInt(1, 1000)),
			},
		},
	}
}
