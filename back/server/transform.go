package server

import (
	"bufio"
	"bytes"
	"strings"

	log "github.com/sirupsen/logrus"
)

func DataToKVSlices(d *bytes.Buffer) Series {
	series := Series{}
	scanner := bufio.NewScanner(d)
	scanner.Split(bufio.ScanLines)
	// get header
	scanner.Scan()
	h := scanner.Text()
	header := strings.Split(h, ",")
	log.Debugf("test data header: %s", header)

	for scanner.Scan() {
		text := strings.Split(scanner.Text(), ",")
		for hIdx, h := range header {
			if _, ok := series[h]; !ok {
				series[h] = make([]string, 0)
			}
			series[h] = append(series[h], text[hIdx])
		}
	}
	return series
}
