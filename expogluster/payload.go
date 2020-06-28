package expogluster

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/prometheus/common/log"
)

type VolumeInfoJSON struct {
	CliOutput struct {
		OpRet    string `json:"opRet"`
		OpErrno  string `json:"opErrno"`
		OpErrstr string `json:"opErrstr"`
		VolInfo  struct {
			Volumes struct {
				Volume []struct {
					BrickCount      string `json:"brickCount"`
					DistCount       string `json:"distCount"`
					ArbiterCount    string `json:"arbiterCount"`
					RedundancyCount string `json:"redundancyCount"`
					Name            string `json:"name"`
					StatusStr       string `json:"statusStr"`
					StripeCount     string `json:"stripeCount"`
					Transport       string `json:"transport"`
					Options         struct {
						Option []struct {
							Name  string `json:"name"`
							Value string `json:"value"`
						} `json:"option"`
					} `json:"options"`
					SnapshotCount string `json:"snapshotCount"`
					ReplicaCount  string `json:"replicaCount"`
					DisperseCount string `json:"disperseCount"`
					Bricks        struct {
						Brick []struct {
							UUID      string `json:"-uuid"`
							Name      string `json:"name"`
							HostUUID  string `json:"hostUuid"`
							IsArbiter string `json:"isArbiter"`
						} `json:"brick"`
					} `json:"bricks"`
					OptCount string `json:"optCount"`
					ID       string `json:"id"`
					Status   string `json:"status"`
					Type     string `json:"type"`
					TypeStr  string `json:"typeStr"`
				} `json:"volume"`
				Count string `json:"count"`
			} `json:"volumes"`
		} `json:"volInfo"`
	} `json:"cliOutput"`
}

// Volume element of "gluster volume info" command
type Volume struct {
	Name       string  `json:"name"`
	ID         string  `json:"id"`
	Status     int     `json:"status"`
	StatusStr  string  `json:"statusStr"`
	BrickCount int     `json:"brickCount"`
	Bricks     []Brick `json:"bricks"`
	DistCount  int     `json:"distCount"`
}

// Brick element of "gluster volume info" command
type Brick struct {
	UUID      string `json:"brick>uuid"`
	Name      string `json:"brick>name"`
	HostUUID  string `json:"brick>hostUuid"`
	IsArbiter int    `json:"brick>isArbiter"`
}

// VolumeListJSON struct represents cliOutput element of "gluster volume list" command
type VolumeListJSON struct {
	CliOutput struct {
		OpRet    string `json:"opRet"`
		OpErrno  string `json:"opErrno"`
		OpErrstr string `json:"opErrstr"`
		VolList  struct {
			Volume []string `json:"volume"`
			Count  string   `json:"count"`
		} `json:"volList"`
	} `json:"cliOutput"`
}

// PeerStatusJSON struct represents cliOutput element of "gluster peer status" command
type PeerStatusJSON struct {
	CliOutput struct {
		OpRet      string `json:"opRet"`
		OpErrno    string `json:"opErrno"`
		OpErrstr   string `json:"opErrstr"`
		PeerStatus struct {
			Peer []struct {
				UUID      string `json:"uuid"`
				Hostname  string `json:"hostname"`
				Hostnames struct {
					Hostname string `json:"hostname"`
				} `json:"hostnames"`
				Connected string `json:"connected"`
				State     string `json:"state"`
				StateStr  string `json:"stateStr"`
			} `json:"peer"`
		} `json:"peerStatus"`
	} `json:"cliOutput"`
}

// Hostnames element of "gluster peer status" command
// type Hostnames struct {
// 	Hostname string `json:"hostname"`
// }

// VolumeProfileJSON struct represents cliOutput element of "gluster volume {volume} profile" command
type VolumeProfileJSON struct {
	// JSONName   string     `json:"cliOutput"`
	OpRet      int        `json:"opRet"`
	OpErrno    int        `json:"opErrno"`
	OpErrstr   string     `json:"opErrstr"`
	VolProfile VolProfile `json:"volProfile"`
}

// VolProfile element of "gluster volume {volume} profile" command
type VolProfile struct {
	Volname    string         `json:"volname"`
	BrickCount int            `json:"brickCount"`
	Brick      []BrickProfile `json:"brick"`
}

// BrickProfile struct for element brick of "gluster volume {volume} profile" command
type BrickProfile struct {
	//JSONName string `json:"brick"`
	BrickName       string          `json:"brickName"`
	CumulativeStats CumulativeStats `json:"cumulativeStats"`
}

// CumulativeStats element of "gluster volume {volume} profile" command
type CumulativeStats struct {
	FopStats   FopStats `json:"fopStats"`
	Duration   int      `json:"duration"`
	TotalRead  int      `json:"totalRead"`
	TotalWrite int      `json:"totalWrite"`
}

// FopStats element of "gluster volume {volume} profile" command
type FopStats struct {
	Fop []Fop `json:"fop"`
}

// Fop is struct for FopStats
type Fop struct {
	Name       string  `json:"name"`
	Hits       int     `json:"hits"`
	AvgLatency float64 `json:"avgLatency"`
	MinLatency float64 `json:"minLatency"`
	MaxLatency float64 `json:"maxLatency"`
}

