package services

type Service struct {
	OrderService *orderService
}

type Depends struct {
	OrderStorage orderStorage
	OrderCache   orderCache
}

func New(depends Depends) *Service {
	return &Service{
		OrderService: newOrderService(depends.OrderStorage, depends.OrderCache),
	}
}
