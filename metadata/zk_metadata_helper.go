package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ashriths/go-graph/system"
	"github.com/samuel/go-zookeeper/zk"
)

var ZK_TIMEOUT = 30 * time.Second

func (self *ZkMetadataMapper) getChildren(path string) ([]string, error) {
	self.connect()
	children, _, err := self.Connection.Children(path)
	return children, err
}

func (self *ZkMetadataMapper) connect() {
	var err error
	if self.Connection == nil { // If no Connection currently exists
		self.Connection, _, err = zk.Connect(self.ZkAddrs, ZK_TIMEOUT)
		must(err)
	} else if _, _, err = self.Connection.Get("/"); err != nil { // If Connection exists, but is faulty
		self.Connection.Close()
		self.Connection, _, err = zk.Connect(self.ZkAddrs, ZK_TIMEOUT)
		must(err)
	}
}

func (self *ZkMetadataMapper) getZnodeData(znodePath string) (map[string]interface{}, *zk.Stat, error) {
	//Establish Connection of zookeeper
	self.connect()

	//Fetch and unmarshal data for znode
	data, statusInfo, err := self.Connection.Get(znodePath)
	if err != nil {
		system.Logf("Error while getting znode at path: %s", znodePath)
		return make(map[string]interface{}), statusInfo, err
	}

	var dat map[string]interface{}
	err = json.Unmarshal(data, &dat)
	if err != nil {
		system.Logf("Error while unmarshalling data for znode: %s", znodePath)
		return make(map[string]interface{}), statusInfo, err
	}

	return dat, statusInfo, nil
}

func (self *ZkMetadataMapper) setZnodeData(znodePath string, data interface{}, versionID int32) error {
	//Establish Connection of zookeeper
	self.connect()

	//Marshal and set data for node
	str, err := json.Marshal(data)
	if err != nil {
		system.Logf("Error while marshalling data for znode: %s", znodePath)
		return err
	}

	_, err = self.Connection.Set(znodePath, str, versionID)

	if err != nil {
		system.Logf("Error while setting znode: %s with data: %s", znodePath, str)
		return err
	}

	return nil
}

func (self *ZkMetadataMapper) checkZnodeExists(znodePath string) (bool, error) {
	self.connect()
	exists, _, err := self.Connection.Exists(znodePath)
	if err != nil {
		system.Logf("Failed to check if znode %s exists", znodePath)
		return false, err
	}
	return exists, nil
}

func (self *ZkMetadataMapper) createZnode(znodePath string, data interface{}) error {
	// Establish Connection to zookeeper
	self.connect()

	// Set partitionID for element
	//znodePath := path.Join("", typeGRAPH, graphID.String(), znodeType, elementID.String())
	exists, err := self.checkZnodeExists(znodePath)
	if err != nil {
		system.Logf("Error trying to check if znode: %s exists", znodePath)
		return err
	}

	if exists == true {
		system.Logf("znode %s already exists", znodePath)
		return errors.New("znode already exists")
	}

	strdata, err := json.Marshal(data)
	if err != nil {
		system.Logf("Error while marshalling data for znode: %s", znodePath)
		return err
	}

	_, err = self.Connection.Create(znodePath, strdata, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		system.Logf("Error while creating znode: %s", znodePath)
		return err
	}

	return nil
}

// Creates znode if it doesn't exist. Doesn't fail if it already exists
func (self *ZkMetadataMapper) createZnodeIfNotExists(znodePath string, data interface{}) error {
	// Establish Connection to zookeeper
	self.connect()

	exists, err := self.checkZnodeExists(znodePath)
	if err != nil {
		system.Logf("Error trying to check if znode: %s exists", znodePath)
		return err
	}

	if exists == true {
		return nil
	}

	strdata, err := json.Marshal(data)
	if err != nil {
		system.Logf("Error while marshalling data for znode: %s", znodePath)
		return err
	}

	_, err = self.Connection.Create(znodePath, strdata, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		system.Logf("Error while creating znode: %s", znodePath)
		return err
	}

	return nil
}

func (self *ZkMetadataMapper) deleteZnode(znodePath string) error {
	self.connect()
	exists, err := self.checkZnodeExists(znodePath)
	if err != nil {
		system.Logf("Failed to check if vertex exists for znode: %s", znodePath)
		return err
	}

	if !exists {
		system.Logf("Znode does not exists for znode: %s", znodePath)
		return fmt.Errorf("Vertex does not exist")
	}

	err = self.Connection.Delete(znodePath, DEFAULTVERSION)
	if err != nil {
		system.Logf("Failed to delete Znode for path: %s", znodePath)
		return err
	}
	return nil
}

func (self *ZkMetadataMapper) changePartitionInfo(znodePath string, value int32) error {
	data, statusInfo, err := self.getZnodeData(znodePath)
	if err != nil {
		system.Logf("Error while retrieving data stored at path: %s", znodePath)
		return nil
	}
	newValue := data["count"].(int32) + value
	newData := map[string]string{"elementCount": string(newValue)}
	err = self.setZnodeData(znodePath, newData, statusInfo.Version)
	return nil
}
