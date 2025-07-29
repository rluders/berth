package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
	"github.com/rluders/berth/mocks/client"
	"github.com/stretchr/testify/mock"
	"io"
	"reflect"
	"strings"
	"testing"
)

// mockReadCloser is a mock implementation of io.ReadCloser for testing
type mockReadCloser struct {
	io.Reader
}

func (m *mockReadCloser) Close() error {
	return nil
}

func newMockReadCloser(s string) io.ReadCloser {
	return &mockReadCloser{strings.NewReader(s)}
}

func TestNewContainerService(t *testing.T) {
	type args struct {
		client dockerClient.APIClient
	}
	mockClient := client.NewMockAPIClient(t)
	tests := []struct {
		name string
		args args
		want ContainerService
	}{
		{
			name: "creates new container service",
			args: args{
				client: mockClient,
			},
			want: &dockerContainerService{
				client: mockClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContainerService(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContainerService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerContainerService_ContainerInspect(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx         context.Context
		containerID string
	}

	mockClient := client.NewMockAPIClient(t)

	// Setup successful container inspect
	successResp := container.InspectResponse{}
	mockClient.EXPECT().ContainerInspect(mock.Anything, "container123").Return(successResp, nil)

	// Setup failed container inspect
	mockClient.EXPECT().ContainerInspect(mock.Anything, "invalid-container").Return(container.InspectResponse{}, fmt.Errorf("container not found"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    container.InspectResponse
		wantErr bool
	}{
		{
			name: "successful container inspect",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:         context.Background(),
				containerID: "container123",
			},
			want:    successResp,
			wantErr: false,
		},
		{
			name: "failed container inspect",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:         context.Background(),
				containerID: "invalid-container",
			},
			want:    container.InspectResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerContainerService{
				client: tt.fields.client,
			}
			got, err := s.ContainerInspect(tt.args.ctx, tt.args.containerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerInspect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContainerInspect() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerContainerService_ContainerLogs(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx         context.Context
		containerID string
		options     container.LogsOptions
	}

	mockClient := client.NewMockAPIClient(t)
	successLogs := newMockReadCloser("container logs")

	// Setup successful container logs
	mockClient.EXPECT().ContainerLogs(mock.Anything, "container123", mock.AnythingOfType("container.LogsOptions")).Return(successLogs, nil)

	// Setup failed container logs
	mockClient.EXPECT().ContainerLogs(mock.Anything, "invalid-container", mock.AnythingOfType("container.LogsOptions")).Return(nil, fmt.Errorf("container not found"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
		{
			name: "successful container logs retrieval",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:         context.Background(),
				containerID: "container123",
				options: container.LogsOptions{
					ShowStdout: true,
					ShowStderr: true,
				},
			},
			want:    successLogs,
			wantErr: false,
		},
		{
			name: "failed container logs retrieval",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:         context.Background(),
				containerID: "invalid-container",
				options: container.LogsOptions{
					ShowStdout: true,
					ShowStderr: true,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerContainerService{
				client: tt.fields.client,
			}
			got, err := s.ContainerLogs(tt.args.ctx, tt.args.containerID, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContainerLogs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerContainerService_ListContainers(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx     context.Context
		options container.ListOptions
	}

	mockClient := client.NewMockAPIClient(t)
	successList := []container.Summary{}

	// Setup successful container list
	mockClient.EXPECT().ContainerList(mock.Anything, container.ListOptions{All: true}).Return(successList, nil)

	// Setup failed container list
	mockClient.EXPECT().ContainerList(mock.Anything, container.ListOptions{}).Return(nil, fmt.Errorf("failed to list containers"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []container.Summary
		wantErr bool
	}{
		{
			name: "successful container list",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:     context.Background(),
				options: container.ListOptions{All: true},
			},
			want:    successList,
			wantErr: false,
		},
		{
			name: "failed container list",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:     context.Background(),
				options: container.ListOptions{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerContainerService{
				client: tt.fields.client,
			}
			got, err := s.ListContainers(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListContainers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListContainers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerContainerService_RemoveContainer(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx         context.Context
		containerID string
		options     container.RemoveOptions
	}

	mockClient := client.NewMockAPIClient(t)

	// Setup successful container removal
	mockClient.EXPECT().ContainerRemove(mock.Anything, "container123", container.RemoveOptions{Force: true}).Return(nil)

	// Setup failed container removal
	mockClient.EXPECT().ContainerRemove(mock.Anything, "invalid-container", container.RemoveOptions{}).Return(fmt.Errorf("container not found"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful container removal",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:         context.Background(),
				containerID: "container123",
				options:     container.RemoveOptions{Force: true},
			},
			wantErr: false,
		},
		{
			name: "failed container removal",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:         context.Background(),
				containerID: "invalid-container",
				options:     container.RemoveOptions{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerContainerService{
				client: tt.fields.client,
			}
			if err := s.RemoveContainer(tt.args.ctx, tt.args.containerID, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("RemoveContainer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_dockerContainerService_StartContainer(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx         context.Context
		containerID string
		options     container.StartOptions
	}

	mockClient := client.NewMockAPIClient(t)

	// Setup successful container start
	mockClient.EXPECT().ContainerStart(mock.Anything, "container123", container.StartOptions{}).Return(nil)

	// Setup failed container start
	mockClient.EXPECT().ContainerStart(mock.Anything, "invalid-container", container.StartOptions{}).Return(fmt.Errorf("container not found"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful container start",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:         context.Background(),
				containerID: "container123",
				options:     container.StartOptions{},
			},
			wantErr: false,
		},
		{
			name: "failed container start",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:         context.Background(),
				containerID: "invalid-container",
				options:     container.StartOptions{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerContainerService{
				client: tt.fields.client,
			}
			if err := s.StartContainer(tt.args.ctx, tt.args.containerID, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("StartContainer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_dockerContainerService_StopContainer(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx         context.Context
		containerID string
		options     container.StopOptions
	}

	mockClient := client.NewMockAPIClient(t)

	// Setup successful container stop
	mockClient.EXPECT().ContainerStop(mock.Anything, "container123", container.StopOptions{}).Return(nil)

	// Setup failed container stop
	mockClient.EXPECT().ContainerStop(mock.Anything, "invalid-container", container.StopOptions{}).Return(fmt.Errorf("container not found"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful container stop",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:         context.Background(),
				containerID: "container123",
				options:     container.StopOptions{},
			},
			wantErr: false,
		},
		{
			name: "failed container stop",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:         context.Background(),
				containerID: "invalid-container",
				options:     container.StopOptions{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerContainerService{
				client: tt.fields.client,
			}
			if err := s.StopContainer(tt.args.ctx, tt.args.containerID, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("StopContainer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
