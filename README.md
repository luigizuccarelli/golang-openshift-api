# A simple golang openshift api call interface

## Example usage

```bash

./build/api -a https://api.crc.testing:6443/apis/machineconfiguration.openshift.io/v1/machineconfigs -l info -f '.items[] | select(.metadata.name | test("^[0-9]{2}-master-kubelet$|^[-0-9]{2}-worker-kubelet$"))' -regex VersionTLS1[2-9]{1}

```

