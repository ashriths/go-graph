package locator

import (
	"fmt"
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Locator interface {
	FindPartition(element graph.Element) (uuid.UUID, error)
	FindBackend(element graph.Element, zkConn *metadata.ZkMetadataMapper, numBackends int) (string, error)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

type RandomLocator struct {
	Metadata metadata.Metadata
}

func (randomLocator *RandomLocator) FindBackend(element graph.Element,
	zkConn *metadata.ZkMetadataMapper, numBackends int) ([]string, error) {

	//panic("implement me")
	partitions, err := zkConn.GetAllPartitions()
	if err != nil {
		system.Logln("Failed to get partitions from zookeeper")
		return make([]string, 0), err
	}
	var partitionID uuid.UUID
	if len(partitions) == 0 {
		//No partitions, create a new one
		partitionID, err = uuid.NewUUID()
		if err != nil {
			system.Logln("Failed to generate UUID")
			return make([]string, 0), err
		}
		err = zkConn.CreatePartition(element.GetGraphID(), partitionID)
		if err != nil {
			system.Logln("Failed to create partition")
			return make([]string, 0), err
		}
	} else {
		randInd := random(0, len(partitions))
		partitionID, err = uuid.Parse(partitions[randInd])
		if err != nil {
			system.Logln("Failed to parse string to UUID")
			return make([]string, 0), err
		}
	}
	for i := 0; i < numBackends; i++ {

	}
}

type SizeBalancedLocator struct {
}

func (sizeBalancedLocator *SizeBalancedLocator) FindBackend(element graph.Element,
	zkConn *metadata.ZkMetadataMapper, numBackends int) (string, error) {
	panic("implement me")
}

type CCLocator struct {
}

func (ccLocator *CCLocator) FindBackend(element graph.Element,
	zkConn *metadata.ZkMetadataMapper, numBackends int) (string, error) {
	panic("implement me")
}
