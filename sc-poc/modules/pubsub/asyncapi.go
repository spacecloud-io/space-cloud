package pubsub

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/spacecloud-io/space-cloud/managers/apis"
)

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

func (a *App) generateASyncAPIDoc() *AsyncAPI {
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
	for channelPath, channelObj := range channels.Channels {
		// Producer channel
		asyncapi.AddChannel(channelPath+"/producer", ChannelItem{
			Publish: &Operation{
				ID: "producerPublish" + getID(channelObj.Name),
				Message: MessageOneOrMany{
					MessageEntity: MessageEntity{
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
											"type": "object",
											"additionalProperties": map[string]interface{}{
												"type": "string",
											},
										},
										"requireAck": map[string]interface{}{
											"type": "boolean",
										},
										"payload": channelObj.Payload.Schema,
									},
									"required": []string{"id", "payload"},
								},
							},
							"required": []string{"event", "data"},
						},
					},
				},
			},
			Subscribe: &Operation{
				ID: "producerSubscribe" + getID(channelObj.Name),
				Message: MessageOneOrMany{
					MessageEntity: MessageEntity{
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
									"required": []string{"id", "ack"},
								},
							},
							"required": []string{"event", "data"},
						},
					},
				},
			},
		})

		// Consumer channel
		asyncapi.AddChannel(channelPath+"/consumer", ChannelItem{
			Publish: &Operation{
				ID: "consumerPublish" + getID(channelObj.Name),
				Message: MessageOneOrMany{
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
										"required": []string{"mode", "capacity", "autoack", "format"},
									},
								},
								"required": []string{"event", "data"},
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
										"required": []string{"id", "ack"},
									},
								},
								"required": []string{"event", "data"},
							},
						},
					},
				},
			},
			Subscribe: &Operation{
				ID: "consumerSubscribe" + getID(channelObj.Name),
				Message: MessageOneOrMany{
					MessageEntity: MessageEntity{
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
											"type": "object",
											"additionalProperties": map[string]interface{}{
												"type": "string",
											},
										},
										"payload": channelObj.Payload.Schema,
									},
									"required": []string{"id", "payload"},
								},
							},
							"required": []string{"event", "data"},
						},
					},
				},
			},
		})
	}

	if channels.Components != nil {
		asyncapi.Components.Schemas = make(map[string]interface{})
		for k, v := range channels.Components.Schemas {
			asyncapi.Components.Schemas[k] = v
		}
	}

	return &asyncapi
}

func (a *App) exposeAsyncAPIDoc() *apis.API {
	return &apis.API{
		Name: "asyncapi",
		Path: "/v1/api/asyncapi.json",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			enc.Encode(a.asyncapiDoc)
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
