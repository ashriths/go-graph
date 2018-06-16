package locator

import (
	"github.com/ashriths/go-graph/common"
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/metadata"
	"github.com/ashriths/go-graph/storage"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
	"math/rand"
	"time"
	"errors"
)

const (
	ELEMENTS_PER_PARTITION = 10
)

type Locator interface {
	FindPartition(element graph.ElementInterface) (uuid.UUID, error)
	RelocateConnectedElements(element graph.ElementInterface) (uuid.UUID, error)
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

func (randomLocator *RandomLocator) RelocateConnectedElements(element graph.ElementInterface) (uuid.UUID, error) {
	//panic("implement me")
	if graph.GetElementType(element) == graph.EDGE {
		edge := element.(graph.EdgeInterface)
		err, src := edge.GetSrcVertex()
		if err != nil {
			system.Logln("Failed to fetch source vertex for edge: ", edge.GetUUID().String())
			return uuid.New(), err
		}
		partitionID, _, err := randomLocator.Metadata.GetVertexLocation(src.GetGraphId(), edge.GetUUID())
		if err != nil {
			system.Logln("Failed to fetch partitionid for vertex: ", src.GetUUID().String())
			return uuid.New(), err
		}
		return *partitionID, nil
	}
	return uuid.New(), errors.New("Element not of type EDGE")
}

func (randomLocator *RandomLocator) createPartition(element graph.ElementInterface) (uuid.UUID, error) {
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
		if count == common.REPLICATION_FACTOR {
			break
		}

		data, err := randomLocator.Metadata.GetBackendInformation(backendId)
		if err != nil {
			system.Logln("Failed to fetch backend Info")
			return uuid.New(), err
		}
		backendAddr := data["address"]
		stClient := storage.NewStorageClient(backendAddr.(string))
		uuidList := [2]uuid.UUID{element.GetGraphId(), partitionID}
		var succ bool
		err = stClient.RegisterToHostPartition(uuidList[:], &succ)
		if err == nil {
			count += 1
		}
	}
	if count < common.REPLICATION_FACTOR {
		system.Logln("Failed to replicate to ", common.REPLICATION_FACTOR, " backends")
		return uuid.New(), err
	}
	//err = randomLocator.Metadata.SetPartitionInformation(element.GetGraphId(), partitionID)
	if err != nil {
		system.Logln("Failed to update element count in partition: ", partitionID)
		return uuid.New(), err
	}
	return partitionID, nil
}

func (randomLocator *RandomLocator) FindPartition(element graph.ElementInterface) (uuid.UUID, error) {

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
			if data["vertexCount"].(float64) < ELEMENTS_PER_PARTITION {
				partitionsWithSpaceArr = append(partitionsWithSpaceArr, partitionUUID)
			}
		}
		if len(partitionsWithSpaceArr) == 0 {
			return randomLocator.createPartition(element)
		}
		randInd := random(0, len(partitionsWithSpaceArr))
		partitionID = partitionsWithSpaceArr[randInd]
	}
	//err = randomLocator.Metadata.SetPartitionInformation(element.GetGraphId(), partitionID)
	if err != nil {
		system.Logln("Failed to update element count in partition: ", partitionID)
		return uuid.New(), err
	}
	return partitionID, nil
}

type CCLocator struct {
	Metadata metadata.Metadata
	StClient []*storage.StorageClient
}

func (ccLocator *CCLocator) FindPartition(element graph.ElementInterface) (uuid.UUID, error) {
	//panic("implement me")
	partitions, err := ccLocator.Metadata.GetAllPartitions(element.GetGraphId())
	if err != nil {
		system.Logln("Failed to get partitions from zookeeper")
		return uuid.New(), err
	}
	var partitionID uuid.UUID
	if len(partitions) == 0 {
		//No partitions, create a new one
		return ccLocator.createPartition(element)
	} else {
		var partitionsWithSpaceArr []uuid.UUID
		for _, partition := range partitions {
			partitionUUID, _ := uuid.Parse(partition)
			data, _ := ccLocator.Metadata.GetPartitionInformation(element.GetGraphId(), partitionUUID)
			if data["elementCount"].(int) < ELEMENTS_PER_PARTITION {
				partitionsWithSpaceArr = append(partitionsWithSpaceArr, partitionUUID)
			}
		}
		if len(partitionsWithSpaceArr) == 0 {
			return ccLocator.createPartition(element)
		}
		randInd := random(0, len(partitionsWithSpaceArr))
		partitionID = partitionsWithSpaceArr[randInd]
	}
	//err = randomLocator.Metadata.SetPartitionInformation(element.GetGraphId(), partitionID)
	if err != nil {
		system.Logln("Failed to update element count in partition: ", partitionID)
		return uuid.New(), err
	}
	return partitionID, nil
}

