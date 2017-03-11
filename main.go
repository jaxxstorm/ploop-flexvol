package main

import (
	"os"
	"syscall"

	"github.com/avagin/ploop-flexvol/volume"
	"github.com/jaxxstorm/flexvolume"
	"github.com/kolyshkin/goploop-cli"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ploop flexvolume"
	app.Usage = "Mount ploop volumes in kubernetes using the flexvolume driver"
	app.Commands = flexvolume.Commands(Ploop{})
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Lee Briggs",
		},
	}
	app.Version = "0.1a"
	app.Run(os.Args)
}

type Ploop struct{}

func (p Ploop) Init() flexvolume.Response {
	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Ploop is available",
	}
}

func (p Ploop) Attach(options map[string]string) flexvolume.Response {

	if options["volumePath"] == "" {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: "Must specify a volume path",
		}
	}

	if options["volumeId"] == "" {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: "Must specify a volume id",
		}
	}

	if _, err := os.Stat(options["volumePath"] + "/" + options["volumeId"] + "/" + "DiskDescriptor.xml"); err == nil {
		return flexvolume.Response{
			Status:  flexvolume.StatusSuccess,
			Message: "Volume already exists",
			Device:  options["volumePath"] + "/" + options["volumeId"],
		}
	}

	if err := volume.Create(options); err != nil {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: err.Error(),
		}
	}

	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Successfully attached the ploop volume",
		Device:  options["volumePath"] + "/" + options["volumeId"] + "/" + options["volumeId"],
	}
}

func (p Ploop) Detach(device string) flexvolume.Response {

	if err := ploop.UmountByDevice(device); err != nil {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: "Unable to detach ploop volume",
			Device:  device,
		}
	}
	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Successfully detached the ploop volume",
		Device:  device,
	}
}

func (p Ploop) Mount(target, device string, options map[string]string) flexvolume.Response {
	// make the target directory we're going to mount to
	err := os.MkdirAll(target, 0755)
	if err != nil {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: err.Error(),
			Device:  device,
		}
	}

	// open the disk descriptor first
	volume, err := ploop.Open(options["volumePath"] + "/" + options["volumeId"] + "/" + "DiskDescriptor.xml")
	if err != nil {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: err.Error(),
			Device:  device,
		}
	}
	defer volume.Close()

	if m, _ := volume.IsMounted(); !m {
		// If it's mounted, let's mount it!

		readonly := false
		if options["kubernetes.io/readwrite"] == "ro" {
			readonly = true
		}

		mp := ploop.MountParam{Target: target, Readonly: readonly}

		dev, err := volume.Mount(&mp)
		if err != nil {
			return flexvolume.Response{
				Status:  flexvolume.StatusFailure,
				Message: err.Error(),
				Device:  dev,
			}
		}

		return flexvolume.Response{
			Status:  flexvolume.StatusSuccess,
			Message: "Successfully mounted the ploop volume",
			Device:  device,
		}
	} else {

		return flexvolume.Response{
			Status:  flexvolume.StatusSuccess,
			Message: "Ploop volume already mounted",
			Device:  device,
		}

	}
}

func (p Ploop) Unmount(mount string) flexvolume.Response {

	if err := syscall.Unmount(mount, 0x2); err != nil {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: "Unable to unmount ploop volume from pod",
			Device:  mount,
		}
	}

	if err := os.Remove(mount); err != nil {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: "Unable to remove stale directory from pod",
			Device:  mount,
		}
	}

	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Successfully unmounted the ploop volume",
		Device:  mount,
	}
}
