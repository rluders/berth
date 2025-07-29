package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/api/types/volume"
	dockerClient "github.com/docker/docker/client"
	"github.com/rluders/berth/mocks/client"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewSystemService(t *testing.T) {
	type args struct {
		client dockerClient.APIClient
	}

	mockClient := client.NewMockAPIClient(t)

	tests := []struct {
		name string
		args args
		want SystemService
	}{
		{
			name: "creates new system service",
			args: args{
				client: mockClient,
			},
			want: &dockerSystemService{
				client: mockClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSystemService(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSystemService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerSystemService_ContainersPrune(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx          context.Context
		pruneFilters filters.Args
	}

	// Success case
	successMock := client.NewMockAPIClient(t)
	successReport := container.PruneReport{
		ContainersDeleted: []string{"container1", "container2"},
		SpaceReclaimed:    1024 * 1024 * 10, // 10MB
	}
	successMock.EXPECT().ContainersPrune(mock.Anything, mock.AnythingOfType("filters.Args")).Return(successReport, nil)

	// Error case
	errorMock := client.NewMockAPIClient(t)
	errorMock.EXPECT().ContainersPrune(mock.Anything, mock.AnythingOfType("filters.Args")).Return(container.PruneReport{}, fmt.Errorf("failed to prune containers"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    container.PruneReport
		wantErr bool
	}{
		{
			name: "success case",
			fields: fields{
				client: successMock,
			},
			args: args{
				ctx:          context.Background(),
				pruneFilters: filters.Args{},
			},
			want:    successReport,
			wantErr: false,
		},
		{
			name: "error case",
			fields: fields{
				client: errorMock,
			},
			args: args{
				ctx:          context.Background(),
				pruneFilters: filters.Args{},
			},
			want:    container.PruneReport{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerSystemService{
				client: tt.fields.client,
			}
			got, err := s.ContainersPrune(tt.args.ctx, tt.args.pruneFilters)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContainersPrune() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContainersPrune() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerSystemService_DiskUsage(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx     context.Context
		options types.DiskUsageOptions
	}

	// Success case
	successMock := client.NewMockAPIClient(t)
	successUsage := types.DiskUsage{
		LayersSize: 1024 * 1024 * 100, // 100MB
		Images:     []*image.Summary{},
		Containers: []*container.Summary{},
		Volumes:    []*volume.Volume{},
	}
	successMock.EXPECT().DiskUsage(mock.Anything, mock.AnythingOfType("types.DiskUsageOptions")).Return(successUsage, nil)

	// Error case
	errorMock := client.NewMockAPIClient(t)
	errorMock.EXPECT().DiskUsage(mock.Anything, mock.AnythingOfType("types.DiskUsageOptions")).Return(types.DiskUsage{}, fmt.Errorf("failed to get disk usage"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.DiskUsage
		wantErr bool
	}{
		{
			name: "success case",
			fields: fields{
				client: successMock,
			},
			args: args{
				ctx:     context.Background(),
				options: types.DiskUsageOptions{},
			},
			want:    successUsage,
			wantErr: false,
		},
		{
			name: "error case",
			fields: fields{
				client: errorMock,
			},
			args: args{
				ctx:     context.Background(),
				options: types.DiskUsageOptions{},
			},
			want:    types.DiskUsage{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerSystemService{
				client: tt.fields.client,
			}
			got, err := s.DiskUsage(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiskUsage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiskUsage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerSystemService_ImagesPrune(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx          context.Context
		pruneFilters filters.Args
	}

	// Success case
	successMock := client.NewMockAPIClient(t)
	successReport := image.PruneReport{
		ImagesDeleted: []image.DeleteResponse{
			{Deleted: "image1"},
			{Deleted: "image2"},
		},
		SpaceReclaimed: 1024 * 1024 * 50, // 50MB
	}
	successMock.EXPECT().ImagesPrune(mock.Anything, mock.AnythingOfType("filters.Args")).Return(successReport, nil)

	// Error case
	errorMock := client.NewMockAPIClient(t)
	errorMock.EXPECT().ImagesPrune(mock.Anything, mock.AnythingOfType("filters.Args")).Return(image.PruneReport{}, fmt.Errorf("failed to prune images"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    image.PruneReport
		wantErr bool
	}{
		{
			name: "success case",
			fields: fields{
				client: successMock,
			},
			args: args{
				ctx:          context.Background(),
				pruneFilters: filters.Args{},
			},
			want:    successReport,
			wantErr: false,
		},
		{
			name: "error case",
			fields: fields{
				client: errorMock,
			},
			args: args{
				ctx:          context.Background(),
				pruneFilters: filters.Args{},
			},
			want:    image.PruneReport{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerSystemService{
				client: tt.fields.client,
			}
			got, err := s.ImagesPrune(tt.args.ctx, tt.args.pruneFilters)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImagesPrune() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImagesPrune() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerSystemService_Info(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx context.Context
	}

	// Success case
	successMock := client.NewMockAPIClient(t)
	successInfo := system.Info{
		ID:                "test-id",
		Containers:        5,
		ContainersRunning: 3,
		ContainersPaused:  0,
		ContainersStopped: 2,
		Images:            10,
	}
	successMock.EXPECT().Info(mock.Anything).Return(successInfo, nil)

	// Error case
	errorMock := client.NewMockAPIClient(t)
	errorMock.EXPECT().Info(mock.Anything).Return(system.Info{}, fmt.Errorf("failed to get info"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    system.Info
		wantErr bool
	}{
		{
			name: "success case",
			fields: fields{
				client: successMock,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    successInfo,
			wantErr: false,
		},
		{
			name: "error case",
			fields: fields{
				client: errorMock,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    system.Info{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerSystemService{
				client: tt.fields.client,
			}
			got, err := s.Info(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Info() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerSystemService_NetworksPrune(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx          context.Context
		pruneFilters filters.Args
	}

	// Success case
	successMock := client.NewMockAPIClient(t)
	successReport := network.PruneReport{
		NetworksDeleted: []string{"network1", "network2"},
	}
	successMock.EXPECT().NetworksPrune(mock.Anything, mock.AnythingOfType("filters.Args")).Return(successReport, nil)

	// Error case
	errorMock := client.NewMockAPIClient(t)
	errorMock.EXPECT().NetworksPrune(mock.Anything, mock.AnythingOfType("filters.Args")).Return(network.PruneReport{}, fmt.Errorf("failed to prune networks"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    network.PruneReport
		wantErr bool
	}{
		{
			name: "success case",
			fields: fields{
				client: successMock,
			},
			args: args{
				ctx:          context.Background(),
				pruneFilters: filters.Args{},
			},
			want:    successReport,
			wantErr: false,
		},
		{
			name: "error case",
			fields: fields{
				client: errorMock,
			},
			args: args{
				ctx:          context.Background(),
				pruneFilters: filters.Args{},
			},
			want:    network.PruneReport{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerSystemService{
				client: tt.fields.client,
			}
			got, err := s.NetworksPrune(tt.args.ctx, tt.args.pruneFilters)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworksPrune() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworksPrune() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerSystemService_VolumesPrune(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx          context.Context
		pruneFilters filters.Args
	}

	// Success case
	successMock := client.NewMockAPIClient(t)
	successReport := volume.PruneReport{
		VolumesDeleted: []string{"volume1", "volume2"},
		SpaceReclaimed: 1024 * 1024 * 20, // 20MB
	}
	successMock.EXPECT().VolumesPrune(mock.Anything, mock.AnythingOfType("filters.Args")).Return(successReport, nil)

	// Error case
	errorMock := client.NewMockAPIClient(t)
	errorMock.EXPECT().VolumesPrune(mock.Anything, mock.AnythingOfType("filters.Args")).Return(volume.PruneReport{}, fmt.Errorf("failed to prune volumes"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    volume.PruneReport
		wantErr bool
	}{
		{
			name: "success case",
			fields: fields{
				client: successMock,
			},
			args: args{
				ctx:          context.Background(),
				pruneFilters: filters.Args{},
			},
			want:    successReport,
			wantErr: false,
		},
		{
			name: "error case",
			fields: fields{
				client: errorMock,
			},
			args: args{
				ctx:          context.Background(),
				pruneFilters: filters.Args{},
			},
			want:    volume.PruneReport{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerSystemService{
				client: tt.fields.client,
			}
			got, err := s.VolumesPrune(tt.args.ctx, tt.args.pruneFilters)
			if (err != nil) != tt.wantErr {
				t.Errorf("VolumesPrune() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VolumesPrune() got = %v, want %v", got, tt.want)
			}
		})
	}
}
