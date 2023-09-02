package zulip

import (
	"bytes"
	"fmt"
	"github.com/TwiN/gatus/v5/alerting/alert"
	"github.com/TwiN/gatus/v5/client"
	"github.com/TwiN/gatus/v5/core"
	"io"
	"net/http"
	"net/url"
)

type AlertProvider struct {
	Email      string `yaml:"email"`
	APIKey     string `yaml:"api-key"`
	APIURL     string `yaml:"api-url"`
	StreamName string `yaml:"stream-name"`
	TopicName  string `yaml:"topic-name"`

	ClientConfig *client.Config `yaml:"client,omitempty"`

	DefaltAlert *alert.Alert `yaml:"default-alert,omitempty"`
}

func (provider *AlertProvider) IsValid() bool {
	if provider.ClientConfig == nil {
		provider.ClientConfig = client.GetDefaultConfig()
	}
	return len(provider.Email) > 0 && len(provider.APIKey) > 0 && len(provider.APIURL) > 0 && len(provider.StreamName) > 0 && len(provider.TopicName) > 0
}

func (provider *AlertProvider) Send(endpoint *core.Endpoint, alert *alert.Alert, result *core.Result, resolved bool) error {
	buffer := bytes.NewBufferString(provider.buildRequestBody(endpoint, alert, result, resolved))

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/messages", provider.APIURL), buffer)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(provider.Email, provider.APIKey)

	response, err := client.GetHTTPClient(provider.ClientConfig).Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 399 {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("call to provider alert returned status code %d: %s", response.StatusCode, string(body))
	}
	return err
}

func (provider *AlertProvider) buildRequestBody(endpoint *core.Endpoint, alert *alert.Alert, result *core.Result, resolved bool) string {
	var message, results string

	if resolved {
		message = fmt.Sprintf("An alert for *%s* has been resolved:\n\n Healthcheck passing successfully %d time(s) in a row", endpoint.DisplayName(), alert.SuccessThreshold)
	} else {
		message = fmt.Sprintf("An alert for *%s* has been triggered:\n\n Healthcheck failed %d time(s) in a row", endpoint.DisplayName(), alert.FailureThreshold)
	}
	for _, conditionResult := range result.ConditionResults {
		var prefix string
		if conditionResult.Success {
			prefix = "✅"
		} else {
			prefix = "❌"
		}
		results += fmt.Sprintf("%s - `%s`\n", prefix, conditionResult.Condition)
	}
	var text string
	if len(alert.GetDescription()) > 0 {
		text = fmt.Sprintf("⛑ *Gatus* \n%s \n*Description* \n_%s_  \n\n*Condition results*\n%s", message, alert.GetDescription(), results)
	} else {
		text = fmt.Sprintf("⛑ *Gatus* \n%s \n*Condition results*\n%s", message, results)
	}

	data := url.Values{}
	data.Set("type", "stream")
	data.Set("to", provider.StreamName)
	data.Set("topic", provider.TopicName)
	data.Set("content", text)

	return data.Encode()
}

func (provider *AlertProvider) GetDefaultAlert() *alert.Alert {
	return provider.DefaltAlert
}
