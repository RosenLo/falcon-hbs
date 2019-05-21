package cmdb

import (
	"fmt"

	"log"

	"github.com/RosenLo/falcon-hbs/g"
	"github.com/RosenLo/toolkits/common"
	"github.com/RosenLo/toolkits/http/requests"
)

type Service struct {
	endpoint string
}

func (s *Service) ReportStatus(hostInfo map[string]interface{}) {
	url := fmt.Sprintf("%s/api/v3/host/add/agent", s.endpoint)
	headers := map[string]string{
		"HTTP_BLUEKING_SUPPLIER_ID": "0",
	}
	body := map[string]interface{}{
		"host_info": hostInfo,
	}
	data, err := requests.Call("POST", url, headers, nil, body, nil)
	if err != nil {
		log.Println(err)
	}
	var ret interface{}
	if err := common.ToJSON(data, &ret); err != nil {
		log.Println(err)
	}
}

func NewService() *Service {
	return &Service{endpoint: g.Config().CMDB.Url}
}

func ReportStatus(body map[string]interface{}) {
	cmdb := NewService()
	cmdb.ReportStatus(body)
}
