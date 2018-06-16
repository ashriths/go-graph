package metadata

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/ashriths/go-graph/system"

	"strings"

	"math/rand"

	"github.com/ashriths/go-graph/common"
	"github.com/google/uuid"
	"github.com/samuel/go-zookeeper/zk"
)

//Constants to be used
const (
	DEFAULTVERSION = -1
	ROOT           = "/"
	BACKEND        = "backends"
	GRAPH          = "graphs"
	PARTITION      = "partitions"
	VERTEX         = "vertices"
	EDGE           = "edges"

	BACKEND_PREFIX = "back-"
	EMPTY_DATA     = "{}"
)

type ZkMetadataMapper struct {
	Connection *zk.Conn
	ZkAddrs    []string
	err        error
}

func NewZkMetadataMapper(ZkAddrs []string) *ZkMetadataMapper {
	zkMapper := ZkMetadataMapper{
		ZkAddrs: ZkAddrs,
	}
	zkMapper.Initialize()
	return &zkMapper
}

//func NewMetadataMapper() *ZkMetadataMapper {
//	return &ZkMetadataMapper{
//		ZkAddrs: []string{"169.228.66.172:21810", "169.228.66.170:21810", "169.228.66.171:21810"},
//	}
//}

func (self *ZkMetadataMapper) Initialize() error {
	var err error
	err = self.createZnodeIfNotExists(path.Join(ROOT, GRAPH), EMPTY_DATA)
	if err != nil {
		system.Logf("Error while creating \"graphs\" znode")
		return err
	}
	err = self.createZnodeIfNotExists(path.Join(ROOT, BACKEND), EMPTY_DATA)
	if err != nil {
		system.Logf("Error while creating \"backends\" znode")
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) CreatePartition(graphID uuid.UUID, partitionId uuid.UUID) error {
	data := map[string]interface{}{"partitionId": partitionId.String(), "vertexCount": float64(0)}
	return self.createZnode(path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionId.String()), data)
}

