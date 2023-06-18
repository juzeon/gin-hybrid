package etclient

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"strings"
)

func GetKeyLastSegment(fullKey string) string {
	arr := strings.Split(fullKey, "/")
	return arr[len(arr)-1]
}
func ConvertStringToLeaseID(str string) clientv3.LeaseID {
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return clientv3.LeaseID(i)
}
