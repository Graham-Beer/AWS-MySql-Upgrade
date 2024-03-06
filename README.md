# AWS-MySql-Upgrade
Using the Go AWS SDK v2, this code creates a snapshot, a db Parameter group and modify's the MySql Database version

### Command usage
Parameters:
  - instance-identifier: Database name
  - Profile: Profile name from ~/.aws/credentials file

./build/upgrade-db -instance-identifier database-1 -profile sandbox

### Console display
<img width="896" alt="Capture" src="https://github.com/Graham-Beer/AWS-MySql-Upgrade/assets/12196171/ff7ea93e-aa83-493e-af02-45dd0d16eaa7">
