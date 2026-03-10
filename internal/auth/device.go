package auth

import (
	"context"
	"fmt"
)

// DeviceClient handles the Device Authorization Flow.
// NOTE: Device authorization is not yet supported by the MYN backend.
// This is a stub per Spec §2.4.
type DeviceClient struct {
	BaseURL string
}

// NewDeviceClient creates a new device auth client.
func NewDeviceClient(baseURL string) *DeviceClient {
	return &DeviceClient{BaseURL: baseURL}
}

// Authorize initiates the device authorization flow.
// Returns ErrNotImplemented until the backend supports it.
func (c *DeviceClient) Authorize(ctx context.Context) error {
	return fmt.Errorf("device authorization flow is not yet supported by the MYN backend\n" +
		"Use 'mynow login' (OAuth) or 'mynow login --api-key' instead")
}
