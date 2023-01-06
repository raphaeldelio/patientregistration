package pingenerator

//go:generate mockgen -source=pinGenerator.go -destination=pinGeneratorMocks.go -package=pingenerator PinGenerator

type PinGenerator interface {
	GeneratePin() (int, error)
}
