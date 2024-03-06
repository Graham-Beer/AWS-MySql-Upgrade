package database

import (
	"time"

	"Kony/v2/pkg/helper"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func WaitForResource(client *rds.Client, identifier string, GetMessage func(), GetFn func(param *rds.DescribeDBInstancesOutput)) error {
	GetMessage()

	for {

		db, err := helper.GetDBInstance(client, identifier)
		if err != nil {
			return err
		}
		status := helper.GetDbStatus(db)

		helper.ProcessStatus(status, client, identifier, helper.GetDBInstance, helper.GetDbStatus)

		if status == "available" {
			helper.IsAvailable()
			GetFn(db)
			break
		}
		time.Sleep(10 * time.Second)
	}

	return nil
}
