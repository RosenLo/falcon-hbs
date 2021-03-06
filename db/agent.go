// Copyright 2018 RosenLo

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

/**
 * This code was originally worte by Xiaomi, Inc. modified by RosenLo.
**/

package db

import (
	"fmt"
	"log"

	"github.com/RosenLo/falcon-hbs/g"
	"github.com/open-falcon/falcon-plus/common/model"
)

func UpdateAgent(agentInfo *model.AgentUpdateInfo) {
	sql := ""
	if g.Config().Hosts == "" {
		sql = fmt.Sprintf(
			"insert into host(hostname, ip, agent_version, plugin_version) values ('%s', '%s', '%s', '%s') on duplicate key update hostname='%s', agent_version='%s', plugin_version='%s'",
			agentInfo.ReportRequest.Hostname,
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
			agentInfo.ReportRequest.Hostname,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
		)
	} else {
		// sync, just update
		sql = fmt.Sprintf(
			"update host set hostname='%s', agent_version='%s', plugin_version='%s' where ip='%s'",
			agentInfo.ReportRequest.Hostname,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
			agentInfo.ReportRequest.IP,
		)
	}

	log.Println("exec update sql: ", sql)
	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}

}

func UpdateCMDBGroup(agentInfo *model.AgentUpdateInfo) {
	cmdbGroup := g.Config().CMDBGroup
	var hostId int64 = -1
	var groupId int64 = -1

	err := DB.QueryRow("SELECT id FROM grp WHERE grp_name = ?", cmdbGroup).Scan(&groupId)
	if err != nil {
		log.Println("get group name fail", err)
		return
	}
	if groupId <= 0 {
		log.Println("group not found, gourp: ", cmdbGroup)
		return
	}

	err = DB.QueryRow("SELECT id FROM host WHERE hostname = ?", agentInfo.ReportRequest.Hostname).Scan(&hostId)
	if err != nil {
		log.Println("get group name fail", err)
		return
	}
	if groupId <= 0 {
		log.Println("host not found, host: ", cmdbGroup)
		return
	}

	var id int64 = -1
	sql := fmt.Sprintf("SELECT grp_id FROM grp_host WHERE grp_id = %d AND host_id = %d", groupId, hostId)
	err = DB.QueryRow(sql).Scan(&id)
	if err != nil {
		log.Println("get host group empty, sql", sql)
	}
	if id != -1 {
		log.Println("host exist in cmdb group")
		return
	}

	sql = fmt.Sprintf("INSERT INTO grp_host(grp_id, host_id) VALUES (%d, %d)", groupId, hostId)
	_, err = DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}
}
