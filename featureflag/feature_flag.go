package featureflag

import (
	"context"

	"github.com/checkr/goflagr"
)

type FeatureFlag interface {
	EvalThenExecute(flag, variant string, featureFunc FeatureFunc) (ok bool, err error)
	Eval(flag, variant string) (ok bool, err error)
}

type featureFlag struct {
	apiClient *goflagr.APIClient
	ctx       context.Context
}

type FeatureFlagConfig struct {
	Username string
	Password string
	BaseURL  string
	Debug    bool
}

type FeatureFunc func() error

func NewFlagrClient(baseURL, username, password string, ctx context.Context) (*goflagr.APIClient, context.Context) {
	basicAuth := goflagr.BasicAuth{
		UserName: username,
		Password: password,
	}

	ctx = context.WithValue(ctx, goflagr.ContextBasicAuth, basicAuth)

	conf := goflagr.NewConfiguration()
	conf.BasePath = baseURL

	return goflagr.NewAPIClient(conf), ctx
}

func NewFeatureFlag(conf *FeatureFlagConfig, ctx context.Context) FeatureFlag {
	apiClient, ctx := NewFlagrClient(conf.BaseURL, conf.Username, conf.Password, ctx)
	return &featureFlag{
		apiClient: apiClient,
		ctx:       ctx,
	}
}

func (ff *featureFlag) evaluationApiService(apiClient *goflagr.APIClient) *goflagr.EvaluationApiService {
	return apiClient.EvaluationApi
}

func (ff *featureFlag) EvalThenExecute(flag, variant string, featureFunc FeatureFunc) (ok bool, err error) {
	ok, err = ff.Eval(flag, variant)

	if err != nil {
		return false, err
	}

	if ok {
		err = featureFunc()
	}

	return
}

func (ff *featureFlag) Eval(flag, variant string) (ok bool, err error) {

	evalApi := ff.evaluationApiService(ff.apiClient)
	evalContext := goflagr.EvalContext{
		FlagKey: flag,
	}

	evalResult, _, err := evalApi.PostEvaluation(ff.ctx, evalContext)
	if err != nil {
		return false, err
	}

	if evalResult.VariantKey == variant {
		return true, nil
	}

	return false, nil
}
