package metrics

import "errors"

var (
	ErrNilConfig = errors.New("nil config")
)

type configurer interface {
}

type metricService struct {
}

func (s *metricService) GetDifficulties() uint {
	return 10 //TODO:  Make senseble
}

func NewMetricService(
	cfg configurer,
) (*metricService, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	s := &metricService{}

	return s, nil
}
