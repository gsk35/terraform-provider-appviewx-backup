# terraform-provider-appviewx

## Build
```
> cd ../terraform-provider-appviewx
> go build -o terraform-provider-appviewx
```

## Enable logs  ( TRACE, DEBUG, INFO, WARN or ERROR )
```
> export TF_LOG=TRACE
```

## Sample ' version.tf ' file
```
terraform {
  required_providers {
    appviewx = {
      version = "0.2"
      source  = "appviewx.com/provider/appviewx"
    }
  }
}
```
## Sample ' appviewx.tf'  file
```
provider "appviewx"{
  appviewx_username="USER_NAME"
	appviewx_password="PASSWORD"
	appviewx_environment_is_https=true 
	appviewx_environment_ip="APPVIEWX_HOST_NAME"
	appviewx_environment_port="APPVIEWX_PORT_NUMBER"
}

resource "appviewx_automation" "newcert"{
 payload= <<EOF
 {
  "payload" : {
    "data" : {
      "input" : {
        "requestData" : [ {
          "sequenceNo" : 1,
          "scenario" : "scenario",
          "fieldInfo" : {
            "commonname" : "www.app1.company.com",
            "email" : "name@company.com"
          }
        } ]
      },
      "task_action" : 1
    },
    "header" : {
      "workflowName" : "Generate New CSR"
    }
  }
}
EOF
action_id= "visualworkflow-submit-request"

  }
```

## Steps to run
```
> Keep the .tf files in the current folder

> keep the "terraform-provider-appviewx" binary file under "~/.terraform.d/plugins/appviewx.com/provider/appviewx/0.2/linux_386"   ( linux_386 is sample, need to change based on the installation system architecture )

> Run the following commands, to reset and trigger the request

	rm -rf ./terraform.tfstate;
	terraform init;
	terraform apply;
```