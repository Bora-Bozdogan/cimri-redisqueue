package handlers

import (
	"cimrique-redis/internal/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type servicesInterface interface {
	UnpackRequest(body []byte) (models.Request, int)
	EnqueueHigh(req models.Request) error
	EnqueueMed(req models.Request) error
	EnqueueLow(req models.Request) error
	IncrementRequestCount()
	IncrementValidRequestCount()
}

type Handler struct {
	service servicesInterface
}

func NewHandler(service servicesInterface) Handler {
	return Handler{service: service}
}

func (h Handler) HandleEnqueue(c *fiber.Ctx) error {
	//metric
	h.service.IncrementRequestCount()
	
	//unpack model and score
	req, score := h.service.UnpackRequest(c.Body())

	//check for invalid score (100> or 1000<)
	if score < 100 || score > 1000 {
		return c.Status(400).SendString("Invalid score")
	}

	//based on score, add request to the correct que using rds.LPUSH("que_name", instance)
	//high, med, low ques
	var err error
	var que string
	if score >= 800 {
		err = h.service.EnqueueHigh(req)
		que = "high"
	} else if score >= 500 {
		err = h.service.EnqueueMed(req)
		que = "med"
	} else {
		err = h.service.EnqueueLow(req)
		que = "low"
	}
	
	if err != nil {
		return c.Status(400).SendString("Couldn't enqueue request")
	}

	//metric
	h.service.IncrementValidRequestCount()

	//alert the workers by sending a message
	return c.Status(200).SendString(fmt.Sprintf("Enqueued request score %d on queue %s", score, que))

	//also write scoring search the que logic
}
