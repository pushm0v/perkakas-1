package log

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	gock "gopkg.in/h2non/gock.v1"
)

type LogTestSuite struct {
	suite.Suite
	Logger       *Logger
	host         string
	endpoint     string
	url          string
	mockResponse map[string]interface{}
	responseTime string
	logID        string
}

func (suite *LogTestSuite) SetupTest() {
	suite.Logger = NewLogger()
	suite.logID = "3891967c-8589-42e0-a493-4fb6a0287992" // should be dynamic uuid. This static uuid was for testing purpose
	suite.host = "https://programming-quotes-api.herokuapp.com"
	suite.endpoint = "quotes"
	suite.url = fmt.Sprintf("%s/%s", suite.host, suite.endpoint)
	suite.responseTime = "2019-09-30T10:25:14+07:00"
	suite.mockResponse = map[string]interface{}{
		"id":     "a1234-abcd",
		"quotes": "Everything the light touches, is our kingdom",
	}
}

func (suite *LogTestSuite) TestLog() {
	defer gock.Off()

	gock.New(suite.host).
		Get(suite.endpoint).
		Reply(200).
		AddHeader("X-Test-Response", "Hello").
		AddHeader("Time", suite.responseTime).
		JSON(suite.mockResponse)

	resp, err := http.Get(suite.url)
	assert.Nil(suite.T(), err, "Nil expected")

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	fields := Field{
		LogID:          suite.logID,
		Endpoint:       suite.url,
		Method:         "GET",
		RequestBody:    nil,
		RequestHeader:  nil,
		ResponseBody:   string(b),
		ResponseHeader: resp.Header,
		ErrorMessage:   "Error",
	}
	suite.Logger.Set(fields)
	suite.Logger.Log(ErrorLevel, "testlog")
	assert.Equal(suite.T(), true, gock.IsDone(), "must be equal")
}

func (suite *LogTestSuite) TestLogBuilder() {
	defer gock.Off()

	gock.New(suite.host).
		Get(suite.endpoint).
		Reply(200).
		AddHeader("X-Test-Response", "Hello").
		AddHeader("Time", suite.responseTime).
		JSON(suite.mockResponse)

	resp, err := http.Get(suite.url)
	assert.Nil(suite.T(), err, "Nil expected")

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	suite.Logger.SetLevel(InfoLevel)
	suite.Logger.
		SetLogID(suite.logID).
		SetEndpoint(suite.url).
		SetMethod("GET").
		SetRequestBody(nil).
		SetRequestHeaders(nil).
		SetResponseBody(string(b)).
		SetResponseHeaders(resp.Header).
		SetErrorMessage(errors.New("Error in code 2123123"))
	suite.Logger.Log(InfoLevel, "testlog")
	assert.Equal(suite.T(), true, gock.IsDone(), "must be equal")
}

func (suite *LogTestSuite) TestSetLevel() {
	suite.Logger.SetLevel(InfoLevel)
	assert.Equal(suite.T(), log.Level(uint32(InfoLevel)), suite.Logger.logger.GetLevel())
}

func (suite *LogTestSuite) TestSetLogID() {
	suite.Logger.SetLogID(suite.logID)
	assert.NotEqual(suite.T(), "", suite.Logger.field.LogID)
}

func (suite *LogTestSuite) TestSetEndpoint() {
	suite.Logger.SetEndpoint(suite.endpoint)
	assert.Equal(suite.T(), suite.endpoint, suite.Logger.field.Endpoint)
}

func (suite *LogTestSuite) TestSetMethod() {
	suite.Logger.SetMethod("GET")
	assert.Equal(suite.T(), "GET", suite.Logger.field.Method)
}

func (suite *LogTestSuite) TestSetRequestBody() {
	body := `{"greet": "hello world!"}`
	suite.Logger.SetRequestBody(body)
	assert.Equal(suite.T(), body, suite.Logger.field.RequestBody)
}

func (suite *LogTestSuite) TestSetRequestHeader() {
	header := `"Content-Type": "application/json"`
	suite.Logger.SetRequestHeaders(header)
	assert.Equal(suite.T(), header, suite.Logger.field.RequestHeader)
}

func (suite *LogTestSuite) TestSetResponseBody() {
	body := `{"greet": "hello world!"}`
	suite.Logger.SetResponseBody(body)
	assert.Equal(suite.T(), body, suite.Logger.field.ResponseBody)
}

func (suite *LogTestSuite) TestSetResponseHeader() {
	header := `"Content-Type": "application/json"`
	suite.Logger.SetResponseHeaders(header)
	assert.Equal(suite.T(), header, suite.Logger.field.ResponseHeader)
}

func (suite *LogTestSuite) TestSetErrorMessage() {
	err := "Internal server error"
	suite.Logger.SetErrorMessage(err)
	assert.Equal(suite.T(), err, suite.Logger.field.ErrorMessage)
}

func (suite *LogTestSuite) TestSetErrorTypeErrorMessage() {
	err := errors.New("Internal server error")
	suite.Logger.SetErrorMessage(err)
	assert.Equal(suite.T(), err, suite.Logger.field.ErrorMessage)
}

func TestLogTestSuite(t *testing.T) {
	suite.Run(t, new(LogTestSuite))
}
