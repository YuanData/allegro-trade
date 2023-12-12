package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/util"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Membername string `json:"membername"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	member, err := processor.store.GetMember(ctx, payload.Membername)
	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Membername:   member.Membername,
		Email:      member.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Subject"
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email??id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf("%s %s", member.NameEntire, verifyUrl)
	to := []string{member.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", member.Email).Msg("processed task")
	return nil
}
