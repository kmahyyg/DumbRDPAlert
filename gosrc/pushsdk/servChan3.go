package pushsdk

type sc3PushContent struct {
	providerName PushProvider
}

func (sc3p *sc3PushContent) Provider() PushProvider {
	return sc3p.providerName
}

func (sc3p *sc3PushContent) FromGeneral(g GeneralPushContent) (PushContent, error) {
	//TODO implement me
	panic("implement me")
}

func (sc3p *sc3PushContent) ToBytes() []byte {
	//TODO implement me
	panic("implement me")
}

func (sc3p *sc3PushContent) SetPushProvider() {
	sc3p.providerName = ServChan3
}

func (sc3p *sc3PushContent) Init() {

}

func (sc3p *sc3PushContent) AcceptExtParamSettings(d any) {

}

type sc3PushProvider struct {
	// ServChan 3 for mobile universal
	ProviderServerURL string                     `json:"serverURL" validate:"url,required"`
	ExtraParams       *sc3PushProviderExtraParam `json:"extParams" validate:"omitempty"`
}

func (s sc3PushProvider) VerifyConfig() error {
	//TODO implement me
	panic("implement me")
}

func (s sc3PushProvider) TransformToSpecificPushContent(g GeneralPushContent) (PushContent, error) {
	//TODO implement me
	panic("implement me")
}

func (s sc3PushProvider) SendPushContent(p PushContent) (*PushResponse, error) {
	//TODO implement me
	panic("implement me")
}

type sc3PushResponse struct {
}

type sc3PushProviderExtraParam struct {
	CustomPushTags []string `json:"customPushTags" validate:"omitempty"`
}
