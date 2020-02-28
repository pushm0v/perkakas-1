package featureflag

import (
	"context"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type FeatureFlagTestSuite struct {
	suite.Suite
	fflag  FeatureFlag
	config *FeatureFlagConfig
}

func (suite *FeatureFlagTestSuite) SetupTest() {
	ctx := context.Background()
	suite.config = &FeatureFlagConfig{
		BaseURL:  "http://some-url/api/v1",
		Username: "",
		Password: "",
	}
	suite.fflag = NewFeatureFlag(suite.config, ctx)
}

func (suite *FeatureFlagTestSuite) TestEval() {

	defer gock.Off()

	mockResult := `{
	  "evalContext": {
		"enableDebug": true,
		"entityContext": "",
		"entityID": "randomly_generated_607811211",
		"flagID": 5,
		"flagKey": "some-flag"
	  },
	  "evalDebugLog": {
		"segmentDebugLogs": [
		  {
			"msg": "matched all constraints. rollout yes. {BucketNum:543 DistributionArray:{VariantIDs:[9] PercentsAccumulated:[1000]} VariantID:9 RolloutPercent:100}",
			"segmentID": 5
		  }
		]
	  },
	  "flagID": 5,
	  "flagKey": "some-flag",
	  "flagSnapshotID": 39,
	  "segmentID": 5,
	  "timestamp": "2020-02-27T04:25:32Z",
	  "variantAttachment": {},
	  "variantID": 9,
	  "variantKey": "some-variant"
	}`

	gock.New(suite.config.BaseURL).
		Post("evaluation").
		Reply(200).
		JSON(mockResult).
		SetHeader("Content-Type", "application/json")

	flagString := "some-flag"
	ok, err := suite.fflag.Eval(flagString, "some-variant")

	assert.Equal(suite.T(), true, gock.IsDone(), "Must be equal")
	assert.True(suite.T(), ok, "Result OK should be true")
	assert.Nil(suite.T(), err, "Error should be nil")
}

func (suite *FeatureFlagTestSuite) TestFailEval() {
	flagString := "some-flag"
	ok, err := suite.fflag.Eval(flagString, "some-variant")

	assert.False(suite.T(), ok, "Result OK should be false")
	assert.NotNil(suite.T(), err, "Error should be not nil")
}

func (suite *FeatureFlagTestSuite) TestEvalThenExecute() {

	defer gock.Off()

	mockResult := `{
	  "evalContext": {
		"enableDebug": true,
		"entityContext": "",
		"entityID": "randomly_generated_607811211",
		"flagID": 5,
		"flagKey": "some-flag"
	  },
	  "evalDebugLog": {
		"segmentDebugLogs": [
		  {
			"msg": "matched all constraints. rollout yes. {BucketNum:543 DistributionArray:{VariantIDs:[9] PercentsAccumulated:[1000]} VariantID:9 RolloutPercent:100}",
			"segmentID": 5
		  }
		]
	  },
	  "flagID": 5,
	  "flagKey": "some-flag",
	  "flagSnapshotID": 39,
	  "segmentID": 5,
	  "timestamp": "2020-02-27T04:25:32Z",
	  "variantAttachment": {},
	  "variantID": 9,
	  "variantKey": "some-variant"
	}`

	gock.New(suite.config.BaseURL).
		Post("evaluation").
		Reply(200).
		JSON(mockResult).
		SetHeader("Content-Type", "application/json")

	var funcCalled = false

	ok, err := suite.fflag.EvalThenExecute("some-flag", "some-variant", func() error {
		funcCalled = true
		return nil
	})

	assert.Equal(suite.T(), true, gock.IsDone(), "Must be equal")
	assert.True(suite.T(), funcCalled, "Feature func should be called")
	assert.True(suite.T(), ok, "Result OK should be true")
	assert.Nil(suite.T(), err, "Error should be nil")
}

func (suite *FeatureFlagTestSuite) TestFailEvalThenExecute() {
	flagString := "some-flag"
	ok, err := suite.fflag.EvalThenExecute(flagString, "some-variant", func() error {
		return nil
	})

	assert.False(suite.T(), ok, "Result OK should be false")
	assert.NotNil(suite.T(), err, "Error should be not nil")
}

func (suite *FeatureFlagTestSuite) TestNonExistVariant() {
	defer gock.Off()

	mockResult := `{
	  "evalContext": {
		"enableDebug": true,
		"entityContext": "",
		"entityID": "randomly_generated_607811211",
		"flagID": 5,
		"flagKey": "some-flag"
	  },
	  "evalDebugLog": {
		"segmentDebugLogs": [
		  {
			"msg": "matched all constraints. rollout yes. {BucketNum:543 DistributionArray:{VariantIDs:[9] PercentsAccumulated:[1000]} VariantID:9 RolloutPercent:100}",
			"segmentID": 5
		  }
		]
	  },
	  "flagID": 5,
	  "flagKey": "some-flag",
	  "flagSnapshotID": 39,
	  "segmentID": 5,
	  "timestamp": "2020-02-27T04:25:32Z",
	  "variantAttachment": {},
	  "variantID": 9,
	  "variantKey": "some-variant"
	}`

	gock.New(suite.config.BaseURL).
		Post("evaluation").
		Reply(200).
		JSON(mockResult).
		SetHeader("Content-Type", "application/json")

	flagString := "some-flag"
	ok, err := suite.fflag.EvalThenExecute(flagString, "some-variant-not-exist", func() error {
		return nil
	})

	assert.Equal(suite.T(), true, gock.IsDone(), "Must be equal")
	assert.False(suite.T(), ok, "Result OK should be false")
	assert.Nil(suite.T(), err, "Error should be nil")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(FeatureFlagTestSuite))
}
