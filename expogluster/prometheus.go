package expogluster

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

const (
	namespace  = "gluster"
	allVolumes = "_all"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of Gluster successful.",
		nil, nil,
	)

	volumesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volumes_available"),
		"How many volumes were up at the last query.",
		nil, nil,
	)

	volumeStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_status"),
		"Status code of requested volume.",
		[]string{"volume"}, nil,
	)

	nodeSizeFreeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "node_size_bytes_bytes"),
		"Free bytes reported for each node on each instance. Labels are to distinguish origins",
		[]string{"hostname", "path", "volume"}, nil,
	)

	nodeSizeTotalBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "node_size_bytes_total"),
		"Total bytes reported for each node on each instance. Labels are to distinguish origins",
		[]string{"hostname", "path", "volume"}, nil,
	)

	nodeInodesTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "node_inodes_total"),
		"Total inodes reported for each node on each instance. Labels are to distinguish origins",
		[]string{"hostname", "path", "volume"}, nil,
	)

	nodeInodesFree = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "node_inodes_free"),
		"Free inodes reported for each node on each instance. Labels are to distinguish origins",
		[]string{"hostname", "path", "volume"}, nil,
	)

	brickCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_available"),
		"Number of bricks available at last query.",
		[]string{"volume"}, nil,
	)

	brickDuration = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_duration_seconds_total"),
		"Time running volume brick in seconds.",
		[]string{"volume", "brick"}, nil,
	)

	brickDataRead = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_data_read_bytes_total"),
		"Total amount of bytes of data read by brick.",
		[]string{"volume", "brick"}, nil,
	)

	brickDataWritten = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_data_written_bytes_total"),
		"Total amount of bytes of data written by brick.",
		[]string{"volume", "brick"}, nil,
	)

	brickFopHits = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_fop_hits_total"),
		"Total amount of file operation hits.",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	brickFopLatencyAvg = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_fop_latency_avg"),
		"Average fileoperations latency over total uptime",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	brickFopLatencyMin = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_fop_latency_min"),
		"Minimum fileoperations latency over total uptime",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	brickFopLatencyMax = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "brick_fop_latency_max"),
		"Maximum fileoperations latency over total uptime",
		[]string{"volume", "brick", "fop_name"}, nil,
	)

	peersConnected = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peers_connected"),
		"Is peer connected to gluster cluster.",
		nil, nil,
	)

	healInfoFilesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "heal_info_files_count"),
		"File count of files out of sync, when calling 'gluster v heal VOLNAME info",
		[]string{"volume"}, nil)

	volumeWriteable = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_writeable"),
		"Writes and deletes file in Volume and checks if it is writeable",
		[]string{"volume", "mountpoint"}, nil)

	mountSuccessful = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mount_successful"),
		"Checks if mountpoint exists, returns a bool value 0 or 1",
		[]string{"volume", "mountpoint"}, nil)

	quotaHardLimit = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_hardlimit"),
		"Quota hard limit (bytes) in a volume",
		[]string{"path", "volume"}, nil)

	quotaSoftLimit = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_softlimit"),
		"Quota soft limit (bytes) in a volume",
		[]string{"path", "volume"}, nil)

	quotaUsed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_used"),
		"Current data (bytes) used in a quota",
		[]string{"path", "volume"}, nil)

	quotaAvailable = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_available"),
		"Current data (bytes) available in a quota",
		[]string{"path", "volume"}, nil)

	quotaSoftLimitExceeded = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_softlimit_exceeded"),
		"Is the quota soft-limit exceeded",
		[]string{"path", "volume"}, nil)

	quotaHardLimitExceeded = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volume_quota_hardlimit_exceeded"),
		"Is the quota hard-limit exceeded",
		[]string{"path", "volume"}, nil)
)

