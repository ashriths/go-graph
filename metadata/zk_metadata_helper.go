package metadata

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/ashriths/go-graph/system"
)

var ZK_TIMEOUT = 2 * time.Minute

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

func (self *ZkMetadataMapper) getZnodeData(znodePath string) (map[string]interface{}, error) {
	//Establish Connection of zookeeper
	self.connect()

	//Fetch and unmarshal data for znode
	data, _, err := self.Connection.Get(znodePath)
	if err != nil {
		system.Logf("Error while getting znode at path: %s", znodePath)
		return make(map[string]interface{}), err
	}

	var dat map[string]interface{}
	err = json.Unmarshal(data, &dat)
	if err != nil {
		system.Logf("Error while unmarshalling data for znode: %s", znodePath)
		return make(map[string]interface{}), err
	}

	return dat, nil
}

func (self *ZkMetadataMapper) setZnodeData(znodePath string, data interface{}) error {
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
