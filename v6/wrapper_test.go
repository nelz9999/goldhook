package goldhook_test

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	"github.com/launchdarkly/go-sdk-common/v3/ldreason"
	"github.com/launchdarkly/go-sdk-common/v3/lduser"
	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"
	ld "github.com/launchdarkly/go-server-sdk/v6"

	"github.com/nelz9999/goldhook/v6"
)

func TestObservable(t *testing.T) {
	now := fmt.Sprintf("%d", time.Now().UnixNano())
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	ctxkey := fmt.Sprintf("ctx-%s", now)

	client, err := ld.MakeCustomClient("", ld.Config{Offline: true}, 0)
	if err != nil {
		t.Fatalf("unexpected: %v\n", err)
	}

	observed := map[string]bool{}
	hook := goldhook.ObserverFunc(
		func(
			ctx context.Context,
			key string,
			ldctx ldcontext.Context,
			callsiteDefault ldvalue.Value,
			_ time.Duration,
			detail ldreason.EvaluationDetail,
			_ error,
		) {
			t.Logf("%s: %v\n", key, detail.Value)
			// since the LDClient is offline, all results are
			// expected to fall back to the callsiteDefault
			if !callsiteDefault.Equal(detail.Value) {
				t.Errorf("result - expected %v; got %v\n", callsiteDefault, detail.Value)
			}
			// verify we got an expected user
			if !strings.Contains(ldctx.Key(), now) {
				t.Errorf("user - expected %q in %q\n", now, ldctx.Key())
			}
			if ctxval, _ := ctx.Value(ctxkey).(string); !strings.Contains(ctxval, now) {
				t.Errorf("context - expected %q in %q\n", now, ctxval)
			}
			observed[key] = true
		},
	)

	hooked, err := goldhook.NewEvaluator(
		context.Background(),
		// this implicitly tests that the Evaluator interface
		// correctly describes the LDClient
		client,
		// this tests the SubscriberFunc adapter
		hook,
	)
	if err != nil {
		t.Fatalf("unexpected: %v\n", err)
	}

	expected := []string{}
	keyFn := func(prefix string) (string, string) {
		ks := fmt.Sprintf("%s-simple-%s", prefix, now)
		kd := fmt.Sprintf("%s-detail-%s", prefix, now)
		expected = append(expected, ks, kd)
		return ks, kd
	}

	ctx := context.WithValue(context.Background(), ctxkey, fmt.Sprintf("ctxvalue-%s", now))
	user := lduser.NewAnonymousUser(now)

	// test Bool
	ks, kd := keyFn("bool")
	dBool := rnd.Intn(2) == 0
	hooked.WithContext(ctx).BoolVariation(ks, user, dBool)
	hooked.WithContext(ctx).BoolVariationDetail(kd, user, dBool)

	// test Float
	ks, kd = keyFn("float")
	dFloat := rnd.Float64()
	hooked.WithContext(ctx).Float64Variation(ks, user, dFloat)
	hooked.WithContext(ctx).Float64VariationDetail(kd, user, dFloat)

	// test Int
	ks, kd = keyFn("int")
	dInt := rnd.Int()
	hooked.WithContext(ctx).IntVariation(ks, user, dInt)
	hooked.WithContext(ctx).IntVariationDetail(kd, user, dInt)

	// test JSON
	ks, kd = keyFn("json")
	dJSON := ldvalue.Parse([]byte(fmt.Sprintf(`{"key":"%s"}`, now)))
	hooked.WithContext(ctx).JSONVariation(ks, user, dJSON)
	hooked.WithContext(ctx).JSONVariationDetail(kd, user, dJSON)

	// test String
	ks, kd = keyFn("string")
	dStr := fmt.Sprintf("result-%s", now)
	hooked.WithContext(ctx).StringVariation(ks, user, dStr)
	hooked.WithContext(ctx).StringVariationDetail(kd, user, dStr)

	// make sure all the expecteed keys were seen in the hook
	for _, exp := range expected {
		if !observed[exp] {
			t.Errorf("expected %q in %v\n", exp, observed)
		}
	}
}
