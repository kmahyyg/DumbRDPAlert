package pushsdk

import (
	"encoding/json"
	"rdpalert/utils"
)

// barkPushContent is an instance of https://github.com/Finb/bark-server/blob/master/docs/API_V2.md
// current version is v2.2.0 for server side, updated at 2025/04/05
type barkPushContent struct {
	// required
	Title     string `json:"title" validate:"required"`
	Body      string `json:"body" validate:"required"`
	DeviceKey string `json:"device_key" validate:"required"`
	// SubTitle is optional, but for UX, I'll make it required
	SubTitle string `json:"subtitle" validate:"required"`
	// from here, all options below are OPTIONAL
	Level barkiOSNotificationLevel `json:"level,omitempty"`
	Badge int                      `json:"badge,omitempty" validate:"omitempty,gt=0"`
	// AutomaticallyCopy must be "1" in string
	AutomaticallyCopy string `json:"automaticallyCopy,omitempty" validate:"omitempty,oneof='0' '1'"`
	// Copy store the value to be copied when user click button
	Copy string `json:"copy,omitempty"`
	// alert Sound, from https://github.com/Finb/Bark/tree/master/Sounds
	Sound barkAlertSound `json:"sound,omitempty"`
	// Notification Icon URL
	Icon string `json:"icon,omitempty"`
	// Group Seperation
	Group string `json:"group,omitempty"`
	// IsArchive must be string, "1" or "0"
	IsArchive string `json:"isArchive,omitempty" validate:"omitempty,oneof='0' '1'"`
	// Jump to URL when clicked
	URL string `json:"url,omitempty" validate:"omitempty,url"`

	// provider info
	providerName PushProvider
}

func (bpct *barkPushContent) FromGeneral(g GeneralPushContent) (PushContent, error) {
	//TODO implement me
	panic("implement me")
}

func (bpct *barkPushContent) ToBytes() []byte {
	//TODO implement me
	panic("implement me")
}

func (bpct *barkPushContent) Init() {
	bpct.AutomaticallyCopy = "1"
	bpct.IsArchive = "1"
	bpct.providerName = BarkForiOS
	bpct.Group = "Sec_RdpAlert"
	bpct.Level = ActiveNotification
}

func (bpct *barkPushContent) Provider() PushProvider {
	return bpct.providerName
}

func (bpct *barkPushContent) SetPushProvider() {
	bpct.providerName = BarkForiOS
}

func (bpct *barkPushContent) AcceptExtParamSettings(d any) {
	d1 := d.(*barkPushProviderExtraParams)
	if d1.IOSNotificationLevel != "" {
		bpct.Level = d1.IOSNotificationLevel
	}
	if d1.NotificationGroup != "" {
		bpct.Group = d1.NotificationGroup
	}
}

type barkPushProvider struct {
	ProviderServerURL string                       `json:"serverURL" validate:"url,required"`
	ExtraParams       *barkPushProviderExtraParams `json:"extParams" validate:"omitempty"`
}

func (b barkPushProvider) VerifyConfig() error {
	err1 := verifier.Struct(b)
	if err1 != nil {
		return err1
	}
	err2 := verifier.Struct(b.ExtraParams)
	if err2 != nil {
		return err2
	}
	return nil
}

func (b barkPushProvider) TransformToSpecificPushContent(g GeneralPushContent) (PushContent, error) {
	//TODO implement me
	panic("implement me")
}

func (b barkPushProvider) SendPushContent(p PushContent) (*PushResponse, error) {
	gLogger, err := utils.GetLoggerInstance()
	if err != nil {
		return nil, err
	}
	pData := p.(*barkPushContent)
	body := pData.ToBytes()
	respData, _, err := SendHttpPostJSON(b.ProviderServerURL, body)
	pushResp := &PushResponse{}
	err = json.Unmarshal(respData, pushResp)
	if err != nil {
		return nil, err
	}
	gLogger.Info("pushResp: ", pushResp.String())
	return pushResp, nil
}

type barkPushProviderExtraParams struct {
	// Bark for ios only
	NotificationGroup    string                   `json:"notificationGrp,omitempty" validate:"omitempty"`
	DeviceKeys           []string                 `json:"deviceKeys" validate:"required"`
	IOSNotificationLevel barkiOSNotificationLevel `json:"iOSNotificationLvl,omitempty" validate:"omitempty,oneof='active' 'passive' 'timeSensitive' 'critical'"`
}

type barkiOSNotificationLevel string

const (
	ActiveNotification        barkiOSNotificationLevel = "active"
	TimeSensitiveNotification barkiOSNotificationLevel = "timeSensitive"
	SilentNotification        barkiOSNotificationLevel = "passive"
	CriticalNotification      barkiOSNotificationLevel = "critical"
)

type barkAlertSound string

const (
	Alarm              barkAlertSound = "alarm.caf"
	Anticipate         barkAlertSound = "anticipate.caf"
	Bell               barkAlertSound = "bell.caf"
	Birdsong           barkAlertSound = "birdsong.caf"
	Bloom              barkAlertSound = "bloom.caf"
	Calypso            barkAlertSound = "calypso.caf"
	Chime              barkAlertSound = "chime.caf"
	Choo               barkAlertSound = "choo.caf"
	Descent            barkAlertSound = "descent.caf"
	Electronic         barkAlertSound = "electronic.caf"
	Fanfare            barkAlertSound = "fanfare.caf"
	Glass              barkAlertSound = "glass.caf"
	Gotosleep          barkAlertSound = "gotosleep.caf"
	Healthnotification barkAlertSound = "healthnotification.caf"
	Horn               barkAlertSound = "horn.caf"
	Ladder             barkAlertSound = "ladder.caf"
	Mailsent           barkAlertSound = "mailsent.caf"
	Minuet             barkAlertSound = "minuet.caf"
	Multiwayinvitation barkAlertSound = "multiwayinvitation.caf"
	Newmail            barkAlertSound = "newmail.caf"
	Newsflash          barkAlertSound = "newsflash.caf"
	Noir               barkAlertSound = "noir.caf"
	Paymentsuccess     barkAlertSound = "paymentsuccess.caf"
	Shake              barkAlertSound = "shake.caf"
	Sherwoodforest     barkAlertSound = "sherwoodforest.caf"
	Silence            barkAlertSound = "silence.caf"
	Spell              barkAlertSound = "spell.caf"
	Suspense           barkAlertSound = "suspense.caf"
	Telegraph          barkAlertSound = "telegraph.caf"
	Tiptoes            barkAlertSound = "tiptoes.caf"
	Typewriters        barkAlertSound = "typewriters.caf"
	Update             barkAlertSound = "update.caf"
)
