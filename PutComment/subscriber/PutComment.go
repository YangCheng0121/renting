package subscriber

import (
	"context"
	"github.com/asim/go-micro/v3/util/log"
	pb "renting/PutComment/proto"
)

type PutComment struct{}

func (e *PutComment) Handle(ctx context.Context, msg *pb.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *pb.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
