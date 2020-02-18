package main

import (
	"context"
	"fmt"

	lvm "github.com/mwennrich/csi-driver-lvm/pkg/lvm"
	"github.com/urfave/cli/v2"
	"k8s.io/klog"
)

const (
	linearType  = "linear"
	stripedType = "striped"
	mirrorType  = "mirror"
)

func createLVCmd() *cli.Command {
	return &cli.Command{
		Name: "createlv",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  flagLVName,
				Usage: "Required. Specify lv name.",
			},
			&cli.Uint64Flag{
				Name:  flagLVSize,
				Usage: "Required. The size of the lv in MiB",
			},
			&cli.StringFlag{
				Name:  flagVGName,
				Usage: "Required. the name of the volumegroup",
			},
			&cli.StringFlag{
				Name:  flagLVMType,
				Usage: "Required. type of lvs, can be either striped or mirrored",
			},
			&cli.StringSliceFlag{
				Name:  flagDevicesPattern,
				Usage: "Required. the patterns of the physical volumes to use.",
			},
		},
		Action: func(c *cli.Context) error {
			if err := createLV(c); err != nil {
				klog.Fatalf("Error creating lv: %v", err)
				return err
			}
			return nil
		},
	}
}

func createLV(c *cli.Context) error {
	lvName := c.String(flagLVName)
	if lvName == "" {
		return fmt.Errorf("invalid empty flag %v", flagLVName)
	}
	lvSize := c.Uint64(flagLVSize)
	if lvSize == 0 {
		return fmt.Errorf("invalid empty flag %v", flagLVSize)
	}
	vgName := c.String(flagVGName)
	if vgName == "" {
		return fmt.Errorf("invalid empty flag %v", flagVGName)
	}
	devicesPattern := c.StringSlice(flagDevicesPattern)
	if len(devicesPattern) == 0 {
		return fmt.Errorf("invalid empty flag %v", flagDevicesPattern)
	}
	lvmType := c.String(flagLVMType)
	if lvmType == "" {
		return fmt.Errorf("invalid empty flag %v", flagLVMType)
	}

	klog.Infof("create lv %s size:%d vg:%s devicespattern:%s  type:%s", lvName, lvSize, vgName, devicesPattern, lvmType)

	// TODO
	// createVG could get called once at the start of the nodeserver
	output, err := lvm.CreateVG(vgName, devicesPattern)
	if err != nil {
		return fmt.Errorf("unable to create vg: %v output:%s", err, output)
	}

	output, err = lvm.CreateLVS(context.Background(), vgName, lvName, lvSize, lvmType)
	if err != nil {
		return fmt.Errorf("unable to create lv: %v output:%s", err, output)
	}
	return nil
}

// TODO
// move everything below to lvm package
// ephemeral volumes can be created directly on the node without a provisioner pod,
// so these functions are needed there too anyway
