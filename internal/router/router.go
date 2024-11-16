package router

import (
	"kong-mock-service/internal/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	serviceController controller.ServiceController
}

func NewRouter(serviceController controller.ServiceController) *Router {
	return &Router{
		serviceController: serviceController,
	}
}

func (r *Router) SetupRouter() *gin.Engine {
	e := gin.Default()

	//Define API
	e.GET("/api/v1/services", r.serviceController.GetAllServices)
	e.GET("/api/v1/service/:id", r.serviceController.GetServiceById)

	return e
}
