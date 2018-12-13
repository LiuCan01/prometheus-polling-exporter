package metric

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"os/exec"
	"strings"
)

var ceph_status_code int = 0
var ceph_status_checks_end string = ""

//this function needs to be add in cron
func Get_ceph_health_metric() {
	src_replace := []string{
		"Degraded data redundancy:",
		"Degraded data redundancy (low space)",
		"Reduced data availability:",
		"backfillfull osd(s)",
		"nearfull osd(s)",
		" full osd(s)",
		"osds down"}

	dst_replacement := []string{
		"Degraded data redundancy",
		"Degraded data redundancy (low space)",
		"Reduced data availability",
		"Exist backfillfull osd(s)",
		"Exist nearfull osd(s)",
		"Exist full osd(s)",
		"Exist osds down"}

	var std_out []byte
	var ceph_health string
	var err error
	var cmd *exec.Cmd
	cmd = exec.Command("ceph","health")
	if std_out, err = cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ceph_health = string(std_out)
	ceph_status := "HEALTH_OK"
	ceph_status_checks_end = ""
	status_index := strings.Index(ceph_health, " ")
        if status_index != -1 {	
		ceph_status = ceph_health[0:status_index] 
	}
	if "HEALTH_WARN" == ceph_status {
		ceph_status_code = 1
	}else if "HEALTH_ERR" == ceph_status {
		ceph_status_code = 2
	}else{
		ceph_status_code = 0
		return 
	}

	ceph_status_checks := ceph_health[status_index+1:]

	sub_str := strings.Split(ceph_status_checks, ";")

	for i:=0; i < len(sub_str); i++{
		flag_tmp := 0
		for j:=0; j < len(src_replace); j++{
			if strings.Contains(sub_str[i], src_replace[j]){
				flag_tmp = 1
				ceph_status_checks_end +=  dst_replacement[j]
				break
			}
		}
		if flag_tmp == 0 {
			ceph_status_checks_end += strings.TrimSpace(sub_str[i])
		}
		if i != len(sub_str)-1 {
			ceph_status_checks_end += "; "
		}
	}
}

type ClusterManager struct {
	OOMCountDesc *prometheus.Desc
	// ... many more fields
}

// Simulate prepare the data
func (c *ClusterManager) ReallyExpensiveAssessmentOfTheSystemState() (
	oomCountByHost map[string]int) {

	oomCountByHost = map[string]int{
		//one metric name with different lables
		ceph_status_checks_end : ceph_status_code,
	}
	return
}

// Describe simply sends the two Descs in the struct to the channel.
func (c *ClusterManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.OOMCountDesc
}

func (c *ClusterManager) Collect(ch chan<- prometheus.Metric) {
	oomCountByHost := c.ReallyExpensiveAssessmentOfTheSystemState()
	for host, oomCount := range oomCountByHost {
		ch <- prometheus.MustNewConstMetric(
			c.OOMCountDesc,
			prometheus.CounterValue,
			float64(oomCount),
			host,
		)
	}
}

//This function need to be register in main.go
// NewClusterManager creates the Descs OOMCountDesc. Note
// that the zone is set as a ConstLabel. (It's different in each instance of the
// ClusterManager, but constant over the lifetime of an instance.) Then there is
// a variable label "host", since we want to partition the collected metrics by
// host. Since all Descs created in this way are consistent across instances,
// with a guaranteed distinction by the "zone" label, we can register different
// ClusterManager instances with the same registry.
func NewClusterManager() *ClusterManager {
	return &ClusterManager{
		OOMCountDesc: prometheus.NewDesc(
			"ecms_ceph_health_status",
			"",
			[]string{"checks"},
			prometheus.Labels{"ceph": ""},
		),
	}
}


