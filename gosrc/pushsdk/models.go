package pushsdk

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

var verifier = validator.New()

type PushConfig struct {
	PushMethods []string `json:"pushMethods" validate:"required"`
	IsDryRun    bool     `json:"isDryRun" validate:"boolean,required"`

	// ServChan 3 for mobile universal
	SC3ServerURL   string   `json:"servChan3ServerURL" validate:"url,omitempty"`
	CustomPushTags []string `json:"servChan3CustomPushTags" validate:"omitempty"`
	// Bark for ios only
	BarkDeviceKeys           []string                 `json:"barkDeviceKeys"`
	BarkServerURL            string                   `json:"barkServerURL" validate:"url,omitempty"`
	BarkiOSNotificationLevel BarkiOSNotificationLevel `json:"barkiOSNotificationLevel,omitempty" validate:"omitempty,oneof='active' 'passive' 'timeSensitive'"`
}

// PushResponse represent HTTP Response Data from PushNotification Service Provider
type PushResponse struct {
	Code      int
	Message   string
	Timestamp int64
}

func (p *PushResponse) String() string {
	t1 := time.Unix(p.Timestamp, 0)
	return fmt.Sprintf("Response %d: %s at time %s", p.Code, p.Message, t1.Format(time.RFC3339))
}

// pusher is general instance for storing push config and responsible for further data transfer
type pusher struct {
	Config   *PushConfig
	Contents []*BarkPushContent
}

func NewPusher(conf *PushConfig) (*pusher, error) {
	err := verifier.Struct(conf)
	if err != nil {
		return nil, err
	}
	return &pusher{
		Config:   conf,
		Contents: nil,
	}, nil
}
