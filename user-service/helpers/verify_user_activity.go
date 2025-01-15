package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	shared "github.com/kannanmohan1994/common-structs"
)

func VerifyUserActivity(ctx context.Context, req *shared.VerifyUserRequest) (*shared.VerifyUserResponse, error) {
	heartbeat := StartHeartbit(ctx, time.Second)
	defer heartbeat.Stop()

	log.Println("verifying user...", req.UserID)
	time.Sleep(2 * time.Second)

	tm := time.Now().Unix()
	var err error
	if tm%2 == 0 {
		//err = &shared.BusinessError{Err: "verification failed"}
		//err = errors.New("some issue that would get solved on retrying")
	}
	//err = &shared.BusinessError{Err: "verification failed"}
	// err = errors.New("some issue that would get solved on retrying")

	fmt.Println("error: ", err)

	return &shared.VerifyUserResponse{}, err
}
