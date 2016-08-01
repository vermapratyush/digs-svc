package models

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