// HealInfoBrick is a struct of HealInfoBricks
type HealInfoBrick struct {
	JSONName        string `json:"brick"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	NumberOfEntries string `json:"numberOfEntries"`
}

// HealInfoBricks is a struct of HealInfo
type HealInfoBricks struct {
	JSONName string          `json:"bricks"`
	Brick    []HealInfoBrick `json:"brick"`
}

// HealInfo is a struct of VolumenHealInfoJSON
type HealInfo struct {
	//JSONName string         `json:"healInfo"`
	Bricks HealInfoBricks `json:"bricks"`
}

// VolumeHealInfoJSON struct represents cliOutput element of "gluster volume {volume} heal info" command
type VolumeHealInfoJSON struct {
	//JSONName string   `json:"cliOutput"`
	OpRet    int      `json:"opRet"`
	OpErrno  int      `json:"opErrno"`
	OpErrstr string   `json:"opErrstr"`
	HealInfo HealInfo `json:"healInfo"`
}

// VolumeHealInfoJSONUnmarshall unmarshalls heal info of gluster cluster
func VolumeHealInfoJSONUnmarshall(cmdOutBuff io.Reader) (VolumeHealInfoJSON, error) {
	var vol VolumeHealInfoJSON
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	err = json.Unmarshal(b, &vol)
	if err != nil {
		log.Error(err)
	}
	return vol, nil
}

func skipRoot(jsonBlob []byte) json.RawMessage {
	var root map[string]json.RawMessage

	if err := json.Unmarshal(jsonBlob, &root); err != nil {
		panic(err)
	}
	for _, v := range root {
		return v
	}
	return nil
}

// VolumeListJSONUnmarshall unmarshalls bytes to VolumeListJSON struct
func VolumeListJSONUnmarshall(cmdOutBuff io.Reader) (VolumeListJSON, error) {
	var vol VolumeListJSON
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}

	err = json.Unmarshal(skipRoot(b), &vol)
	return vol, err
}

// VolumeInfoJSONUnmarshall unmarshalls bytes to VolumeInfoJSON struct
func VolumeInfoJSONUnmarshall(cmdOutBuff io.Reader) (VolumeInfoJSON, error) {
	var vol VolumeInfoJSON
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}

	err = json.Unmarshal(b, &vol)
	return vol, err
}

// PeerStatusJSONUnmarshall unmarshalls bytes to PeerStatusJSON struct
func PeerStatusJSONUnmarshall(cmdOutBuff io.Reader) (PeerStatusJSON, error) {
	var vol PeerStatusJSON
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	err = json.Unmarshal(b, &vol)
	return vol, err
}

// VolumeProfileGvInfoCumulativeJSONUnmarshall unmarshalls cumulative profile of gluster volume profile
func VolumeProfileGvInfoCumulativeJSONUnmarshall(cmdOutBuff io.Reader) (VolumeProfileJSON, error) {
	var vol VolumeProfileJSON
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	err = json.Unmarshal(b, &vol)
	return vol, err
}

// VolumeStatusJSON JSON type of "gluster volume status"
type VolumeStatusJSON struct {
	CliOutput struct {
		OpRet     string `json:"opRet"`
		OpErrno   string `json:"opErrno"`
		OpErrstr  string `json:"opErrstr"`
		VolStatus struct {
			Volumes struct {
				Volume []struct {
					VolName   string `json:"volName"`
					NodeCount string `json:"nodeCount"`
					Node      []struct {
						Hostname    string `json:"hostname"`
						Peerid      string `json:"peerid"`
						SizeTotal   string `json:"sizeTotal"`
						MntOptions  string `json:"mntOptions"`
						Path        string `json:"path"`
						Pid         string `json:"pid"`
						FsName      string `json:"fsName"`
						InodesFree  string `json:"inodesFree"`
						Device      string `json:"device"`
						BlockSize   string `json:"blockSize"`
						InodesTotal string `json:"inodesTotal"`
						Status      string `json:"status"`
						Port        string `json:"port"`
						Ports       struct {
							Rdma string `json:"rdma"`
							TCP  string `json:"tcp"`
						} `json:"ports"`
						SizeFree string `json:"sizeFree"`
					} `json:"node"`
				} `json:"volume"`
			} `json:"volumes"`
		} `json:"volStatus"`
	} `json:"cliOutput"`
}

// VolumeStatusAllDetailJSONUnmarshall reads bytes.buffer and returns unmarshalled json
func VolumeStatusAllDetailJSONUnmarshall(cmdOutBuff io.Reader) (VolumeStatusJSON, error) {
	var vol VolumeStatusJSON
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return vol, err
	}
	err = json.Unmarshal(b, &vol)
	return vol, err
}

type Quota struct {
	Quota   string `json:"quota"`
	Volume  string `json:"volume"`
	Subdir  string `json:"subdir"`
	Percent string `json:"percent"`
}

// // VolumeQuotaJSON JSON type of "gluster volume quota list"
type VolumeQuotaJSON struct {
	CliOutput struct {
		OpRet    string `json:"opRet"`
		OpErrno  string `json:"opErrno"`
		OpErrstr string `json:"opErrstr"`
		VolQuota struct {
			Limit []struct {
				HardLimit        string `json:"hard_limit"`
				SoftLimitPercent string `json:"soft_limit_percent"`
				SoftLimitValue   string `json:"soft_limit_value"`
				UsedSpace        string `json:"used_space"`
				AvailSpace       string `json:"avail_space"`
				SlExceeded       string `json:"sl_exceeded"`
				HlExceeded       string `json:"hl_exceeded"`
				Path             string `json:"path"`
			} `json:"limit"`
		} `json:"volQuota"`
	} `json:"cliOutput"`
}

// VolumeQuotaListJSONUnmarshall function parse "gluster volume quota list" JSON output
func VolumeQuotaListJSONUnmarshall(cmdOutBuff io.Reader) (VolumeQuotaJSON, error) {
	var volQuotaJSON VolumeQuotaJSON
	b, err := ioutil.ReadAll(cmdOutBuff)
	if err != nil {
		log.Error(err)
		return volQuotaJSON, err
	}
	err = json.Unmarshal(b, &volQuotaJSON)
	return volQuotaJSON, err
}
