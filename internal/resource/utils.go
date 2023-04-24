package config

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	internaltypes "terraform-provider-pingaccess/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingaccess-go-client"
)

// Get BasicAuth context with a username and password
func BasicAuthContext(ctx context.Context, username, password string) context.Context {
	return context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: username,
		Password: password,
	})
}

// Get a BasicAuth context from a ProviderConfiguration
func ProviderBasicAuthContext(ctx context.Context, providerConfig internaltypes.ProviderConfiguration) context.Context {
	return BasicAuthContext(ctx, providerConfig.Username, providerConfig.Password)
}

// Error from PA API
type pingAccessError struct {
	Schemas []string `json:"schemas"`
	Status  string   `json:"status"`
	Detail  string   `json:"detail"`
}

// Report an HTTP error
func ReportHttpError(ctx context.Context, diagnostics *diag.Diagnostics, errorSummary string, err error, httpResp *http.Response) {
	httpErrorPrinted := false
	var internalError error
	if httpResp != nil {
		body, internalError := io.ReadAll(httpResp.Body)
		if internalError == nil {
			tflog.Debug(ctx, "Error HTTP response body: "+string(body))
			var paError pingAccessError
			internalError = json.Unmarshal(body, &paError)
			if internalError == nil {
				diagnostics.AddError(errorSummary, err.Error()+" - Detail: "+string(body))
				httpErrorPrinted = true
			}
		}
	}
	if !httpErrorPrinted {
		if internalError != nil {
			tflog.Warn(ctx, "Failed to unmarshal HTTP response body: "+internalError.Error())
		}
		diagnostics.AddError(errorSummary, err.Error())
	}
}

// Write out messages from the Config API response to tflog
func logMessages(ctx context.Context, messages *client.APIResponse) {
	if messages == nil {
		return
	}

	// for _, message := range messages.Message {
	tflog.Warn(ctx, "Configuration API Notification: "+messages.Message)
	// }

	// for _, action := range messages.RequiredActions {
	// 	actionJson, err := action.MarshalJSON()
	// 	if err != nil {
	// 		tflog.Warn(ctx, "Configuration API RequiredAction: "+string(actionJson))
	// 	}
	// }
}

// // Read messages from the Configuration API response
func ReadMessages(ctx context.Context, messages *client.APIResponse) types.Set {
	// Report any notifications from the Config API
	var Message types.Set
	if messages != nil {
		Message = internaltypes.GetStringSet(messages.Message)
		logMessages(ctx, messages)
	} else {
		Message, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	return Message
}
