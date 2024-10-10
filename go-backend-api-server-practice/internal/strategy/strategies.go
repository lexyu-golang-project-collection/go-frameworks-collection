package strategy

import (
	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/model"
)

type CarBuildStrategy interface {
	BuildCar() *model.Car
}

type SportCarStrategy struct{}

func (scs *SportCarStrategy) BuildCar() *model.Car {
	return &model.Car{Brand: "Sport-Ferrari", Model: "F8 Tributo"}
}

type FamilyCarStrategy struct{}

func (fcs *FamilyCarStrategy) BuildCar() *model.Car {
	return &model.Car{Brand: "Family-Toyota", Model: "Camry"}
}
