package realtime

// feedData is the format to send realtime data
type feedData struct {
	QueryID   string                 `json:"id" mapstructure:"id" structs:"id"`
	Find      map[string]interface{} `json:"find" structs:"find"`
	Type      string                 `json:"type" structs:"type"`
	Payload   interface{}            `json:"payload" structs:"payload"`
	TimeStamp int64                  `json:"time" structs:"time"`
	Group     string                 `json:"group" structs:"group"`
	DBType    string                 `json:"dbType" structs:"dbType"`
	TypeName  string                 `json:"__typename,omitempty" structs:"__typename,omitempty"`
}
