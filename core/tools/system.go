package core

import (
	"os"
	"regexp"
	"slack-wails/lib/structs"
	"slack-wails/lib/utils"
	"strings"

	"gopkg.in/yaml.v2"
)

var AuthPatchs = []structs.AuthPatch{
	{MS: "MS17-017", Patch: "KB4013081", Description: "GDIPaletteObjectsLocalPrivilegeEscalation", System: "windows7/8", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS17-017"},
	{MS: "MS17-010", Patch: "KB4013389", Description: "WindowsKernelModeDrivers", System: "windows7/2008/2003/XP", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS17-010"},
	{MS: "MS16-135", Patch: "KB3199135", Description: "WindowsKernelModeDrivers", System: "2016", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-135"},
	{MS: "MS16-111", Patch: "KB3186973", Description: "kernelapi", System: "Windows1010586(32/64)/8.1", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-111"},
	{MS: "MS16-098", Patch: "KB3178466", Description: "KernelDriver", System: "Win8.1", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-098"},
	{MS: "MS16-075", Patch: "KB3164038", Description: "HotPotato", System: "2003/2008/7/8/2012", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-075"},
	{MS: "MS16-034", Patch: "KB3143145", Description: "KernelDriver", System: "2008/7/8/10/2012", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-034"},
	{MS: "MS16-032", Patch: "KB3143141", Description: "SecondaryLogonHandle", System: "2008/7/8/10/2012", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-032"},
	{MS: "MS16-016", Patch: "KB3136041", Description: "WebDAV", System: "2008/Vista/7", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-016"},
	{MS: "MS16-014", Patch: "K3134228", Description: "remotecodeexecution", System: "2008/Vista/7", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS16-014"},
	{MS: "MS15-097", Patch: "KB3089656", Description: "remotecodeexecution", System: "win8.1/2012", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-097"},
	{MS: "MS15-076", Patch: "KB3067505", Description: "RPC", System: "2003/2008/7/8/2012", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-076"},
	{MS: "MS15-077", Patch: "KB3077657", Description: "ATM", System: "XP/Vista/Win7/Win8/2000/2003/2008/2012", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-077"},
	{MS: "MS15-061", Patch: "KB3057839", Description: "KernelDriver", System: "2003/2008/7/8/2012", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-061"},
	{MS: "MS15-051", Patch: "KB3057191", Description: "WindowsKernelModeDrivers", System: "2003/2008/7/8/2012", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-051"},
	{MS: "MS15-015", Patch: "KB3031432", Description: "KernelDriver", System: "Win7/8/8.1/2012/RT/2012R2/2008R2", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-015"},
	{MS: "MS15-010", Patch: "KB3036220", Description: "KernelDriver", System: "2003/2008/7/8", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-010"},
	{MS: "MS15-001", Patch: "KB3023266", Description: "KernelDriver", System: "2008/2012/7/8", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS15-001"},
	{MS: "MS14-070", Patch: "KB2989935", Description: "KernelDriver", System: "2003", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-070"},
	{MS: "MS14-068", Patch: "KB3011780", Description: "DomainPrivilegeEscalation", System: "2003/2008/2012/7/8", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-068"},
	{MS: "MS14-058", Patch: "KB3000061", Description: "Win32k.sys", System: "2003/2008/2012/7/8", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-058"},
	{MS: "MS14-066", Patch: "KB2992611", Description: "WindowsSchannelAllowingremotecodeexecution", System: "VistaSP2/7SP1/8/Windows8.1/2003SP2/2008SP2/2008R2SP1/2012/2012R2/WindowsRT/WindowsRT8.1", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-066"},
	{MS: "MS14-040", Patch: "KB2975684", Description: "AFDDriver", System: "2003/2008/2012/7/8", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-040"},
	{MS: "MS14-002", Patch: "KB2914368", Description: "NDProxy", System: "2003/XP", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS14-002"},
	{MS: "MS13-053", Patch: "KB2850851", Description: "win32k.sys", System: "XP/Vista/2003/2008/win7", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS13-053"},
	{MS: "MS13-046", Patch: "KB2840221", Description: "dxgkrnl.sys", System: "Vista/2003/2008/2012/7", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS13-046"},
	{MS: "MS13-005", Patch: "KB2778930", Description: "KernelModeDriver", System: "2003/2008/2012/win7/8", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS13-005"},
	{MS: "MS12-042", Patch: "KB2972621", Description: "ServiceBus", System: "2008/2012/win7", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS12-042"},
	{MS: "MS12-020", Patch: "KB2671387", Description: "RDP", System: "2003/2008/7/XP", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS12-020"},
	{MS: "MS11-080", Patch: "KB2592799", Description: "AFD.sys", System: "2003/XP", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS11-080"},
	{MS: "MS11-062", Patch: "KB2566454", Description: "NDISTAPI", System: "2003/XP", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS11-062"},
	{MS: "MS11-046", Patch: "KB2503665", Description: "AFD.sys", System: "2003/2008/7/XP", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS11-046"},
	{MS: "MS11-011", Patch: "KB2393802", Description: "kernelDriver", System: "2003/2008/7/XP/Vista", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS11-011"},
	{MS: "MS10-092", Patch: "KB2305420", Description: "TaskScheduler", System: "2008/7", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS10-092"},
	{MS: "MS10-059", Patch: "KB982799", Description: "ACL-Churraskito", System: "2008/7/Vista", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS10-059"},
	{MS: "MS10-048", Patch: "KB2160329", Description: "win32k.sys", System: "XPSP2&SP3/2003SP2/VistaSP1&SP2/2008Gold&SP2&R2/Win7", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS10-048"},
	{MS: "MS10-015", Patch: "KB977165", Description: "KiTrap0D", System: "2003/2008/7/XP", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS10-015"},
	{MS: "MS10-012", Patch: "KB971468", Description: "SMBClientTrans2stackoverflow", System: "Windows7/2008R2", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS10-012"},
	{MS: "MS09-050", Patch: "KB975517", Description: "RemoteCodeExecution", System: "2008/Vista", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS09-050"},
	{MS: "MS09-020", Patch: "KB970483", Description: "IIS6.0", System: "IIS5.1and6.0", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS09-020"},
	{MS: "MS09-012", Patch: "KB959454", Description: "Chimichurri", System: "Vista/win7/2008/Vista", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS09-012"},
	{MS: "MS08-068", Patch: "KB957097", Description: "RemoteCodeExecution", System: "2000/XP", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS08-068"},
	{MS: "MS08-067", Patch: "KB958644", Description: "RemoteCodeExecution", System: "Windows2000/XP/Server2003/Vista/Server2008", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS08-067"},
	{MS: "MS08-066", Patch: "KB956803", Description: "AFD.sys", System: "Windows2000/XP/Server2003", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS08-066"},
	{MS: "MS08-025", Patch: "KB941693", Description: "Win32.sys", System: "XP/2003/2008/Vista", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS08-025"},
	{MS: "MS06-040", Patch: "KB921883", Description: "RemoteCodeExecution", System: "2003/xp/2000", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS06-040"},
	{MS: "MS05-039", Patch: "KB899588", Description: "PnPService", System: "Win9X/ME/NT/2000/XP/2003", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS05-039"},
	{MS: "MS03-026", Patch: "KB823980", Description: "BufferOverrunInRPCInterface", System: "/NT/2000/XP/2003", Reference: "https://github.com/SecWiki/windows-kernel-exploits/tree/master/MS03-026"},
}

// 杀软识别
func (t *Tools) AntivirusIdentify(tasklist string) []structs.AntivirusResult {
	var antivirus []structs.AntivirusResult
	file := utils.HomeDir() + "/slack/config/antivirues.yaml"
	yamlData, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	data := make(map[string][]string)
	yaml.Unmarshal(yamlData, &data)
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
						antivirus = append(antivirus, structs.AntivirusResult{
							Process: p,
							Pid:     pid,
							Name:    name,
						})
					}
				}
			}
		}
	}
	return antivirus
}

// 补丁识别
func (t *Tools) PatchIdentify(systeminfo string) []structs.AuthPatch {
	var result []structs.AuthPatch
	for _, ap := range AuthPatchs {
		if !strings.Contains(systeminfo, ap.Patch) {
			result = append(result, ap)
		}
	}
	return result
}
