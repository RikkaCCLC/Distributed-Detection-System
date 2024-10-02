package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"open-devops/src/common"
)

var (
	// 耗时统计的
	// 刷索引
	IndexFlushDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_index_flush_last_duration_seconds",
		Help: "Duration of index flush  ",
	}, []string{common.LABEL_RESOURCE_TYPE})

	// 公有云同步
	PublicCloudSyncDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "public_cloud_sync_last_duration_seconds",
		Help: "Duration of public cloud sync",
	}, []string{common.LABEL_RESOURCE_TYPE})

	// 查stree_path表耗时
	GetGPAFromDbDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "get_gpa_from_db_last_duration_seconds",
		Help: "get_gpa_from_db_last_duration_seconds",
	})
	// 统计耗时
	ResourceLastStatisticsDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resource_last_statistics_last_duration_seconds",
		Help: "resource_last_statistics_last_duration_seconds",
	})

	// 全局资源统计
	ResourceNumCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_num_count",
		Help: "Num of resource",
	}, []string{common.LABEL_RESOURCE_TYPE})

	ResourceNumRegionCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_num_region_count",
		Help: "Num of resource with region tag",
	}, []string{common.LABEL_RESOURCE_TYPE, common.LABEL_REGION})

	ResourceNumCloudProviderCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_num_cloud_provider_count",
		Help: "Num of resource with cloud_provider tag",
	}, []string{common.LABEL_RESOURCE_TYPE, common.LABEL_CLOUD_PROVIDER})

	ResourceNumClusterCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_num_cluster_count",
		Help: "Num of resource with cluster tag",
	}, []string{common.LABEL_RESOURCE_TYPE, common.LABEL_CLUSTER})

	// gpa 通用资源统计
	GPACount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gpa_count",
		Help: "Num gpas",
	})

	GPAAllNumCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_all_num_count",
		Help: "Num gpa of all",
	}, []string{common.LABEL_GPA_NAME, common.LABEL_RESOURCE_TYPE})

	GPAAllNumRegionCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_all_region_num_count",
		Help: "Num gpa of all with tag region",
	}, []string{common.LABEL_GPA_NAME, common.LABEL_RESOURCE_TYPE, common.LABEL_REGION})
	GPAAllNumCloudProviderCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_all_cloud_provider_num_count",
		Help: "Num gpa of all with tag cloud_provider",
	}, []string{common.LABEL_GPA_NAME, common.LABEL_RESOURCE_TYPE, common.LABEL_CLOUD_PROVIDER})
	GPAAllNumClusterCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_all_cluster_num_count",
		Help: "Num gpa of all with tag cluster",
	}, []string{common.LABEL_GPA_NAME, common.LABEL_RESOURCE_TYPE, common.LABEL_CLUSTER})
	GPAAllNumInstanceTypeCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_all_instance_type_num_count",
		Help: "Num gpa of all with tag instance_type",
	}, []string{common.LABEL_GPA_NAME, common.LABEL_RESOURCE_TYPE, common.LABEL_INSTANCE_TYPE})

	// host特殊的
	GPAHostCpuCores = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_host_cpu_cores",
		Help: "Num gpa cpu cores of ecs",
	}, []string{common.LABEL_GPA_NAME})

	GPAHostMemGbs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_host_mem_gbs",
		Help: "Num gpa mem gbs of ecs",
	}, []string{common.LABEL_GPA_NAME})
	GPAHostDiskGbs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_host_disk_gbs",
		Help: "Num gpa disk gbs of ecs",
	}, []string{common.LABEL_GPA_NAME})
)

func NewMetrics() {
	// last耗时
	prometheus.DefaultRegisterer.MustRegister(IndexFlushDuration)
	prometheus.DefaultRegisterer.MustRegister(PublicCloudSyncDuration)
	prometheus.DefaultRegisterer.MustRegister(GetGPAFromDbDuration)
	prometheus.DefaultRegisterer.MustRegister(ResourceLastStatisticsDuration)
	// 全局资源统计
	prometheus.DefaultRegisterer.MustRegister(ResourceNumCount)
	prometheus.DefaultRegisterer.MustRegister(ResourceNumRegionCount)
	prometheus.DefaultRegisterer.MustRegister(ResourceNumCloudProviderCount)
	prometheus.DefaultRegisterer.MustRegister(ResourceNumClusterCount)

	// gpa 通用资源统计
	prometheus.DefaultRegisterer.MustRegister(GPACount)
	prometheus.DefaultRegisterer.MustRegister(GPAAllNumCount)
	prometheus.DefaultRegisterer.MustRegister(GPAAllNumRegionCount)
	prometheus.DefaultRegisterer.MustRegister(GPAAllNumCloudProviderCount)
	prometheus.DefaultRegisterer.MustRegister(GPAAllNumClusterCount)
	prometheus.DefaultRegisterer.MustRegister(GPAAllNumInstanceTypeCount)

	// host 特殊
	prometheus.DefaultRegisterer.MustRegister(GPAHostCpuCores)
	prometheus.DefaultRegisterer.MustRegister(GPAHostMemGbs)
	prometheus.DefaultRegisterer.MustRegister(GPAHostDiskGbs)

}
