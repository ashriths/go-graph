package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/ashriths/go-graph/system"
	"path"

	"github.com/google/uuid"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

//Constants to be used
const (
	versionID = -1
	ROOT      = "/"
	BACKEND   = "backends"
	GRAPH     = "graphs"
	PARTITION = "partitions"
	VERTEX    = "vertices"
	EDGE      = "edges"

	BACKEND_PREFIX = "back-"
	EMPTY_DATA     = "{}"
)

type ZkMetadataMapper struct {
	Connection *zk.Conn
	ZkAddrs    []string
	err        error
	Watches    map[string]<-chan zk.Event
}

func NewZkMetadataMapper(ZkAddrs []string) *ZkMetadataMapper {
	zkMapper := ZkMetadataMapper{
		ZkAddrs: ZkAddrs,
		Watches: make(map[string]<-chan zk.Event),
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

func (self *ZkMetadataMapper) CreatePartition(graphID uuid.UUID, partitionID uuid.UUID) error {
	data := map[string]string{"partitionID": partitionID.String()}
	return self.createZnode(path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionID.String()), data)
}

func (self *ZkMetadataMapper) CreateVertex(graphID uuid.UUID, partitionID uuid.UUID, vertexID uuid.UUID) error {
	data := map[string]string{"partitionID": partitionID.String()}
	return self.createZnode(path.Join(ROOT, GRAPH, graphID.String(), VERTEX, vertexID.String()), data)
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

func (self *ZkMetadataMapper) GetPartitionInformation(graphID uuid.UUID, partitionID uuid.UUID) (map[string]interface{}, error) {
	data, err := self.getZnodeData(path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionID.String()))
	if err != nil {
		system.Logf("Error while retrieving data stored at %s partition", partitionID.String())
		return nil, err
	}
	return data, nil
}

func (self *ZkMetadataMapper) SetPartitionInformation(graphID uuid.UUID, partitionID uuid.UUID, data interface{}) error {
	var exists bool
	var err error

	znodePath := path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionID.String())
	exists, err = self.checkZnodeExists(znodePath)
	if err != nil {
		system.Logf("Error while checking if %s partition exists", partitionID.String())
		return err
	}
	if exists != true {
		system.Logf("%s partition does not exist", partitionID.String())
		return fmt.Errorf("Partition does not exist")
	}

	err = self.setZnodeData(znodePath, data)
	if err != nil {
		system.Logf("Error while setting data at %s partition", partitionID.String())
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) GetBackendInformation(backendID string) (map[string]interface{}, error) {
	data, err := self.getZnodeData(path.Join(ROOT, BACKEND, backendID))
	if err != nil {
		system.Logf("Error while retrieving data stored at %s backend", backendID)
		return nil, err
	}
	return data, nil
}

// GetBackendsForPartition : Add backends to Partitions
func (self *ZkMetadataMapper) GetBackendsForPartition(graphID uuid.UUID, partitionID uuid.UUID) ([]string, error) {
	children, err := self.getChildren(path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionID.String()))
	if err != nil {
		system.Logf("Error while getting %s paritionID", partitionID)
		return nil, err
	}
	return children, nil
}

func (self *ZkMetadataMapper) AddBackendToPartition(graphID uuid.UUID, partitionID uuid.UUID, backendID string) ([]string, error) {
	var err error
	var liveBackends []string
	watch := make(<-chan zk.Event)
	err = self.createZnode(path.Join(ROOT, GRAPH, graphID.String(), PARTITION, partitionID.String(), backendID), EMPTY_DATA)
	if err != nil {
		system.Logf("Error while adding %s backend to %s partition", backendID, partitionID.String())
		return nil, err
	}
	backendNode := path.Join(ROOT, GRAPH, graphID.String(), partitionID.String())

	// pass the full path upto node whose children have to be watched
	liveBackends, watch, err = self.GetWatchOnChildren(backendNode)
	if err != nil {
		system.Logf("Error while setting a watch on %s", backendID)
		return nil, err
	}
	self.Watches[partitionID.String()] = watch
	go self.startWatchingPartition(partitionID, watch)
	return liveBackends, err
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
	var partitionID string
	var children []string
	var err error

	// get partitionID from graphID and vertexID
	znodePath := path.Join(ROOT, GRAPH, graphID.String(), VERTEX, vertexID.String())
	data, err := self.getZnodeData(znodePath)
	if err != nil {
		system.Logf("Error while getting %s vertex", vertexID.String())
		return nil, nil, err
	}
	partitionID = data["partitionID"].(string)

	//  get backends from paritionID
	znodePath = path.Join(ROOT, GRAPH, graphID.String(), "partitions", partitionID)
	children, err = self.getChildren(znodePath)
	if err != nil {
		system.Logf("Error while getting %s paritionID", partitionID)
		return nil, nil, err
	}

	partitionUUID, err := uuid.Parse(partitionID)
	if err != nil {
		system.Logf("Error while parsing %s paritionID", partitionID)
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

func (self *ZkMetadataMapper) SetVertexLocation(graphID uuid.UUID, partitionID uuid.UUID, vertexID uuid.UUID) error {
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

	err = self.setZnodeData(znodePath, map[string]string{"partitionID": partitionID.String()})
	if err != nil {
		system.Logf("Error while setting %s vertex with %s paritionID", vertexID.String(), partitionID.String())
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) GetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID) (*uuid.UUID, []string, error) {
	znodePath := path.Join(ROOT, GRAPH, graphID.String(), EDGE, edgeID.String())
	data, err := self.getZnodeData(znodePath)
	if err != nil {
		system.Logf("Failed to fetch data for znode: %s", znodePath)
		return nil, make([]string, 0), err
	}

	srcID_str := data["srcID"].(string)
	srcID, err := uuid.FromBytes([]byte(srcID_str))
	if err != nil {
		return nil, make([]string, 0), err
	}

	return self.GetVertexLocation(graphID, srcID)
}
func (self *ZkMetadataMapper) SetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID, srcID uuid.UUID) error {
	znodePath := path.Join(ROOT, GRAPH, graphID.String(), EDGE, edgeID.String())
	err := self.setZnodeData(znodePath, map[string]string{"srcID": srcID.String()})
	if err != nil {
		system.Logf("Failed to set data for znode: %s", znodePath)
		return err
	}

	return nil
}

func (self *ZkMetadataMapper) startWatchingPartition(partitionID uuid.UUID, watch <-chan zk.Event) {
	for {
		evt := <-watch
		system.Logln("Watch fired for ", partitionID.String(), evt.Err)
		//TODO: Call replication here
	}
}

var _ Metadata = new(ZkMetadataMapper)
