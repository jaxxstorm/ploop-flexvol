# ploop-flexvol

A [FlexVolume](https://github.com/kubernetes/kubernetes/blob/master/examples/volumes/flexvolume/README.md) driver for kubernetes which allows you to mount [Ploop](https://openvz.org/Man/ploop.8) volumes to your kubernetes pods.

## Status

Kubernetes FlexVolumes are currently in Alpha state, so this plugin is as well. Use it at your own risk.

## Using

### Build

This project uses [glide](http://glide.readthedocs.io/en/latest/) so the easiest way to install your dependencies is using that.

Run `glide up` to install your dependencies.

```
$ glide up
[INFO]  Downloading dependencies. Please wait...
[INFO]  --> Fetching updates for github.com/jaxxstorm/flexvolume.
[INFO]  --> Fetching updates for github.com/urfave/cli.
[INFO]  --> Fetching updates for github.com/kolyshkin/goploop-cli.
[INFO]  --> Fetching updates for github.com/dustin/go-humanize.
[INFO]  --> Detected semantic version. Setting version for github.com/urfave/cli to v1.19.1.
[INFO]  Resolving imports
[INFO]  Downloading dependencies. Please wait...
[INFO]  Setting references for remaining imports
[INFO]  Exporting resolved dependencies...
[INFO]  --> Exporting github.com/dustin/go-humanize
[INFO]  --> Exporting github.com/jaxxstorm/flexvolume
[INFO]  --> Exporting github.com/kolyshkin/goploop-cli
[INFO]  --> Exporting github.com/urfave/cli
[INFO]  Replacing existing vendor dependencies
[INFO]  Versions did not change. Skipping glide.lock update.
[INFO]  Project relies on 4 dependencies.
```

Then build the binary:

```
go build -o ploop main.go
```

### Installing

In order to use the flexvolume driver, you'll need to install it on every node you want to use ploop on in the kubelet `volume-plugin-dir`. By default this is `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/`

You need a directory for the volume driver vendor, so create it:

```
mkdir -p /usr/libexec/kubernetes/kubelet-plugins/volume/exec/jaxxstorm~ploop
```

Then drop the binary in there:

```
mv ploop /usr/libexec/kubernetes/kubelet-plugins/volume/exec/jaxxstorm~ploop/ploop
```

You can now use ploops as usual!

### Pod Config

An example pod config would look like this:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-ploop
spec:
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: test
      mountPath: /data
    ports:
    - containerPort: 80
  nodeSelector:
    os: parallels # make sure you label your nodes to be ploop compatible 
  volumes:
  - name: test
    flexVolume:
      driver: "jaxxstorm/ploop" # this must match your vendor dir
      options:
        volumeId: "golang-ploop-test"
        size: "10G"
        volumePath: "/vstorage/storage_pool/kubernetes"
```

This will create a ploop volume `/vstorage/storage_pool/kubernetes/golang-ploop-test`. The block device which will be mounted will be at `/vstorage/storage_pool/kubernetes/golang-ploop-test/golang-ploop-test` and the `DiskDescriptor.xml` will be located at /vstorage/storage_pool/kubernetes/golang-ploop-test/DiskDescriptior.xml`

You can verify the ploop volume was created by finding the node where your pod was scheduled by running `ploop list`:

```
# ploop list
ploop18115  /vstorage/storage_pool/kubernetes/golang-ploop-test/golang-ploop-test
```



