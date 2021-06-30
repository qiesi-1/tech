package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) SaveTech(tech TechSample) (string, error) {

	eventID := "AddObject"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(tech)
	if err != nil {
		return "", fmt.Errorf("将对象进行序列化时发行错误: %v", err)
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addTech", Args: [][]byte{b, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

func (t *ServiceSetup) UpdateTech(tech TechSample) (string, error) {
	eventID := "UpdateObject"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(tech)
	if err != nil {
		return "", fmt.Errorf("将对象进行序列化时发行错误: %v", err)
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "updateTech", Args: [][]byte{b, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

func (t *ServiceSetup) FindByNameAndMphone(name, mphone string) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryTechByNameAndMphone", Args: [][]byte{[]byte(name), []byte(mphone)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

func (t *ServiceSetup) FindByEntityID(entityID string) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryTechByID", Args: [][]byte{[]byte(entityID)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}
