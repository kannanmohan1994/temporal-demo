package helpers

import (
	"time"

	shared "github.com/kannanmohan1994/common-structs"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Driverworkflow struct {
	Client client.Client
}

func (w *Driverworkflow) AssignDriverWorkflow(ctx workflow.Context, req *shared.AssignDriverWorkflowRequest) (
	resp *shared.AssignDriverWorkflowResponse, err error) {

	logger := workflow.GetLogger(ctx)
	logger.Info("assigning driver...", "user-id", req.UserID)

	steps := shared.StepsStore{
		shared.Step{Name: "VERIFY_USER", Progress: "not_started", Error: nil},
		shared.Step{Name: "MATCH_DRIVER", Progress: "not_started", Error: nil},
		shared.Step{Name: "CONFIRM_DRIVE", Progress: "not_started", Error: nil}}

	return w.flow(ctx, req, steps)
}

func (w *Driverworkflow) flow(ctx workflow.Context, req *shared.AssignDriverWorkflowRequest, steps shared.StepsStore) (
	resp *shared.AssignDriverWorkflowResponse, err error) {

	defer func() {
		// sends notification on failure
		if err != nil {
			steps.FailPresentStage()
			ctx = w.WithActivityOptions(ctx, shared.TaskQueueNotifyUserActivity)
			workflow.ExecuteActivity(ctx, shared.NotifyUserActivity, &shared.NotifyUserRequest{
				UserID: req.UserID,
				Steps:  steps,
			}).Get(ctx, nil)
		}
	}()

	cancelCtx, _ := workflow.WithCancel(ctx)
	if err = w.MoveToNextStageAndNotify(cancelCtx, req.UserID, steps); err != nil {
		return nil, err
	}

	var userDetails *shared.FetchUserDetailsResponse
	var verificationDetails *shared.VerifyUserResponse
	{
		// user details fetch and verification activities can happen in parallel
		ctx1 := w.WithActivityOptions(cancelCtx, shared.TaskQueueFetchUserDetailsActivity)
		future := workflow.ExecuteActivity(ctx1, shared.FetchUserDetailsActivity, &shared.FetchUserDetailsRequest{
			UserID: req.UserID,
		})

		ctx2 := w.WithActivityOptions(cancelCtx, shared.TaskQueueVerifyUserActivity)
		future2 := workflow.ExecuteActivity(ctx2, shared.VerifyUserActivity, &shared.VerifyUserRequest{
			UserID: req.UserID,
		})

		// waits for executon of 2 parallel processes
		err = future.Get(ctx1, &userDetails)
		if err != nil {
			// wil come here only if retry attempt exhausted, timed-out or business error
			return
		}
		err = future2.Get(ctx2, &verificationDetails)
		if err != nil {
			return
		}
	}

	for {
		if err = w.MoveToNextStageAndNotify(cancelCtx, req.UserID, steps); err != nil {
			return nil, err
		}

		var matchDriver *shared.MatchDriverResponse
		ctx = w.WithActivityOptionsTimeout(cancelCtx, shared.TaskQueueMatchDriverActivity, 1*time.Second)
		err = workflow.ExecuteActivity(ctx, shared.MatchDriverActivity, &shared.MatchDriverRequest{
			UserID:       req.UserID,
			UserLocation: req.Location,
		}).Get(cancelCtx, &matchDriver)
		if err != nil {
			return nil, err
		}

		if err = w.MoveToNextStageAndNotify(cancelCtx, req.UserID, steps); err != nil {
			return nil, err
		}

		var waitDriver *shared.WaitDriverResponse
		ctx = w.WithActivityOptions(cancelCtx, shared.TaskQueueWaitDriverActivity)
		err = workflow.ExecuteActivity(ctx, shared.WaitDriverActivity, &shared.WaitDriverRequest{
			UserID:       req.UserID,
			UserLocation: req.Location,
		}).Get(cancelCtx, &waitDriver)
		if err != nil {
			return nil, err
		}

		if waitDriver.Accept {
			break
		}
		steps.MoveToStep(1)
	}

	if err = w.MoveToNextStageAndNotify(cancelCtx, req.UserID, steps); err != nil {
		return nil, err
	}

	return &shared.AssignDriverWorkflowResponse{}, nil
}

func (w *Driverworkflow) MoveToNextStageAndNotify(ctx workflow.Context, userID string, steps shared.StepsStore) error {
	steps.StartNextStage()
	ctx = w.WithActivityOptions(ctx, shared.TaskQueueNotifyUserActivity)
	return workflow.ExecuteActivity(ctx, shared.NotifyUserActivity, &shared.NotifyUserRequest{
		UserID: userID,
		Steps:  steps,
	}).Get(ctx, nil)
}

func (w *Driverworkflow) WithActivityOptions(ctx workflow.Context, TaskQueueName string) workflow.Context {
	ao := workflow.ActivityOptions{
		HeartbeatTimeout:       time.Second, // will help in identifying oom's
		TaskQueue:              TaskQueueName,
		ScheduleToCloseTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Minute * 5,
			NonRetryableErrorTypes: []string{"BusinessError"},
		},
	}
	return workflow.WithActivityOptions(ctx, ao)
}

func (w *Driverworkflow) WithActivityOptionsTimeout(ctx workflow.Context, TaskQueueName string, Timeout time.Duration) workflow.Context {
	ao := workflow.ActivityOptions{
		HeartbeatTimeout:       time.Second, // will help in identifying oom's
		TaskQueue:              TaskQueueName,
		ScheduleToCloseTimeout: Timeout,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Minute * 5,
			NonRetryableErrorTypes: []string{"BusinessError"},
		},
	}
	return workflow.WithActivityOptions(ctx, ao)
}
