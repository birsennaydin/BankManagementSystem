package database

import (
	"fmt"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitCassandra() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "bank_app"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("Cassandra connection error: %v", err))
	}

	fmt.Println("✅ Cassandra connected")
	Session = session
}
