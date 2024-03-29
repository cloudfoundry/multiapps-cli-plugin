// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new operations API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for operations API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
ExecuteOperationAction Executes a particular action over Multi-Target Application operation

*/
func (a *Client) ExecuteOperationAction(params *ExecuteOperationActionParams, authInfo runtime.ClientAuthInfoWriter) (*ExecuteOperationActionAccepted, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewExecuteOperationActionParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "ExecuteOperationAction",
		Method:             "POST",
		PathPattern:        "/operations/{operationId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ExecuteOperationActionReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ExecuteOperationActionAccepted), nil

}

/*
GetMta Retrieves Multi-Target Application in a space

*/
func (a *Client) GetMta(params *GetMtaParams, authInfo runtime.ClientAuthInfoWriter) (*GetMtaOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMtaParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetMta",
		Method:             "GET",
		PathPattern:        "/mtas/{mta_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetMtaReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetMtaOK), nil

}

/*
GetMtaFiles Retrieves all Multi-Target Application files

*/
func (a *Client) GetMtaFiles(params *GetMtaFilesParams, authInfo runtime.ClientAuthInfoWriter) (*GetMtaFilesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMtaFilesParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetMtaFiles",
		Method:             "GET",
		PathPattern:        "/files",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetMtaFilesReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetMtaFilesOK), nil

}

/*
GetMtaOperation Retrieves Multi-Target Application operation

*/
func (a *Client) GetMtaOperation(params *GetMtaOperationParams, authInfo runtime.ClientAuthInfoWriter) (*GetMtaOperationOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMtaOperationParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetMtaOperation",
		Method:             "GET",
		PathPattern:        "/operations/{operationId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetMtaOperationReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetMtaOperationOK), nil

}

/*
GetMtaOperationLogContent Retrieves the log content for Multi-Target Application operation

*/
func (a *Client) GetMtaOperationLogContent(params *GetMtaOperationLogContentParams, authInfo runtime.ClientAuthInfoWriter) (*GetMtaOperationLogContentOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMtaOperationLogContentParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetMtaOperationLogContent",
		Method:             "GET",
		PathPattern:        "/operations/{operationId}/logs/{logId}/content",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetMtaOperationLogContentReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetMtaOperationLogContentOK), nil

}

/*
GetMtaOperationLogs Retrieves the logs Multi-Target Application operation

*/
func (a *Client) GetMtaOperationLogs(params *GetMtaOperationLogsParams, authInfo runtime.ClientAuthInfoWriter) (*GetMtaOperationLogsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMtaOperationLogsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetMtaOperationLogs",
		Method:             "GET",
		PathPattern:        "/operations/{operationId}/logs",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetMtaOperationLogsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetMtaOperationLogsOK), nil

}

/*
GetMtaOperations Retrieves Multi-Target Application operations

*/
func (a *Client) GetMtaOperations(params *GetMtaOperationsParams, authInfo runtime.ClientAuthInfoWriter) (*GetMtaOperationsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMtaOperationsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetMtaOperations",
		Method:             "GET",
		PathPattern:        "/operations",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetMtaOperationsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetMtaOperationsOK), nil

}

/*
GetMtas Retrieves all Multi-Target Applications in a space

*/
func (a *Client) GetMtas(params *GetMtasParams, authInfo runtime.ClientAuthInfoWriter) (*GetMtasOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMtasParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetMtas",
		Method:             "GET",
		PathPattern:        "/mtas",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetMtasReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetMtasOK), nil

}

/*
GetOperationActions Retrieves available actions for Multi-Target Application operation

*/
func (a *Client) GetOperationActions(params *GetOperationActionsParams, authInfo runtime.ClientAuthInfoWriter) (*GetOperationActionsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetOperationActionsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetOperationActions",
		Method:             "GET",
		PathPattern:        "/operations/{operationId}/actions",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetOperationActionsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetOperationActionsOK), nil

}

/*
StartMtaOperation Starts execution of a Multi-Target Application operation

*/
func (a *Client) StartMtaOperation(params *StartMtaOperationParams, authInfo runtime.ClientAuthInfoWriter) (*StartMtaOperationAccepted, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewStartMtaOperationParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "StartMtaOperation",
		Method:             "POST",
		PathPattern:        "/operations",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &StartMtaOperationReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*StartMtaOperationAccepted), nil

}

/*
GetCsrfToken Retrieves a csrf-token header

*/
func (a *Client) GetCsrfToken(params *GetCsrfTokenParams, authInfo runtime.ClientAuthInfoWriter) (*GetCsrfTokenNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetCsrfTokenParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getCsrfToken",
		Method:             "GET",
		PathPattern:        "/csrf-token",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetCsrfTokenReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetCsrfTokenNoContent), nil

}

/*
GetInfo Retrieve information about the Deploy Service application

*/
func (a *Client) GetInfo(params *GetInfoParams, authInfo runtime.ClientAuthInfoWriter) (*GetInfoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetInfoParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getInfo",
		Method:             "GET",
		PathPattern:        "/info",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetInfoReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetInfoOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
