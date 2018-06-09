package metadata

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"errors"

	"github.com/google/uuid"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	versionID     = -1
	typeBACKEND   = "backends"
	typeGRAPH     = "graph"
	typePARTITION = "partitions"
	typeVERTEX    = "vertices"
	typeEDGE      = "edges"
)

type ZkMetadataMapper struct {
	connection *zk.Conn
	err        error
}

func (self *ZkMetadataMapper) CreatePartition(graphID uuid.UUID, partitionID uuid.UUID) error {
	data := map[string]string{"partitionID": partitionID.String()}
	return self.createZnode(path.Join("", graphID.String(), typePARTITION, partitionID.String()), data)
}

func (self *ZkMetadataMapper) CreateVertex(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error {
	data := map[string]string{"partitionID": partitionID.String()}
	return self.createZnode(path.Join("", graphID.String(), typeVERTEX, vertexID.String()), data)
}

func (self *ZkMetadataMapper) CreateEdge(graphID uuid.UUID, edgeID uuid.UUID, srcID uuid.UUID) error {
	data := map[string]string{"srcID": srcID.String()}
	return self.createZnode(path.Join("", graphID.String(), typeEDGE, edgeID.String()), data)
}

func (self *ZkMetadataMapper) CreateBackend(backendAddr string) (string, error) {
	var err error
	var data []byte
	var backendID string
	conn := connect(self.connection, self.err)
	znodePath := path.Join("", typeBACKEND)
	data, err = json.Marshal(backendAddr)
	if err != nil {
		fmt.Printf("Error while Marshalling the backendAddr %s", backendAddr)
		return backendID, err
	}
	backendID, err = conn.CreateProtectedEphemeralSequential(znodePath, data, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Printf("Error while creating %s backend node", backendAddr)
		return backendID, err
	}
	return backendID, nil
}

func (self *ZkMetadataMapper) AddBackendToPartition(graphID uuid.UUID, partitionID uuid.UUID, backendID string) ([]string, <-chan zk.Event, error) {
	var err error
	var liveBackends []string
	watch := make(<-chan zk.Event)
	err = self.createZnode(path.Join("", graphID.String(), typePARTITION, partitionID.String(), backendID), "")
	if err != nil {
		fmt.Printf("Error while adding %s backend to %s partition", backendID, partitionID.String())
		return nil, nil, err
	}
	backendNode := path.Join("", graphID.String(), partitionID.String())

	// pass the full path upto node whose children have to be watched
	liveBackends, watch, err = self.GetWatchOnChildren(backendNode)
	if err != nil {
		fmt.Printf("Error while setting a watch on %s", backendID)
		return nil, watch, err
	}
	return liveBackends, watch, err
}

// @param : backendNode - full path of node whose children need to be watched
// @return : list of all alive children,
func (self *ZkMetadataMapper) GetWatchOnChildren(backendNode string) ([]string, <-chan zk.Event, error) {
	conn := connect(self.connection, self.err)
	snapshot, _, watch, err := conn.ChildrenW(backendNode)
	return snapshot, watch, err
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (self *ZkMetadataMapper) GetVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) ([]string, error) {
	var partitionID string
	var backendIDs, children []string
	var err error

	// get partitionID from graphID and vertexID
	znodePath := path.Join("", typeGRAPH, graphID.String(), typeVERTEX, vertexID.String())
	data, err := self.getZnodeData(znodePath)
	if err != nil {
		fmt.Printf("Error while getting %s vertex", vertexID.String())
	}
	partitionID = data["partitionID"].(string)

	//  get backends from paritionID
	znodePath = path.Join("", typeGRAPH, graphID.String(), "partitions", partitionID)
	data, err = self.getZnodeData(znodePath)
	if err != nil {
		fmt.Printf("Error while getting %s paritionID", partitionID)
		return nil, err
	}

	backendIDs = append([]string{data["Primary"].(string)}, backendIDs...)
	backendIDs = append(backendIDs, data["Secondaries"].([]string)...)

	for _, backendID := range backendIDs {
		znodePath = path.Join("", typeBACKEND, backendID)
		data, err = self.getZnodeData(znodePath)

		if err != nil {
			fmt.Printf("Error while getting backend address of %s backendID", backendID)
			return nil, err
		}
		children = append(children, data["addr"].(string))
	}

	return children, err
}

func (self *ZkMetadataMapper) SetVertexLocation(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error {
	//conn := connect(self.connection, self.err)
	//var buffer bytes.Buffer
	var exists bool
	var err error

	znodePath := path.Join("", typeGRAPH, graphID.String(), typeVERTEX, vertexID.String())
	exists, err = self.checkZnodeExists(znodePath)
	if err != nil {
		fmt.Printf("Error while checking if %s vertex exists", vertexID.String())
		return err
	}
	if exists != true {
		fmt.Printf("%s vertex does not exist", vertexID.String())
		return fmt.Errorf("Vertex does not exist")
	}

	err = self.setZnodeData(znodePath, map[string]string{"partitionID": partitionID.String()})
	if err != nil {
		fmt.Printf("Error while setting %s vertex with %s paritionID", vertexID.String(), partitionID.String())
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) GetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID) ([]string, error) {
	znodePath := path.Join("", typeGRAPH, graphID.String(), typeEDGE, edgeID.String())
	data, err := self.getZnodeData(znodePath)
	if err != nil {
		fmt.Printf("Failed to fetch data for znode: %s", znodePath)
		return make([]string, 0), err
	}

	srcID_str := data["srcID"].(string)
	srcID, err := uuid.FromBytes([]byte(srcID_str))
	if err != nil {
		return make([]string, 0), err
	}

	return self.GetVertexLocation(graphID, srcID)
}
func (self *ZkMetadataMapper) SetEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID, srcID uuid.UUID) error {
	znodePath := path.Join("", typeGRAPH, graphID.String(), typeEDGE, edgeID.String())
	err := self.setZnodeData(znodePath, map[string]string{"srcID": srcID.String()})
	if err != nil {
		fmt.Printf("Failed to set data for znode: %s", znodePath)
		return err
	}

	return nil
}

func connect(connection *zk.Conn, err error) *zk.Conn {
	Addrs := []string{"169.228.66.172:21810", "169.228.66.170:21810", "169.228.66.171:21810"}
	if connection == nil { // If no connection currently exists
		connection, _, err = zk.Connect(Addrs, time.Second)
		must(err)
	} else if _, _, err = connection.Get("/"); err != nil { // If connection exists, but is faulty
		connection.Close()
		connection, _, err = zk.Connect(Addrs, time.Second)
		must(err)
	}
	return connection
}

func (self *ZkMetadataMapper) getZnodeData(znodePath string) (map[string]interface{}, error) {
	//Establish connection of zookeeper
	conn := connect(self.connection, self.err)

	//Fetch and unmarshal data for znode
	data, _, err := conn.Get(znodePath)
	if err != nil {
		fmt.Printf("Error while getting znode at path: %s", znodePath)
		return make(map[string]interface{}), err
	}

	var dat map[string]interface{}
	err = json.Unmarshal(data, &dat)
	if err != nil {
		fmt.Printf("Error while unmarshalling data for znode: %s", znodePath)
		return make(map[string]interface{}), err
	}

	return dat, nil
}

func (self *ZkMetadataMapper) setZnodeData(znodePath string, data interface{}) error {
	//Establish connection of zookeeper
	conn := connect(self.connection, self.err)

	//Marshal and set data for node
	str, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error while marshalling data for znode: %s", znodePath)
		return err
	}

	_, err = conn.Set(znodePath, str, versionID)

	if err != nil {
		fmt.Printf("Error while setting znode: %s with data: %s", znodePath, str)
		return err
	}

	return nil
}

func (self *ZkMetadataMapper) checkZnodeExists(znodePath string) (bool, error) {
	conn := connect(self.connection, self.err)

	exists, _, err := conn.Exists(znodePath)
	if err != nil {
		fmt.Printf("Failed to check if znode %s exists", znodePath)
		return false, err
	}
	return exists, nil
}

func (self *ZkMetadataMapper) createZnode(znodePath string, data interface{}) error {
	// Establish connection to zookeeper
	conn := connect(self.connection, self.err)

	// Set partitionID for element
	//znodePath := path.Join("", typeGRAPH, graphID.String(), znodeType, elementID.String())
	exists, err := self.checkZnodeExists(znodePath)
	if err != nil {
		fmt.Printf("Error trying to check if znode: %s exists", znodePath)
		return err
	}

	if exists == true {
		fmt.Printf("znode %s already exists", znodePath)
		return errors.New("znode already exists")
	}

	strdata, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error while marshalling data for znode: %s", znodePath)
		return err
	}

	_, err = conn.Create(znodePath, strdata, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Printf("Error while creating znode: %s", znodePath)
		return err
	}

	return nil
}

var _ Metadata = new(ZkMetadataMapper)