// Describe all the metrics exported by Gluster exporter. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- volumeStatus
	ch <- volumesCount
	ch <- brickCount
	ch <- brickDuration
	ch <- brickDataRead
	ch <- brickDataWritten
	ch <- peersConnected
	ch <- nodeSizeFreeBytes
	ch <- nodeSizeTotalBytes
	ch <- brickFopHits
	ch <- brickFopLatencyAvg
	ch <- brickFopLatencyMin
	ch <- brickFopLatencyMax
	ch <- healInfoFilesCount
	ch <- volumeWriteable
	ch <- mountSuccessful
	ch <- quotaHardLimit
	ch <- quotaSoftLimit
	ch <- quotaUsed
	ch <- quotaAvailable
	ch <- quotaSoftLimitExceeded
	ch <- quotaHardLimitExceeded
}

// Collect collects all the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	volumeInfo, err := ExecVolumeInfo()
	// Couldn't parse xml, so something is really wrong and up=0
	if err != nil {
		log.Errorf("couldn't parse json volume info: %v", err)
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0,
		)
	}

	// use OpErrno as indicator for up
	if i, _ := strconv.Atoi(volumeInfo.CliOutput.OpErrno); i != 0 {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 1.0,
		)
	}

	if len(volumeInfo.CliOutput.VolInfo.Volumes.Volume) != 0 {
		if i, _ := strconv.Atoi(volumeInfo.CliOutput.VolInfo.Volumes.Count); i != 0 {
			ch <- prometheus.MustNewConstMetric(
				volumesCount, prometheus.GaugeValue, float64(i),
			)
		}

	}

	for _, volume := range volumeInfo.CliOutput.VolInfo.Volumes.Volume {
		if e.Volumes[0] == allVolumes || ContainsVolume(e.Volumes, volume.Name) {

			if i, _ := strconv.Atoi(volume.BrickCount); i != 0 {
				ch <- prometheus.MustNewConstMetric(
					brickCount, prometheus.GaugeValue, float64(i), volume.Name,
				)
			}

			if i, _ := strconv.Atoi(volume.Status); i != 0 {
				ch <- prometheus.MustNewConstMetric(
					volumeStatus, prometheus.GaugeValue, float64(i), volume.Name,
				)
			}

		}
	}

	// reads gluster peer status
	peerStatus, peerStatusErr := ExecPeerStatus()
	if peerStatusErr != nil {
		log.Errorf("couldn't parse json of peer status: %v", peerStatusErr)
	}
	count := 0
	for range peerStatus.CliOutput.PeerStatus.Peer {
		count++
	}
	ch <- prometheus.MustNewConstMetric(
		peersConnected, prometheus.GaugeValue, float64(count),
	)

	// reads profile info
	if e.Profile {
		for _, volume := range volumeInfo.CliOutput.VolInfo.Volumes.Volume {
			if e.Volumes[0] == allVolumes || ContainsVolume(e.Volumes, volume.Name) {
				volumeProfile, execVolProfileErr := ExecVolumeProfileGvInfoCumulative(volume.Name)
				if execVolProfileErr != nil {
					log.Errorf("Error while executing or marshalling gluster profile output: %v", execVolProfileErr)
				}
				for _, brick := range volumeProfile.VolProfile.Brick {
					if strings.HasPrefix(brick.BrickName, e.Hostname) {
						ch <- prometheus.MustNewConstMetric(
							brickDuration, prometheus.CounterValue, float64(brick.CumulativeStats.Duration), volume.Name, brick.BrickName,
						)

						ch <- prometheus.MustNewConstMetric(
							brickDataRead, prometheus.CounterValue, float64(brick.CumulativeStats.TotalRead), volume.Name, brick.BrickName,
						)

						ch <- prometheus.MustNewConstMetric(
							brickDataWritten, prometheus.CounterValue, float64(brick.CumulativeStats.TotalWrite), volume.Name, brick.BrickName,
						)
						for _, fop := range brick.CumulativeStats.FopStats.Fop {
							ch <- prometheus.MustNewConstMetric(
								brickFopHits, prometheus.CounterValue, float64(fop.Hits), volume.Name, brick.BrickName, fop.Name,
							)

							ch <- prometheus.MustNewConstMetric(
								brickFopLatencyAvg, prometheus.GaugeValue, fop.AvgLatency, volume.Name, brick.BrickName, fop.Name,
							)

							ch <- prometheus.MustNewConstMetric(
								brickFopLatencyMin, prometheus.GaugeValue, fop.MinLatency, volume.Name, brick.BrickName, fop.Name,
							)

							ch <- prometheus.MustNewConstMetric(
								brickFopLatencyMax, prometheus.GaugeValue, fop.MaxLatency, volume.Name, brick.BrickName, fop.Name,
							)
						}
					}
				}
			}
		}
	}

	// executes gluster status all detail
	volumeStatusAll, err := ExecVolumeStatusAllDetail()
	if err != nil {
		log.Errorf("couldn't parse json of peer status: %v", err)
	}

	for _, vol := range volumeStatusAll.CliOutput.VolStatus.Volumes.Volume {
		for _, node := range vol.Node {
			if i, _ := strconv.Atoi(node.SizeTotal); i != 0 {
				ch <- prometheus.MustNewConstMetric(
					nodeSizeTotalBytes, prometheus.CounterValue, float64(i), node.Hostname, node.Path, vol.VolName,
				)
			}
			if i, _ := strconv.Atoi(node.SizeFree); i != 0 {
				ch <- prometheus.MustNewConstMetric(
					nodeSizeFreeBytes, prometheus.GaugeValue, float64(i), node.Hostname, node.Path, vol.VolName,
				)
			}
			if i, _ := strconv.Atoi(node.InodesTotal); i != 0 {
				ch <- prometheus.MustNewConstMetric(
					nodeInodesTotal, prometheus.CounterValue, float64(i), node.Hostname, node.Path, vol.VolName,
				)
			}

			if i, _ := strconv.Atoi(node.InodesFree); i != 0 {
				ch <- prometheus.MustNewConstMetric(
					nodeInodesFree, prometheus.GaugeValue, float64(i), node.Hostname, node.Path, vol.VolName,
				)
			}

		}
	}
	vols := e.Volumes
	if vols[0] == allVolumes {
		log.Warn("no Volumes were given.")
		volumeList, volumeListErr := ExecVolumeList()
		if volumeListErr != nil {
			log.Error(volumeListErr)
		}

		vols = volumeList.CliOutput.VolList.Volume
	}

	for _, vol := range vols {
		filesCount, volumeHealErr := ExecVolumeHealInfo(vol)
		if volumeHealErr == nil {
			ch <- prometheus.MustNewConstMetric(
				healInfoFilesCount, prometheus.CounterValue, float64(filesCount), vol,
			)
		}
	}

	mountBuffer, execMountCheckErr := ExecMountCheck()
	if execMountCheckErr != nil {
		log.Error(execMountCheckErr)
	} else {
		mounts, err := parseMountOutput(mountBuffer.String())
		if err != nil {
			log.Error(err)
			if len(mounts) > 0 {
				for _, mount := range mounts {
					ch <- prometheus.MustNewConstMetric(
						mountSuccessful, prometheus.GaugeValue, float64(0), mount.volume, mount.mountPoint,
					)
				}
			}
		} else {
			for _, mount := range mounts {
				ch <- prometheus.MustNewConstMetric(
					mountSuccessful, prometheus.GaugeValue, float64(1), mount.volume, mount.mountPoint,
				)

				isWriteable, err := ExecTouchOnVolumes(mount.mountPoint)
				if err != nil {
					log.Error(err)
				}
				if isWriteable {
					ch <- prometheus.MustNewConstMetric(
						volumeWriteable, prometheus.GaugeValue, float64(1), mount.volume, mount.mountPoint,
					)
				} else {
					ch <- prometheus.MustNewConstMetric(
						volumeWriteable, prometheus.GaugeValue, float64(0), mount.volume, mount.mountPoint,
					)
				}
			}
		}
	}

	if e.Quota {
		for _, volume := range volumeInfo.CliOutput.VolInfo.Volumes.Volume {
			if e.Volumes[0] == allVolumes || ContainsVolume(e.Volumes, volume.Name) {
				volumeQuotaJSON, err := ExecVolumeQuotaList(volume.Name)
				if err != nil {
					log.Error("Cannot create quota metrics if quotas are not enabled in your gluster server")
				} else {
					for _, limit := range volumeQuotaJSON.CliOutput.VolQuota.Limit {
						if i, err := strconv.Atoi(limit.HardLimit); err != nil {
							ch <- prometheus.MustNewConstMetric(
								quotaHardLimit,
								prometheus.CounterValue,
								float64(i),
								limit.Path,
								volume.Name,
							)

						}
						if i, err := strconv.Atoi(limit.SoftLimitValue); err != nil {
							ch <- prometheus.MustNewConstMetric(
								quotaSoftLimit,
								prometheus.CounterValue,
								float64(i),
								limit.Path,
								volume.Name,
							)

						}

						if i, err := strconv.Atoi(limit.UsedSpace); err != nil {
							ch <- prometheus.MustNewConstMetric(
								quotaUsed,
								prometheus.CounterValue,
								float64(i),
								limit.Path,
								volume.Name,
							)

						}

						if i, err := strconv.Atoi(limit.AvailSpace); err != nil {
							ch <- prometheus.MustNewConstMetric(
								quotaAvailable,
								prometheus.CounterValue,
								float64(i),
								limit.Path,
								volume.Name,
							)

						}

						slExceeded := 0.0
						if limit.SlExceeded != "No" {
							slExceeded = 1.0
						}
						ch <- prometheus.MustNewConstMetric(
							quotaSoftLimitExceeded,
							prometheus.CounterValue,
							slExceeded,
							limit.Path,
							volume.Name,
						)

						hlExceeded := 0.0
						if limit.HlExceeded != "No" {
							hlExceeded = 1.0
						}
						ch <- prometheus.MustNewConstMetric(
							quotaHardLimitExceeded,
							prometheus.CounterValue,
							hlExceeded,
							limit.Path,
							volume.Name,
						)
					}
				}
			}
		}
	}
}

