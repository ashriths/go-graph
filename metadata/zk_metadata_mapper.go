package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	versionID = -1
)

type ZkMetadataMapper struct {
	connection *zk.Conn
	err        error
}

func must(err error) {
	if err != nil {
		panic(err)
	}
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
	data, _, err = conn.Get(znodePath)
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

func setZnodeData(znodePath string, data interface{}) error {
	//Establish connection of zookeeper
	con := connect(self.connection, self.err)
}

func (self *ZkMetadataMapper) createVertexZnode(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error {
	return createElementZnode(graphID, vertexID, partitionID, "vertices")
}

func (self *ZkMetadataMapper) createEdgeZnode(graphID uuid.UUID, edgeID uuid.UUID, partitionID uuid.UUID) error {
	return createElementZnode(graphID, vertexID, partitionID, "edges")
}

func (self *ZkMetadataMapper) createElementZnode(graphID uuid.UUID, elementID uuid.UUID, partitionID uuid.UUID, znodeType string) error {
	// Establish connection to zookeeper
	conn := connect(self.connection, self.err)

	// Set partitionID for element
	znodePath = path.Join("/graph", graphID.String(), znodeType, elementID)
	_, err = conn.Set(znodePath, []byte(partitionID.String()), versionID)

	if err != nil {
		fmt.Printf("Error while setting %s vertex with %s paritionID", vertexID.String(), partitionID.String())
		return err
	}

	return nil
}

func (self *ZkMetadataMapper) getVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) ([]string, error) {

	var buffer bytes.Buffer
	var partitionID string
	var backendIDs, children []string
	var data []byte
	var err error
	var dat map[string]interface{}

	// establish connection to zookeeper ensemble
	conn := connect(self.connection, self.err)

	// get partitionID from graphID and vertexID
	buffer.WriteString("/graph")
	buffer.WriteString("/")
	buffer.WriteString(graphID.String())
	buffer.WriteString("/")
	buffer.WriteString("vertices")
	buffer.WriteString("/")
	buffer.WriteString(vertexID.String())
	data, _, err = conn.Get(buffer.String())
	if err != nil {
		fmt.Printf("Error while getting %s vertex", vertexID.String())
		return nil, err
	}
	partitionID = string(data)

	//  get backends from paritionID
	buffer.Reset()
	buffer.WriteString("/graph/")
	buffer.WriteString(graphID.String())
	buffer.WriteString("/")
	buffer.WriteString("partitions")
	buffer.WriteString("/")
	buffer.WriteString(partitionID)
	data, _, err = conn.Get(buffer.String())
	if err != nil {
		fmt.Printf("Error while getting %s paritionID", partitionID)
		return nil, err
	}
	err = json.Unmarshal(data, &dat)
	if err != nil {
		fmt.Printf("Error while unmarshalling %s paritionID's backend locations", partitionID)
		return nil, err
	}

	backendIDs = append([]string{dat["Primary"].(string)}, backendIDs...)
	backendIDs = append(backendIDs, dat["Secondaries"].([]string)...)

	for _, backendID := range backendIDs {
		buffer.Reset()
		buffer.WriteString("/backends")
		buffer.WriteString("/")
		buffer.WriteString(backendID)
		data, _, err = conn.Get(buffer.String())
		if err != nil {
			fmt.Printf("Error while getting backend address of %s backendID", backendID)
			return nil, err
		}
		children = append(children, string(data))
	}

	return children, err
}

func (self *ZkMetadataMapper) setVertexLocation(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error {
	conn := connect(self.connection, self.err)
	var buffer bytes.Buffer
	var exists bool
	var err error
	buffer.WriteString("/graph")
	buffer.WriteString("/")
	buffer.WriteString(graphID.String())
	buffer.WriteString("/")
	buffer.WriteString("vertices")
	buffer.WriteString("/")
	buffer.WriteString(vertexID.String())
	exists, _, err = conn.Exists(buffer.String())
	if err != nil {
		fmt.Printf("Error while checking if %s vertex exists", vertexID.String())
		return err
	}
	if exists != true {
		fmt.Printf("%s vertex does not exist", vertexID.String())
		return fmt.Errorf("Vertex does not exist")
	}
	_, err = conn.Set(buffer.String(), []byte(partitionID.String()), versionID)
	if err != nil {
		fmt.Printf("Error while setting %s vertex with %s paritionID", vertexID.String(), partitionID.String())
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) getEdgeLocation(graphID uuid.UUID, edgeID uuid.UUID) []string {
	conn := connect(self.connection, self.err)

	znodePath = path.Join("/graph", graphID.String(), "vertices", edgeID)
	var data map[string]string

}
func (self *ZkMetadataMapper) setEdgeLocation(nodeID uuid.UUID, backend string) {
	conn := connect(self.connection, self.err)
}

var _ Metadata = new(ZkMetadataMapper)