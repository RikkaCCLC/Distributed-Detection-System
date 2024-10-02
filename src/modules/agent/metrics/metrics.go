package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"open-devops/src/models"
)

func CreateMetrics(ss []*models.LogStrategy) map[string]*prometheus.GaugeVec {
	mmap := map[string]*prometheus.GaugeVec{}
	for _, s := range ss {
		labels := []string{}
		for k := range s.Tags {
			labels = append(labels, k)
		}
		m := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: s.MetricName,
			Help: s.MetricHelp,
		}, labels)
		mmap[s.MetricName] = m

	}
	return mmap
}

func StartMetricWeb(addr string) error {

	http.Handle("/metrics", promhttp.Handler())
	srv := http.Server{Addr: addr}
	err := srv.ListenAndServe()
	return err

}
