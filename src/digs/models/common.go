package models

import (
	"github.com/deckarep/golang-set"
	"github.com/astaxie/beego"
	"digs/logger"
)

func (this *UserAccount) OneToOneGroupId() []string {
	gids := make([]string, 0)
	for gid, visibility := range(this.GroupMember) {
		if visibility == 0 {
			gids = append(gids, gid)
		}
	}
	return gids
}

func (this *UserAccount) MultiGroupId() []string {
	gids := make([]string, 0)
	for gid, visibility := range(this.GroupMember) {
		if visibility == 1 {
			gids = append(gids, gid)
		}
	}
	return gids
}

func GetOneToOneCommonId(uid1, uid2 *UserAccount) string {
	gid1Set := mapset.NewSet()
	for _, gid := range(uid1.OneToOneGroupId()) {
		gid1Set.Add(gid)
	}
	gid2Set := mapset.NewSet()
	for _, gid := range(uid2.OneToOneGroupId()) {
		gid2Set.Add(gid)
	}
	commonGid := gid1Set.Intersect(gid2Set).ToSlice()
	if len(commonGid) != 1 {
		if len(commonGid) > 1 {
			logger.Critical("DataCorrupt|UID1=", uid1.UID, "|UID2=", uid2.UID, "|Duplicate 1-1 GID")
		}
		return ""
	} else {
		return commonGid[0].(string)
	}
}