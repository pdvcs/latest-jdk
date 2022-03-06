package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/itchyny/gojq"
)

const PROGRAM_VERSION = "0.1.0"

type ProgramArgs struct {
	release int
	arch    string
	os      string
	lts     bool
}

func main() {
	args := setupArgs()
	if args.release == 0 {
		apiGetLatestRelease(args)
	}
	url := apiGetPackageUrl(args.arch, args.os, apiResponseBytes(apiEndpoint(args.release)))
	fmt.Println(url)
}

func apiGetLatestRelease(args *ProgramArgs) {
	const apiUrl = "https://api.adoptium.net/v3/info/available_releases"
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

func apiGetPackageUrl(arch, os string, apiRespnse []byte) string {
	const jqQuery = `.[] | .binary | select(.image_type == "jdk") | select(.architecture == "$ARCH") | select(.os == "$OS") | .package.link`
	queryStr := strings.Replace(jqQuery, "$ARCH", arch, 1)
	queryStr = strings.Replace(queryStr, "$OS", os, 1)
	return queryForString(queryStr, apiRespnse)
}

func apiResponseBytes(url string) []byte {
	userAgent := "latest-jdk/" + PROGRAM_VERSION
	fetcher := http.Client{
		Timeout: time.Second * 4,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	check(err)

	req.Header.Set("User-Agent", userAgent)
	res, err := fetcher.Do(req)
	check(err)
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	check(err)

	return body
}

func queryForString(jqQuery string, jsonBytes []byte) string {
	var unmarshaledJson interface{}
	e := json.Unmarshal(jsonBytes, &unmarshaledJson)
	check(e)

	queryRunner, err := gojq.Parse(jqQuery)
	check(err)

	result := "error: could not understand API response"
	iter := queryRunner.Run(unmarshaledJson)
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
	}
	return result
}

func setupArgs() *ProgramArgs {

	if len(os.Args) == 2 && os.Args[1] == "--help" {
		fmt.Println("Use the '-h' command-line option for help.")
		os.Exit(0)
	}

	jdkRelease := flag.Int("release", 0, "The major JDK release; if not provided, the latest is found and used.")
	jdkArch := flag.String("arch", "x64", "The JDK target machine architecture")
	jdkOS := flag.String("os", runtime.GOOS, "The JDK target OS")
	ltsRelease := flag.Bool("lts", false, "Get the latest LTS release")
	showHelp := flag.Bool("h", false, "Show help/usage")
	showVersion := flag.Bool("v", false, "Show version info")
	flag.Parse()

	if *showHelp {
		flag.PrintDefaults()
		fmt.Println("\nTry `curl -LO $(latest-jdk)` to download the latest JDK.")
		os.Exit(0)
	}
	if *showVersion {
		fmt.Printf("latest-jdk v%v\n", PROGRAM_VERSION)
		os.Exit(0)
	}

	args := ProgramArgs{release: *jdkRelease, arch: *jdkArch, os: *jdkOS, lts: *ltsRelease}
	return &args
}

func check(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
