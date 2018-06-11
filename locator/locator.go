package locator

import (
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/storage"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const (
	REPLICATIONFACTOR      = 3
	ELEMENTS_PER_PARTITION = 10
)

type Locator interface {
	FindPartition(element graph.ElementInterface) (uuid.UUID, error)
	//FindBackend(element graph.Element, zkConn *metadata.ZkMetadataMapper, numBackends int) (string, error)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

type RandomLocator struct {
	Metadata metadata.Metadata
	StClient []*storage.StorageClient
}

func (randomLocator *RandomLocator) createPartition(element graph.Element) (uuid.UUID, error) {
	partitionID := uuid.New()
	err := randomLocator.Metadata.CreatePartition(element.GetGraphId(), partitionID)
	if err != nil {
		system.Logln("Failed to create partition")
		return uuid.New(), err
	}
	backends, err := randomLocator.Metadata.GetAllBackends()
	if err != nil {
		system.Logln("Failed to fetch backends")
		return uuid.New(), err
	}

	count := 0
	for _, backendId := range backends {
		if count == REPLICATIONFACTOR {
			break
		}

		data, err := randomLocator.Metadata.GetBackendInformation(backendId)
		if err != nil {
			system.Logln("Failed to fetch backend Info")
			return uuid.New(), err
		}
		backendAddr := data["address"]
		stClient := storage.NewStorageClient(backendAddr.(string))
		uuidList := [2]uuid.UUID{element.GraphUUID, partitionID}
		var succ bool
		err = stClient.RegisterToHostPartition(uuidList[:], &succ)
		if err == nil {
			count += 1
		}
	}
	if count < REPLICATIONFACTOR {
		system.Logln("Failed to replicate to ", REPLICATIONFACTOR, " backends")
		return uuid.New(), err
	}
	err = randomLocator.Metadata.SetPartitionInformation(element.GetGraphId(), partitionID)
	if err != nil {
		system.Logln("Failed to update element count in partition: ", partitionID)
		return uuid.New(), err
	}
	return partitionID, nil
}

func (randomLocator *RandomLocator) FindPartition(element graph.Element) (uuid.UUID, error) {

	partitions, err := randomLocator.Metadata.GetAllPartitions(element.GetGraphId())
	if err != nil {
		system.Logln("Failed to get partitions from zookeeper")
		return uuid.New(), err
	}
	var partitionID uuid.UUID
	if len(partitions) == 0 {
		//No partitions, create a new one
		return randomLocator.createPartition(element)
	} else {
		var partitionsWithSpaceArr []uuid.UUID
		for _, partition := range partitions {
			partitionUUID, _ := uuid.Parse(partition)
			data, _ := randomLocator.Metadata.GetPartitionInformation(element.GetGraphId(), partitionUUID)
			if data["elementCount"].(int) < ELEMENTS_PER_PARTITION {
				partitionsWithSpaceArr = append(partitionsWithSpaceArr, partitionUUID)
			}
		}
		if len(partitionsWithSpaceArr) == 0 {
			return randomLocator.createPartition(element)
		}
		randInd := random(0, len(partitionsWithSpaceArr))
		partitionID = partitionsWithSpaceArr[randInd]
	}
	err = randomLocator.Metadata.SetPartitionInformation(element.GetGraphId(), partitionID)
	if err != nil {
		system.Logln("Failed to update element count in partition: ", partitionID)
		return uuid.New(), err
	}
	return partitionID, nil
}

type SizeBalancedLocator struct {
}

func (sizeBalancedLocator *SizeBalancedLocator) FindPartition(element graph.Element) (uuid.UUID, error) {
	panic("implement me")
}

type CCLocator struct {
}

func (ccLocator *CCLocator) FindPartition(element graph.Element) (uuid.UUID, error) {
	panic("implement me")
}
