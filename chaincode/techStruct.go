package main

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