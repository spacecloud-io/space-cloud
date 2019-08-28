package model

import (
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// RaftCommand is the object passed as a raft entry
type RaftCommand struct {
	Kind      utils.RaftCommandType   `json:"kind"`
	ID        string                  `json:"projectId"`
	Project   *config.Project         `json:"project"`
	Deploy    *config.Deploy          `json:"deploy"`
	Operation *config.OperationConfig `json:"operation"`
	Static    *config.Static          `json:"static"`
}
