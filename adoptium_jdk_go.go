package adoptium_jdk_go

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/itchyny/gojq"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const ProgramVersion = "0.1.1"
const apiUrl = "https://api.adoptium.net/v3/info/available_releases"

type ProgramArgs struct {
	release      int
	arch         string
	os           string
	lts          bool
	showVer      bool
	listReleases bool
}

func main() {

	args := setupArgs()
	if args.release == 0 {

		apiLatestRelease(args)
	}

	var result string
	if args.showVer {

		result = apiLatestVersion(apiResponseBytes(apiEndpoint(args.release)))

	} else if args.listReleases {

		result = ApiListReleases()

	} else {

		result = apiPackageUrl(args.arch, args.os, apiResponseBytes(apiEndpoint(args.release)))
	}
	fmt.Println(result)
}

func ApiListReleases() string {

	query := ".available_releases"
	return queryForString(query, apiResponseBytes(apiUrl))
}

func apiLatestRelease(args *ProgramArgs) {
	query := ".most_recent_feature_release"
	if args.lts {
		query = ".most_recent_lts"
	}
	fRel, err := strconv.ParseFloat(queryForString(query, apiResponseBytes(apiUrl)), 64)
	check(err)
	rel := int(fRel)
	args.release = rel
}

func apiEndpoint(release int) string {
	const api = "https://api.adoptium.net/v3/assets/latest/$RELEASE/hotspot?vendor=eclipse"
	return strings.Replace(api, "$RELEASE", strconv.Itoa(release), 1)
}

func apiPackageUrl(arch, os string, apiResponse []byte) string {
	const jqQuery = `.[] | .binary | select(.image_type == "jdk") | select(.architecture == "$ARCH") | select(.os == "$OS") | .package.link`
	queryStr := strings.Replace(jqQuery, "$ARCH", arch, 1)
	queryStr = strings.Replace(queryStr, "$OS", os, 1)
	return queryForString(queryStr, apiResponse)
}

func apiLatestVersion(apiResponse []byte) string {
	version := queryForString(`. | first | .version.openjdk_version`, apiResponse)
	r := strings.Index(version, "+")
	if r == -1 {
		r = strings.Index(version, "-")
	}
	if r == -1 {
		return "error"
	} else {
		return version[0:r]
	}
}

func apiResponseBytes(url string) []byte {

	res, err := http.Get(url)
	check(err)

	body, err := io.ReadAll(res.Body)
	check(err)

	return body
}

func queryForString(jqQuery string, jsonBytes []byte) string {

	var unmarshalledJson interface{}
	e := json.Unmarshal(jsonBytes, &unmarshalledJson)
	check(e)

	queryRunner, err := gojq.Parse(jqQuery)
	check(err)

	result := "error: could not understand API response"
	iter := queryRunner.Run(unmarshalledJson)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			check(err)
		}
		if num, ok := v.(float64); ok {
			result = fmt.Sprintf("%f", num)
		}
		if s, ok := v.(string); ok {
			result = s
		}
		if interfaceArray, ok := v.([]interface{}); ok {

			result = ""
			for index, interfaceArrayItem := range interfaceArray {

				packageUrl := apiPackageUrl(arch(), runtime.GOOS, apiResponseBytes(apiEndpoint(int(interfaceArrayItem.(float64)))))
				if index == 0 {

					result += packageUrl

				} else {

					result += "\n" + packageUrl
				}
			}
		}
	}
	return result
}

func setupArgs() *ProgramArgs {

	if len(os.Args) == 2 && os.Args[1] == "--help" {
		fmt.Println("Use the '-h' command-line option for help.")
		os.Exit(0)
	}

	jdkRelease := flag.Int("release", 0, "The major JDK release; if not provided, the latest is found and used.")
	jdkArch := flag.String("arch", arch(), "The JDK target machine architecture")
	jdkOS := flag.String("os", runtime.GOOS, "The JDK target OS")
	ltsRelease := flag.Bool("lts", false, "Get the latest LTS release")
	showVer := flag.Bool("jv", false, "Print the JDK version only, not the URL")
	showHelp := flag.Bool("h", false, "Show help/usage")
	showVersion := flag.Bool("v", false, "Show version info")
	listReleases := flag.Bool("listReleases", false, "List Available Releases")
	flag.Parse()

	if *showHelp {
		flag.PrintDefaults()
		fmt.Println("\nTry `curl -LO $(latest-jdk)` to download the latest JDK.")
		os.Exit(0)
	}
	if *showVersion {
		fmt.Printf("latest-jdk v%v\n", ProgramVersion)
		os.Exit(0)
	}

	args := ProgramArgs{release: *jdkRelease, arch: *jdkArch, os: *jdkOS, lts: *ltsRelease, showVer: *showVer, listReleases: *listReleases}
	return &args
}

func arch() string {
	if runtime.GOARCH == "amd64" {
		return "x64"
	} else {
		return runtime.GOARCH
	}
}

func check(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