func (self *ZkMetadataMapper) CreateVertex(graphID uuid.UUID, partitionId uuid.UUID, vertexID uuid.UUID) error {
	data := map[string]string{"partitionId": partitionId.String()}
	err := self.createZnode(path.Join(ROOT, GRAPH, graphID.String(), VERTEX, vertexID.String()), data)
	if err != nil {
		return err
	}
	err = self.IncrementElementCount(graphID, partitionId)
	if err != nil {
		system.Logf("Failed to increment element count at parition :%s", partitionId.String())
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) CreateEdge(graphID uuid.UUID, edgeID uuid.UUID, srcID uuid.UUID) error {
	data := map[string]string{"srcID": srcID.String()}
	return self.createZnode(path.Join(ROOT, GRAPH, graphID.String(), EDGE, edgeID.String()), data)
}

func (self *ZkMetadataMapper) CreateBackend(backendAddr string) (string, error) {
	var err error
	var data []byte
	var backendID string
	self.connect()
	znodePath := path.Join(ROOT, BACKEND, BACKEND_PREFIX)
	dataMap := map[string]string{"address": backendAddr}
	data, err = json.Marshal(dataMap)
	if err != nil {
		system.Logf("Error while Marshalling the backendAddr %s", backendAddr)
		return backendID, err
	}
	//backendID, err = conn.CreateProtectedEphemeralSequential(znodePath, data, zk.WorldACL(zk.PermAll))
	backendID, err = self.Connection.Create(znodePath, data, zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		system.Logf("Error while creating %s backend node", backendAddr)
		return backendID, err
	}
	backendID = strings.Split(backendID, "/")[2]
	return backendID, nil
}

//CreateGraph : Gets all graphs in the systems
func (self *ZkMetadataMapper) GetGraphs(graphIds *[]uuid.UUID) error {
	self.connect()
	znodePath := path.Join(ROOT, GRAPH)
	graphIdStrs, _, err := self.Connection.Children(znodePath)
	if err != nil {
		return err
	}
	for _, i := range graphIdStrs {
		u, _ := uuid.Parse(i)
		*graphIds = append(*graphIds, u)
	}
	return nil
}

//CreateGraph : creates a graph Znode
func (self *ZkMetadataMapper) CreateGraph(graphID uuid.UUID, data interface{}) error {
	err := self.createZnode(path.Join(ROOT, GRAPH, graphID.String()), data)
	if err != nil {
		system.Logf("Error while creating %s graph", graphID.String())
		return err
	}
	err = self.createZnode(path.Join(ROOT, GRAPH, graphID.String(), EDGE), EMPTY_DATA)
	if err != nil {
		system.Logf("Error while creating \"edges\" node under %s graph", graphID.String())
		return err
	}
	err = self.createZnode(path.Join(ROOT, GRAPH, graphID.String(), VERTEX), EMPTY_DATA)
	if err != nil {
		system.Logf("Error while creating \"vertices\" node under %s graph", graphID.String())
		return err
	}
	err = self.createZnode(path.Join(ROOT, GRAPH, graphID.String(), PARTITION), EMPTY_DATA)
	if err != nil {
		system.Logf("Error while creating \"partitions\" node under %s graph", graphID.String())
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) GetAllBackends() ([]string, error) {
	children, err := self.getChildren(path.Join(ROOT, BACKEND))
	if err != nil {
		system.Logf("Error while getting all backend IDs")
		return nil, err
	}
	return children, nil
}

func (self *ZkMetadataMapper) GetAllPartitions(graphID uuid.UUID) ([]string, error) {
	children, err := self.getChildren(path.Join(ROOT, GRAPH, graphID.String(), PARTITION))
	if err != nil {
		system.Logf("Error while getting all partitions of %s graph", graphID.String())
		return nil, err
	}
	return children, nil
}

func (self *ZkMetadataMapper) GetPartitionInformation(graphID uuid.UUID, partitionId uuid.UUID) (map[string]interface{}, error) {
	data, _, err := self.getZnodeData(path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionId.String()))
	if err != nil {
		system.Logf("Error while retrieving data stored at %s partition", partitionId.String())
		return nil, err
	}
	return data, nil
}

func (self *ZkMetadataMapper) UpdateVertexInformation(graphID uuid.UUID, vertexID uuid.UUID, key interface{}, value interface{}) error {
	var err error
	var exists bool
	var curData map[string]interface{}
	var statusInfo *zk.Stat

	znodePath := path.Join(ROOT, GRAPH, graphID.String(), VERTEX, vertexID.String())
	for {
		exists, err = self.checkZnodeExists(znodePath)
		if err != nil {
			system.Logf("Error while checking if %s vertex exists", vertexID.String())
			return err
		}
		if exists != true {
			system.Logf("%s vertex does not exist", vertexID.String())
			return fmt.Errorf("Partition does not exist")
		}
		curData, statusInfo, err = self.getZnodeData(znodePath)
		if err != nil {
			system.Logf("Error while retrieving data stored at %s vertex", vertexID.String())
			return err
		}
		curData[key.(string)] = value
		err = self.setZnodeData(znodePath, curData, statusInfo.Version)
		if err != nil {
			system.Logf("Error while updating data at %s vertex. Trying again..", vertexID.String())
			//err = self.setZnodeData(znodePath, curData, statusInfo.Version)
		} else {
			system.Logf("Succesfully updated data at vertex: %s", vertexID.String())
			break
		}
	}
	return nil
}

func (self *ZkMetadataMapper) UpdateEdgeInformation(graphID uuid.UUID, edgeID uuid.UUID, key interface{}, value interface{}) error {
	var err error
	var exists bool
	var curData map[string]interface{}
	var statusInfo *zk.Stat

	znodePath := path.Join(ROOT, GRAPH, graphID.String(), EDGE, edgeID.String())
	for {
		exists, err = self.checkZnodeExists(znodePath)
		if err != nil {
			system.Logf("Error while checking if %s edge exists", edgeID.String())
			return err
		}
		if exists != true {
			system.Logf("%s edge does not exist", edgeID.String())
			return fmt.Errorf("Partition does not exist")
		}
		curData, statusInfo, err = self.getZnodeData(znodePath)
		if err != nil {
			system.Logf("Error while retrieving data stored at %s edge", edgeID.String())
			return err
		}
		curData[key.(string)] = value.(string)
		err = self.setZnodeData(znodePath, curData, statusInfo.Version)
		if err != nil {
			system.Logf("Error while updating data at %s edge. Trying again..", edgeID.String())
		} else {
			system.Logf("Succesfully updated data at edge: %s", edgeID.String())
			break
		}
	}
	return nil
}

func (self *ZkMetadataMapper) DeleteVertexInformation(graphID uuid.UUID, vertexID uuid.UUID, key interface{}) error {
	var err error
	var exists bool
	var curData map[string]interface{}
	var statusInfo *zk.Stat

	znodePath := path.Join(ROOT, GRAPH, graphID.String(), VERTEX, vertexID.String())
	for {
		exists, err = self.checkZnodeExists(znodePath)
		if err != nil {
			system.Logf("Error while checking if %s vertex exists", vertexID.String())
			return err
		}
		if exists != true {
			system.Logf("%s vertex does not exist", vertexID.String())
			return fmt.Errorf("Partition does not exist")
		}
		curData, statusInfo, err = self.getZnodeData(znodePath)
		if err != nil {
			system.Logf("Error while retrieving data stored at %s vertex", vertexID.String())
			return err
		}
		delete(curData, key.(string))
		err = self.setZnodeData(znodePath, curData, statusInfo.Version)
		if err != nil {
			system.Logf("Error while deleting data at %s vertex. Trying again..", vertexID.String())
			//err = self.setZnodeData(znodePath, curData, statusInfo.Version)
		} else {
			system.Logf("Succesfully deleted data at vertex: %s", vertexID.String())
			break
		}
	}
	return nil
}

func (self *ZkMetadataMapper) GetEdgeInformation(graphID uuid.UUID, edgeID uuid.UUID) (map[string]interface{}, error) {
	var err error
	var exists bool
	var curData map[string]interface{}

	znodePath := path.Join(ROOT, GRAPH, graphID.String(), EDGE, edgeID.String())
	exists, err = self.checkZnodeExists(znodePath)
	if err != nil {
		system.Logf("Error while checking if %s edge exists", edgeID.String())
		return nil, err
	}
	if exists != true {
		system.Logf("%s edge does not exist", edgeID.String())
		return nil, fmt.Errorf("Edge does not exist")
	}
	curData, _, err = self.getZnodeData(znodePath)
	if err != nil {
		system.Logf("Error while retrieving data stored at %s edge", edgeID.String())
		return nil, err
	}
	return curData, err
}

func (self *ZkMetadataMapper) GetVertexInformation(graphID uuid.UUID, vertexID uuid.UUID) (map[string]interface{}, error) {
	var err error
	var exists bool
	var curData map[string]interface{}

	znodePath := path.Join(ROOT, GRAPH, graphID.String(), VERTEX, vertexID.String())
	exists, err = self.checkZnodeExists(znodePath)
	if err != nil {
		system.Logf("Error while checking if %s vertex exists", vertexID.String())
		return nil, err
	}
	if exists != true {
		system.Logf("%s vertex does not exist", vertexID.String())
		return nil, fmt.Errorf("Edge does not exist")
	}
	curData, _, err = self.getZnodeData(znodePath)
	if err != nil {
		system.Logf("Error while retrieving data stored at %s vertex", vertexID.String())
		return nil, err
	}
	return curData, err
}

func (self *ZkMetadataMapper) SetPartitionInformation(graphID uuid.UUID, partitionId uuid.UUID, data interface{}) error {
	var exists bool
	var err error

	znodePath := path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionId.String())
	exists, err = self.checkZnodeExists(znodePath)
	if err != nil {
		system.Logf("Error while checking if %s partition exists", partitionId.String())
		return err
	}
	if exists != true {
		system.Logf("%s partition does not exist", partitionId.String())
		return fmt.Errorf("Partition does not exist")
	}

	err = self.setZnodeData(znodePath, data, DEFAULTVERSION)
	if err != nil {
		system.Logf("Error while setting data at %s partition", partitionId.String())
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) GetBackendInformation(backendID string) (map[string]interface{}, error) {
	data, _, err := self.getZnodeData(path.Join(ROOT, BACKEND, backendID))
	if err != nil {
		system.Logf("Error while retrieving data stored at %s backend", backendID)
		return nil, err
	}
	return data, nil
}

// GetBackendsForPartition : Add backends to Partitions
func (self *ZkMetadataMapper) GetBackendsForPartition(graphID uuid.UUID, partitionId uuid.UUID) ([]string, error) {
	children, err := self.getChildren(path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionId.String()))
	if err != nil {
		system.Logln("Error while getting paritionID", partitionId, err)
		return nil, err
	}
	return children, nil
}

