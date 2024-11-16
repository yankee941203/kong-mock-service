package controller

import (
	"fmt"
	"kong-mock-service/internal/model"
	"kong-mock-service/internal/repository"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceController interface {
	GetAllServices(ctx *gin.Context)
	GetServiceById(ctx *gin.Context)
}

type ServiceControllerImp struct {
	serviceRepository repository.ServiceRepository
}

func NewServiceControllerImp(repo repository.ServiceRepository) *ServiceControllerImp {
	return &ServiceControllerImp{
		serviceRepository: repo,
	}
}

func (c *ServiceControllerImp) GetAllServices(ctx *gin.Context) {
	temp := "/api/v1/services?page=%d&per_page=%d"
	name := ctx.Query("name")

	// Check page or page_size is valid or not
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	pnum, err := strconv.Atoi(page)
	if err != nil {
		log.Printf("[ERROR]: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
	psnum, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("[ERROR]: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
	}

	// Check sort_by field is valid or not
	sortBy := ctx.DefaultQuery("sort_by", "ID")
	v := reflect.ValueOf(model.ServiceInfo{})
	if !v.FieldByName(sortBy).IsValid() {
		msg := fmt.Sprintf("[ERROR]: Sort by field %s is invalid", sortBy)
		log.Println(msg)
		ctx.JSON(http.StatusBadRequest, msg)
	}

	ses, err := c.serviceRepository.GetAllServicesWithImages()
	if err != nil {
		log.Printf("[ERROR]: %v", err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	sis := []model.ServiceInfo{}
	for _, se := range ses {
		si := ConvertServicesRntityToInfo(se)
		sis = append(sis, si)
	}
	resInfo := model.ResponseInfo{}
	resInfo.Data = sis
	resInfo.Meta = model.MetaInfo{}
	if name != "" {
		FilterServiceResponseByName(&resInfo, name)
		temp = fmt.Sprintf("%s&name=%s", temp, name)
	}
	SortServiceResponseByFieldName(&resInfo, sortBy)
	temp = fmt.Sprintf("%s&sort_by=%s", temp, sortBy)
	PagingServicesResponse(&resInfo, pnum, psnum, temp)
	ctx.JSON(http.StatusOK, resInfo)
}

func (c *ServiceControllerImp) GetServiceById(ctx *gin.Context) {
	id := ctx.Param("id")
	idnum, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("[ERROR]: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
	se, err := c.serviceRepository.GetServiceByIdWithImages(uint(idnum))
	if err != nil {
		log.Printf("[ERROR]: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
	si := ConvertServicesRntityToInfo(se)
	ctx.JSON(http.StatusOK, si)
}

func FilterServiceResponseByName(ri *model.ResponseInfo, name string) {
	sis := ri.Data
	fsis := []model.ServiceInfo{}
	for _, si := range sis {
		if si.Name == name {
			fsis = append(fsis, si)
		}
	}
	ri.Data = fsis
}

func SortServiceResponseByFieldName(ri *model.ResponseInfo, fieldName string) {
	sis := ri.Data
	sort.Slice(sis, func(i, j int) bool {
		valI := reflect.ValueOf(sis[i]).FieldByName(fieldName).Interface()
		valJ := reflect.ValueOf(sis[j]).FieldByName(fieldName).Interface()
		switch valI.(type) {
		case int:
			return valI.(int) < valJ.(int)
		case float64:
			return valI.(float64) < valJ.(float64)
		case string:
			return valI.(string) < valJ.(string)
		default:
			return false
		}
	})
}

func PagingServicesResponse(ri *model.ResponseInfo, page int, pageSize int, linktemp string) {
	ri.Meta.Page = page
	ri.Meta.PerPage = pageSize
	ti := len(ri.Data)
	ri.Meta.TotalItems = ti
	tp := ti / pageSize
	if ti%pageSize > 0 {
		tp += 1
	}
	ri.Meta.TotalPages = tp
	pp := page - 1
	np := page + 1
	lp := tp
	if pp > 0 {
		ri.Meta.PrevPage = fmt.Sprintf(linktemp, pp, pageSize)
	}
	if np <= lp {
		ri.Meta.NextPage = fmt.Sprintf(linktemp, np, pageSize)
	}
	ri.Meta.LastPage = fmt.Sprintf(linktemp, lp, pageSize)

	sis := ri.Data
	start := pageSize * (page - 1)
	end := start + pageSize
	ri.Data = sis[start:end]
}

func ConvertServicesRntityToInfo(e model.ServiceEntity) model.ServiceInfo {
	avs := []string{}
	imgs := e.Images
	for _, img := range imgs {
		av := fmt.Sprintf("%s:%s", img.Name, img.Version)
		avs = append(avs, av)
	}
	return model.ServiceInfo{
		ID:               e.ID,
		Name:             e.Name,
		Description:      e.Description,
		AvaliableVersion: avs,
	}
}
