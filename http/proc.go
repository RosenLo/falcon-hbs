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

package http

import (
	"fmt"
	"net/http"

	"github.com/RosenLo/falcon-hbs/cache"
	"github.com/open-falcon/falcon-plus/common/model"
)

func configProcRoutes() {
	http.HandleFunc("/expressions", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, cache.ExpressionCache.Get())
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		m := make(map[string]interface{})
		cache.Agents.RLock()

		m["data"] = cache.Agents.M
		m["total"] = len(cache.Agents.M)
		RenderDataJson(w, m)
		cache.Agents.RUnlock()
	})

	http.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*model.Host, len(cache.MonitoredHosts.Get()))
		for k, v := range cache.MonitoredHosts.Get() {
			data[fmt.Sprint(k)] = v
		}
		m := make(map[string]interface{})
		m["data"] = data
		m["total"] = len(data)
		RenderDataJson(w, m)
	})

	http.HandleFunc("/strategies", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*model.Strategy, len(cache.Strategies.GetMap()))
		for k, v := range cache.Strategies.GetMap() {
			data[fmt.Sprint(k)] = v
		}
		RenderDataJson(w, data)
	})

	http.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*model.Template, len(cache.TemplateCache.GetMap()))
		for k, v := range cache.TemplateCache.GetMap() {
			data[fmt.Sprint(k)] = v
		}
		RenderDataJson(w, data)
	})

	http.HandleFunc("/plugins/", func(w http.ResponseWriter, r *http.Request) {
		hostname := r.URL.Path[len("/plugins/"):]
		RenderDataJson(w, cache.GetPlugins(hostname))
	})

}
