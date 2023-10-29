# iam-disable

Use this tool to disable an IAM user's console access AND each of their access keys.  

NOTE: make it easy tio disable additional  components like 
 - SSH public keys for AWS CodeCommit
 - HTTPS Git credentials for AWS CodeCommit

## Usage

### Help
To see the help text, run the tool with the --help flag

```
./build/current/linux/amd64/iam-disable [ -h | --help | help ]
Usage:
  iam-disable [target file path]

If run with no arguments, discover IAM user accounts and create a targets file if one doesn't
exist. example: 0123456789_targets.txt

create/overwrite a report file. example: 0123456789_report.txt

If run with one argument (target file path), disable the IAM users in the targets file.
```

### Discover IAM users
To discover the IAM users in your account, set your AWS credentials and run to tool with no arguments:
NOTE: It's normal to see errors about users not having a console profile.  This is expected for users that do not have console access.
```
./build/current/linux/amd64/iam-disable
{"level":"info","version":"00574cf2e6f9121e32f01a330533e09c1e55cd1d","account":"0123456789","arn":"arn:aws:iam::0123456789:user/nmarks","user":"AIAIAIAIAI","mode":"discover","time":"2023-10-29T06:39:14-04:00","message":"starting"}
{"level":"error","version":"00574cf2e6f9121e32f01a330533e09c1e55cd1d","account":"0123456789","arn":"arn:aws:iam::0123456789:user/nmarks","user":"AIAIAIAIAI","mode":"discover","error":"operation error IAM: GetLoginProfile, https response error StatusCode: 404, RequestID: 7f4a8c4c-3e6d-4bdb-b19e-2025237eda20, NoSuchEntity: Login Profile for User FAKEUSER1 cannot be found.","time":"2023-10-29T06:39:14-04:00","message":"error checking for console profile: FAKEUSER1"}
{"level":"error","version":"00574cf2e6f9121e32f01a330533e09c1e55cd1d","account":"0123456789","arn":"arn:aws:iam::0123456789:user/nmarks","user":"AIAIAIAIAI","mode":"discover","error":"operation error IAM: GetLoginProfile, https response error StatusCode: 404, RequestID: 869432bc-8acb-4ad9-b6e8-95438ad6560d, NoSuchEntity: Login Profile for User FAKEUSER2 cannot be found.","time":"2023-10-29T06:39:16-04:00","message":"error checking for console profile: FAKEUSER2"}

```

The toole will create two files:
 - 0123456789_report.txt
 - 0123456789_targets.txt

Review the report file. Edit the targets file to remove any users you do not want to disable. 


### Disable IAM users
To disable IAM users, edit the targets list to include only the users you want to disable.  Then run the tool with the targets file as the only argument.

NOTE: you'll be prompted to enter 4 random characters to confirm the disable operation.  This is to prevent accidental execution of the tool.
```
./build/current/linux/amd64/iam-disable 0123456789_targets.txt
```