func (self *ZkMetadataMapper) AddBackendToPartition(graphID uuid.UUID, partitionId uuid.UUID, backendID string) ([]string, interface{}, error) {
	var err error
	var liveBackends []string
	watch := make(<-chan zk.Event)
	err = self.createZnode(path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionId.String(), backendID), EMPTY_DATA)
	if err != nil {
		system.Logf("Error while adding %s backend to %s partition", backendID, partitionId.String())
		return nil, nil, err
	}
	backendNode := path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionId.String())

	// pass the full path upto node whose children have to be watched
	liveBackends, watch, err = self.GetWatchOnChildren(backendNode)
	if err != nil {
		system.Logln("Error while setting a watch on ", backendID, err)
		return nil, nil, err
	}
	return liveBackends, watch, err
}

// @param : backendNode - full path of node whose children need to be watched
// @return : list of all alive children,
func (self *ZkMetadataMapper) GetWatchOnChildren(backendNode string) ([]string, <-chan zk.Event, error) {
	self.connect()
	snapshot, _, watch, err := self.Connection.ChildrenW(backendNode)
	return snapshot, watch, err
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (self *ZkMetadataMapper) GetVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) (*uuid.UUID, []string, error) {
	var partitionId string
	var children []string
	var err error

	// get partitionId from graphID and vertexID
	znodePath := path.Join(ROOT, GRAPH, graphID.String(), VERTEX, vertexID.String())
	data, _, err := self.getZnodeData(znodePath)
	if err != nil {
		system.Logf("Error while getting %s vertex", vertexID.String())
		return nil, nil, err
	}
	partitionId = data["partitionId"].(string)

	//  get backends from paritionID
	znodePath = path.Join(ROOT, GRAPH, graphID.String(), "partitions", partitionId)
	children, err = self.getChildren(znodePath)
	if err != nil {
		system.Logf("Error while getting %s paritionID", partitionId)
		return nil, nil, err
	}

	partitionUUID, err := uuid.Parse(partitionId)
	if err != nil {
		system.Logf("Error while parsing %s paritionID", partitionId)
		return nil, nil, err
	}
	//sort.Strings(children)
	// backendIDs = append([]string{data["Primary"].(string)}, backendIDs...)
	// backendIDs = append(backendIDs, data["Secondaries"].([]string)...)

	// for _, backendID := range backendIDs {
	// 	znodePath = path.Join("", BACKEND, backendID)
	// 	data, err = self.getZnodeData(znodePath)

	// 	if err != nil {
	// 		system.Logf()("Error while getting backend address of %s backendID", backendID)
	// 		return nil, err
	// 	}
	// 	children = append(children, data["addr"].(string))
	// }

	return &partitionUUID, children, err
}

