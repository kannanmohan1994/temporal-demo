package helpers

import (
	"context"
	"log"
	"time"

	shared "github.com/kannanmohan1994/common-structs"
)

func FetchUserDetailsActivity(ctx context.Context, req *shared.FetchUserDetailsRequest) (*shared.FetchUserDetailsResponse, error) {
	heartbeat := StartHeartbit(ctx, time.Second)
	defer heartbeat.Stop()

	log.Println("fetching user details...", req.UserID)
	time.Sleep(time.Second)
	return &shared.FetchUserDetailsResponse{}, nil
}
