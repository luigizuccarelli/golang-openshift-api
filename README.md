# A simple golang openshift api call interface

## Example usage

```bash

export BASE_URL=https://api.crc.testing:6443
export HOME_DIR=/home/<user>

# build the binary
$ make build

./build/api -l debug -c input-rules.json

```

## Example input-rules.json file

```json
[
  {
    "name": "check disk_encryption",
    "api": "/apis/machine.openshift.io/v1beta1/machinesets?limit=500",
    "filter": "[.items[] | .spec.template.spec.providerSpec.value.blockDevices[0].ebs.encrypted] | map(. == true)",
    "regex": "true",
    "occurence": "all",
    "operand": "eq"
  },
  {
    "name": "check fips",
    "api": "/apis/machineconfiguration.openshift.io/v1/machineconfigs",
    "filter": "[.items[] | select(.metadata.name | test(\"^[0-9]{2}-worker$|^[0-9]{2}-master$\"))]|map(.spec.fips == false)",
    "regex": "true",
    "occurence": "2",
    "operand": "eq"
  }
]

```