func (self *ZkMetadataMapper) SetVertexLocation(graphID uuid.UUID, partitionId uuid.UUID, vertexID uuid.UUID) error {
	//conn := connect(self.Connection, self.err)
	//var buffer bytes.Buffer
	var exists bool
	var err error

	znodePath := path.Join(ROOT, GRAPH, graphID.String(), VERTEX, vertexID.String())
	exists, err = self.checkZnodeExists(znodePath)
	if err != nil {
		system.Logf("Error while checking if %s vertex exists", vertexID.String())
		return err
	}
	if exists != true {
		system.Logf("%s vertex does not exist", vertexID.String())
		return fmt.Errorf("Vertex does not exist")
	}

	err = self.setZnodeData(znodePath, map[string]string{"partitionId": partitionId.String()}, DEFAULTVERSION)
	if err != nil {
		system.Logf("Error while setting %s vertex with %s paritionID", vertexID.String(), partitionId.String())
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) GetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID) (*uuid.UUID, []string, error) {
	znodePath := path.Join(ROOT, GRAPH, graphID.String(), EDGE, edgeID.String())
	data, _, err := self.getZnodeData(znodePath)
	if err != nil {
		system.Logf("Failed to fetch data for znode: %s", znodePath)
		return nil, make([]string, 0), err
	}

	srcIdStr := data["srcID"].(string)
	srcID, err := uuid.Parse(srcIdStr)
	if err != nil {
		return nil, make([]string, 0), err
	}

	return self.GetVertexLocation(graphID, srcID)
}
func (self *ZkMetadataMapper) SetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID, srcID uuid.UUID) error {
	znodePath := path.Join(ROOT, GRAPH, graphID.String(), EDGE, edgeID.String())
	err := self.setZnodeData(znodePath, map[string]string{"srcID": srcID.String()}, DEFAULTVERSION)
	if err != nil {
		system.Logf("Failed to set data for znode: %s", znodePath)
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) DeleteVertex(graphID uuid.UUID, vertexID uuid.UUID) error {
	var partitionId *uuid.UUID
	znodePath := path.Join(ROOT, GRAPH, graphID.String(), VERTEX, vertexID.String())

	partitionId, _, err := self.GetVertexLocation(graphID, vertexID)
	if err != nil {
		system.Logf("Failed to get partition of vertex: %s", vertexID.String())
		return err
	}
	err = self.deleteZnode(znodePath)
	if err != nil {
		system.Logf("Failed to delete znode for vertex: %s", vertexID.String())
		return err
	}

	err = self.DecrementElementCount(graphID, *partitionId)
	if err != nil {
		system.Logf("Failed to decrement Element count at partition: %s", (*partitionId).String())
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) DeleteEdge(graphID uuid.UUID, edgeID uuid.UUID) error {
	znodePath := path.Join(ROOT, GRAPH, graphID.String(), EDGE, edgeID.String())
	return self.deleteZnode(znodePath)
}

func (self *ZkMetadataMapper) IncrementElementCount(graphID uuid.UUID, partitionId uuid.UUID) error {
	znodePath := path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionId.String())
	return self.changePartitionInfo(znodePath, 1)
}

func (self *ZkMetadataMapper) DecrementElementCount(graphID uuid.UUID, partitionId uuid.UUID) error {
	znodePath := path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionId.String())
	return self.changePartitionInfo(znodePath, -1)
}

func (self *ZkMetadataMapper) FindNewBackendForPartition(graphID uuid.UUID, partitionId uuid.UUID) (string, error) {
	backends, err := self.GetAllBackends()
	var back string
	if err != nil {
		return "", err
	}
	currBacks, err := self.GetBackendsForPartition(graphID, partitionId)
	if err != nil {
		return "", err
	}
	if len(backends) < common.REPLICATION_FACTOR {
		return "", fmt.Errorf("Not enough  backends")
	}
	for {
		ind := rand.Intn(len(backends) - 1)
		back = backends[ind]
		new := true
		for _, v := range currBacks {
			if back == v {
				new = false
			}
		}
		if new {
			break
		}
	}
	return back, nil
}

var _ Metadata = new(ZkMetadataMapper)
