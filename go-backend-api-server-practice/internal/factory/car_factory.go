package factory

import (
	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/strategy"
	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/types"
)

func GetCarBuildStrategy(carType types.CarType) strategy.CarBuildStrategy {
	switch carType {
	case types.SPORT_CAR:
		return &strategy.SportCarStrategy{}
	case types.FAMILY_CAR:
		return &strategy.FamilyCarStrategy{}
	default:
		return nil
	}
}
