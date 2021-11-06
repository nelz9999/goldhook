package goldhook

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/launchdarkly/go-sdk-common.v2/ldreason"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

type ObservedEvaluator struct {
	client Evaluator
	hooks  []Observer
	ctx    context.Context
}

func NewEvaluator(client Evaluator, subscribers ...Observer) (*ObservedEvaluator, error) {
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
		ctx:    context.Background(),
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
	key string,
	user lduser.User,
	callsiteDefault ldvalue.Value,
	elapsed time.Duration,
	detail ldreason.EvaluationDetail,
	evalErr error,
) {
	for _, h := range oe.hooks {
		h.Observe(oe.ctx, key, user, callsiteDefault, elapsed, detail, evalErr)
	}
}

func (oe *ObservedEvaluator) BoolVariation(key string, user lduser.User, defaultVal bool) (bool, error) {
	_, detail, err := oe.BoolVariationDetail(key, user, defaultVal)
	return detail.Value.BoolValue(), err
}

func (oe *ObservedEvaluator) BoolVariationDetail(key string, user lduser.User, defaultVal bool) (bool, ldreason.EvaluationDetail, error) {
	start := time.Now()
	_, detail, err := oe.client.BoolVariationDetail(key, user, defaultVal)
	oe.notifyHooks(key, user, ldvalue.Bool(defaultVal), time.Since(start), detail, err)
	return detail.Value.BoolValue(), detail, err
}

func (oe *ObservedEvaluator) Float64Variation(key string, user lduser.User, defaultVal float64) (float64, error) {
	_, detail, err := oe.Float64VariationDetail(key, user, defaultVal)
	return detail.Value.Float64Value(), err
}

func (oe *ObservedEvaluator) Float64VariationDetail(key string, user lduser.User, defaultVal float64) (float64, ldreason.EvaluationDetail, error) {
	start := time.Now()
	_, detail, err := oe.client.Float64VariationDetail(key, user, defaultVal)
	oe.notifyHooks(key, user, ldvalue.Float64(defaultVal), time.Since(start), detail, err)
	return detail.Value.Float64Value(), detail, err
}

func (oe *ObservedEvaluator) IntVariation(key string, user lduser.User, defaultVal int) (int, error) {
	_, detail, err := oe.IntVariationDetail(key, user, defaultVal)
	return detail.Value.IntValue(), err
}

func (oe *ObservedEvaluator) IntVariationDetail(key string, user lduser.User, defaultVal int) (int, ldreason.EvaluationDetail, error) {
	start := time.Now()
	_, detail, err := oe.client.IntVariationDetail(key, user, defaultVal)
	oe.notifyHooks(key, user, ldvalue.Int(defaultVal), time.Since(start), detail, err)
	return detail.Value.IntValue(), detail, err
}

func (oe *ObservedEvaluator) JSONVariation(key string, user lduser.User, defaultVal ldvalue.Value) (ldvalue.Value, error) {
	_, detail, err := oe.JSONVariationDetail(key, user, defaultVal)
	return detail.Value, err
}

func (oe *ObservedEvaluator) JSONVariationDetail(key string, user lduser.User, defaultVal ldvalue.Value) (ldvalue.Value, ldreason.EvaluationDetail, error) {
	start := time.Now()
	_, detail, err := oe.client.JSONVariationDetail(key, user, defaultVal)
	oe.notifyHooks(key, user, defaultVal, time.Since(start), detail, err)
	return detail.Value, detail, err
}

func (oe *ObservedEvaluator) StringVariation(key string, user lduser.User, defaultVal string) (string, error) {
	_, detail, err := oe.StringVariationDetail(key, user, defaultVal)
	return detail.Value.StringValue(), err
}

func (oe *ObservedEvaluator) StringVariationDetail(key string, user lduser.User, defaultVal string) (string, ldreason.EvaluationDetail, error) {
	start := time.Now()
	_, detail, err := oe.client.StringVariationDetail(key, user, defaultVal)
	oe.notifyHooks(key, user, ldvalue.String(defaultVal), time.Since(start), detail, err)
	return detail.Value.StringValue(), detail, err
}