type mount struct {
	mountPoint string
	volume     string
}

// ParseMountOutput pares output of system execution 'mount'
func parseMountOutput(mountBuffer string) ([]mount, error) {
	mounts := make([]mount, 0, 2)
	mountRows := strings.Split(mountBuffer, "\n")
	for _, row := range mountRows {
		trimmedRow := strings.TrimSpace(row)
		if len(row) > 3 {
			mountColumns := strings.Split(trimmedRow, " ")
			mounts = append(mounts, mount{mountPoint: mountColumns[2], volume: mountColumns[0]})
		}
	}
	return mounts, nil
}

// ContainsVolume checks a slice if it contains an element
func ContainsVolume(slice []string, element string) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}

// NewExporter initialises exporter
func NewExporter(hostname, glusterExecPath, volumesString string, profile bool, quota bool) (*Exporter, error) {
	if len(glusterExecPath) < 1 {
		log.Fatalf("Gluster executable path is wrong: %v", glusterExecPath)
	}
	volumes := strings.Split(volumesString, ",")
	if len(volumes) < 1 {
		log.Warnf("No volumes given. Proceeding without volume information. Volumes: %v", volumesString)
	}

	return &Exporter{
		Hostname: hostname,
		Volumes:  volumes,
		Profile:  profile,
		Quota:    quota,
	}, nil
}

func init() {
	prometheus.MustRegister(version.NewCollector("gluster_exporter"))
}
