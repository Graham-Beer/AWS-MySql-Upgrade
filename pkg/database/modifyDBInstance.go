package database

import (
	"Kony/v2/pkg/helper"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func ModifyDBInstance(client *rds.Client, instanceIdentifier, newEngineVersion, parameterGroupName string) error {
	params := &rds.ModifyDBInstanceInput{
		AllowMajorVersionUpgrade: aws.Bool(true),
		ApplyImmediately:         aws.Bool(true),
		DBInstanceIdentifier:     &instanceIdentifier,
		EngineVersion:            &newEngineVersion,
		DBParameterGroupName:     &parameterGroupName,
	}

	_, err := client.ModifyDBInstance(context.Background(), params)

	for {
		output, err := helper.GetDBInstance(client, instanceIdentifier)
		if err != nil {
			return err
		}
		status := helper.GetDbStatus(output)
		if status == "upgrading" {
			break
		}

		time.Sleep(100 * time.Second)
	}

	return err
}
