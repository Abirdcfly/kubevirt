package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	libvirt "libvirt.org/libvirt-go"

	k6tv1 "kubevirt.io/client-go/api/v1"
	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/stats"
	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/statsconv"
	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/statsconv/util"
)

type fakeCollector struct {
}

func (fc fakeCollector) Describe(_ chan<- *prometheus.Desc) {
}

//Collect needs to report all metrics to see it in docs
func (fc fakeCollector) Collect(ch chan<- prometheus.Metric) {
	ps := prometheusScraper{ch: ch}

	libstatst, err := util.LoadStats()
	if err != nil {
		panic(err)
	}

	in := &libstatst[0]
	inMem := []libvirt.DomainMemoryStat{}
	out := stats.DomainStats{}
	ident := statsconv.DomainIdentifier(&fakeIdentifier{})

	if err = statsconv.Convert_libvirt_DomainStats_to_stats_DomainStats(ident, in, inMem, &out); err != nil {
		panic(err)
	}

	out.Memory.ActualBalloonSet = true
	out.Memory.UnusedSet = true
	out.Memory.AvailableSet = true
	out.Memory.RSSSet = true
	out.Memory.SwapInSet = true
	out.CPUMapSet = true

	vmi := k6tv1.VirtualMachineInstance{
		Status: k6tv1.VirtualMachineInstanceStatus{
			Phase: k6tv1.Running,
		},
	}
	ps.Report("test", &vmi, &out)
	updateVMIsPhase("test", []*k6tv1.VirtualMachineInstance{&vmi}, ch)
}

type fakeIdentifier struct {
}

func (*fakeIdentifier) GetName() (string, error) {
	return "test", nil
}

func (*fakeIdentifier) GetUUIDString() (string, error) {
	return "uuid", nil
}

func RegisterFakeCollector() {
	prometheus.MustRegister(fakeCollector{})
}
