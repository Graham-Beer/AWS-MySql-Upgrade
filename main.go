package main

import (
	"Kony/v2/pkg/database"
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// Define variables instead of repeating strings
var (
	profile                     string
	InstanceIdentifier          string
	SnapshotIdentifier          = "cutover"
	NewEngineVersion            = "8.0.36"
	DbParameterGroupName        = "mysql8-capture"
	DbParameterGroupFamily      = "mysql8.0"
	DbParameterGroupDescription = "MySQL 8 parameter group"
)

func init() {
	flag.StringVar(&InstanceIdentifier, "instance-identifier", "", "RDS instance identifier")
	flag.StringVar(&profile, "profile", "default", "AWS profile name (defaults to 'default')")
	flag.Parse()
}

func main() {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion("eu-west-1"))
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return
	}

	// Create an RDS client
	client := rds.NewFromConfig(cfg)

	// Print header information
	fmt.Println("**Database Upgrade Report**")
	fmt.Println("")
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	fmt.Printf("**Upgrade Start Time: %s**\n", formattedTime)

	// Create a single function for common error handling
	handleError := func(err error, operation string) {
		if err != nil {
			fmt.Printf("Error during %s: %v\n", operation, err)
			return
		}
	}

	// Format the date in the desired format (YYYYMMDD)
	formattedDate := currentTime.Format("060104")

	SnapshotFormattedName := fmt.Sprintf("%s-%s", SnapshotIdentifier, formattedDate)

	// Create database snapshot
	snapshot, err := database.CreateDBSnapshot(client, InstanceIdentifier, SnapshotFormattedName)
	handleError(err, "creating snapshot")
	if err != nil {
		return
	}

	// Print snapshot information
	fmt.Println("**Snapshot Information:**")
	fmt.Printf("  * Snapshot ID: %s\n", *snapshot.DBSnapshotIdentifier)
	fmt.Printf("  * DB Engine: %s\n", *snapshot.Engine)
	fmt.Printf("  * Engine Version: %s\n", *snapshot.EngineVersion)

	// Wait for snapshot completion
	fmt.Println("")
	fmt.Printf("**Snapshot of [%s] Started**\n", InstanceIdentifier)
	err = database.WaitForSnapshotCompletion(client, SnapshotFormattedName)
	handleError(err, "waiting for snapshot completion")
	if err != nil {
		return
	}

	// Create DB parameter group
	err = database.CreateDBParameterGroup(client, DbParameterGroupName, DbParameterGroupFamily, DbParameterGroupDescription)
	handleError(err, "creating DB parameter group")
	if err != nil {
		return
	}
	fmt.Println("2. DB parameter group created successfully")

	// Modify DB instance
	err = database.ModifyDBInstance(client, InstanceIdentifier, NewEngineVersion, DbParameterGroupName)
	handleError(err, "modifying DB instance")
	if err != nil {
		return
	}
	fmt.Println("3. DB instance engine version modified")

	// Wait for DB upgrade
	err = database.WaitForDbUpgrade(client, InstanceIdentifier)
	handleError(err, "upgrading DB instance")
	if err != nil {
		return
	}
}
