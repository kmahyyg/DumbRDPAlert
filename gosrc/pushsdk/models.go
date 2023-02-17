package pushsdk

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

var verifier = validator.New()

type pusher struct {
	Config   *PushConfig
	l        *DumbLogger
	Contents []*PushContent
}

func NewPusher(conf *PushConfig, logger *DumbLogger) (*pusher, error) {
	err := verifier.Struct(conf)
	if err != nil {
		return nil, err
	}
	return &pusher{
		Config:   conf,
		l:        logger,
		Contents: nil,
	}, nil
}

type PushConfig struct {
	DeviceKeys        []string          `json:"deviceKeys" validate:"required"`
	ServerAddress     string            `json:"serverAddress" validate:"url,required"`
	NotificationLevel NotificationLevel `json:"notificationLevel,omitempty" validate:"oneof='active' 'passive' 'timeSensitive',omitempty"`
}

// PushContent is an instance of https://github.com/Finb/bark-server/blob/master/docs/API_V2.md
type PushContent struct {
	// required
	Title     string `json:"title" validate:"required"`
	Body      string `json:"body" validate:"required"`
	Category  string `json:"category" validate:"len=0"`
	DeviceKey string `json:"device_key" validate:"required"`
	// optional, from here
	Level NotificationLevel `json:"level,omitempty"`
	Badge int               `json:"badge,omitempty" validate:"gt=0,omitempty"`
	// must be "1"
	AutomaticallyCopy string `json:"automaticallyCopy,omitempty" validate:"oneof='0' '1',omitempty"`
	// value to be copied
	Copy string `json:"copy,omitempty"`
	// alert sound, from https://github.com/Finb/Bark/tree/master/Sounds
	Sound AlertSound `json:"sound,omitempty"`
	// Notification Icon URL
	Icon string `json:"icon,omitempty"`
	// Group Seperation
	Group string `json:"group,omitempty"`
	// isArchive must be string, "1" or "0"
	IsArchive string `json:"isArchive,omitempty" validate:"oneof='0' '1',omitempty"`
	// Jump to URL when clicked
	URL string `json:"URL,omitempty" validate:"url,omitempty"`
}

func (pct *PushContent) Init() {
	pct.AutomaticallyCopy = "1"
	pct.IsArchive = "0"
	pct.Category = ""
	pct.Level = ActiveNotification
}

type PushResponse struct {
	Code      int
	Message   string
	Timestamp int64
}

func (p *PushResponse) String() string {
	t1 := time.Unix(p.Timestamp, 0)
	return fmt.Sprintf("Response %d: %s at time %s", p.Code, p.Message, t1.Format(time.RFC3339))
}

type NotificationLevel string

const (
	ActiveNotification        NotificationLevel = "active"
	TimeSensitiveNotification NotificationLevel = "timeSensitive"
	SilentNotification        NotificationLevel = "passive"
)

type AlertSound string

const (
	Alarm              AlertSound = "alarm.caf"
	Anticipate         AlertSound = "anticipate.caf"
	Bell               AlertSound = "bell.caf"
	Birdsong           AlertSound = "birdsong.caf"
	Bloom              AlertSound = "bloom.caf"
	Calypso            AlertSound = "calypso.caf"
	Chime              AlertSound = "chime.caf"
	Choo               AlertSound = "choo.caf"
	Descent            AlertSound = "descent.caf"
	Electronic         AlertSound = "electronic.caf"
	Fanfare            AlertSound = "fanfare.caf"
	Glass              AlertSound = "glass.caf"
	Gotosleep          AlertSound = "gotosleep.caf"
	Healthnotification AlertSound = "healthnotification.caf"
	Horn               AlertSound = "horn.caf"
	Ladder             AlertSound = "ladder.caf"
	Mailsent           AlertSound = "mailsent.caf"
	Minuet             AlertSound = "minuet.caf"
	Multiwayinvitation AlertSound = "multiwayinvitation.caf"
	Newmail            AlertSound = "newmail.caf"
	Newsflash          AlertSound = "newsflash.caf"
	Noir               AlertSound = "noir.caf"
	Paymentsuccess     AlertSound = "paymentsuccess.caf"
	Shake              AlertSound = "shake.caf"
	Sherwoodforest     AlertSound = "sherwoodforest.caf"
	Silence            AlertSound = "silence.caf"
	Spell              AlertSound = "spell.caf"
	Suspense           AlertSound = "suspense.caf"
	Telegraph          AlertSound = "telegraph.caf"
	Tiptoes            AlertSound = "tiptoes.caf"
	Typewriters        AlertSound = "typewriters.caf"
	Update             AlertSound = "update.caf"
)
