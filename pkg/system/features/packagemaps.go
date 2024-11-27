package features

import "github.com/kairos-io/kairos-init/pkg/values"

// Immutability is a map of packages to install for each distro.
// so we can deal with stupid different names between distros.

// Immutability that are named the same across all distros
var commonPackages = []string{
	"curl",
	"file",
	"gawk",
	"iptables",
	"less",
	"nano",
	"sudo",
	"tar",
	"zstd",
	"rsync",
	"systemd",
	"lvm2",
	"jq",
	"dosfstools",
	"e2fsprogs",
	"parted",
}

// we need:
// - grub2 for the bootloader
// - linux-image-generic for the kernel
// - dracut for the initrd to generate one with immucore on it
// - dosfstools for the fat32 partition :(
// - e2fsprogs for the other partitions formatting only
// Ideally for ubuntu:
// Get the actual version of the kernel from linux-image0-generic
// Install ONLY that kernel image
// That saves 400Mb as it doesnt bring any other stuff like firmware and extra modules
// Then we have that as an extra feature or whatever so we can install it if needed (uki slim vs fat)
// curl is needed for livenet, which in turn is needed for kairos-network
var immucorePackages = map[values.Distro][]string{
	values.Ubuntu: {
		"dbus", "dracut", "dracut-network", "dracut-live", "dosfstools", "e2fsprogs", "isc-dhcp-common",
		"isc-dhcp-client", "lvm2", "curl", "parted", "fdisk", "gdisk", "rsync", "cryptsetup", "ca-certificates",
		"systemd-sysv", "cloud-guest-utils", "gawk",
	},
}

var kernelPackages = map[values.Distro][]string{
	values.Ubuntu: {"linux-image-generic"},
}

var basePackages = map[values.Distro][]string{
	values.Debian: {"grub2"},
	values.Ubuntu: append([]string{
		"gdisk",
		"fdisk",
		"ca-certificates",
		"conntrack",
		"console-data",
		"cloud-guest-utils",
		"cryptsetup",
		"debianutils",
		"gettext",
		"haveged",
		"iproute2",
		"iputils-ping",
		"krb5-locales",
		"nbd-client",
		"nfs-common",
		"open-iscsi",
		"open-vm-tools",
		"openssh-server",
		"systemd-timesyncd",
		"systemd-container",
		"ubuntu-advantage-tools",
		"xz-utils",
		"tpm2-tools",
	}, commonPackages...),
	values.RedHat: {"grub2"},
	values.Fedora: {"grub2"},
	values.Alpine: {"grub2"},
	values.Arch:   {"grub2"},
}

// grubPackages is a map of packages to install for each distro and architecture.
var grubPackages = map[values.Distro]map[values.Architecture][]string{
	values.Ubuntu: {
		values.ArchAMD64: {
			"grub2",
			"grub-efi-amd64-bin",
			"grub-efi-amd64-signed",
			"grub-pc-bin",
		},
		values.ArchARM64: {
			"grub-efi-arm64",
			"grub-efi-arm64-bin",
			"grub-efi-arm64-signed",
		},
	},
}

var systemdPackages = map[values.Distro]map[values.Architecture][]string{
	values.Ubuntu: {
		values.ArchAMD64: {
			"systemd",
		},
		values.ArchARM64: {
			"systemd",
		},
	},
}
