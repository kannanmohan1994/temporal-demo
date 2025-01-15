package helpers

import (
	"context"
	"log"
	"time"

	shared "github.com/kannanmohan1994/common-structs"
)

func MatchDriverActivity(ctx context.Context, req *shared.MatchDriverRequest) (*shared.MatchDriverResponse, error) {
	heartbeat := StartHeartbit(ctx, time.Second)
	defer heartbeat.Stop()

	log.Println("matching driver...", req.UserID, req.UserLocation)
	//time.Sleep(2 * time.Second)

	return &shared.MatchDriverResponse{}, nil
}
