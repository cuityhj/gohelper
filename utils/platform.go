package utils

import "strings"

const UNKNOWN = "unknown"

func ParsePlatform(userAgentStr string) string {
	uaLower := strings.ToLower(userAgentStr)
	switch {
	case strings.Contains(uaLower, "ios") || strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "ipod"):
		return "iOS"
	case strings.Contains(uaLower, "android"):
		return "Android"
	case strings.Contains(uaLower, "windows"):
		return "Windows"
	case strings.Contains(uaLower, "macintosh") || strings.Contains(uaLower, "mac os"):
		return "MacOS"
	case strings.Contains(uaLower, "ubuntu"):
		return "Ubuntu"
	case strings.Contains(uaLower, "centos"):
		return "CentOS"
	case strings.Contains(uaLower, "fedora"):
		return "Fedora"
	case strings.Contains(uaLower, "linux"):
		return "Linux"
	case strings.Contains(uaLower, "chromeos") || strings.Contains(uaLower, "chromiumos"):
		return "ChromeOS"
	case strings.Contains(uaLower, "freebsd"):
		return "FreeBSD"
	case strings.Contains(uaLower, "netbsd"):
		return "NetBSD"
	case strings.Contains(uaLower, "openbsd"):
		return "OpenBSD"
	case strings.Contains(uaLower, "harmony"):
		return "HarmonyOS"
	case strings.Contains(uaLower, "velaos") || strings.Contains(uaLower, "vela"):
		return "VelaOS"
	case strings.Contains(uaLower, "blueos"):
		return "BlueOS"
	case strings.Contains(uaLower, "pantanal"):
		return "Pantanal"
	case strings.Contains(uaLower, "tizen"):
		return "Tizen"
	case strings.Contains(uaLower, "playstation"):
		return "PlayStation"
	case strings.Contains(uaLower, "xbox"):
		return "XBox"
	case strings.Contains(uaLower, "nintendo"):
		return "Nintendo"
	default:
		return UNKNOWN
	}
}
