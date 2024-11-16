package main

import (
	"kong-mock-service/internal/controller"
	"kong-mock-service/internal/repository"
	"kong-mock-service/internal/router"
)

func main() {
	sr := repository.NewServiceRepositoryDbImp()
	sr.InitDb()
	sc := controller.NewServiceControllerImp(sr)
	ro := router.NewRouter(sc)
	r := ro.SetupRouter()
	r.Run(":8000")
}
