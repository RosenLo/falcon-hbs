// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"fmt"
	"log"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
)

func UpdateAgent(agentInfo *model.AgentUpdateInfo) {
	sql := ""
	if g.Config().Hosts == "" {
		sql = fmt.Sprintf(
			"insert into host(hostname, ip, agent_version, plugin_version) values ('%s', '%s', '%s', '%s') on duplicate key update ip='%s', agent_version='%s', plugin_version='%s'",
			agentInfo.ReportRequest.Hostname,
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
		)
	} else {
		// sync, just update
		sql = fmt.Sprintf(
			"update host set ip='%s', agent_version='%s', plugin_version='%s' where hostname='%s'",
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
			agentInfo.ReportRequest.Hostname,
		)
	}

	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}

}

func UpdateCMDBGroup(agentInfo *model.AgentUpdateInfo) {
	cmdb_group := g.Config().CMDBGroup
	var host_id int64 = -1
	var group_id int64 = -1

	err := DB.QueryRow("SELECT grp_id FROM grp WHERE grp_name = ?", cmdb_group).Scan(&group_id)
	if err != nil {
		log.Println("get group name fail", err)
		return
	}
	if group_id <= 0 {
		log.Println("group not found, gourp: ", cmdb_group)
		return
	}

	err = DB.QueryRow("SELECT id FROM host WHERE hostname = ?", agentInfo.ReportRequest.Hostname).Scan(&host_id)
	if err != nil {
		log.Println("get group name fail", err)
		return
	}
	if group_id <= 0 {
		log.Println("host not found, host: ", cmdb_group)
		return
	}

	var id int64 = -1
	sql := fmt.Sprintf("SELECT grp_id FROM grp_host WHERE grp_id = %d AND host_id = %d", group_id, host_id)
	err = DB.QueryRow(sql).Scan(&id)
	if err != nil {
		log.Println("get host group empty, sql", sql)
	}
	if id != -1 {
		return
	}

	sql = fmt.Sprintf("INSERT INTO grp_host(grp_id, host_id) VALUES (%d, %d)", group_id, host_id)
	_, err = DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}
}
