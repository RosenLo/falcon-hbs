// Copyright 2018-2019 RosenLo

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

package cache

// 每个agent心跳上来的时候立马更新一下数据库是没必要的
// 缓存起来，每隔一个小时写一次DB
// 提供http接口查询机器信息，排查重名机器的时候比较有用

import (
	"log"
	"sync"
	"time"

	"github.com/RosenLo/falcon-hbs/db"
	"github.com/RosenLo/falcon-hbs/g"
	"github.com/RosenLo/falcon-hbs/util/cmdb"
	"github.com/open-falcon/falcon-plus/common/model"
)

type SafeAgents struct {
	sync.RWMutex
	M map[string]*model.AgentUpdateInfo
}

var Agents = NewSafeAgents()

func NewSafeAgents() *SafeAgents {
	return &SafeAgents{M: make(map[string]*model.AgentUpdateInfo)}
}

func (this *SafeAgents) Put(req *model.AgentReportRequest) {
	val := &model.AgentUpdateInfo{
		LastUpdate:    time.Now().Unix(),
		ReportRequest: req,
	}

	if agentInfo, exists := this.Get(req.IP); !exists ||
		agentInfo.ReportRequest.AgentVersion != req.AgentVersion ||
		agentInfo.ReportRequest.Hostname != req.Hostname ||
		agentInfo.ReportRequest.PluginVersion != req.PluginVersion {

		go db.UpdateAgent(val)
		go db.UpdateCMDBGroup(val)
		if val.ReportRequest.HostInfo != nil {
			go cmdb.ReportStatus(val.ReportRequest.HostInfo)
		}
	}
	this.Lock()
	defer this.Unlock()
	this.M[req.IP] = val
}

func (this *SafeAgents) Get(ip string) (*model.AgentUpdateInfo, bool) {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[ip]
	return val, exists
}

func (this *SafeAgents) Delete(ip string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, ip)
}

func (this *SafeAgents) Keys() []string {
	this.RLock()
	defer this.RUnlock()
	count := len(this.M)
	keys := make([]string, count)
	i := 0
	for ip := range this.M {
		keys[i] = ip
		i++
	}
	return keys
}

func (this *SafeAgents) GetMap() map[string]*model.AgentUpdateInfo {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func DeleteStaleAgents() {
	duration := time.Minute * time.Duration(1)
	for {
		time.Sleep(duration)
		deleteStaleAgents()
	}
}

func deleteStaleAgents() {
	// 十分钟都没有心跳的Agent，从内存中干掉
	before := time.Now().Unix() - g.Config().Interval
	keys := Agents.Keys()
	count := len(keys)
	if count == 0 {
		return
	}

	for i := 0; i < count; i++ {
		curr, _ := Agents.Get(keys[i])
		if curr.LastUpdate < before {
			curr.ReportRequest.HostInfo["online"] = false
			go cmdb.ReportStatus(curr.ReportRequest.HostInfo)
			Agents.Delete(curr.ReportRequest.IP)
			log.Println("delete the host from cache, host: ", curr.ReportRequest.IP)
		}
	}
}
