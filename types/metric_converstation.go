package qtypes

import (
	"sort"
	"fmt"
	"strings"
)


func (m *Metric) ConvertToOpenTSDBHandler(targetForm string) (datapoint string) {
	//orders dimensions so datapoint keeps consistent name
	var keys []string
	dimensions := m.GetDimensions(make(map[string]string))
	for k := range dimensions {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	datapoint = fmt.Sprintf("put %s", m.Name)
	var dims []string
	for _, key := range keys {
		dims = append(dims, fmt.Sprintf("%s=%s", key, dimensions[key]))
	}
	if len(dimensions) == 0 {
		datapoint = fmt.Sprintf("%s %d %f\n", datapoint, m.GetTime().Unix(), m.Value)
	} else {
		switch targetForm {
		case "influxdb":
			datapoint = fmt.Sprintf("%s %d %f %s\n", datapoint, m.GetTime().Unix(), m.Value, strings.Join(dims[:], " "))
		default:
			datapoint = fmt.Sprintf("%s %s %d %f\n", datapoint, strings.Join(dims[:], ","), m.GetTime().Unix(), m.Value)
		}

	}
	return datapoint
}
