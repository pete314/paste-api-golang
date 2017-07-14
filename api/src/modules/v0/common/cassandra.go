//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description: collection of methods for Cassandra
package common

import (
	"github.com/gocql/gocql"
	"sync"
)

type CassandraCluster struct {
	cluster *gocql.ClusterConfig
}

var (
	ccInstance    *CassandraCluster
	flywInstances map[string]*CassandraCluster
	once          sync.Once
)

//NewCluster - connects to new cluster
func NewCluster(ks string) *CassandraCluster {
	once.Do(func() {
		initBStruct()
	})

	if nil != ccInstance {
		clstr := gocql.NewCluster("127.0.0.1")
		clstr.Keyspace = ks
		clstr.Consistency = gocql.Quorum
		ccInstance = &CassandraCluster{cluster: clstr}
	}

	return ccInstance
}

//NewSession - create new session
func NewSession(ks string) (*gocql.Session, error) {
	v, ok := flywInstances[ks]
	if !ok {
		v = NewCluster(ks)
		flywInstances[ks] = v
	}

	return v.cluster.CreateSession()
}

func initBStruct() {

}
