package initialize


func init() {
	//初始化配置
	InitConf()

	//初始化服务
	InitService()

	//初始化管道
	InitConnection()

	//初始化多个worker
	InitializeServiceLink()
}
