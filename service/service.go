package service

type Service struct {
}

// New init
func New() (s *Service) {
	s = &Service{
	}
	return s
}