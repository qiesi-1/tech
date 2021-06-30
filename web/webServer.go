package web

import (
	"tech/web/controller"
	"fmt"
	"net/http"
)



func WebStart(app controller.WebApplication){

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 路由
	http.HandleFunc("/", app.Index)
	http.HandleFunc("/index", app.Index)

	http.HandleFunc("/savePage", app.SavePage)	// 进入添加信息页面
	http.HandleFunc("/save", app.SaveTech)	// 保存信息

	http.HandleFunc("/updatePage", app.UpdatePage)	// 进入修改信息页面
	http.HandleFunc("/update", app.UpdateTech)	// 修改信息


	http.HandleFunc("/findbyname", app.FindByName)	// 根据名称及联系电话查询信息

	http.HandleFunc("/findHistoryPage", app.FindHistoryPage)
	http.HandleFunc("/findById", app.FindByEntityID)	// 根据ID查询信息(历史记录信息)

	fmt.Println("启动Web服务, 监听端口为: 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("启动Web服务失败: %s\n", err.Error())
		return
	}

}
