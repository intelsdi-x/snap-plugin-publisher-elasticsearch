/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	Namespace   = "Namespace"
	Config      = "Config"
	Data        = "Data"
	Description = "Description"
	Timestamp   = "Timestamp"
	Unit        = "Unit"
	Version     = "Version"
)

type ElasticsearchPublisher struct {
}

func (els ElasticsearchPublisher) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewStringRule([]string{""}, "protocol", false, plugin.SetDefaultString("http"))
	policy.AddNewStringRule([]string{""}, "address", false, plugin.SetDefaultString("localhost"))
	policy.AddNewIntRule([]string{""}, "port", false, plugin.SetDefaultInt(9200))
	policy.AddNewStringRule([]string{""}, "index", false, plugin.SetDefaultString("log"))
	policy.AddNewStringRule([]string{""}, "type", false, plugin.SetDefaultString("message"))
	policy.AddNewBoolRule([]string{""}, "index_timestamp", false, plugin.SetDefaultBool(true))
	policy.AddNewStringRule([]string{""}, "publish_fields", false,
		plugin.SetDefaultString("Namespace|Config|Data|Description|Timestamp|Unit|Version"))
	return *policy, nil
}

func (els ElasticsearchPublisher) Publish(mts []plugin.Metric, cfg plugin.Config) error {
	url, err := createURL(cfg)
	publish_fields, err := readFieldsToPublish(cfg)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(mts))
	c := make(chan error, len(mts))
	for _, metric := range mts {
		go func(metric plugin.Metric) {
			defer wg.Done()
			c <- sendRequest(metric, url, publish_fields)
		}(metric)
	}
	wg.Wait()
	close(c)
	for status := range c {
		if status != nil {
			log.WithFields(log.Fields{
				"_block": "Publish",
				"status": status,
			}).Errorf("Publishing failed")
			return status
		}
	}
	return nil
}

func sendRequest(metric plugin.Metric, url string, fieldsToPublish []string) error {
	//Change metric structure to map instead of serialize it to JSON.
	//It is caused by serializing Namespace as a structure not as string. It is not readable in database
	metricToPublish := make(map[string]interface{})
	for _, field := range fieldsToPublish {
		metricToPublish[field] = getValue(metric, field)
	}

	for k, v := range metric.Tags {
		metricToPublish[k] = v
	}
	data, err := json.Marshal(metricToPublish)
	if err != nil {
		return err
	}
	r := bytes.NewReader(data)
	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return err
	}
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		resp.Body.Close()
		return nil
	}
	log.WithFields(log.Fields{
		"_block": "Publish",
		"url":    url,
	}).Errorf("Unable to publish metric %v, status: %v", strings.Join(metric.Namespace.Strings(), "/"), resp.Status)
	resp.Body.Close()
	return fmt.Errorf("Unable to publish metric %v, status: %v", strings.Join(metric.Namespace.Strings(), "/"), resp.Status)
}

func getValue(metric plugin.Metric, field string) interface{} {
	switch field {
	case Namespace:
		return strings.Join(metric.Namespace.Strings(), "/")
	case Config:
		return metric.Config
	case Data:
		return metric.Data
	case Description:
		return metric.Description
	case Timestamp:
		return metric.Timestamp
	case Unit:
		return metric.Unit
	case Version:
		return metric.Version
	default:
		return nil
	}
}

func createURL(cfg plugin.Config) (string, error) {
	//Protocol specification (HTTP, HTTPS)
	protocol, err := cfg.GetString("protocol")
	if err != nil {
		return "", err
	}
	//Elasticsearch address
	address, err := cfg.GetString("address")
	if err != nil {
		return "", err
	}
	//Elasticsearch port
	port, err := cfg.GetInt("port")
	if err != nil {
		return "", err
	}
	//Add timestamp to index, default value is true
	index_timestamp := true
	index_timestamp, err = cfg.GetBool("index_timestamp")
	if err != nil {
		return "", err
	}
	//Index to publish data
	pref, err := cfg.GetString("index")
	if err != nil {
		return "", err
	}
	index := ""
	//If index_timestamp is set, index will be extended with timestamp e.x. (log-2016.10.16)
	if index_timestamp {
		now := time.Now()
		index = fmt.Sprintf("%v-%v", pref, now.Format("2006.01.02"))
	} else {
		index = pref
	}
	//Message type in database
	messageType := "message"
	messageType, err = cfg.GetString("type")
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%v://%v:%v/%v/%v", protocol, address, port, index, messageType)
	return url, nil
}

func readFieldsToPublish(cfg plugin.Config) ([]string, error) {
	acceptedFields := []string{Namespace, Config, Data, Version, Timestamp, Unit}
	var fieldsToPublish []string
	fieldsToPublish = make([]string, 0)
	configValue, err := cfg.GetString("publish_fields")
	if err != nil {
		return []string{""}, err
	}
	configValues := strings.Split(configValue, "|")
	for _, v := range configValues {
		if listContains(v, acceptedFields) {
			fieldsToPublish = append(fieldsToPublish, v)
		} else {
			log.WithFields(log.Fields{
				"_block":           "Publish",
				"field_to_publish": v,
			}).Warnf("Unsupported field passed ", v)
		}
	}
	return fieldsToPublish, nil
}

func listContains(s string, list []string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}
