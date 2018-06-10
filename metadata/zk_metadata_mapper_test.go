package metadata_test

import (
	"go-graph/common"
	"go-graph/metadata"
	"path"
	"testing"

	"github.com/google/uuid"
)

func DoInitialize(m *metadata.ZkMetadataMapper) error {
	return m.Initialize()
}

func DoCreateGraph(m *metadata.ZkMetadataMapper, graphID uuid.UUID) error {
	return m.CreateGraph(graphID)
}

func DoCreatePartition(m *metadata.ZkMetadataMapper, graphID uuid.UUID, partitionID uuid.UUID) error {
	return m.CreatePartition(graphID, partitionID)
}

func DoCreateVertex(m *metadata.ZkMetadataMapper, graphID uuid.UUID, partitionID uuid.UUID, vertexID uuid.UUID) error {
	return m.CreateVertex(graphID, partitionID, vertexID)
}

func TestZkMetadataMapper_Initialize(t *testing.T) {
	var succ bool
	var err error
	m := metadata.NewMetadataMapper()
	err = m.Initialize()
	common.Assert(err == nil, t)

	succ, _, err = m.Connection.Exists(path.Join("", metadata.GRAPH))
	common.Assert(succ == true, t)
	common.Assert(err == nil, t)

	succ, _, err = m.Connection.Exists(path.Join("", metadata.BACKEND))
	common.Assert(succ == true, t)
	common.Assert(err == nil, t)
}

func TestZkMetadataMapper_CreateGraph(t *testing.T) {

}
