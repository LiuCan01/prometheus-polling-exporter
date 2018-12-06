package metric

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"os/exec"
)

var f1 int = 6

//this function needs to be add in cron
func HaConfig() {
	var whoami []byte
	var err error
	var cmd *exec.Cmd
	//The command you want to exec
	cmd = exec.Command("pwd")
	if whoami, err = cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	f1 = f1 + 1
	//res := string(whoami)
	fmt.Println(string(whoami))
}

type ClusterManager struct {
	Zone         string
	OOMCountDesc *prometheus.Desc
	// ... many more fields
}

// Simulate prepare the data
func (c *ClusterManager) ReallyExpensiveAssessmentOfTheSystemState() (
	oomCountByHost map[string]int) {

	oomCountByHost = map[string]int{
		//one metric name with different lables
		"foo.example.org": f1,
		"bar.example.org": 2001,
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
func NewClusterManager(zone string) *ClusterManager {
	return &ClusterManager{
		Zone: zone,
		OOMCountDesc: prometheus.NewDesc(
			"clustermanager_oom_crashes_total",
			"Number of OOM crashes.",
			[]string{"host"},
			prometheus.Labels{"zone": zone},
		),
	}
}


