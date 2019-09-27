package log

import (
	"context"
	"net/http"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type LogTestSuite struct {
	suite.Suite
	Logger *Logger
	url    string
}

func (suite *LogTestSuite) SetupTest() {
	ctx := context.WithValue(context.Background(), "log_id", uuid.NewV1().String())
	suite.Logger = NewLogger(ctx)
	suite.url = "https://programming-quotes-api.herokuapp.com/quotes"
}

func (suite *LogTestSuite) TestLog() {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, suite.url, nil)
	if err != nil {
		suite.FailNow(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		suite.FailNow(err.Error())
	}
	defer resp.Body.Close()

	// b, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	suite.FailNow(err.Error())
	// }

	// suite.Logger.SetResponseBody(string(b))
	fields := Field{
		Endpoint:      req.URL.String(),
		Method:        req.Method,
		RequestBody:   nil,
		RequestHeader: nil,
		// ResponseBody:   string(b),
		ResponseHeader: resp.Header,
		ErrorMessage:   "error",
	}
	suite.Logger.Set(fields)
	suite.Logger.Log(ErrorLevel, "testlog")
	suite.Logger.Log(ErrorLevel, "testlog1234")
}

func TestLogTestSuite(t *testing.T) {
	suite.Run(t, new(LogTestSuite))
}
