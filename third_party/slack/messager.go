package slack

import (
	"TrackMaster/pkg"
	"TrackMaster/third_party"
	"encoding/json"
)

type payload struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

var Messager = &third_party.ThirdPartyDataFetcher{
	Host:    "http://alert-trigger.infra.svc.cluster.local:5556",
	Path:    "/alert/slack",
	Query:   nil,
	OnError: nil,
}

func SendMessage(message string) *pkg.Error {
	// channel 暂时写死是C02RWG9QGDP
	pl := payload{
		Channel: "C02RWG9QGDP",
		Message: message,
	}
	reqBody, err1 := json.Marshal(pl)
	if err1 != nil {
		return pkg.NewError(pkg.ServerError, err1.Error())
	}

	_, err := Messager.PatchData("POST", "", reqBody)
	if err != nil {
		return err
	}

	return nil
}
