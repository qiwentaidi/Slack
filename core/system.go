package core

import (
	"os"
	"regexp"
	"slack-wails/lib/util"
	"strings"

	"gopkg.in/yaml.v2"
)

type AuthPatch struct {
	Cve               string
	Patch             string
	Description       string
	System            string
	VulnerabilityPath string
}

var AuthPatchs = []AuthPatch{
	{"MS17-017", "KB4013081", "GDIPaletteObjectsLocalPrivilegeEscalation", "windows7/8", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS17-017"},
	{"MS17-010", "KB4013389", "WindowsKernelModeDrivers", "windows7/2008/2003/XP", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS17-010"},
	{"MS16-135", "KB3199135", "WindowsKernelModeDrivers", "2016", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-135"},
	{"MS16-111", "KB3186973", "kernelapi", "Windows1010586(32/64)/8.1", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-111"},
	{"MS16-098", "KB3178466", "KernelDriver", "Win8.1", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-098"},
	{"MS16-075", "KB3164038", "HotPotato", "2003/2008/7/8/2012", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-075"},
	{"MS16-034", "KB3143145", "KernelDriver", "2008/7/8/10/2012", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-034"},
	{"MS16-032", "KB3143141", "SecondaryLogonHandle", "2008/7/8/10/2012", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-032"},
	{"MS16-016", "KB3136041", "WebDAV", "2008/Vista/7", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-016"},
	{"MS16-014", "K3134228", "remotecodeexecution", "2008/Vista/7", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-014"},
	{"MS15-097", "KB3089656", "remotecodeexecution", "win8.1/2012", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-097"},
	{"MS15-076", "KB3067505", "RPC", "2003/2008/7/8/2012", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-076"},
	{"MS15-077", "KB3077657", "ATM", "XP/Vista/Win7/Win8/2000/2003/2008/2012", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-077"},
	{"MS15-061", "KB3057839", "KernelDriver", "2003/2008/7/8/2012", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-061"},
	{"MS15-051", "KB3057191", "WindowsKernelModeDrivers", "2003/2008/7/8/2012", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-051"},
	{"MS15-015", "KB3031432", "KernelDriver", "Win7/8/8.1/2012/RT/2012R2/2008R2", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-015"},
	{"MS15-010", "KB3036220", "KernelDriver", "2003/2008/7/8", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-010"},
	{"MS15-001", "KB3023266", "KernelDriver", "2008/2012/7/8", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-001"},
	{"MS14-070", "KB2989935", "KernelDriver", "2003", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-070"},
	{"MS14-068", "KB3011780", "DomainPrivilegeEscalation", "2003/2008/2012/7/8", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-068"},
	{"MS14-058", "KB3000061", "Win32k.sys", "2003/2008/2012/7/8", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-058"},
	{"MS14-066", "KB2992611", "WindowsSchannelAllowingremotecodeexecution", "VistaSP2/7SP1/8/Windows8.1/2003SP2/2008SP2/2008R2SP1/2012/2012R2/WindowsRT/WindowsRT8.1", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-066"},
	{"MS14-040", "KB2975684", "AFDDriver", "2003/2008/2012/7/8", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-040"},
	{"MS14-002", "KB2914368", "NDProxy", "2003/XP", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-002"},
	{"MS13-053", "KB2850851", "win32k.sys", "XP/Vista/2003/2008/win7", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS13-053"},
	{"MS13-046", "KB2840221", "dxgkrnl.sys", "Vista/2003/2008/2012/7", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS13-046"},
	{"MS13-005", "KB2778930", "KernelModeDriver", "2003/2008/2012/win7/8", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS13-005"},
	{"MS12-042", "KB2972621", "ServiceBus", "2008/2012/win7", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS12-042"},
	{"MS12-020", "KB2671387", "RDP", "2003/2008/7/XP", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS12-020"},
	{"MS11-080", "KB2592799", "AFD.sys", "2003/XP", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS11-080"},
	{"MS11-062", "KB2566454", "NDISTAPI", "2003/XP", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS11-062"},
	{"MS11-046", "KB2503665", "AFD.sys", "2003/2008/7/XP", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS11-046"},
	{"MS11-011", "KB2393802", "kernelDriver", "2003/2008/7/XP/Vista", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS11-011"},
	{"MS10-092", "KB2305420", "TaskScheduler", "2008/7", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS10-092"},
	{"MS10-059", "KB982799", "ACL-Churraskito", "2008/7/Vista", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS10-059"},
	{"MS10-048", "KB2160329", "win32k.sys", "XPSP2&SP3/2003SP2/VistaSP1&SP2/2008Gold&SP2&R2/Win7", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS10-048"},
	{"MS10-015", "KB977165", "KiTrap0D", "2003/2008/7/XP", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS10-015"},
	{"MS10-012", "KB971468", "SMBClientTrans2stackoverflow", "Windows7/2008R2", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS10-012"},
	{"MS09-050", "KB975517", "RemoteCodeExecution", "2008/Vista", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS09-050"},
	{"MS09-020", "KB970483", "IIS6.0", "IIS5.1and6.0", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS09-020"},
	{"MS09-012", "KB959454", "Chimichurri", "Vista/win7/2008/Vista", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS09-012"},
	{"MS08-068", "KB957097", "RemoteCodeExecution", "2000/XP", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS08-068"},
	{"MS08-067", "KB958644", "RemoteCodeExecution", "Windows2000/XP/Server2003/Vista/Server2008", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS08-067"},
	{"MS08-066", "KB956803", "AFD.sys", "Windows2000/XP/Server2003", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS08-066"},
	{"MS08-025", "KB941693", "Win32.sys", "XP/2003/2008/Vista", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS08-025"},
	{"MS06-040", "KB921883", "RemoteCodeExecution", "2003/xp/2000", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS06-040"},
	{"MS05-039", "KB899588", "PnPService", "Win9X/ME/NT/2000/XP/2003", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS05-039"},
	{"MS03-026", "KB823980", "BufferOverrunInRPCInterface", "/NT/2000/XP/2003", "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS03-026"},
}

func AntivirusIdentify(tasklist string) ([][]string, error) {
	var AntivirusResult [][]string
	yamlData, err := os.ReadFile(util.HomeDir() + "/slack/antivirues.yaml")
	if err != nil {
		return nil, err
	}
	data := make(map[string][]string)
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, err
	}
	// 换行读取每列，然后做指纹比对
	for _, line := range strings.Split(tasklist, "\n") {
		re := regexp.MustCompile(`^(\S+)\s+(\d+)`)
		result := re.FindStringSubmatch(line)
		if len(result) >= 3 {
			processName := result[1]
			pid := result[2]
			for name, process := range data {
				for _, p := range process {
					if processName == p {
						AntivirusResult = append(AntivirusResult, []string{p, pid, name})
					}
				}
			}
		}
	}
	return AntivirusResult, nil
}

func Patch(systeminfo string) [][]string {
	var AuthPatchResult [][]string
	for _, ap := range AuthPatchs {
		if !strings.Contains(systeminfo, ap.Patch) {
			AuthPatchResult = append(AuthPatchResult, []string{ap.Cve, ap.Patch, ap.Description, ap.System, ap.VulnerabilityPath})
		}
	}
	return AuthPatchResult
}
