package service

import (
	"context"
	"fmt"
	imageTypes "github.com/docker/docker/api/types/image"
	dockerClient "github.com/docker/docker/client"
	"github.com/rluders/berth/mocks/client"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewImageService(t *testing.T) {
	type args struct {
		client dockerClient.APIClient
	}
	mockClient := client.NewMockAPIClient(t)
	tests := []struct {
		name string
		args args
		want ImageService
	}{
		{
			name: "creates new image service",
			args: args{
				client: mockClient,
			},
			want: &dockerImageService{
				client: mockClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewImageService(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewImageService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerImageService_ImageList(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx     context.Context
		options imageTypes.ListOptions
	}

	mockClient := client.NewMockAPIClient(t)

	// Setup successful image list
	successList := []imageTypes.Summary{}
	mockClient.EXPECT().ImageList(mock.Anything, imageTypes.ListOptions{All: true}).Return(successList, nil)

	// Setup failed image list
	mockClient.EXPECT().ImageList(mock.Anything, imageTypes.ListOptions{}).Return(nil, fmt.Errorf("failed to list images"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []imageTypes.Summary
		wantErr bool
	}{
		{
			name: "successful image list",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:     context.Background(),
				options: imageTypes.ListOptions{All: true},
			},
			want:    successList,
			wantErr: false,
		},
		{
			name: "failed image list",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:     context.Background(),
				options: imageTypes.ListOptions{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerImageService{
				client: tt.fields.client,
			}
			got, err := s.ImageList(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImageList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dockerImageService_ImageRemove(t *testing.T) {
	type fields struct {
		client dockerClient.APIClient
	}
	type args struct {
		ctx     context.Context
		imageID string
		options imageTypes.RemoveOptions
	}

	mockClient := client.NewMockAPIClient(t)

	// Setup successful image removal
	successResp := []imageTypes.DeleteResponse{{
		Deleted:  "sha256:1234567890",
		Untagged: "test:latest",
	}}
	mockClient.EXPECT().ImageRemove(mock.Anything, "image123", imageTypes.RemoveOptions{Force: true}).Return(successResp, nil)

	// Setup failed image removal
	mockClient.EXPECT().ImageRemove(mock.Anything, "invalid-image", imageTypes.RemoveOptions{}).Return(nil, fmt.Errorf("image not found"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []imageTypes.DeleteResponse
		wantErr bool
	}{
		{
			name: "successful image removal",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:     context.Background(),
				imageID: "image123",
				options: imageTypes.RemoveOptions{Force: true},
			},
			want:    successResp,
			wantErr: false,
		},
		{
			name: "failed image removal",
			fields: fields{
				client: mockClient,
			},
			args: args{
				ctx:     context.Background(),
				imageID: "invalid-image",
				options: imageTypes.RemoveOptions{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dockerImageService{
				client: tt.fields.client,
			}
			got, err := s.ImageRemove(tt.args.ctx, tt.args.imageID, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageRemove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImageRemove() got = %v, want %v", got, tt.want)
			}
		})
	}
}
