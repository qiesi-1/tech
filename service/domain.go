package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)

type TechSample struct {
	ObjectType	string	`json:"docType"`
	Name	string	`json:"Name"`
	EntityID	string	`json:"EntityID"`
	Gender	string	`json:"Gender"`
	Mphone	string	`json:"Mphone"`

	Historys []HistoryItem	// 当前对象的历史记录
}


type HistoryItem struct {
	TxId	string
	TechSample TechSample
}

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

func registerEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {
	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
		return nil
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的EventID接收到相应的链码事件(%s)", eventID)
	}
	//return nil
}