package main

import (
	"Kony/v2/pkg/database"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// Define variables instead of repeating strings
var (
	InstanceIdentifier          = "database-1"
	SnapshotIdentifier          = "cutover-test-1"
	NewEngineVersion            = "8.0.36"
	DbParameterGroupName        = "mysql8"
	DbParameterGroupFamily      = "mysql8.0"
	DbParameterGroupDescription = "MySQL 8 parameter group"
)

func main() {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
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

	// Create database snapshot
	snapshot, err := database.CreateDBSnapshot(client, InstanceIdentifier, SnapshotIdentifier)
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
	fmt.Printf("**Snapshot of [%s] Started**\n", InstanceIdentifier)
	err = database.WaitForSnapshotCompletion(client, SnapshotIdentifier)
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
