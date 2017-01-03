/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package elasticsearch

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestElasticsearchInvalidConfig(t *testing.T) {
	el := ElasticsearchPublisher{}
	err := el.Publish(nil, nil)
	Convey("Check if Publish behave properly without proper config", t, func() {
		So(err, ShouldNotBeNil)
	})
}

func TestElasticsearchPublishOK(t *testing.T) {
	el := ElasticsearchPublisher{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var mts []plugin.Metric
	mts = make([]plugin.Metric, 1)
	url := strings.Split(server.URL, "/")
	addr := strings.Split(url[2], ":")[0]
	port, _ := strconv.ParseInt(strings.Split(url[2], ":")[1], 10, 64)
	mts[0] = plugin.Metric{Namespace: plugin.NewNamespace("A"), Data: 7}
	cfg := plugin.Config{"protocol": "http", "address": addr, "port": port, "index": "database", "type": "test",
		"index_timestamp": true, "publish_fields": "Namespace|Data"}
	err := el.Publish(mts, cfg)
	Convey("Check if Publish method accepts http status 200 (OK) in return", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestElasticsearchPublishCreated(t *testing.T) {
	el := ElasticsearchPublisher{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	var mts []plugin.Metric
	mts = make([]plugin.Metric, 1)
	url := strings.Split(server.URL, "/")
	addr := strings.Split(url[2], ":")[0]
	port, _ := strconv.ParseInt(strings.Split(url[2], ":")[1], 10, 64)
	mts[0] = plugin.Metric{Namespace: plugin.NewNamespace("A"), Data: 7}
	cfg := plugin.Config{"protocol": "http", "address": addr, "port": port, "index": "database", "type": "test",
		"index_timestamp": true, "publish_fields": "Namespace|Data"}
	err := el.Publish(mts, cfg)
	Convey("Check if Publish method accepts http status 201 (Created) in return", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestElasticsearchPublishDynamicMetric(t *testing.T) {
	el := ElasticsearchPublisher{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	var mts []plugin.Metric
	mts = make([]plugin.Metric, 1)
	url := strings.Split(server.URL, "/")
	addr := strings.Split(url[2], ":")[0]
	port, _ := strconv.ParseInt(strings.Split(url[2], ":")[1], 10, 64)
	mts[0] = plugin.Metric{Namespace: plugin.NewNamespace("A").AddDynamicElement("B", "test").AddStaticElement("C"), Data: 7}
	cfg := plugin.Config{"protocol": "http", "address": addr, "port": port, "index": "database", "type": "test",
		"index_timestamp": true, "publish_fields": "Namespace|Data"}
	err := el.Publish(mts, cfg)
	Convey("Check if Publish method accepts http status 201 (Created) in return", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestElasticsearchPublishBadRequest(t *testing.T) {
	el := ElasticsearchPublisher{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	var mts []plugin.Metric
	mts = make([]plugin.Metric, 1)
	url := strings.Split(server.URL, "/")
	addr := strings.Split(url[2], ":")[0]
	port, _ := strconv.ParseInt(strings.Split(url[2], ":")[1], 10, 64)
	mts[0] = plugin.Metric{Namespace: plugin.NewNamespace("A"), Data: 7}
	cfg := plugin.Config{"protocol": "http", "address": addr, "port": port, "index": "database", "type": "test",
		"index_timestamp": true, "publish_fields": "Namespace|Data"}
	err := el.Publish(mts, cfg)
	Convey("Check if Publish method returns error on http status 400 (Bad Request) in return", t, func() {
		So(err, ShouldNotBeNil)
	})
}

func TestElasticsearchPublishNotFoud(t *testing.T) {
	el := ElasticsearchPublisher{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	var mts []plugin.Metric
	mts = make([]plugin.Metric, 1)
	url := strings.Split(server.URL, "/")
	addr := strings.Split(url[2], ":")[0]
	port, _ := strconv.ParseInt(strings.Split(url[2], ":")[1], 10, 64)
	mts[0] = plugin.Metric{Namespace: plugin.NewNamespace("A"), Data: 7}
	cfg := plugin.Config{"protocol": "http", "address": addr, "port": port, "index": "database", "type": "test",
		"index_timestamp": true, "publish_fields": "Namespace|Data"}
	err := el.Publish(mts, cfg)
	Convey("Check if Publish method returns error on http status 404 (Not Found) in return", t, func() {
		So(err, ShouldNotBeNil)
	})
}

func TestElasticsearchPublishWithoutServer(t *testing.T) {
	el := ElasticsearchPublisher{}

	var mts []plugin.Metric
	mts = make([]plugin.Metric, 1)
	mts[0] = plugin.Metric{Namespace: plugin.NewNamespace("A"), Data: 7}
	cfg := plugin.Config{"protocol": "http", "address": "localhost", "port": "9200", "index": "database", "type": "test",
		"index_timestamp": true, "publish_fields": "Namespace|Data"}
	err := el.Publish(mts, cfg)
	Convey("Check if Publish method handle server error", t, func() {
		So(err, ShouldNotBeNil)
	})
}
