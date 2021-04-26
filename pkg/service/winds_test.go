package service

import (
	"testing"

	"github.com/psihachina/windfarms-backend/pkg/repository"
)

func TestWindsService_Create(t *testing.T) {
	type fields struct {
		repo             repository.Winds
		windfarmsService Windfarms
	}
	type args struct {
		userID     string
		windfarmID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &WindsService{
				repo:             tt.fields.repo,
				windfarmsService: tt.fields.windfarmsService,
			}
			if err := s.Create(tt.args.userID, tt.args.windfarmID); (err != nil) != tt.wantErr {
				t.Errorf("WindsService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
