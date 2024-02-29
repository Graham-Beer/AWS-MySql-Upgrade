package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type DescribeDB interface {
	*rds.DescribeDBSnapshotsOutput | *rds.DescribeDBInstancesOutput
}

type ResourceGetter[T DescribeDB] func(client *rds.Client, identifier string) (T, error)
type getstatus[T DescribeDB] func(db T) string

func ProcessStatus[T DescribeDB](
	status string,
	client *rds.Client,
	identifier string,
	getter ResourceGetter[T],
	getstate getstatus[T],
) (string, error) {
	switch status {
	case "creating", "modifying", "upgrading":
		fmt.Printf("%s..", status)
		for status == "creating" || status == "modifying" || status == "upgrading" {
			time.Sleep(10 * time.Second)
			fmt.Print("..")
			resource, err := getter(client, identifier)
			if err != nil {
				return "", err
			}
			status = getstate(resource)
		}
		return status, nil
	case "available":
		return "", nil
	default:
		fmt.Printf("\nUnknown status: %s\n", status)
		return "", nil
	}
}

func PrintUpgradeSuccess(created string) {
	fmt.Println("]")
	fmt.Println("")
	fmt.Println("**Upgrade Steps:**")
	fmt.Printf("1. Snapshot created successfully %s.\n", created)
}

func IsAvailable() {
	fmt.Printf("%s", "]")
	fmt.Println("")
}

func DbComplete(engineVersion string) {
	completed := time.Now()
	formattedCompleteTime := completed.Format("2006-01-02 15:04:05")
	fmt.Println("")
	fmt.Printf("**Database version check: %s**\n", engineVersion)
	fmt.Printf("**Database upgrade complete at: %s**\n", formattedCompleteTime)
}

func DbUpgradeMessage() {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	fmt.Printf("4. Database upgrade started at: %s\n", formattedTime)
	fmt.Println("")
	fmt.Println("**Upgrade Progress:**")
	fmt.Printf("  Status: %s", "[")
}

func GetDBInstance(client *rds.Client, identifier string) (*rds.DescribeDBInstancesOutput, error) {
	return client.DescribeDBInstances(context.Background(), &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: &identifier,
	})
}

func GetDBSnapshot(client *rds.Client, identifier string) (*rds.DescribeDBSnapshotsOutput, error) {
	return client.DescribeDBSnapshots(context.Background(), &rds.DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: &identifier,
	})
}

func GetDbStatus(db *rds.DescribeDBInstancesOutput) string {
	return *db.DBInstances[0].DBInstanceStatus
}

func GetDbVersion(db *rds.DescribeDBInstancesOutput) string {
	return *db.DBInstances[0].EngineVersion
}

func GetDBSnapshotStatus(snap *rds.DescribeDBSnapshotsOutput) string {
	return *snap.DBSnapshots[0].Status
}

func GetSnapshotCreateTime(snap *rds.DescribeDBSnapshotsOutput) string {
	created := *snap.DBSnapshots[0].SnapshotCreateTime
	return created.Format("2006-01-02 15:04:05")
}
