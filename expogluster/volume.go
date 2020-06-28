package expogluster

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	// "github.com/google/martian/log"

	"github.com/prometheus/common/log"
)

//ExecMountCheck checks mount point
func ExecMountCheck() (*bytes.Buffer, error) {
	stdoutBuffer := &bytes.Buffer{}
	mountCmd := exec.Command("mount", "-t", "fuse.glusterfs")

	mountCmd.Stdout = stdoutBuffer

	return stdoutBuffer, mountCmd.Run()
}

//ExecTouchOnVolumes checks mountpoint permission
func ExecTouchOnVolumes(mountpoint string) (bool, error) {
	testFileName := fmt.Sprintf("%v/%v_%v", mountpoint, "gluster_mount.test", time.Now())
	_, createErr := os.Create(testFileName)
	if createErr != nil {
		return false, createErr
	}
	removeErr := os.Remove(testFileName)
	if removeErr != nil {
		return false, removeErr
	}
	return true, nil
}

// ExecVolumeInfo executes "gluster volume info" at the local machine and
// returns VolumeInfoJSON struct and error
func ExecVolumeInfo() (VolumeInfoJSON, error) {
	bytesBuffer, cmdErr := gluster("volume", "info")
	if cmdErr != nil {
		return VolumeInfoJSON{}, cmdErr
	}
	volumeInfo, err := VolumeInfoJSONUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling json: %v", err)
		return volumeInfo, err
	}

	return volumeInfo, nil
}

// returns VolumeList struct and error
func ExecVolumeList() (*VolumeListJSON, error) {
	bytesBuffer, cmdErr := gluster("volume", "list")
	if cmdErr != nil {
		return &VolumeListJSON{}, cmdErr
	}
	volumeList, err := VolumeListJSONUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling json: %v", err)
		return &VolumeListJSON{}, err
	}
	return &volumeList, nil
}

// ExecPeerStatus executes "gluster peer status" at the local machine and
// returns PeerStatus struct and error
func ExecPeerStatus() (*PeerStatusJSON, error) {
	bytesBuffer, cmdErr := gluster("peer", "status")
	if cmdErr != nil {
		return &PeerStatusJSON{}, cmdErr
	}
	peerStatus, err := PeerStatusJSONUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling json: %v", err)
		return &peerStatus, err
	}

	return &peerStatus, nil
}

// ExecVolumeProfileGvInfoCumulative executes "gluster volume {volume] profile info cumulative" at the local machine and
// returns VolumeInfoJSON struct and error
func ExecVolumeProfileGvInfoCumulative(volumeName string) (*VolumeProfileJSON, error) {
	args := []string{"volume", "profile", volumeName, "info", "cumulative"}
	bytesBuffer, cmdErr := gluster(args...)
	if cmdErr != nil {
		return &VolumeProfileJSON{}, cmdErr
	}
	volumeProfile, err := VolumeProfileGvInfoCumulativeJSONUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling json: %v", err)
		return &volumeProfile, err
	}
	return &volumeProfile, nil
}

// ExecVolumeStatusAllDetail executes "gluster volume status all detail" at the local machine
// returns VolumeStatusJSON struct and error
func ExecVolumeStatusAllDetail() (*VolumeStatusJSON, error) {
	args := []string{"volume", "status", "all", "detail"}
	bytesBuffer, cmdErr := gluster(args...)
	if cmdErr != nil {
		return &VolumeStatusJSON{}, cmdErr
	}
	volumeStatus, err := VolumeStatusAllDetailJSONUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling json: %v", err)
		return &volumeStatus, err
	}

	log.Infof("%v", volumeStatus.CliOutput.VolStatus.Volumes.Volume[1].Node[0].SizeFree)
	return &volumeStatus, nil
}

// ExecVolumeHealInfo executes volume heal info on host system and processes input
// returns (int) number of unsynced files
func ExecVolumeHealInfo(volumeName string) (int, error) {
	entriesOutOfSync := 0
	bytesBuffer, cmdErr := gluster("volume", "heal", volumeName, "info")
	if cmdErr != nil {
		return -1, cmdErr
	}
	healInfo, err := VolumeHealInfoJSONUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling json: %v", err)
		return -1, err
	}

	for _, brick := range healInfo.HealInfo.Bricks.Brick {
		var count int
		var err error
		count, err = strconv.Atoi(brick.NumberOfEntries)
		if err != nil {
			log.Errorf("Something went wrong while parsing brick info: %v", err)
			return -1, err
		}
		entriesOutOfSync += count
	}
	return entriesOutOfSync, nil
}

// ExecVolumeQuotaList executes volume quota list on host system and processes input
// returns QuotaList structs and errors
func ExecVolumeQuotaList(volumeName string) (VolumeQuotaJSON, error) {

	result, cmdErr := gluster("volume", "quota", volumeName, "list")
	if cmdErr != nil {
		return VolumeQuotaJSON{}, cmdErr
	}
	volumeQuota, err := VolumeQuotaListJSONUnmarshall(result)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling json: %v", err)
		return volumeQuota, err
	}
	return volumeQuota, nil
}
