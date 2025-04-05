package pushsdk

// BarkPushContent is an instance of https://github.com/Finb/bark-server/blob/master/docs/API_V2.md
// current version is v2.2.0 for server side, updated at 2025/04/05
type BarkPushContent struct {
	// required
	Title     string `json:"title" validate:"required"`
	Body      string `json:"body" validate:"required"`
	DeviceKey string `json:"device_key" validate:"required"`
	// SubTitle is optional, but for UX, I'll make it required
	SubTitle string `json:"subtitle" validate:"required"`
	// from here, all options below are OPTIONAL
	Level BarkiOSNotificationLevel `json:"level,omitempty"`
	Badge int                      `json:"badge,omitempty" validate:"omitempty,gt=0"`
	// AutomaticallyCopy must be "1" in string
	AutomaticallyCopy string `json:"automaticallyCopy,omitempty" validate:"omitempty,oneof='0' '1'"`
	// Copy store the value to be copied when user click button
	Copy string `json:"copy,omitempty"`
	// alert Sound, from https://github.com/Finb/Bark/tree/master/Sounds
	Sound BarkAlertSound `json:"sound,omitempty"`
	// Notification Icon URL
	Icon string `json:"icon,omitempty"`
	// Group Seperation
	Group string `json:"group,omitempty"`
	// IsArchive must be string, "1" or "0"
	IsArchive string `json:"isArchive,omitempty" validate:"omitempty,oneof='0' '1'"`
	// Jump to URL when clicked
	URL string `json:"url,omitempty" validate:"omitempty,url"`
}

func (bpct *BarkPushContent) Init() {
	bpct.AutomaticallyCopy = "1"
	bpct.IsArchive = "1"
	bpct.Level = ActiveNotification
}

type BarkiOSNotificationLevel string

const (
	ActiveNotification        BarkiOSNotificationLevel = "active"
	TimeSensitiveNotification BarkiOSNotificationLevel = "timeSensitive"
	SilentNotification        BarkiOSNotificationLevel = "passive"
	CriticalNotification      BarkiOSNotificationLevel = "critical"
)

type BarkAlertSound string

const (
	Alarm              BarkAlertSound = "alarm.caf"
	Anticipate         BarkAlertSound = "anticipate.caf"
	Bell               BarkAlertSound = "bell.caf"
	Birdsong           BarkAlertSound = "birdsong.caf"
	Bloom              BarkAlertSound = "bloom.caf"
	Calypso            BarkAlertSound = "calypso.caf"
	Chime              BarkAlertSound = "chime.caf"
	Choo               BarkAlertSound = "choo.caf"
	Descent            BarkAlertSound = "descent.caf"
	Electronic         BarkAlertSound = "electronic.caf"
	Fanfare            BarkAlertSound = "fanfare.caf"
	Glass              BarkAlertSound = "glass.caf"
	Gotosleep          BarkAlertSound = "gotosleep.caf"
	Healthnotification BarkAlertSound = "healthnotification.caf"
	Horn               BarkAlertSound = "horn.caf"
	Ladder             BarkAlertSound = "ladder.caf"
	Mailsent           BarkAlertSound = "mailsent.caf"
	Minuet             BarkAlertSound = "minuet.caf"
	Multiwayinvitation BarkAlertSound = "multiwayinvitation.caf"
	Newmail            BarkAlertSound = "newmail.caf"
	Newsflash          BarkAlertSound = "newsflash.caf"
	Noir               BarkAlertSound = "noir.caf"
	Paymentsuccess     BarkAlertSound = "paymentsuccess.caf"
	Shake              BarkAlertSound = "shake.caf"
	Sherwoodforest     BarkAlertSound = "sherwoodforest.caf"
	Silence            BarkAlertSound = "silence.caf"
	Spell              BarkAlertSound = "spell.caf"
	Suspense           BarkAlertSound = "suspense.caf"
	Telegraph          BarkAlertSound = "telegraph.caf"
	Tiptoes            BarkAlertSound = "tiptoes.caf"
	Typewriters        BarkAlertSound = "typewriters.caf"
	Update             BarkAlertSound = "update.caf"
)
