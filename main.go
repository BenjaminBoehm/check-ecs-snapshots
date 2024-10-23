package main

import (
	"log"
	"time"

	check "github.com/NETWAYS/go-check"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/backups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v2/snapshots"
)

func main() {
	config := check.NewConfig()
	config.Name = "check-ecs-snapshots"
	config.Readme = `Check for Elastic Cloud Server (ECS) snapshot age on Open Telekom Cloud (OTC)`
	config.Version = "1.0.0"

	critical := config.FlagSet.IntP("critical", "c", 60, "critical threshold for age of snapshot in days")
	warning := config.FlagSet.IntP("warning", "w", 30, "warning threshold for age of snapshot in days")
	cbr := config.FlagSet.Bool("cbr", false, "check cbr backups instead of snapshots")

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

	var bckps []backups.Backup
	var snaps []snapshots.Snapshot
	var forLength int

	if *cbr {

		cbrClient, err := openstack.NewCBRService(providerClient, golangsdk.EndpointOpts{})

		if err != nil {
			log.Fatalf("Problem with cbr Service: %s", err)
		}

		bckps, err = backups.List(cbrClient, backups.ListOpts{})

		forLength = len(bckps)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		storageClient, err := openstack.NewBlockStorageV2(providerClient, golangsdk.EndpointOpts{})

		if err != nil {
			log.Fatalf("Problem with evs Service: %s", err)
		}

		pager := snapshots.List(storageClient, snapshots.ListOpts{})

		if err != nil {
			log.Fatal(err)
		}
		pages, err := pager.AllPages()

		if err != nil {
			log.Fatal(err)
		}

		snaps, err = snapshots.ExtractSnapshots(pages)

		if err != nil {
			log.Fatal(err)
		}

		forLength = len(snaps)

	}
	currentTime := time.Now()
	layout := "2006-01-02T15:04:05.999999"

	criticalCount := 0
	warningCount := 0
	okCount := 0

	for i := 0; i < forLength; i++ {

		var createdAt time.Time

		switch {
		case len(bckps) > 0:
			createdAtParsed, err := time.Parse(layout, bckps[i].CreatedAt)

			if err != nil {
				log.Fatal("Could not parse time:", err, createdAtParsed, currentTime)
			}

			createdAt = createdAtParsed
		case len(snaps) > 0:
			createdAt = snaps[i].CreatedAt
		default:
			log.Fatal("Can't check backups.")
		}

		duration := currentTime.Sub(createdAt)
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
		check.Exitf(check.Critical, "%d Snapshots older than %d days", criticalCount, *critical)
	} else if warningCount > 0 {
		check.Exitf(check.Warning, "%d Snapshots older than %d days", warningCount, *warning)
	} else {
		check.Exitf(check.OK, "No snapshots (%d) older than %d days", okCount, *warning)
	}

}
