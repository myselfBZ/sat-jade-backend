package feedback

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/myselfBZ/sat-jade/internal/llm"
	"github.com/myselfBZ/sat-jade/internal/queries/feedbacks"
	"github.com/myselfBZ/sat-jade/internal/store"
)

type GenerateParams struct {
	Overview *store.ResultOverview
	ResultId int32
	UserId   uuid.UUID
}

func New(st store.FeedbackRepository, llm llm.LLM) Service {
	return &FeedbackService{
		store: st,
		llm: llm,
	}
}

type Service interface {
	Generate(ctx context.Context, params *GenerateParams) (*store.Feedback, error)
}

type FeedbackService struct {
	store store.FeedbackRepository
	llm   llm.LLM
}

func (f *FeedbackService) Generate(ctx context.Context, params *GenerateParams) (*store.Feedback, error) {
	fb, err := f.store.Get(ctx, params.ResultId)
	if err != nil {
		if !errors.Is(err, store.ErrRecordNotFound) {
			return nil, err
		}

		res, err := f.llm.GeneratePracticeOverview(&llm.PracticeOverviewParams{
			CorrectAnswers: params.Overview.CorrectAnswers,
			Mistakes:       params.Overview.MistakesByDomain,
		})

		if err != nil {
			return nil, err
		}

		fb, err := f.store.Create(ctx, feedbacks.CreateParams{
			ResultID: params.ResultId,
			UserID:   pgtype.UUID{Bytes: params.UserId, Valid: true},
			Header:   res.Overview,
			Body:     strings.Join(res.Suggesttions, "\n"),
			Footer:   res.Motivation,
		})

		if err != nil {
			return nil, err
		}

		return fb, nil
	}

	return fb, nil
}
