package models

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"open-devops/src/common"
	"strings"
)

func ResourceQuery(resourceType string, matchIds []uint64, logger log.Logger, limit, offset int) (interface{}, error) {
	ids := ""
	for _, id := range matchIds {
		ids += fmt.Sprintf("%d,", id)
	}

	ids = strings.TrimRight(ids, ",")
	inSql := fmt.Sprintf("id in (%s) ", ids)
	level.Info(logger).Log("msg", "ResourceQuery.sql.show", "resourceType", resourceType, "inSql", inSql,
		"limit", limit,
		"offset", offset,
	)

	var (
		res interface{}
		err error
	)
	switch resourceType {
	case common.RESOURCE_HOST:
		res, err = ResourceHostGetManyWithLimit(limit, offset, inSql)
	case common.RESOURCE_RDS:

	}
	return res, err
}

func ResourceHostGetManyWithLimit(limit, offset int, where string, args ...interface{}) ([]ResourceHost, error) {
	var objs []ResourceHost
	err := DB["stree"].Where(where, args...).Limit(limit, offset).Find(&objs)
	if err != nil {
		return nil, err
	}
	return objs, nil

}

func ResourceHostGetMany(where string, args ...interface{}) ([]ResourceHost, error) {
	var objs []ResourceHost
	err := DB["stree"].Where(where, args...).Find(&objs)
	if err != nil {
		return nil, err
	}
	return objs, nil

}
