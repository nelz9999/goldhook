package goldhook

import (
	"context"
	"fmt"
	"time"

	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	"github.com/launchdarkly/go-sdk-common/v3/ldreason"
	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"
)

type ObservedEvaluator struct {
	client EvaluatorCtx
	hooks  []Observer
	ctx    context.Context
}

func NewEvaluator(ctx context.Context, client EvaluatorCtx, subscribers ...Observer) (*ObservedEvaluator, error) {
	if client == nil {
		return nil, fmt.Errorf("client must not be nil")
	}
	// we don't check that there are 1+ subscribers vis Postel's Law
	// but we do check for nil subscribers which could cause later panics
	for _, h := range subscribers {
		if h == nil {
			return nil, fmt.Errorf("subscribers must not be nil")
		}
	}
	return &ObservedEvaluator{
		client: client,
		hooks:  subscribers,
		ctx:    ctx,
	}, nil
}

func (oe *ObservedEvaluator) WithContext(c context.Context) Evaluator {
	return &ObservedEvaluator{
		client: oe.client,
		hooks:  oe.hooks,
		ctx:    c,
	}
}

func (oe *ObservedEvaluator) notifyHooks(
	ctx context.Context,
	key string,
	ldctx ldcontext.Context,
	callsiteDefault ldvalue.Value,
	elapsed time.Duration,
	detail ldreason.EvaluationDetail,
	evalErr error,
) {
	for _, h := range oe.hooks {
		h.Observe(ctx, key, ldctx, callsiteDefault, elapsed, detail, evalErr)
	}
}

/* * * BOOL * * */

func (oe *ObservedEvaluator) BoolVariation(key string, ldctx ldcontext.Context, defaultVal bool) (bool, error) {
	_, detail, err := oe.BoolVariationDetailCtx(oe.ctx, key, ldctx, defaultVal)
	return detail.Value.BoolValue(), err
}

func (oe *ObservedEvaluator) BoolVariationDetail(key string, ldctx ldcontext.Context, defaultVal bool) (bool, ldreason.EvaluationDetail, error) {
	_, detail, err := oe.BoolVariationDetailCtx(oe.ctx, key, ldctx, defaultVal)
	return detail.Value.BoolValue(), detail, err
}

func (oe *ObservedEvaluator) BoolVariationCtx(ctx context.Context, key string, ldctx ldcontext.Context, defaultVal bool) (bool, error) {
	_, detail, err := oe.BoolVariationDetailCtx(ctx, key, ldctx, defaultVal)
	return detail.Value.BoolValue(), err
}

func (oe *ObservedEvaluator) BoolVariationDetailCtx(ctx context.Context, key string, ldctx ldcontext.Context, defaultVal bool) (bool, ldreason.EvaluationDetail, error) {
	start := time.Now()
	_, detail, err := oe.client.BoolVariationDetailCtx(ctx, key, ldctx, defaultVal)
	oe.notifyHooks(ctx, key, ldctx, ldvalue.Bool(defaultVal), time.Since(start), detail, err)
	return detail.Value.BoolValue(), detail, err
}

/* * * FLOAT * * */

func (oe *ObservedEvaluator) Float64Variation(key string, ldctx ldcontext.Context, defaultVal float64) (float64, error) {
	_, detail, err := oe.Float64VariationDetailCtx(oe.ctx, key, ldctx, defaultVal)
	return detail.Value.Float64Value(), err
}

func (oe *ObservedEvaluator) Float64VariationDetail(key string, ldctx ldcontext.Context, defaultVal float64) (float64, ldreason.EvaluationDetail, error) {
	_, detail, err := oe.Float64VariationDetailCtx(oe.ctx, key, ldctx, defaultVal)
	return detail.Value.Float64Value(), detail, err
}

func (oe *ObservedEvaluator) Float64VariationCtx(ctx context.Context, key string, ldctx ldcontext.Context, defaultVal float64) (float64, error) {
	_, detail, err := oe.Float64VariationDetailCtx(ctx, key, ldctx, defaultVal)
	return detail.Value.Float64Value(), err
}

func (oe *ObservedEvaluator) Float64VariationDetailCtx(ctx context.Context, key string, ldctx ldcontext.Context, defaultVal float64) (float64, ldreason.EvaluationDetail, error) {
	start := time.Now()
	_, detail, err := oe.client.Float64VariationDetailCtx(ctx, key, ldctx, defaultVal)
	oe.notifyHooks(ctx, key, ldctx, ldvalue.Float64(defaultVal), time.Since(start), detail, err)
	return detail.Value.Float64Value(), detail, err
}

