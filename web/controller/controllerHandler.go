package controller

import (
	"tech/service"
	"encoding/json"
	"net/http"
)

func (app *WebApplication) Index(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "index.html", nil)
}

// 进入添加信息页面
func (app *WebApplication) SavePage(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "addInfo.html", nil)
}

// 保存信息
func (app *WebApplication) SaveTech(w http.ResponseWriter, r *http.Request)  {
	data := struct {
		Tech service.TechSample
		Msg string
		Flag bool
	}{}

	flag := r.FormValue("f")
	if flag == "f" {
		ShowView(w, r, "addInfo.html", data)
		return
	}

	tech := service.TechSample{
		Name: r.FormValue("name"),
		EntityID: r.FormValue("entityID"),
		Gender: r.FormValue("gender"),
		Mphone: r.FormValue("mphone"),
	}

	txid, err := app.Service.SaveTech(tech)

	if err != nil {
		data.Msg = err.Error()
	} else {
		data.Tech = tech
		data.Msg = "数据保存成功, TxID: " + txid
	}
	data.Flag = true

	ShowView(w, r, "addInfo.html", data)
}

// 进入修改信息页面
func (app *WebApplication) UpdatePage(w http.ResponseWriter, r *http.Request)  {

}

// 修改信息
func (app *WebApplication) UpdateTech(w http.ResponseWriter, r *http.Request)  {

}

// 根据名称及联系电话查询信息
func (app *WebApplication) FindByName(w http.ResponseWriter, r *http.Request)  {

}

func (app *WebApplication) FindHistoryPage(w http.ResponseWriter, r *http.Request)  {
	ShowView(w, r, "queryHistory.html", nil)
}

// 根据ID查询信息
func (app *WebApplication) FindByEntityID(w http.ResponseWriter, r *http.Request)  {
	id := r.FormValue("entityID")
	result, err := app.Service.FindByEntityID(id)
	var tech service.TechSample
	json.Unmarshal(result, &tech)

	data := struct {
		Tech service.TechSample
		Msg string
		Flag bool
		History bool
	}{
		Tech:tech,
		Msg:"",
		Flag:false,
		History:true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "historyResult.html", data)

}


