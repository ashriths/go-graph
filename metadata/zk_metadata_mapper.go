package metadata

import (
	"encoding/json"
	"bytes"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	versionID = 0
)

type ZkMetadataMapper struct {
	connection *zk.Conn
	err        error
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}P

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

func (self *ZkMetadataMapper) getVertexLocation(graphID uuid.UUID, vertexID uuid.UUID) (error, []string) {
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
	buffer.WriteString(graphID.String())
	buffer.WriteString("vertices")
	buffer.WriteString(vertexID.String())
	data, _, err = conn.Get(buffer.String())
	if err != nil {
		fmt.Printf("Error while getting %s vertex", vertexID.String())
		return err, nil
	}
	partitionID = string(data)

	//  get backends from paritionID
	buffer.Reset()
	buffer.WriteString("/graph")
	buffer.WriteString(graphID.String())
	buffer.WriteString("partitions")
	buffer.WriteString(partitionID)
	data, _, err = conn.Get(buffer.String())
	if err != nil {
		fmt.Printf("Error while getting %s paritionID", partitionID)
		return err, nil
	}
	err = json.Unmarshal(data, &dat)
	if err != nil {
		fmt.Printf("Error while unmarshalling %s paritionID's backend locations", partitionID)
		return err, nil
	}

	backendIDs = append([]string{dat["Primary"]},dat["Secondaries"])
	
	for _,backend_id := range backendIDs {
		buffer.Reset()
		buffer.WriteString("/backends")
		buffer.WriteString("/")
		buffer.WriteString(backend_id)
		data, _, err = conn.Get(buffer.String())
		if err != nil {
			fmt.Printf("Error while getting backend address of %s backendID", backend_id)
			return err, nil
		}
		children = append(children, string(data))
	}

	return nil, children
}

func (self *ZkMetadataMapper) setVertexLocation(graphID uuid.UUID, vertexID uuid.UUID, partitionID uuid.UUID) error {
	conn := connect(self.connection, self.err)
	var buffer bytes.Buffer
	var exists bool
	var err error
	buffer.WriteString("/graph")
	buffer.WriteString(graphID.String())
	buffer.WriteString("vertices")
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

func (self *ZkMetadataMapper) getEdgeLocation(nodeID uuid.UUID, backend string) []string {
	conn := connect(self.connection, self.err)
}
func (self *ZkMetadataMapper) setEdgeLocation(nodeID uuid.UUID, backend string) {
	conn := connect(self.connection, self.err)
}

var _ Metadata = new(ZkMetadataMapper)
