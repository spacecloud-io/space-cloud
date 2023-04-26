package pubsub

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/spacecloud-io/space-cloud/managers/apis"
)

type AsyncAPI struct {
	SpecVersion string             `json:"asyncapi"` // Required
	Info        Info               `json:"info"`     // Required
	Channels    Channels           `json:"channels"` // Required
	Servers     Servers            `json:"servers,omitempty"`
	Components  AsyncAPIComponents `json:"components,omitempty"`
}

type Servers map[string]*ServerItem
type Channels map[string]*ChannelItem

type ChannelItem struct {
	Subscribe *Operation `json:"subscribe,omitempty"`
	Publish   *Operation `json:"publish,omitempty"`
}

type Operation struct {
	Message OneOf  `json:"message,omitempty"`
	ID      string `json:"operationId,omitempty"`
}

type OneOf struct {
	OneOf []MessageEntity `json:"oneOf,omitempty"`
}

type MessageEntity struct {
	Name        string                 `json:"name"`        // Required
	ContentType string                 `json:"contentType"` // Required
	Payload     map[string]interface{} `json:"payload"`     // Required
}

type Info struct {
	Title       string `json:"title"`   // Required
	Version     string `json:"version"` // Required
	Description string `json:"description,omitempty"`
}

// An object representing a Server.
type ServerItem struct {
	URL         string `json:"url"`      // Required.
	Protocol    string `json:"protocol"` // Required.
	Description string `json:"description,omitempty"`
}

type AsyncAPIComponents struct {
	Schemas map[string]interface{} `json:"schemas,omitempty"`
}

func (a *AsyncAPI) AddServer(name string, srv ServerItem) {
	if a.Servers == nil {
		a.Servers = make(map[string]*ServerItem)
	}

	a.Servers[name] = &srv
}

func (a *AsyncAPI) AddChannel(name string, channel ChannelItem) {
	if a.Channels == nil {
		a.Channels = make(map[string]*ChannelItem)
	}

	a.Channels[name] = &channel
}

func (a *App) generateASyncAPIDoc() *apis.API {
	return &apis.API{
		Name: "asyncapi",
		Path: "/v1/api/asyncapi.json",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add general Info
			asyncapi := AsyncAPI{}
			asyncapi.SpecVersion = "2.6.0"
			asyncapi.Info.Title = "SpaceCloud exposed AsyncAPIs"
			asyncapi.Info.Version = "v0.22.0"
			asyncapi.Info.Description = "Specification of all the AsyncAPIs exposed by the pubsub module of SpaceCloud"

			// Add server info
			asyncapi.AddServer("SpaceCloud", ServerItem{
				URL:         "/v1/pubsub/default",
				Protocol:    "websocket",
				Description: "URL for websocket connection",
			})

			// Add channels
			channels := a.Channels()
			for channelName, channelObj := range channels.Channels {
				// Producer channel
				asyncapi.AddChannel(channelName+"/producer", ChannelItem{
					Publish: &Operation{
						ID: "producerPublish" + getID(channelObj.Name),
						Message: OneOf{
							OneOf: []MessageEntity{
								{
									Name:        "Message",
									ContentType: "application/json",
									Payload: map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"event": map[string]interface{}{
												"type": "string",
											},
											"data": map[string]interface{}{
												"type": "object",
												"properties": map[string]interface{}{
													"id": map[string]interface{}{
														"type": "string",
													},
													"metadata": map[string]interface{}{
														"type":                 "object",
														"additionalProperties": true,
													},
													"requireAck": map[string]interface{}{
														"type": "boolean",
													},
													"payload": channelObj.Payload.Schema,
												},
											},
										},
									},
								},
							},
						},
					},
					Subscribe: &Operation{
						ID: "producerSubscribe" + getID(channelObj.Name),
						Message: OneOf{
							OneOf: []MessageEntity{
								{
									Name:        "Acknowledgement",
									ContentType: "application/json",
									Payload: map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"event": map[string]interface{}{
												"type": "string",
											},
											"data": map[string]interface{}{
												"type": "object",
												"properties": map[string]interface{}{
													"id": map[string]interface{}{
														"type": "string",
													},
													"ack": map[string]interface{}{
														"type": "boolean",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				})

				// Consumer channel
				asyncapi.AddChannel(channelName+"/consumer", ChannelItem{
					Publish: &Operation{
						ID: "consumerPublish" + getID(channelObj.Name),
						Message: OneOf{
							OneOf: []MessageEntity{
								{
									Name:        "startSubscribe",
									ContentType: "application/json",
									Payload: map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"event": map[string]interface{}{
												"type": "string",
											},
											"data": map[string]interface{}{
												"type": "object",
												"properties": map[string]interface{}{
													"mode": map[string]interface{}{
														"type": "string",
													},
													"capacity": map[string]interface{}{
														"type": "integer",
													},
													"autoack": map[string]interface{}{
														"type": "boolean",
													},
													"format": map[string]interface{}{
														"type": "string",
													},
												},
											},
										},
									},
								},
								{
									Name:        "Acknowlegement",
									ContentType: "application/json",
									Payload: map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"event": map[string]interface{}{
												"type": "string",
											},
											"data": map[string]interface{}{
												"type": "object",
												"properties": map[string]interface{}{
													"id": map[string]interface{}{
														"type": "string",
													},
													"ack": map[string]interface{}{
														"type": "boolean",
													},
												},
											},
										},
									},
								},
							},
						},
					},
					Subscribe: &Operation{
						ID: "consumerSubscribe" + getID(channelObj.Name),
						Message: OneOf{
							OneOf: []MessageEntity{
								{
									Name:        "Message",
									ContentType: "application/json",
									Payload: map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"event": map[string]interface{}{
												"type": "string",
											},
											"data": map[string]interface{}{
												"type": "object",
												"properties": map[string]interface{}{
													"id": map[string]interface{}{
														"type": "string",
													},
													"metadata": map[string]interface{}{
														"type":                 "object",
														"additionalProperties": true,
													},
													"payload": channelObj.Payload.Schema,
												},
											},
										},
									},
								},
							},
						},
					},
				})
			}

			if channels.Components != nil {
				asyncapi.Components.Schemas = channels.Components.Schemas
			}

			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			enc.Encode(asyncapi)
		}),
	}
}

func getID(name string) string {
	arr := strings.Split(name, "-")
	for i, item := range arr {
		arr[i] = strings.Title(item)
	}

	return strings.Join(arr, "")
}
