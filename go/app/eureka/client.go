package eureka

import (
	"bytes"
	"iQuest/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/HikoQiu/go-eureka-client/eureka"
	"github.com/HikoQiu/go-feign/feign"
)

func init() {

	logFunc := func(level int, format string, a ...interface{}) {

		var logFunc *zerolog.Event
		switch level {
		case eureka.LevelDebug:
			logFunc = log.Debug()
		case eureka.LevelInfo:
			logFunc = log.Info()
		case eureka.LevelError:
			logFunc = log.Error()
		}
		if logFunc != nil {
			funcName, file, line, _ := runtime.Caller(2)
			fullFuncName := runtime.FuncForPC(funcName).Name()
			arr := strings.Split(fullFuncName, "/")
			arrFile := strings.Split(file, "/")

			logFunc.Str("file", arrFile[len(arrFile)-1]).Int("line", line).Str("func", arr[len(arr)-1]).Msgf(format, a...)
		}
	}
	eureka.SetLogger(logFunc)
	feign.SetLogger(logFunc)

	eurekaConfig := eureka.GetDefaultEurekaClientConfig()
	eurekaConfig.UseDnsForFetchingServiceUrls = false
	eurekaConfig.ServiceUrl = map[string]string{
		eureka.DEFAULT_ZONE: config.Viper.GetString("EUREKA_DEFAULT_ZONE"),
	}

	port := strings.Trim(config.Viper.GetString("RESTFUL_PORT"), ":")
	portInt, _ := strconv.Atoi(port)
	eureka.DefaultClient.Config(eurekaConfig).
		Register(config.Viper.GetString("EUREKA_SERVICE_NAME"), portInt).
		Run()

}

// Get 请求
func Get(rawURL string) (*http.Response, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	res, err := feign.DefaultFeign.App(strings.ToUpper(u.Host)).R().SetHeaders(map[string]string{
		"Content-Type": "application/json",
	}).SetQueryString(u.RawQuery).Get(u.Path)
	if err != nil {
		return nil, err
	}
	res.RawResponse.Body = ioutil.NopCloser(bytes.NewReader(res.Body()))
	return res.RawResponse, nil
}

// Post 请求
func Post(rawURL string, body *bytes.Reader) (*http.Response, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	res, err := feign.DefaultFeign.App(strings.ToUpper(u.Host)).R().SetHeaders(map[string]string{
		"Content-Type": "application/json",
	}).SetQueryString(u.RawQuery).SetBody(body).Post(u.Path)

	if err != nil {
		return nil, err
	}
	res.RawResponse.Body = ioutil.NopCloser(bytes.NewReader(res.Body()))
	return res.RawResponse, nil
}
