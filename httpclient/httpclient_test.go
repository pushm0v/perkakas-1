package httpclient

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gojektech/heimdall/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"gopkg.in/h2non/gock.v1"
)

type LogTestSuite struct {
	suite.Suite
	HttpClient   *httpclient.Client
	host         string
	endpoint     string
	url          string
	mockResponse map[string]interface{}
}

func (suite *LogTestSuite) SetupTest() {
	suite.HttpClient = NewHttpClient(nil) // Default configuration
	suite.host = "http://some-test-url.com"
	suite.endpoint = "test"
	suite.url = fmt.Sprintf("%s/%s", suite.host, suite.endpoint)
	suite.mockResponse = map[string]interface{}{
		"test": "a1234-abcd",
	}
}

func (suite *LogTestSuite) TestGetRequest() {
	defer gock.Off()

	gock.New(suite.host).
		Get(suite.endpoint).
		Reply(200).
		JSON(suite.mockResponse)

	resp, err := suite.HttpClient.Get(suite.url, nil)
	assert.Nil(suite.T(), err, "Nil expected")

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	assert.Equal(suite.T(), true, gock.IsDone(), "must be equal")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(LogTestSuite))
}
