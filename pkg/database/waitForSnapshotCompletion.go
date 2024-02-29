package database

import (
	"fmt"

	"Kony/v2/pkg/helper"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func WaitForSnapshotCompletion(client *rds.Client, snapshotIdentifier string) error {
	fmt.Printf("  Status: %s", "[")
	for {

		resp, err := helper.GetDBSnapshot(client, snapshotIdentifier)
		if err != nil {
			return err
		}

		status := helper.GetDBSnapshotStatus(resp)

		helper.ProcessStatus(status, client, snapshotIdentifier, helper.GetDBSnapshot, helper.GetDBSnapshotStatus)

		if status == "available" {
			created := helper.GetSnapshotCreateTime(resp)
			helper.PrintUpgradeSuccess(created)
			break
		}
	}
	return nil
}
