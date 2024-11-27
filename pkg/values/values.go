package values

// BinariesCheck returns the list of expected binaries to be in a kairos system
func BinariesCheck() []string {
	return []string{
		"immucore",
		"kairos-agent",
		"grub-install|grub2-install", // same binary, different names across OSes
	}
}

type Architecture string

const (
	ArchAMD64 Architecture = "amd64"
	ArchARM64 Architecture = "arm64"
)

type Distro string

// Individual distros for when we need to be specific
const (
	Unknown            Distro = "unknown"
	Debian             Distro = "debian"
	Ubuntu             Distro = "ubuntu"
	RedHat             Distro = "redhat"
	RockyLinux         Distro = "rocky"
	AlmaLinux          Distro = "almalinux"
	Fedora             Distro = "fedora"
	Arch               Distro = "arch"
	Alpine             Distro = "alpine"
	OpenSUSELeap       Distro = "opensuse-leap"
	OpenSUSETumbleweed Distro = "opensuse-tumbleweed"
)

const (
	ImmutabilitySentinel = "/etc/kairos/.inmmutability_installed"
	KernelSentinel       = "/etc/kairos/.kernel_installed"
	InitrdSentinel       = "/etc/kairos/.initrd_installed"
)
