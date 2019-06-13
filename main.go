package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/dylannz-sailthru/ad-data-generator/config"
	"github.com/dylannz-sailthru/ad-data-generator/generator"
	"github.com/sirupsen/logrus"
)

const (
	serviceName = "ad-data-generate"
)

var (
	logger logrus.FieldLogger
	c      *http.Client
)

func initLogger() error {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	parsedLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(parsedLevel)
	logger = logrus.WithField("service", serviceName)
	return nil
}

func check(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	check(godotenv.Load(".env"))
	check(initLogger())

	cfg, err := config.NewFromTOML(os.Getenv("CONFIG_FILE"))
	check(err)

	c = &http.Client{}

	logger.Debug("Resetting index...")

	_ = doESHTTP("DELETE", "/data/", http.NoBody)
	_ = doESHTTP("PUT", "/data/", bytes.NewBufferString(cfg.ES.Mapping))

	_ = doESHTTP("GET", "/_cluster/health?wait_for_status=yellow", http.NoBody)
	_ = doESHTTP("DELETE", "/_search/template/hotcloud", http.NoBody)
	_ = doESHTTP("POST", "/_search/template/hotcloud", bytes.NewBufferString(cfg.ES.Query))

	ctx := context.Background()
	generator.GenerateTimeline(ctx, cfg)

	_ = doESHTTP("PUT", "/data/_refresh", http.NoBody)

}

func doESHTTP(method, path string, body io.Reader) []byte {
	return doHTTP(method, os.Getenv("ELASTICSEARCH_HOST")+path, body)
}

func doHTTP(method, url string, body io.Reader) []byte {
	req, err := http.NewRequest(method, url, body)
	if body != http.NoBody {
		req.Header.Set("Content-Type", "application/json")
	}
	check(err)
	res, err := c.Do(req)
	check(err)
	defer res.Body.Close()
	logrus.Infof("%s %s: %d %s", method, url, res.StatusCode, http.StatusText(res.StatusCode))
	if res.StatusCode >= 400 && (method != "DELETE" || res.StatusCode != 404) {
		b, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(b))
	}
	b, err := ioutil.ReadAll(res.Body)
	check(err)
	return b
}
