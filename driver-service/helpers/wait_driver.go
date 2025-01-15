package helpers

import (
	"context"
	"log"
	"math/rand"
	"time"

	shared "github.com/kannanmohan1994/common-structs"
)

func WaitDriverActivity(ctx context.Context, req *shared.WaitDriverRequest) (*shared.WaitDriverResponse, error) {
	heartbeat := StartHeartbit(ctx, time.Second)
	defer heartbeat.Stop()

	log.Println("waiting driver response...", req.UserID, req.UserLocation)
	time.Sleep(2 * time.Second)

	return &shared.WaitDriverResponse{Accept: rand.Intn(100)%2 == 0}, nil
}
