package about

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

func NewESInfo(c *gin.Context) *ESInfo {
	version := strings.TrimLeft(Version, "v")
	userAgent := c.Request.UserAgent()
	if strings.Contains(userAgent, "Elastic") {
		reg := regexp.MustCompile(`([0-9]+\.[0-9]+\.[0-9]+)`)
		matches := reg.FindAllString(userAgent, 1)
		if len(matches) > 0 {
			version = matches[0]
		}
	}
	return &ESInfo{
		Name:        "quicksearch",
		ClusterName: "N/A",
		ClusterUUID: "N/A",
		Version: ESInfoVersion{
			Number:                    version,
			BuildFlavor:               "default",
			BuildHash:                 CommitHash,
			BuildDate:                 BuildDate,
			BuildSnapshot:             false,
			LuceneVersion:             "N/A",
			MinimumWireVersion:        "N/A",
			MinimumIndexCompatibility: "N/A",
		},
		Tagline: "You Know, for Search",
	}
}

func NewESLicense(_ *gin.Context) *ESLicense {
	return &ESLicense{
		Status: "active",
	}
}

func NewESXPack(_ *gin.Context) *ESXPack {
	return &ESXPack{
		Build:    make(map[string]bool),
		Features: make(map[string]bool),
		License: ESLicense{
			Status: "active",
		},
	}
}

type ESInfo struct {
	Name        string        `json:"name"`
	ClusterName string        `json:"cluster_name"`
	ClusterUUID string        `json:"cluster_uuid"`
	Version     ESInfoVersion `json:"version"`
	Tagline     string        `json:"tagline"`
}

type ESInfoVersion struct {
	Number                    string `json:"number"`
	BuildFlavor               string `json:"build_flavor"`
	BuildHash                 string `json:"build_hash"`
	BuildDate                 string `json:"build_date"`
	BuildSnapshot             bool   `json:"build_snapshot"`
	LuceneVersion             string `json:"lucene_version"`
	MinimumWireVersion        string `json:"minimum_wire_version"`
	MinimumIndexCompatibility string `json:"minimum_index_compatibility"`
}

type ESLicense struct {
	Status string `json:"status"`
}

type ESXPack struct {
	Build    map[string]bool `json:"build"`
	Features map[string]bool `json:"features"`
	License  ESLicense       `json:"license"`
}
