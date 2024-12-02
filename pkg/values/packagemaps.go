package values

import (
	"bytes"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
)
import "text/template"

// packagemaps is a map of packages to install for each distro.
// so we can deal with stupid different names between distros.

// The format is usually a map[Distro]map[Architecture][]string
// So we can store the packages for each distro and architecture independently
// Except common packages, which are named the same across all distros
// Packages can be templated, so we can pass a map of parameters to replace in the package name
// So we can transform "linux-image-generic-hwe-{{.VERSION}}" into the proper version for each ubuntu release
// the params are not hardcoded or autogenerated anywhere yet.
// Ideally the System struct should have a method to generate the params for the packages automatically
// based on the distro and version, so we can pass them to the installer without anything from our side.
// Either we set also a Common key for the common packages, or we just duplicate them for both arches if needed
//

// CommonPackages are packages that are named the same across all distros and arches
var CommonPackages = []string{
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

// ImmucorePackages are the minimum set of packages that immucore needs.
// Otherwise you wont be able to build the initrd with immucore on it.
// - dosfstools for the fat32 partition :(
// - e2fsprogs for the other partitions formatting only
// - curl is needed for livenet, which in turn is needed for kairos-network, basically for netboot!
// - parted for the partitioning used by yip
// - fdisk for the partitioning used by yip
// - gdisk for the partitioning used by yip
// - rsync for the rsync command, not sure why we need it. Immucore doesnt rsync anything I think?
// - cryptsetup for the encrypted partitions
// - lvm for rpi3 only :(
// - systemd-sysv, like wtf? we should drop that
// - cloud-guest-utils??? what is that? Drop it
// - gawk for the scripts I guess? I think some network-legacy stuff nedeed it so we need to include it
var ImmucorePackages = map[Distro]map[Architecture]map[string][]string{
	Ubuntu: {
		ArchAMD64: {
			Common: {
				"dbus", "dracut", "dracut-network", "dosfstools", "e2fsprogs", "isc-dhcp-common",
				"isc-dhcp-client", "lvm2", "curl", "parted", "fdisk", "gdisk", "rsync", "cryptsetup",
				"systemd-sysv", "cloud-guest-utils", "gawk",
			},
			">=24.04": {"dracut-live"},
			">=20.04": {"caca"},
		},
		ArchARM64: {},
	},
}

// KernelPackages is a map of packages to install for each distro.
// No arch required here, maybe models will need different packages?
var KernelPackages = map[Distro][]string{
	Ubuntu: {"linux-image-generic-hwe-{{.version}}"}, // This is a template, so we can replace the version with the actual version of the system
}

// BasePackages is a map of packages to install for each distro and architecture.
// This comprises the base packages that are needed for the system to work on a Kairos system
var BasePackages = map[Distro]map[Architecture]map[string][]string{
	Ubuntu: {
		ArchAMD64: {
			Common: {
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
				"dmsetup",
				"mdadm",
				"ncurses-term",
				"networkd-dispatcher",
				"packagekit-tools",
				"publicsuffix",
				"xdg-user-dirs",
				"xxd",
				"zerofree",
			},
			">=24:04": {
				"systemd-resolved",
			},
		},
		ArchARM64: {},
	},
	RedHat: {},
	Fedora: {},
	Alpine: {},
	Arch:   {},
	Debian: {},
}

// GrubPackages is a map of packages to install for each distro and architecture.
// TODO: Check why some packages we only install on amd64 and not on arm64?? Like neovim???
var GrubPackages = map[Distro]map[Architecture][]string{
	Ubuntu: {
		ArchAMD64: {
			"grub2",
			"grub-efi-amd64-bin",
			"grub-efi-amd64-signed",
			"grub-pc-bin",
			"coreutils",
			"grub2-common",
			"kbd",
			"lldpd",
			"neovim",
			"shim-signed",
			"snmpd",
			"squashfs-tools",
			"zfsutils-linux",
		},
		ArchARM64: {
			"grub-efi-arm64",
			"grub-efi-arm64-bin",
			"grub-efi-arm64-signed",
		},
	},
}

// SystemdPackages is a map of packages to install for each distro and architecture for systemd-boot (trusted boot) variants
// TODO: Check why some packages we only install on amd64 and not on arm64?? Like kmod???
var SystemdPackages = map[Distro]map[Architecture]map[string][]string{
	Ubuntu: {
		ArchAMD64: {
			Common: {
				"systemd",
			},
			">=24:04": {
				"iucode-tool",
				"kmod",
				"linux-base",
				"systemd-boot",
			},
		},
		ArchARM64: {
			Common: {
				"systemd",
			},
		},
	},
}

// PackageListToTemplate takes a list of packages and a map of parameters to replace in the package name
// and returns a list of packages with the parameters replaced.
func PackageListToTemplate(packages []string, params map[string]string, l sdkTypes.KairosLogger) ([]string, error) {
	var finalPackages []string
	for _, pkg := range packages {
		var result bytes.Buffer
		tmpl, err := template.New("versionTemplate").Parse(pkg)
		if err != nil {
			l.Logger.Error().Err(err).Str("package", pkg).Msg("Error parsing template.")
			return []string{}, err
		}
		err = tmpl.Execute(&result, params)
		if err != nil {
			l.Logger.Error().Err(err).Str("package", pkg).Msg("Error executing template.")
			return []string{}, err
		}
		finalPackages = append(finalPackages, result.String())
	}
	return finalPackages, nil
}
