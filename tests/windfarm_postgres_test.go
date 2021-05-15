package tests

import (
	"errors"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/pgtype"
	"github.com/psihachina/windfarms-backend/models"
	mockdb "github.com/psihachina/windfarms-backend/pkg/repository/mock"
	"github.com/psihachina/windfarms-backend/utils"
	"github.com/stretchr/testify/assert"
)

func TestWindfarmPostgres_Create(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockdb.NewMockWindfarms(ctrl)

	type args struct {
		UserID   string
		windfarm models.Windfarm
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
				windfarm: randomWindfarm(false),
				UserID:   strconv.Itoa(utils.RandomInt(1, 1000)),
			},
			mock: func(args args) {
				r.EXPECT().
					Create(args.UserID, gomock.Eq(args.windfarm)).
					Times(1).
					Return(args.windfarm.WindfarmID, nil)
			},
		},
		{
			name: "Empty field",
			input: args{
				windfarm: randomWindfarm(true),
				UserID:   strconv.Itoa(utils.RandomInt(1, 1000)),
			},
			mock: func(args args) {
				r.EXPECT().
					Create(args.UserID, gomock.Eq(args.windfarm)).
					Times(1).
					Return("", errors.New("no rows in result set"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.Create(tt.input.UserID, tt.input.windfarm)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func randomWindfarm(emptyName bool) models.Windfarm {
	var name string
	if emptyName {
		name = ""
	} else {
		name = utils.RandomString(6)
	}
	return models.Windfarm{
		WindfarmID:   strconv.Itoa(utils.RandomInt(1, 1000)),
		WindfarmName: name,
		PolygonDB: pgtype.Polygon{
			P: []pgtype.Vec2{
				{
					X: float64(utils.RandomInt(1, 1000)),
					Y: float64(utils.RandomInt(1, 1000)),
				},
				{
					X: float64(utils.RandomInt(1, 1000)),
					Y: float64(utils.RandomInt(1, 1000)),
				},
				{
					X: float64(utils.RandomInt(1, 1000)),
					Y: float64(utils.RandomInt(1, 1000)),
				},
			},
		},
		Longitude:       float64(utils.RandomInt(1, 1000)),
		Latitude:        float64(utils.RandomInt(1, 1000)),
		Capacity:        float64(utils.RandomInt(1, 1000)),
		RangeToCity:     float64(utils.RandomInt(1, 1000)),
		RangeToRoad:     float64(utils.RandomInt(1, 1000)),
		RangeToCityLine: float64(utils.RandomInt(1, 1000)),
		CityLatitude:    float64(utils.RandomInt(1, 1000)),
		CityLongitude:   float64(utils.RandomInt(1, 1000)),
		Description:     utils.RandomString(100),
	}
}
