package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type TechChaincode struct {

}

func (t *TechChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response  {

	return shim.Success(nil)
}

func (t *TechChaincode) Invoke(stub shim.ChaincodeStubInterface)  peer.Response{
	fun, args := stub.GetFunctionAndParameters()
	if fun == "addTech" {
		return t.addTech(stub, args)
	} else if fun == "updateTech" {
		return t.updateTech(stub, args)
	} else if fun == "queryTechByNameAndMphone" {
		return t.queryTechByNameAndMphone(stub, args)
	} else if fun == "queryTechByID" {
		return t.queryTechByID(stub, args)
	}
	return peer.Response(shim.Error("指定的函数名称错误"))
}

func main()  {
	err := shim.Start(new(TechChaincode))
	if err != nil{
		fmt.Printf("启动TechChaincode时发生错误: %s", err)
	}
}
