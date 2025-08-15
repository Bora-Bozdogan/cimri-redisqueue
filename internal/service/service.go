package service

import "cimrique-redis/internal/models"

type clientServiceInterface interface {
	UnpackRequest(body []byte) (models.Request, int)
	EnqueueHigh(req models.Request) error
	EnqueueMed(req models.Request) error
	EnqueueLow(req models.Request) error
}

type metricsInterface interface {
	IncrementRequestCount() 
	IncrementValidRequestCount()
}

type ServicesFuncs struct {
	client clientServiceInterface
	metric metricsInterface
}

func NewServicesFuncs(client clientServiceInterface, metrics metricsInterface) ServicesFuncs {
	return ServicesFuncs{client: client, metric: metrics}
}

func (s ServicesFuncs) UnpackRequest(body []byte) (models.Request, int) {
	return s.client.UnpackRequest(body)
}

func (s ServicesFuncs) EnqueueHigh(req models.Request) error {
	return s.client.EnqueueHigh(req)
}

func (s ServicesFuncs) EnqueueMed(req models.Request) error {
	return s.client.EnqueueMed(req)
}

func (s ServicesFuncs) EnqueueLow(req models.Request) error {
	return s.client.EnqueueLow(req)
}

func (s ServicesFuncs) IncrementRequestCount() {
	s.metric.IncrementRequestCount()
}

func (s ServicesFuncs) IncrementValidRequestCount() {
	s.metric.IncrementValidRequestCount()
}
