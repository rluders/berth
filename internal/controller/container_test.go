package controller

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types/container"
	clientmock "github.com/rluders/berth/mocks/client"
	"github.com/rluders/berth/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListContainers_returnsMappedContainers(t *testing.T) {
	mockClient := clientmock.NewMockAPIClient(t)
	mockClient.EXPECT().
		ContainerList(mock.Anything, container.ListOptions{All: true}).
		Return([]container.Summary{
			{
				ID:      "abcdef123456789",
				Image:   "nginx:latest",
				Command: "nginx -g 'daemon off;'",
				Names:   []string{"/my-nginx"},
				Status:  "Up 2 hours",
				State:   "running",
			},
		}, nil)

	setContainerServiceForTest(service.NewContainerService(mockClient))

	result, err := ListContainers()

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "abcdef123456", result[0].ID) // truncated to 12
	assert.Equal(t, "nginx:latest", result[0].Image)
	assert.Equal(t, "my-nginx", result[0].Names) // leading / stripped
	assert.Equal(t, "running", result[0].State)
}

func TestListContainers_propagatesError(t *testing.T) {
	mockClient := clientmock.NewMockAPIClient(t)
	mockClient.EXPECT().
		ContainerList(mock.Anything, container.ListOptions{All: true}).
		Return(nil, errors.New("connection refused"))

	setContainerServiceForTest(service.NewContainerService(mockClient))

	_, err := ListContainers()
	assert.ErrorContains(t, err, "connection refused")
}

func TestStartContainer_callsServiceWithID(t *testing.T) {
	mockClient := clientmock.NewMockAPIClient(t)
	mockClient.EXPECT().
		ContainerStart(mock.Anything, "abc123", container.StartOptions{}).
		Return(nil)

	setContainerServiceForTest(service.NewContainerService(mockClient))

	err := StartContainer("abc123")
	assert.NoError(t, err)
}

func TestStopContainer_callsServiceWithID(t *testing.T) {
	mockClient := clientmock.NewMockAPIClient(t)
	mockClient.EXPECT().
		ContainerStop(mock.Anything, "abc123", container.StopOptions{}).
		Return(nil)

	setContainerServiceForTest(service.NewContainerService(mockClient))

	err := StopContainer("abc123")
	assert.NoError(t, err)
}

func TestRestartContainer_callsServiceWithID(t *testing.T) {
	mockClient := clientmock.NewMockAPIClient(t)
	mockClient.EXPECT().
		ContainerRestart(mock.Anything, "abc123", container.StopOptions{}).
		Return(nil)

	setContainerServiceForTest(service.NewContainerService(mockClient))

	err := RestartContainer("abc123")
	assert.NoError(t, err)
}

func TestRemoveContainer_callsServiceWithForce(t *testing.T) {
	mockClient := clientmock.NewMockAPIClient(t)
	mockClient.EXPECT().
		ContainerRemove(mock.Anything, "abc123", container.RemoveOptions{Force: true}).
		Return(nil)

	setContainerServiceForTest(service.NewContainerService(mockClient))

	err := RemoveContainer("abc123")
	assert.NoError(t, err)
}
