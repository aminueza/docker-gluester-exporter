package expogluster

import (
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Exporter holds name, path and volumes to be monitored
type Exporter struct {
	Router   *mux.Router
	Hostname string
	Volumes  []string
	Profile  bool
	Quota    bool
}

//NewPromExporter creates a new exporter for prometheus
func NewPromExporter() *Exporter {
	return &Exporter{
		Hostname: getEnv("PROM_HOSTNAME", "0.0.0.0:9189"),
		Volumes:  []string{getEnv("PROM_VOLUMES", "_all")},
		Profile:  parseBool(getEnv("PROM_PROFILE", "false")),
		Quota:    parseBool(getEnv("PROM_QUOTA", "false")),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func parseBool(env string) bool {
	sslbool, _ := strconv.ParseBool(env)
	return sslbool
}
