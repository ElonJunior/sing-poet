package shortcuts

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	C "github.com/sagernet/sing-box/constant"

	"github.com/sagernet/sing-box/poet/constant"
)

func ShowVersion() {
	output := fmt.Sprintf("\n\t%s version %s\nIntro: %s\n", strings.ToUpper(constant.Name), constant.Version, constant.Intro)
	output += "Developed by KaringX @ElonWang\nLicensed under GNU General Public License version 3\nGitHub Repository:\thttps://github.com/KaringX/sing-poet\nDocuments:\thttps://karing.app/sing-poet\n"

	version := "\n\tsing-box version " + C.Version + "\n"
	version += "Environment: " + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH + "\n"

	var tags string
	var revision string

	debugInfo, loaded := debug.ReadBuildInfo()
	if loaded {
		for _, setting := range debugInfo.Settings {
			switch setting.Key {
			case "-tags":
				tags = setting.Value
			case "vcs.revision":
				revision = setting.Value
			}
		}
	}

	if tags != "" {
		version += "Tags: " + tags + "\n"
	}
	if revision != "" {
		version += "Revision: " + revision + "\n"
	}

	if C.CGO_ENABLED {
		version += "CGO: enabled\n"
	} else {
		version += "CGO: disabled\n"
	}

	os.Stdout.WriteString(output + version)
}

// func getGitLastTag() (string, error) {
// 	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return "", err
// 	}
// 	tag := strings.TrimSpace(string(output))
// 	return tag, nil
// }
