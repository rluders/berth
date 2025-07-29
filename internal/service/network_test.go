package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/filters"
	networkTypes "github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
	"github.com/rluders/berth/mocks/client"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewNetworkService(t *testing.T) {
	type args struct {
		client dockerClient.APIClient
	}
	mockClient := client.NewMockAPIClient(t)
	tests := []struct {
		name string
		args args
		want NetworkService
	}{
		{
			name: "creates new network service",
			args: args{
				client: mockClient,
			},
			want: &dockerNetworkService{
				client: mockClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNetworkService(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNetworkService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerNetworkService_NetworkInspect(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx       context.Context
		networkID string
		options   networkTypes.InspectOptions
	}

	mockClient := client.NewMockAPIClient(t)

	// Setup successful network inspect
	successResp := networkTypes.Inspect{
		ID:     "network123",
		Name:   "test-network",
		Driver: "bridge",
		Scope:  "local",
	}
	mockClient.EXPECT().NetworkInspect(mock.Anything, "network123", networkTypes.InspectOptions{}).Return(successResp, nil)

	// Setup failed network inspect
	mockClient.EXPECT().NetworkInspect(mock.Anything, "invalid-network", networkTypes.InspectOptions{}).Return(networkTypes.Inspect{}, fmt.Errorf("network not found"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    networkTypes.Inspect
		wantErr bool
	}{
		{
			name: "successful network inspect",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:       context.Background(),
				networkID: "network123",
				options:   networkTypes.InspectOptions{},
			},
			want:    successResp,
			wantErr: false,
		},
		{
			name: "failed network inspect",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:       context.Background(),
				networkID: "invalid-network",
				options:   networkTypes.InspectOptions{},
			},
			want:    networkTypes.Inspect{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &dockerNetworkService{
				client: tt.fields.client,
			}
			got, err := s.NetworkInspect(tt.args.ctx, tt.args.networkID, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworkInspect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworkInspect() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerNetworkService_NetworkList(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx     context.Context
		options networkTypes.ListOptions
	}

	mockClient := client.NewMockAPIClient(t)

	// Setup successful network list
	successList := []networkTypes.Summary{{
		ID:     "network123",
		Name:   "test-network",
		Driver: "bridge",
		Scope:  "local",
	}}
	mockClient.EXPECT().NetworkList(mock.Anything, networkTypes.ListOptions{}).Return(successList, nil)

	// Setup different options for failed network list
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", "test=true")
	filterArgs.Add("dangling", "true")
	filteredOptions := networkTypes.ListOptions{
		Filters: filterArgs,
	}
	mockClient.EXPECT().NetworkList(mock.Anything, filteredOptions).Return(nil, fmt.Errorf("failed to list networks"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []networkTypes.Summary
		wantErr bool
	}{
		{
			name: "successful network list",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:     context.Background(),
				options: networkTypes.ListOptions{},
			},
			want:    successList,
			wantErr: false,
		},
		{
			name: "failed network list with filters",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:     context.Background(),
				options: filteredOptions,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &dockerNetworkService{
				client: tt.fields.client,
			}
			got, err := s.NetworkList(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworkList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworkList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