func (ccLocator *CCLocator) RelocateConnectedElements(element graph.ElementInterface) (uuid.UUID, error) {
	//panic("implement me")
	if graph.GetElementType(element) == graph.EDGE {
		edge := element.(graph.EdgeInterface)
		err, src := edge.GetSrcVertex()
		if err != nil {
			system.Logln("Failed to fetch source vertex for edge: ", edge.GetUUID().String())
			return uuid.New(), err
		}

		err, dest := edge.GetDestVertex()
		if err != nil {
			system.Logln("Failed to fetch destination vertex for edge: ", edge.GetUUID().String())
			return uuid.New(), err
		}

		err, srcParents := src.GetParentVertices([]string{""})
		if err != nil {
			system.Logln("Failed to fetch parent vertices for: ", src.GetUUID().String())
			return uuid.New(), err
		}

		err, srcChildren := src.GetChildVertices([]string{""})
		if err != nil {
			system.Logln("Failed to fetch child vertices for: ", src.GetUUID().String())
			return uuid.New(), err
		}

		neighbors := append(srcParents, srcChildren...)
		partitionCounts := map[*uuid.UUID]int{}
		partition, _, err := ccLocator.Metadata.GetVertexLocation(dest.GetGraphId(), dest.GetUUID())
		partitionCounts[partition] += 1
		for _, neighbor := range neighbors {
			partition, _, err := ccLocator.Metadata.GetVertexLocation(neighbor.GetGraphId(), neighbor.GetUUID())
			if err != nil {
				system.Logln("Failed to fetch partition ID for vertex : ", neighbor.GetUUID().String())
				return uuid.New(), err
			}
			partitionCounts[partition] += 1
		}

		var bestPartition *uuid.UUID
		var bestPartitionCount = 0
		for part, count := range partitionCounts {
			if count > bestPartitionCount {
				partData, err := ccLocator.Metadata.GetPartitionInformation(edge.GetGraphId(), *part)
				if err != nil {
					system.Logln("Failed to fetch partition data for partition : ", bestPartition.String())
					continue
				}
				if partData["elementCount"].(int) < ELEMENTS_PER_PARTITION {
					bestPartition = part
				}
			}
		}

		if bestPartition == nil {
			system.Logln("Failed to find any partition to relocate element to")
			return uuid.New(), 	errors.New("Failed to find any partition to relocate element to")
		}
		return *bestPartition, nil

	}
	return uuid.New(), errors.New("Element not an Edge")
}

func (ccLocator *CCLocator) createPartition(element graph.ElementInterface) (uuid.UUID, error) {
	partitionID := uuid.New()
	err := ccLocator.Metadata.CreatePartition(element.GetGraphId(), partitionID)
	if err != nil {
		system.Logln("Failed to create partition")
		return uuid.New(), err
	}
	backends, err := ccLocator.Metadata.GetAllBackends()
	if err != nil {
		system.Logln("Failed to fetch backends")
		return uuid.New(), err
	}

	count := 0
	for _, backendId := range backends {
		if count == common.REPLICATION_FACTOR {
			break
		}

		data, err := ccLocator.Metadata.GetBackendInformation(backendId)
		if err != nil {
			system.Logln("Failed to fetch backend Info")
			return uuid.New(), err
		}
		backendAddr := data["address"]
		stClient := storage.NewStorageClient(backendAddr.(string))
		uuidList := [2]uuid.UUID{element.GetGraphId(), partitionID}
		var succ bool
		err = stClient.RegisterToHostPartition(uuidList[:], &succ)
		if err == nil {
			count += 1
		}
	}
	if count < common.REPLICATION_FACTOR {
		system.Logln("Failed to replicate to ", common.REPLICATION_FACTOR, " backends")
		return uuid.New(), err
	}
	//err = randomLocator.Metadata.SetPartitionInformation(element.GetGraphId(), partitionID)
	if err != nil {
		system.Logln("Failed to update element count in partition: ", partitionID)
		return uuid.New(), err
	}
	return partitionID, nil
}
