package goldhook

import (
	"context"
	"time"

	"gopkg.in/launchdarkly/go-sdk-common.v2/ldreason"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

// Subscriber is interested in knowing the details about the results of
// a feature flag evaluation.
type Subscriber interface {
	// Observer is invoked with information concerning a feature flag evaluation
	//
	// If, at the callsite, the consumer did not take the pains to call the
	// evaluation using WithContext chained in beforehand, Observe will receive
	// a default context (i.e. provided by context.Background)
	Observe(
		ctx context.Context,
		key string,
		user lduser.User,
		callsiteDefault ldvalue.Value,
		elapsed time.Duration,
		detail ldreason.EvaluationDetail,
		evalErr error,
	)
}

// SubsriberFunc is a function adapter for the Subscriber interface
type SubscriberFunc func(
	ctx context.Context,
	key string,
	user lduser.User,
	callsiteDefault ldvalue.Value,
	elapsed time.Duration,
	detail ldreason.EvaluationDetail,
	evalErr error,
)

// Observe conforms to the Subscriber interface
func (fn SubscriberFunc) Observe(
	ctx context.Context,
	key string,
	user lduser.User,
	callsiteDefault ldvalue.Value,
	elapsed time.Duration,
	detail ldreason.EvaluationDetail,
	evalErr error,
) {
	fn(ctx, key, user, callsiteDefault, elapsed, detail, evalErr)
}

// Evaluator is the interface that describes the subset of all (most?) methods on an
// LDClient, which a consumer would invoke to retrieve the value of individual flags.
type Evaluator interface {
	BoolVariation(key string, user lduser.User, defaultVal bool) (bool, error)
	BoolVariationDetail(key string, user lduser.User, defaultVal bool) (bool, ldreason.EvaluationDetail, error)
	Float64Variation(key string, user lduser.User, defaultVal float64) (float64, error)
	Float64VariationDetail(key string, user lduser.User, defaultVal float64) (float64, ldreason.EvaluationDetail, error)
	IntVariation(key string, user lduser.User, defaultVal int) (int, error)
	IntVariationDetail(key string, user lduser.User, defaultVal int) (int, ldreason.EvaluationDetail, error)
	JSONVariation(key string, user lduser.User, defaultVal ldvalue.Value) (ldvalue.Value, error)
	JSONVariationDetail(key string, user lduser.User, defaultVal ldvalue.Value) (ldvalue.Value, ldreason.EvaluationDetail, error)
	StringVariation(key string, user lduser.User, defaultVal string) (string, error)
	StringVariationDetail(key string, user lduser.User, defaultVal string) (string, ldreason.EvaluationDetail, error)
}

type ContextualEvaluator interface {
	Evaluator
	WithContext(context.Context) Evaluator
}
