package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"

)

const DOC_TYPE = "techObj"

// 添加信息
// EntityID为 key, HelloWorld为 value
func (t *TechChaincode) addTech(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var tech TechSample
	err := json.Unmarshal([]byte(args[0]), &tech)
	if err != nil {
		return peer.Response(shim.Error("反序列化信息时发生错误，添加失败: "  + err.Error()))
	}
	tech.ObjectType = DOC_TYPE

	// 查询要添加的EntityID是否已存在省略...

	// 保存至账本中
	b, err := json.Marshal(tech)
	if err != nil {
		return shim.Error("序列化数据时发生错误: " + err.Error())
	}
	err = stub.PutState(tech.EntityID, b)
	if err != nil {
		return shim.Error("保存数据时发生错误: " + err.Error())
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}

// 根据EntityID更新信息
func (t *TechChaincode) updateTech(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数错误")
	}

	var tech TechSample
	err := json.Unmarshal([]byte(args[0]), &tech)
	if err != nil {
		return shim.Error("反序列化对象时发生错误: " + err.Error())
	}

	// 根据EntityID查询信息时发生错误
	b, err := stub.GetState(tech.EntityID)
	if err != nil {
		return shim.Error("根据身份证号码查询信息时发生错误")
	}

	if b == nil {
		return shim.Error("根据指定的EntityID没有查询到相应的信息")
	}

	var sample TechSample
	err = json.Unmarshal(b, &sample)
	if err != nil {
		return shim.Error("对查询结果进行序列化时发生错误: " + err.Error())
	}

	sample.Name = tech.Name
	sample.Gender = tech.Gender
	sample.Mphone = tech.Mphone

	// 保存修改之后的信息
	sample.ObjectType = DOC_TYPE
	br, err := json.Marshal(sample)
	stub.PutState(sample.EntityID, br)

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息修改成功"))
}

// 根据姓名及联系电话查询信息
func (t *TechChaincode) queryTechByNameAndMphone(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数不符合要求")
	}
	name := args[0]
	mphone := args[1]

	fmt.Println("========> Name = " + name + ", Mphone = " + mphone)

	// 拼接查询字符串
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"Name\":\"%s\", \"Mphone\":\"%s\"}}", DOC_TYPE, name, mphone)

	fmt.Println("QueryString = " + queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error("查询数据时出现错误: " + err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("根据姓名及联系电话查询信息时发生错误: " + err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	// 查询数据
	result := buffer.Bytes()

	if result == nil {
		return shim.Error("根据姓名及联系电话没有查询到相关的信息")
	}
	return shim.Success(result)

}

// 根据指定的EntityID查询信息(溯源实现)
func (t *TechChaincode) queryTechByID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数不符合要求，必须指定为EntityID")
	}

	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据指定的EntityID查询数据时发生错误: " + err.Error())
	}

	if b == nil {
		return shim.Error("根据指定的Entity没有查询到相关的数据")
	}

	var tech TechSample
	err = json.Unmarshal(b, &tech)
	if err != nil {
		return shim.Error("对查询结果进行反序列化时发生错误: " + err.Error())
	}

	// 获取历史变更数据
	iterator, err := stub.GetHistoryForKey(tech.EntityID)
	if err != nil {
		return shim.Error("根据指定的EntityID查询历史变更记录时发生错误: " + err.Error())
	}
	defer iterator.Close()

	var historys []HistoryItem
	var hisHello TechSample
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取历史数据错误: " + err.Error())
		}

		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisHello)
		if hisData.Value == nil {
			var empty TechSample
			historyItem.TechSample = empty
		}else{
			historyItem.TechSample = hisHello
		}

		historys = append(historys, historyItem)
	}

	tech.Historys = historys

	result, err := json.Marshal(tech)
	if err != nil {
		return shim.Error("序列化查询结果时发生错误: " + err.Error())
	}
	return shim.Success(result)
}
