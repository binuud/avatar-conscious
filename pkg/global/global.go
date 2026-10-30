package global

import "github.com/sirupsen/logrus"

var IsDarwin = false
var IsGpuMetricsAvailable = false
var IsDocker = false

func DumpGlobals() {
	logrus.Infof("IsDarwin : %v", IsDarwin)
	logrus.Infof("IsGpuMetricsAvailable : %v", IsGpuMetricsAvailable)
	logrus.Infof("IsGpuMetricsAvailable : %v", IsGpuMetricsAvailable)
}
