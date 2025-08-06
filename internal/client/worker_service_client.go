package client

import (
	"cimrique-redis/internal/models"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type WorkerServiceClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewWorkerServiceClient(addr string, pass string, num int, protocol int) WorkerServiceClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,     // No password set
		DB:       num,      // Use default DB
		Protocol: protocol, // Connection protocol
	})
	return WorkerServiceClient{client: client, ctx: context.Background()}
}

func (s WorkerServiceClient) UnpackRequest(body []byte) (models.Request, int) {
	//que logic here
	var msg models.QueMessage
	json.Unmarshal(body, &msg)
	req := msg.Message
	score := msg.Score
	return req, score
}

func (s WorkerServiceClient) EnqueueHigh(req models.Request) error {
	obj, err := json.Marshal(req)
	if err != nil {
		return err
	}
	cmd := s.client.LPush(s.ctx, "high", obj)
	_, err = cmd.Result()
	if err != nil {
		return err
	}
	return nil
}

func (s WorkerServiceClient) EnqueueMed(req models.Request) error {
	obj, err := json.Marshal(req)
	if err != nil {
		return err
	}
	cmd := s.client.LPush(s.ctx, "med", obj)
	_, err = cmd.Result()
	if err != nil {
		return err
	}
	return nil
}

func (s WorkerServiceClient) EnqueueLow(req models.Request) error {
	obj, err := json.Marshal(req)
	if err != nil {
		return err
	}
	cmd := s.client.LPush(s.ctx, "low", obj)
	_, err = cmd.Result()
	if err != nil {
		return err
	}
	return nil
}
