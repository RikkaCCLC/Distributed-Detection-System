package models

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"open-devops/src/common"
	"strings"
)

var availableResources = map[string]struct{}{
	"resource_host": {},
}

func CheckResources(resource string) bool {
	_, ok := availableResources[resource]
	return ok
}

func ResourceMount(req *common.ResourceMountReq, logger log.Logger) (int64, error) {
	gpas := strings.Split(req.TargetPath, ".")
	g, p, a := gpas[0], gpas[1], gpas[2]

	ids := ""
	for _, id := range req.ResourceIds {
		ids += fmt.Sprintf("%d,", id)
	}
	ids = strings.TrimRight(ids, ",")

	rawSql := fmt.Sprintf(`update %s set stree_group='%s' ,stree_product='%s' ,stree_app='%s' where id in (%s)`,
		req.ResourceType,
		g,
		p,
		a,
		ids,
	)
	level.Info(logger).Log("msg", "ResourceMount.sql.show", "rawSql", rawSql)
	res, err := DB["stree"].Exec(rawSql)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	return rowsAffected, err

}

func ResourceUnMount(req *common.ResourceMountReq, logger log.Logger) (int64, error) {

	ids := ""
	for _, id := range req.ResourceIds {
		ids += fmt.Sprintf("%d,", id)
	}
	ids = strings.TrimRight(ids, ",")

	rawSql := fmt.Sprintf(`update %s set stree_group='' ,stree_product='' ,stree_app='' where id in (%s)`,
		req.ResourceType,
		ids,
	)
	level.Info(logger).Log("msg", "ResourceUnMount.sql.show", "rawSql", rawSql, "g.p.a", req.TargetPath)
	res, err := DB["stree"].Exec(rawSql)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	return rowsAffected, err

}
