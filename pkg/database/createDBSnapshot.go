package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func CreateDBSnapshot(client *rds.Client, instanceIdentifier, snapshotIdentifier string) (*types.DBSnapshot, error) {
	input := &rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: &instanceIdentifier,
		DBSnapshotIdentifier: &snapshotIdentifier,
	}

	output, err := client.CreateDBSnapshot(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return output.DBSnapshot, nil
}
