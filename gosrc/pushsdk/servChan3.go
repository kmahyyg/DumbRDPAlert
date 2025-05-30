package pushsdk

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type sc3PushContent struct {
	Title         string   `json:"title" validate:"required"`
	Description   string   `json:"desp" validate:"required"`
	TagsStr       string   `json:"tags,omitempty"`
	ShortBriefing string   `json:"short,omitempty"`
	Tags          []string `json:"-"`

	providerName PushProvider
}

func (sc3p *sc3PushContent) Provider() PushProvider {
	return sc3p.providerName
}

func (sc3p *sc3PushContent) FromGeneral(g *GeneralPushContent) (PushContent, error) {
	sc3p.Title = g.Title
	sc3p.ShortBriefing = g.ShortTitle
	sc3p.Description = g.Description
	if len(sc3p.Tags) != 0 {
		sc3p.TagsStr = strings.Join(sc3p.Tags, "|")
	}
	return sc3p, nil
}

func (sc3p *sc3PushContent) ToBytes() ([]byte, error) {
	return json.Marshal(sc3p)
}

func (sc3p *sc3PushContent) SetPushProvider() {
	sc3p.providerName = ServChan3
}

func (sc3p *sc3PushContent) Init() {
	sc3p.SetPushProvider()
}

func (sc3p *sc3PushContent) AcceptExtParamSettings(d any) {
	d1 := d.(*sc3PushProviderExtraParam)
	if len(d1.CustomPushTags) != 0 {
		sc3p.Tags = d1.CustomPushTags
	}
	return
}

type sc3PushProvider struct {
	// ServChan 3 for mobile universal
	ProviderServerURL string                     `json:"serverURL" validate:"url,required"`
	ExtraParams       *sc3PushProviderExtraParam `json:"extParams" validate:"omitempty"`
}

func (s sc3PushProvider) VerifyConfig() error {
	err1 := verifier.Struct(s)
	if err1 != nil {
		return err1
	}
	err2 := verifier.Struct(s.ExtraParams)
	if err2 != nil {
		return err2
	}
	return nil
}

func (s sc3PushProvider) TransformToSpecificPushContent(g *GeneralPushContent) (PushContent, error) {
	sc3p := &sc3PushContent{}
	sc3p.Init()
	sc3p.AcceptExtParamSettings(s.ExtraParams)
	return sc3p.FromGeneral(g)
}

func (s sc3PushProvider) SendPushContent(p PushContent) (*PushResponse, error) {
	pData := p.(*sc3PushContent)
	body, err := pData.ToBytes()
	if err != nil {
		return nil, err
	}
	respData, _, err := SendHttpPostJSON(s.ProviderServerURL, body)
	if err != nil {
		return nil, err
	}
	sc3pr := &sc3PushResponse{}
	err = json.Unmarshal(respData, sc3pr)
	if err != nil {
		return nil, err
	}
	gpr, err := sc3pr.ToGeneralPushResponse()
	if err != nil {
		return nil, err
	}
	return gpr, nil
}

type sc3PushResponse struct {
	Code    int `json:"code"`
	ErrorNo int `json:"errno"`
	Data    struct {
		PushID int `json:"pushid"`
		Meta   struct {
			Android struct {
				MessageIds struct {
					MessageId []string `json:"messageId,omitempty"`
				} `json:"messageIds,omitempty"`
				RequestID string `json:"requestId,omitempty"`
			} `json:"android,omitempty"`
			IOS struct {
				MessageIds struct {
					MessageId []string `json:"messageId,omitempty"`
				} `json:"messageIds,omitempty"`
				RequestID string `json:"requestId,omitempty"`
			} `json:"ios,omitempty"`
		} `json:"meta,omitempty"`
	} `json:"data"`
	Message string `json:"message"`
}

func (sc3pr *sc3PushResponse) ToGeneralPushResponse() (*PushResponse, error) {
	respMsg := fmt.Sprintf("ErrorNo: %d , PushID: %d, OriRespMsg: %s, AliyunPushMsgIDs(Android): %v - Req: %s , AliyunPushMsgIDs(iOS): %v - Req: %s ", sc3pr.ErrorNo, sc3pr.Data.PushID,
		sc3pr.Message, sc3pr.Data.Meta.Android.MessageIds.MessageId, sc3pr.Data.Meta.Android.RequestID, sc3pr.Data.Meta.IOS.MessageIds.MessageId, sc3pr.Data.Meta.IOS.RequestID)
	gpr := &PushResponse{
		Code:      sc3pr.Code,
		Message:   respMsg,
		Timestamp: time.Now().Unix(),
	}
	return gpr, nil
}

type sc3PushProviderExtraParam struct {
	CustomPushTags []string `json:"customPushTags" validate:"omitempty"`
}
