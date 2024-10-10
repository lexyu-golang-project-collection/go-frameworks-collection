package builder

import (
	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/model"
	"github.com/lexyu-golang-project-collection/go-design-patterns/combined/factory_with_strategy/internal/strategy"
)

type CarBuilder struct {
	strategy strategy.CarBuildStrategy
}

func (carBuilder *CarBuilder) ConstrcutCar() *model.Car {
	return carBuilder.strategy.BuildCar()
}
