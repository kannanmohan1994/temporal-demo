package helpers

import (
	"context"
	"log"
	"time"

	shared "github.com/kannanmohan1994/common-structs"
)

func NotifyUserActivity(ctx context.Context, req *shared.NotifyUserRequest) (*shared.NotifyUserResponse, error) {
	heartbeat := StartHeartbit(ctx, time.Second)
	defer heartbeat.Stop()

	log.Println("notifying user details...", req.UserID)
	return &shared.NotifyUserResponse{}, nil
}
