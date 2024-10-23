package main

import (
	"fmt"
	"log"
	"time"

	check "github.com/NETWAYS/go-check"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/backups"
)

func main() {
	config := check.NewConfig()
	config.Name = "check-ecs-snapshots"
	config.Readme = `Check for Elastic Cloud Server (ECS) snapshot age on Open Telekom Cloud (OTC)`
	config.Version = "1.0.0"

	critical := config.FlagSet.IntP("critical", "c", 60, "critical threshold for age of snapshot in days")
	warning := config.FlagSet.IntP("warning", "w", 30, "warning threshold for age of snapshot in days")

	config.ParseArguments()

	env := openstack.NewEnv("OS_", true)
	cloud, err := env.Cloud()

	if err != nil {
		panic(err)
	}

	providerClient, err := openstack.AuthenticatedClientFromCloud(cloud)

	if err != nil {
		log.Fatalf("Can't authenticate: %s", err)
	}

	cbr, err := openstack.NewCBRService(providerClient, golangsdk.EndpointOpts{})

	if err != nil {
		log.Fatalf("Problem with cbr Service: %s", err)
	}

	snaps, err := backups.List(cbr, backups.ListOpts{})

	if err != nil {
		log.Fatal(err)
	}
	currentTime := time.Now()
	layout := "2006-01-02T15:04:05.999999"

	criticalCount := 0
	warningCount := 0
	okCount := 0

	for _, s := range snaps {
		createdAtParsed, err := time.Parse(layout, s.CreatedAt)
		if err != nil {
			fmt.Println("Could not parse time:", err, createdAtParsed, currentTime)
		}

		duration := currentTime.Sub(createdAtParsed)
		daysDiff := int(duration.Hours() / 24)

		if *critical <= daysDiff {
			criticalCount += 1

		} else if *warning <= daysDiff {
			warningCount += 1
		} else {
			okCount += 1
		}
	}

	if criticalCount > 0 {
		check.Exitf(check.Critical, "%d Snapshots older than %d days\n", criticalCount, *critical)
	} else if warningCount > 0 {
		check.Exitf(check.Warning, "%d Snapshots older than %d days\n", warningCount, *warning)
	} else {
		check.Exitf(check.OK, "No snapshots (%d) older than %d days\n", okCount, *warning)
	}

}
