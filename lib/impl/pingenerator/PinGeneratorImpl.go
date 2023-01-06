package pingenerator

import (
	"github.com/MarcGrol/patientregistration/lib/api/pingenerator"
	"math/rand"
)

type pinGeneratorImpl struct {
}

func New() pingenerator.PinGenerator {
	return &pinGeneratorImpl{}
}

func (p pinGeneratorImpl) GeneratePin() (int, error) {
	min := 1000
	max := 9999
	pin := rand.Intn(max-min) + min
	return pin, nil
}
