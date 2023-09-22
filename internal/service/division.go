package service

type DivisionService interface {
	// TODO
}

type divisionService struct {
	Service
}

func NewDivisionService(service Service) DivisionService {
	return &divisionService{Service: service}
}
