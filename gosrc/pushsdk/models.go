package pushsdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"rdpalert/utils"
	"time"
)

var (
	verifier = validator.New()
	//ErrConfigLogicMismatch    = errors.New("config logic mismatch, required item is not in place")
	ErrGPCIsNotSet            = errors.New("general push content is not staged")
	ErrPushMethodNotSupported = errors.New("push method not supported")
)

type PushContent interface {
	Init()                                                  // Init initiate default content and config
	Provider() PushProvider                                 // Provider return provider name for looking up config
	FromGeneral(g *GeneralPushContent) (PushContent, error) // FromGeneral transform from generalPushContent to specific
	ToBytes() ([]byte, error)                               // ToBytes convert data to bytes then for send out
	SetPushProvider()                                       // SetPushProvider ensures provider info is attached
	AcceptExtParamSettings(d any)                           // AcceptExtParamSettings get config from PushProviderImpl and apply
}

type GeneralPushContent struct {
	Title        string         `json:"title"`
	ShortTitle   string         `json:"short_title"`
	Description  string         `json:"description"`
	ExtParams    map[string]any `json:"ext_params"`
	TagsOrGroups []string       `json:"tags_or_groups"`
	providerName PushProvider
}

func (gpc *GeneralPushContent) SetSpecificPushProvider(p PushProvider) {
	gpc.providerName = p
}

// PushConfig stored user-defined required configuration
type PushConfig struct {
	PushMethods map[PushProvider]json.RawMessage `json:"pushMethods" validate:"required"`
	IsDryRun    bool                             `json:"isDryRun" validate:"boolean,required"`
}

func (pc *PushConfig) VerifyConfig() error {
	err := verifier.Struct(pc)
	if err != nil {
		return err
	}
	return nil
}

// PushProvider records legit and supported push service provider
type PushProvider string

const (
	// BarkForiOS stands for Bark, check: https://github.com/Finb/bark-server
	BarkForiOS PushProvider = "bark"
	// ServChan3 stands for ServChan3 offered by EasyChen, check: https://sc3.ft07.com
	ServChan3 PushProvider = "sc3"
	//ServChanTurbo   PushProvider = "sct"     // not supported yet
	//generalProvider PushProvider = "general" // used for internal data structure only
)

// AbstractPushProvider handle all provider specific issues
// Only used for skeleton and code framework, didn't implement any actual methods
type AbstractPushProvider struct {
	ProviderServerURL string `json:"serverURL" validate:"url,required"`
	ExtraParams       any    `json:"extParams" validate:"omitempty"`
	PushProviderImpl
}
type PushProviderImpl interface {
	VerifyConfig() error
	TransformToSpecificPushContent(g *GeneralPushContent) (PushContent, error)
	SendPushContent(p PushContent) (*PushResponse, error)
}

// PushResponse represent HTTP Response Data from PushNotification Service Provider
// Bark shared the same format of response,
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
	Config               *PushConfig
	GeneralContent       *GeneralPushContent
	SpecificPushContents []*PushContent
}

// NewPusher will validate config and instantiate push service
func NewPusher(conf *PushConfig) (*pusher, error) {
	err := verifier.Struct(conf)
	if err != nil {
		return nil, err
	}
	err = conf.VerifyConfig()
	if err != nil {
		return nil, err
	}
	return &pusher{
		Config:               conf,
		SpecificPushContents: []*PushContent{},
		GeneralContent:       nil,
	}, nil
}

func (p *pusher) StageGeneralPushContent(g *GeneralPushContent) {
	p.GeneralContent = g
}

func (p *pusher) SendPush() error {
	gLogger, err := utils.GetLoggerInstance()
	if err != nil {
		return err
	}
	if p.GeneralContent == nil {
		return ErrGPCIsNotSet
	}
	if p.Config.IsDryRun {
		gLogger.Info("Config Is Set To DryRun, No HTTP Request will be sent.")
		return nil
	}
	// loop through each provider and instantiate
	for k, v := range p.Config.PushMethods {
		switch k {
		case BarkForiOS:
			var prv = &barkPushProvider{}
			err = json.Unmarshal(v, prv)
			if err != nil {
				return err
			}
			err = prv.VerifyConfig()
			if err != nil {
				return err
			}
			spc, err := prv.TransformToSpecificPushContent(p.GeneralContent)
			if err != nil {
				gLogger.Error("Failed to transform to specific push content: ", err.Error())
				return err
			}
			spr, err := prv.SendPushContent(spc)
			if err != nil {
				gLogger.Error("Failed to send push content: ", err.Error())
				return err
			}
			gLogger.Info("Push response received: ", spr.String())
		case ServChan3:
			var prv = &sc3PushProvider{}
			err = json.Unmarshal(v, prv)
			if err != nil {
				return err
			}
			err = prv.VerifyConfig()
			if err != nil {
				return err
			}
			spc, err := prv.TransformToSpecificPushContent(p.GeneralContent)
			if err != nil {
				gLogger.Error("Failed to transform to specific push content: ", err.Error())
				return err
			}
			spr, err := prv.SendPushContent(spc)
			if err != nil {
				gLogger.Error("Failed to send push content: ", err.Error())
				return err
			}
			gLogger.Info("Push response received: ", spr.String())
		default:
			return ErrPushMethodNotSupported
		}
	}
	return nil
}