/* * * INT * * */

func (oe *ObservedEvaluator) IntVariation(key string, ldctx ldcontext.Context, defaultVal int) (int, error) {
	_, detail, err := oe.IntVariationDetailCtx(oe.ctx, key, ldctx, defaultVal)
	return detail.Value.IntValue(), err
}

func (oe *ObservedEvaluator) IntVariationDetail(key string, ldctx ldcontext.Context, defaultVal int) (int, ldreason.EvaluationDetail, error) {
	_, detail, err := oe.IntVariationDetailCtx(oe.ctx, key, ldctx, defaultVal)
	return detail.Value.IntValue(), detail, err
}

func (oe *ObservedEvaluator) IntVariationCtx(ctx context.Context, key string, ldctx ldcontext.Context, defaultVal int) (int, error) {
	_, detail, err := oe.IntVariationDetailCtx(ctx, key, ldctx, defaultVal)
	return detail.Value.IntValue(), err
}

func (oe *ObservedEvaluator) IntVariationDetailCtx(ctx context.Context, key string, ldctx ldcontext.Context, defaultVal int) (int, ldreason.EvaluationDetail, error) {
	start := time.Now()
	_, detail, err := oe.client.IntVariationDetailCtx(ctx, key, ldctx, defaultVal)
	oe.notifyHooks(ctx, key, ldctx, ldvalue.Int(defaultVal), time.Since(start), detail, err)
	return detail.Value.IntValue(), detail, err
}

/* * * JSON * * */

func (oe *ObservedEvaluator) JSONVariation(key string, ldctx ldcontext.Context, defaultVal ldvalue.Value) (ldvalue.Value, error) {
	_, detail, err := oe.JSONVariationDetailCtx(oe.ctx, key, ldctx, defaultVal)
	return detail.Value, err
}

func (oe *ObservedEvaluator) JSONVariationDetail(key string, ldctx ldcontext.Context, defaultVal ldvalue.Value) (ldvalue.Value, ldreason.EvaluationDetail, error) {
	_, detail, err := oe.JSONVariationDetailCtx(oe.ctx, key, ldctx, defaultVal)
	return detail.Value, detail, err
}

func (oe *ObservedEvaluator) JSONVariationCtx(ctx context.Context, key string, ldctx ldcontext.Context, defaultVal ldvalue.Value) (ldvalue.Value, error) {
	_, detail, err := oe.JSONVariationDetailCtx(ctx, key, ldctx, defaultVal)
	return detail.Value, err
}

func (oe *ObservedEvaluator) JSONVariationDetailCtx(ctx context.Context, key string, ldctx ldcontext.Context, defaultVal ldvalue.Value) (ldvalue.Value, ldreason.EvaluationDetail, error) {
	start := time.Now()
	_, detail, err := oe.client.JSONVariationDetailCtx(ctx, key, ldctx, defaultVal)
	oe.notifyHooks(ctx, key, ldctx, defaultVal, time.Since(start), detail, err)
	return detail.Value, detail, err
}

/* * * STRING * * */

func (oe *ObservedEvaluator) StringVariation(key string, ldctx ldcontext.Context, defaultVal string) (string, error) {
	_, detail, err := oe.StringVariationDetailCtx(oe.ctx, key, ldctx, defaultVal)
	return detail.Value.StringValue(), err
}

func (oe *ObservedEvaluator) StringVariationDetail(key string, ldctx ldcontext.Context, defaultVal string) (string, ldreason.EvaluationDetail, error) {
	_, detail, err := oe.StringVariationDetailCtx(oe.ctx, key, ldctx, defaultVal)
	return detail.Value.StringValue(), detail, err
}

func (oe *ObservedEvaluator) StringVariationCtx(ctx context.Context, key string, ldctx ldcontext.Context, defaultVal string) (string, error) {
	_, detail, err := oe.StringVariationDetailCtx(ctx, key, ldctx, defaultVal)
	return detail.Value.StringValue(), err
}

func (oe *ObservedEvaluator) StringVariationDetailCtx(ctx context.Context, key string, ldctx ldcontext.Context, defaultVal string) (string, ldreason.EvaluationDetail, error) {
	start := time.Now()
	_, detail, err := oe.client.StringVariationDetailCtx(ctx, key, ldctx, defaultVal)
	oe.notifyHooks(ctx, key, ldctx, ldvalue.String(defaultVal), time.Since(start), detail, err)
	return detail.Value.StringValue(), detail, err
}
