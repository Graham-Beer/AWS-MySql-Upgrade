package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func CreateDBParameterGroup(client *rds.Client, groupName, family, description string) error {
	params := &rds.CreateDBParameterGroupInput{
		DBParameterGroupFamily: &family,
		DBParameterGroupName:   &groupName,
		Description:            &description,
	}

	_, err := client.CreateDBParameterGroup(context.Background(), params)
	return err
}
