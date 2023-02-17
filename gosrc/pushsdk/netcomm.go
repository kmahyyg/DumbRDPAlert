package pushsdk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const postJSONContentType = "application/json; charset=utf-8"

func (p *pusher) SetContents(conts []*PushContent) error {
	for i := 0; i < len(conts); i++ {
		err := verifier.Struct(conts[i])
		if err != nil {
			return err
		}
	}
	p.Contents = conts
	return nil
}

func (p *pusher) SendPushRequests() error {
	for _, v := range p.Contents {
		bodyJson, err := json.Marshal(v)
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(bodyJson)
		resp, err := http.Post(p.Config.ServerAddress, postJSONContentType, buf)
		if err != nil {
			return err
		}
		if resp.StatusCode == 200 {
			respJ, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			respI := &PushResponse{}
			err = json.Unmarshal(respJ, respI)
			if err != nil {
				return err
			}
			p.l.Info(respI.String())
		} else {
			p.l.Error("Push Request to Server Error, status: ", resp.StatusCode)
		}
	}
	return nil
}
