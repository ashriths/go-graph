package metadata

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"fmt"
	"encoding/json"
	"errors"
)

func (self *ZkMetadataMapper) getChildren(path string) ([]string, error) {
	conn := self.connect(self.connection, self.err)
	children, _, err := conn.Children(path)
	return children, err
}

func (self *ZkMetadataMapper) connect(connection *zk.Conn, err error) *zk.Conn {
	//Addrs := []string{"169.228.66.172:21810", "169.228.66.170:21810", "169.228.66.171:21810"}
	Addrs := self.ZkAddrs
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
	conn := self.connect(self.connection, self.err)

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
	conn := self.connect(self.connection, self.err)

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
	conn := self.connect(self.connection, self.err)

	exists, _, err := conn.Exists(znodePath)
	if err != nil {
		fmt.Printf("Failed to check if znode %s exists", znodePath)
		return false, err
	}
	return exists, nil
}

func (self *ZkMetadataMapper) createZnode(znodePath string, data interface{}) error {
	// Establish connection to zookeeper
	conn := self.connect(self.connection, self.err)

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
