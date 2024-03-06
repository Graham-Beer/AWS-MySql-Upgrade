package database

import (
	"Kony/v2/pkg/helper"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func RebootAwsDatabase(client *rds.Client, instanceIdentifier string) error {
	_, err := client.RebootDBInstance(context.Background(), &rds.RebootDBInstanceInput{
		DBInstanceIdentifier: &instanceIdentifier,
	})
	if err != nil {
		return err
	}

	helper.WaitForStatus(client, "rebooting", instanceIdentifier)

	return nil
}
