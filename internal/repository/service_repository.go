package repository

import (
	"fmt"
	"kong-mock-service/internal/model"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DB_PATH = "sources/test.db"

type ServiceRepository interface {
	GetAllServices() ([]model.ServiceEntity, error)
	GetAllServicesWithImages() ([]model.ServiceEntity, error)
	GetServiceByIdWithImages(id uint) (model.ServiceEntity, error)
}

type ServiceRepositoryDbImp struct {
	db *gorm.DB
}

func NewServiceRepositoryDbImp() *ServiceRepositoryDbImp {
	db, err := gorm.Open(sqlite.Open(DB_PATH), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &ServiceRepositoryDbImp{
		db: db,
	}
}

func (r *ServiceRepositoryDbImp) InitDb() {
	r.db.AutoMigrate(&model.ServiceEntity{})
	err := InsertMockData(r.db)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *ServiceRepositoryDbImp) GetAllServices() ([]model.ServiceEntity, error) {
	var services []model.ServiceEntity
	err := r.db.Model(&model.ServiceEntity{}).Find(&services).Error
	return services, err
}

func (r *ServiceRepositoryDbImp) GetAllServicesWithImages() ([]model.ServiceEntity, error) {
	var services []model.ServiceEntity
	err := r.db.Model(&model.ServiceEntity{}).Preload("Images").Find(&services).Error
	return services, err
}

func (r *ServiceRepositoryDbImp) GetServiceByIdWithImages(id uint) (model.ServiceEntity, error) {
	var service model.ServiceEntity
	err := r.db.Model(&model.ServiceEntity{}).Preload("Images").Where("id=?", id).First(&service).Error
	return service, err
}

func InsertMockData(db *gorm.DB) error {
	mimgs := []*model.ImageEntity{}
	name := "pokemon/hudi"
	for i := range 6 {
		img := model.ImageEntity{}
		if i%2 == 1 {
			name = "pokemon/pika"
		}
		img.Name = name
		img.Version = fmt.Sprintf("v%d", i/2)
		mimgs = append(mimgs, &img)
	}
	res := db.Create(mimgs)
	err := res.Error
	if err != nil {
		log.Fatal(err)
	}

	msvcs := []*model.ServiceEntity{}
	for i := range 30 {
		simgs := []*model.ImageEntity{}
		svc := model.ServiceEntity{}
		svc.Name = fmt.Sprintf("kong-test-%d", i)
		svc.Description = "a test service"
		for j, img := range mimgs {
			if j == i%6 {
				simgs = append(simgs, img)
			}
			if j == i%4 {
				simgs = append(simgs, img)
			}
		}
		svc.Images = simgs
		msvcs = append(msvcs, &svc)
	}
	res = db.Create(msvcs)
	err = res.Error
	if err != nil {
		log.Fatal(err)
	}
	return err
}
