package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/volume"
	dockerClient "github.com/docker/docker/client"
	"github.com/rluders/berth/mocks/client"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewVolumeService(t *testing.T) {
	type args struct {
		client dockerClient.APIClient
	}

	mockClient := client.NewMockAPIClient(t)

	tests := []struct {
		name string
		args args
		want VolumeService
	}{
		{
			name: "creates new volume service",
			args: args{
				client: mockClient,
			},
			want: &dockerVolumeService{
				client: mockClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVolumeService(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVolumeService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerVolumeService_VolumeList(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx     context.Context
		options volume.ListOptions
	}

	// Success case
	successMock := client.NewMockAPIClient(t)
	successResponse := volume.ListResponse{
		Volumes: []*volume.Volume{
			{
				Name:       "volume1",
				Driver:     "local",
				Mountpoint: "/var/lib/docker/volumes/volume1",
			},
			{
				Name:       "volume2",
				Driver:     "local",
				Mountpoint: "/var/lib/docker/volumes/volume2",
			},
		},
	}
	successMock.EXPECT().VolumeList(mock.Anything, mock.AnythingOfType("volume.ListOptions")).Return(successResponse, nil)

	// Error case
	errorMock := client.NewMockAPIClient(t)
	errorMock.EXPECT().VolumeList(mock.Anything, mock.AnythingOfType("volume.ListOptions")).Return(volume.ListResponse{}, fmt.Errorf("failed to list volumes"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    volume.ListResponse
		wantErr bool
	}{
		{
			name: "success case",
			fields: fields{
				client: successMock,
			},
			args: args{
				ctx:     context.Background(),
				options: volume.ListOptions{},
			},
			want:    successResponse,
			wantErr: false,
		},
		{
			name: "error case",
			fields: fields{
				client: errorMock,
			},
			args: args{
				ctx:     context.Background(),
				options: volume.ListOptions{},
			},
			want:    volume.ListResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerVolumeService{
				client: tt.fields.client,
			}
			got, err := s.VolumeList(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("VolumeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VolumeList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerVolumeService_VolumeRemove(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx      context.Context
		volumeID string
		force    bool
	}

	// Success case
	successMock := client.NewMockAPIClient(t)
	successMock.EXPECT().VolumeRemove(mock.Anything, "volume1", false).Return(nil)

	// Error case
	errorMock := client.NewMockAPIClient(t)
	errorMock.EXPECT().VolumeRemove(mock.Anything, "invalid-volume", true).Return(fmt.Errorf("volume not found"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success case",
			fields: fields{
				client: successMock,
			},
			args: args{
				ctx:      context.Background(),
				volumeID: "volume1",
				force:    false,
			},
			wantErr: false,
		},
		{
			name: "error case",
			fields: fields{
				client: errorMock,
			},
			args: args{
				ctx:      context.Background(),
				volumeID: "invalid-volume",
				force:    true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerVolumeService{
				client: tt.fields.client,
			}
			if err := s.VolumeRemove(tt.args.ctx, tt.args.volumeID, tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("VolumeRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
