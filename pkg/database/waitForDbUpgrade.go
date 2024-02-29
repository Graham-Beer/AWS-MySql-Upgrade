package database

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/rds"

	"Kony/v2/pkg/helper"
)

func WaitForDbUpgrade(client *rds.Client, dbinstanceIdentifier string) error {
	helper.DbUpgradeMessage()
	
	var EngVer string
	for {

		db, err := helper.GetDBInstance(client, dbinstanceIdentifier)
		if err != nil {
			return err
		}
		status := helper.GetDbStatus(db)

		helper.ProcessStatus(status, client, dbinstanceIdentifier, helper.GetDBInstance, helper.GetDbStatus)

		if status == "available" {
			EngVer = helper.GetDbVersion(db)
			helper.IsAvailable()
			break
		}
		time.Sleep(10 * time.Second)
	}
	helper.DbComplete(EngVer)

	return nil
}
