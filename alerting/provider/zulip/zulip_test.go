package zulip

import (
	"github.com/TwiN/gatus/v5/alerting/alert"
	"github.com/TwiN/gatus/v5/client"
	"github.com/TwiN/gatus/v5/core"
	"testing"
)

func TestAlertProvider_IsValid(t *testing.T) {
	t.Run("invalid-provider", func(t *testing.T) {
		invalidProvider := AlertProvider{}
		if invalidProvider.IsValid() {
			t.Error("provider shouldn't have been valid")
		}
	})

	t.Run("valid-provider", func(t *testing.T) {
		validProvider := AlertProvider{Email: "zulip-bot@emailaddr.fi", APIKey: "zyx57W2v1u123ew11aa1d", APIURL: "https://zulip.instance.fi", StreamName: "gatus", TopicName: "uptime-alerts"}
		if validProvider.ClientConfig != nil {
			t.Error("provider client config should have been nil prior to IsValid() being executed")
		}
		if !validProvider.IsValid() {
			t.Error("provider should've been valid")
		}
		if validProvider.ClientConfig == nil {
			t.Error("provider client config should have been set after IsValid() was executed")
		}
	})
}

func TestAlertProvider_Send(t *testing.T) {
	type fields struct {
		Email        string
		APIKey       string
		APIURL       string
		StreamName   string
		TopicName    string
		ClientConfig *client.Config
		DefaltAlert  *alert.Alert
	}
	type args struct {
		endpoint *core.Endpoint
		alert    *alert.Alert
		result   *core.Result
		resolved bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &AlertProvider{
				Email:        tt.fields.Email,
				APIKey:       tt.fields.APIKey,
				APIURL:       tt.fields.APIURL,
				StreamName:   tt.fields.StreamName,
				TopicName:    tt.fields.TopicName,
				ClientConfig: tt.fields.ClientConfig,
				DefaltAlert:  tt.fields.DefaltAlert,
			}
			if err := provider.Send(tt.args.endpoint, tt.args.alert, tt.args.result, tt.args.resolved); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlertProvider_buildRequestBody(t *testing.T) {
	type fields struct {
		Email        string
		APIKey       string
		APIURL       string
		StreamName   string
		TopicName    string
		ClientConfig *client.Config
		DefaltAlert  *alert.Alert
	}
	type args struct {
		endpoint *core.Endpoint
		alert    *alert.Alert
		result   *core.Result
		resolved bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &AlertProvider{
				Email:        tt.fields.Email,
				APIKey:       tt.fields.APIKey,
				APIURL:       tt.fields.APIURL,
				StreamName:   tt.fields.StreamName,
				TopicName:    tt.fields.TopicName,
				ClientConfig: tt.fields.ClientConfig,
				DefaltAlert:  tt.fields.DefaltAlert,
			}
			if got := provider.buildRequestBody(tt.args.endpoint, tt.args.alert, tt.args.result, tt.args.resolved); got != tt.want {
				t.Errorf("buildRequestBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
