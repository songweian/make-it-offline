package main

import (
	"flag"
	"fmt"
	"make-it-offline/pkg/plugins"
	_ "make-it-offline/repos/grafana"
	_ "make-it-offline/repos/mattermost"
	_ "make-it-offline/repos/mysql"
	_ "make-it-offline/repos/nginx"
	_ "make-it-offline/repos/postgresql"
	_ "make-it-offline/repos/prometheus"
	_ "make-it-offline/repos/redis"
	"os"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <app@version> <os@version> <arch> <formats>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nArguments:")
		fmt.Fprintln(os.Stderr, "  app@version   Application name and version (e.g., nginx@1.21)")
		fmt.Fprintln(os.Stderr, "  os@version    Operating system and version (e.g., ubuntu@20.04)")
		fmt.Fprintln(os.Stderr, "  arch          Architecture (e.g., x86_64, arm64)")
		fmt.Fprintln(os.Stderr, "  formats       Comma-separated list of formats (e.g., docker-compose,rpm)")
		fmt.Fprintln(os.Stderr, "\nExample:")
		fmt.Fprintln(os.Stderr, "  make-it-offline nginx@1.21 ubuntu@20.04 x86_64 docker-compose,rpm")
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 4 {
		flag.Usage()
		os.Exit(1)
	}

	appInfo := args[0]
	osInfo := args[1]
	arch := args[2]
	formats := strings.Split(args[3], ",")

	appName, appVersion := parseInfo(appInfo)
	osName, osVersion := parseInfo(osInfo)

	plugin, ok := plugins.GetPlugin(appName)
	if !ok {
		fmt.Printf("Error: Plugin for %s not found\n", appName)
		os.Exit(1)
	}

	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error: Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	resultPath, err := plugin.Generate(appVersion, osName, osVersion, arch, formats)
	if err != nil {
		fmt.Printf("Error generating package: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated package: %s\n", resultPath)
}

func parseInfo(infoStr string) (string, string) {
	if infoStr == "" {
		return "", ""
	}
	if strings.Contains(infoStr, "@") {
		parts := strings.SplitN(infoStr, "@", 2)
		return parts[0], parts[1]
	}
	return infoStr, "latest"
}
