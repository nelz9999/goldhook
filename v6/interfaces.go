package goldhook

import (
	"context"
	"time"

	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	"github.com/launchdarkly/go-sdk-common/v3/ldreason"
	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"
)

// Observer is interested in knowing the details about the results of
// a feature flag evaluation.
type Observer interface {
	// Observer is invoked with information concerning a feature flag evaluation
	//
	// If, at the callsite, the consumer did not take the pains to call the
	// evaluation using WithContext chained in beforehand, Observe will receive
	// a default context (i.e. provided by context.Background)
	Observe(
		ctx context.Context,
		key string,
		ldctx ldcontext.Context,
		callsiteDefault ldvalue.Value,
		elapsed time.Duration,
		detail ldreason.EvaluationDetail,
		evalErr error,
	)
}

// SubsriberFunc is a function adapter for the Observer interface
type ObserverFunc func(
	ctx context.Context,
	key string,
	ldctx ldcontext.Context,
	callsiteDefault ldvalue.Value,
	elapsed time.Duration,
	detail ldreason.EvaluationDetail,
	evalErr error,
)

// Observe conforms to the Observer interface
func (fn ObserverFunc) Observe(
	ctx context.Context,
	key string,
	ldctx ldcontext.Context,
	callsiteDefault ldvalue.Value,
	elapsed time.Duration,
	detail ldreason.EvaluationDetail,
	evalErr error,
) {
	fn(ctx, key, ldctx, callsiteDefault, elapsed, detail, evalErr)
}

// Evaluator is the interface that describes the subset of all (most?) methods on an
// LDClient, which a consumer would invoke to retrieve the value of individual flags.
type Evaluator interface {
	BoolVariation(key string, ldctx ldcontext.Context, defaultVal bool) (bool, error)
	BoolVariationDetail(key string, ldctx ldcontext.Context, defaultVal bool) (bool, ldreason.EvaluationDetail, error)
	Float64Variation(key string, ldctx ldcontext.Context, defaultVal float64) (float64, error)
	Float64VariationDetail(key string, ldctx ldcontext.Context, defaultVal float64) (float64, ldreason.EvaluationDetail, error)
	IntVariation(key string, ldctx ldcontext.Context, defaultVal int) (int, error)
	IntVariationDetail(key string, ldctx ldcontext.Context, defaultVal int) (int, ldreason.EvaluationDetail, error)
	JSONVariation(key string, ldctx ldcontext.Context, defaultVal ldvalue.Value) (ldvalue.Value, error)
	JSONVariationDetail(key string, ldctx ldcontext.Context, defaultVal ldvalue.Value) (ldvalue.Value, ldreason.EvaluationDetail, error)
	StringVariation(key string, ldctx ldcontext.Context, defaultVal string) (string, error)
	StringVariationDetail(key string, ldctx ldcontext.Context, defaultVal string) (string, ldreason.EvaluationDetail, error)
}

type ContextualEvaluator interface {
	Evaluator
	WithContext(context.Context) Evaluator
}
