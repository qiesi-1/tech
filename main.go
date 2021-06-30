package main

import (
	"tech/sdkInit"
	"tech/service"
	"encoding/json"
	"fmt"
	"os"
	"tech/web"
	"tech/web/controller"
)

const (
	configFile = "config.yaml"
	initialized = false
	SimpleCC = "simplecc"
)

func main()  {
	initInfo := &sdkInit.InitInfo {
		ChannelID: "mychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/tech/fixtures/artifacts/channel.tx",

		OrgAdmin:"Admin",
		OrgName:"Org1",
		OrdererOrgName: "orderer.example.com",

		ChaincodeID: SimpleCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath: "tech/chaincode/",
		UserName: "User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()

	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)


	// Service
	serviceSetup := service.ServiceSetup{
		ChaincodeID: SimpleCC,
		Client: channelClient,
	}

	tech := service.TechSample{
		Name: "Jack",
		EntityID: "123456",
		Gender: "男",
		Mphone: "13100000000",
	}
	msg, err := serviceSetup.SaveTech(tech)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息添加成功, 交易编号为: " + msg)
	}

	result, err := serviceSetup.FindByNameAndMphone("Jack", "13100000000")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var hw service.TechSample
		json.Unmarshal(result, &hw)
		fmt.Println("根据姓名与ID查询信息成功：")
		fmt.Println(hw)
	}

	// 更新信息
	sample := service.TechSample{
		Name: "Jack",
		EntityID: "123456",
		Gender: "男",
		Mphone: "13300000000",
	}
	msg, err = serviceSetup.UpdateTech(sample)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息修改成功, 交易编号为: " + msg)
	}

	// 查询历史记录
	result, err = serviceSetup.FindByEntityID("123456")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var hw service.TechSample
		json.Unmarshal(result, &hw)
		fmt.Println("根据ID查询信息成功：")
		fmt.Println(hw)
	}


	//Web
	app := controller.WebApplication {
		Service: &serviceSetup,
	}
	web.WebStart(app)
}



