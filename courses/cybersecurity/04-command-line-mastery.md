# Chapter 04: Command Line Mastery — The Terminal Is Your Home

*Every security professional lives in the terminal. This chapter makes you fluent in the command line — not just comfortable.*

---

## Why the Terminal?

GUIs (graphical interfaces) are built for normal users. The terminal is built for power users.

- **Automation:** Repeat tasks in seconds with scripts
- **Remote access:** SSH into servers with no GUI
- **Precision:** Exactly what you type is exactly what runs
- **Speed:** 10x faster than clicking through menus
- **Security tools:** Almost every security tool is command-line

If you're not comfortable in the terminal, you're handicapped.

---

## Navigation and Files

```bash
# Where am I?
pwd                         # /home/hacker

# What's here?
ls                          # basic listing
ls -la                      # long format, hidden files, sizes
ls -lhS                     # sorted by size, human readable
ls -lt                      # sorted by time (newest first)
ls /etc/*.conf              # list all .conf files in /etc

# Move around
cd /var/log                 # absolute path
cd logs                     # relative path
cd ..                       # parent directory
cd ~                        # home directory
cd -                        # previous directory (toggle)

# File operations
cp source.txt dest.txt      # copy
cp -r dir/ newdir/          # copy recursively
mv old.txt new.txt          # rename
mv file.txt /tmp/           # move
rm file.txt                 # delete file
rm -rf dir/                 # DANGEROUS: delete dir and all contents
mkdir -p a/b/c              # create nested directories
touch newfile.txt           # create empty file
```

---

## Viewing Files

```bash
cat /etc/hosts              # print whole file
less /var/log/syslog        # scroll through (q to quit, / to search)
head -20 file.txt           # first 20 lines
tail -20 file.txt           # last 20 lines
tail -f /var/log/auth.log   # follow file in real time (crucial for monitoring)
wc -l file.txt              # count lines
wc -c file.txt              # count bytes

# Hex dump — view raw bytes
xxd file.bin | head -20     # see hex + ASCII
hexdump -C file.bin | head  # alternative
od -A x -t x1z file.bin     # another alternative
```

---

## Searching — The Security Power Tools

```bash
# grep — find patterns
grep "password" /etc/config          # basic search
grep -i "password" /etc/config       # case insensitive
grep -r "password" /etc/             # recursive (search all files)
grep -rn "password" /etc/            # with line numbers
grep -v "^#" /etc/config             # exclude comments
grep -E "^root|^admin" /etc/passwd   # regex

# Find files
find / -name "*.conf" 2>/dev/null    # find by name
find / -name "id_rsa" 2>/dev/null    # find SSH private keys
find / -perm -4000 -type f           # SUID files
find / -perm 777 -type f             # world-writable files
find /home -user www-data            # files owned by user
find /tmp -mtime -1                  # modified in last day

# locate — indexed search (faster, not always updated)
locate passwd
locate -i "*.conf"
updatedb                             # update the index

# which / whereis — find executables
which python3                        # /usr/bin/python3
whereis nmap                         # nmap: /usr/bin/nmap /usr/share/man/man1/nmap.1
```

---

## Text Processing — Extracting What You Need

```bash
# awk — column-based text processing
awk '{print $1}' file.txt            # print first column
awk -F: '{print $1}' /etc/passwd     # colon separator, print usernames
awk -F: '$3 == 0' /etc/passwd        # users with UID 0 (root-level!)
awk '{sum+=$3} END {print sum}' file # sum a column

# sed — stream editor (find and replace)
sed 's/old/new/g' file.txt           # replace old with new
sed 's/password=.*/REDACTED/' log    # redact passwords in logs
sed -n '10,20p' file.txt             # print lines 10-20
sed '/^#/d' config.txt               # delete comment lines

# cut — cut columns
cut -d: -f1 /etc/passwd             # first field, colon delimiter
cut -d, -f2,4 data.csv              # fields 2 and 4 from CSV

# sort and uniq — analyze lists
sort ips.txt | uniq -c | sort -rn   # count IPs, most common first
sort -t. -k1,1n -k2,2n ips.txt     # sort IP addresses numerically
```

---

## Pipes — The Unix Philosophy

Pipes (`|`) connect commands: output of one becomes input of the next.

```bash
# Find top IPs brute-forcing SSH
grep "Failed password" /var/log/auth.log \
  | awk '{print $11}' \
  | sort | uniq -c | sort -rn \
  | head -20

# Find the most common processes
ps aux | awk '{print $11}' | sort | uniq -c | sort -rn | head

# Find large files in /var
find /var -type f | xargs ls -lS | head -20

# Search all PHP files for eval() (webshell indicator)
find /var/www -name "*.php" | xargs grep -l "eval(" 2>/dev/null

# Monitor a log file for a pattern
tail -f /var/log/nginx/access.log | grep "POST /wp-login"
```

---

## Networking Commands

```bash
# Interface info
ip addr show                # all interfaces and IPs
ip addr show eth0           # specific interface
ip link show                # physical interface state
ifconfig                    # older equivalent

# Routing
ip route show               # routing table
route -n                    # alternative
ip route get 8.8.8.8       # which interface to reach 8.8.8.8?

# Open ports and connections
ss -tulpn                   # listening TCP/UDP with process
ss -tp                      # active TCP connections with process
netstat -tulpn              # older alternative
lsof -i :80                 # what's using port 80?
lsof -i tcp                 # all TCP sockets

# DNS
nslookup google.com         # basic DNS lookup
dig google.com              # detailed DNS lookup
dig MX google.com           # mail records
dig +short google.com       # just the IP
host -t NS google.com       # nameservers

# Connectivity
ping -c 4 8.8.8.8           # ICMP ping
traceroute google.com       # route to destination
mtr google.com              # live traceroute (press q to exit)
curl -I http://example.com  # HTTP headers only
curl -v https://example.com # verbose HTTP (shows TLS handshake)
wget -q -O- http://api.ipify.org  # get your public IP
nc -zv 192.168.1.1 22       # check if port is open (Netcat)
nc -l 4444                  # listen on port 4444 (simple server)
```

---

## Process Management

```bash
# View processes
ps aux                      # all processes
ps aux | grep nginx         # find nginx
top                         # live view (q to quit)
htop                        # better live view
pstree                      # process tree

# Process details
ls -la /proc/1234/          # process 1234's /proc entry
cat /proc/1234/cmdline      # full command line
cat /proc/1234/environ      # environment variables
ls /proc/1234/fd/           # open file descriptors
cat /proc/1234/maps         # memory map (useful for detecting injection)

# Manage processes
kill 1234                   # send SIGTERM (graceful)
kill -9 1234                # send SIGKILL (force)
kill -15 1234               # same as SIGTERM
killall nginx               # kill by name
pkill -f "python script.py" # kill by pattern

# Background jobs
sleep 100 &                 # run in background
jobs                        # list background jobs
fg 1                        # bring job 1 to foreground
Ctrl+Z                      # suspend current job
nohup command &             # keep running after logout
```

---

## Permissions and Users

```bash
# Who am I?
whoami                      # username
id                          # uid, gid, groups

# Change permissions
chmod 644 file.txt          # rw-r--r--
chmod +x script.sh          # add execute
chmod go-w file.txt         # remove write for group/others
chmod 4755 binary           # set SUID

# Change ownership
chown root:www-data file    # change owner and group
chown -R hacker /home/hacker  # recursive

# Privilege escalation tools
sudo -l                     # what can I run as sudo?
su - root                   # switch to root (needs password)
sudo bash                   # root shell (if sudo allowed)
sudo -u www-data bash       # shell as another user
```

---

## Shell Tricks for Speed

```bash
# History
history                     # command history
!42                         # run command #42 from history
!!                          # repeat last command
!grep                       # run last command starting with grep
Ctrl+R                      # search history interactively

# Tab completion
cd /etc/sys<TAB>            # autocomplete
git com<TAB><TAB>           # show completions

# Redirects
command > file.txt          # redirect stdout (overwrite)
command >> file.txt         # redirect stdout (append)
command 2> errors.txt       # redirect stderr
command 2>&1                # redirect stderr to stdout
command > /dev/null 2>&1    # suppress all output

# Variables
TARGET="192.168.1.1"
for port in 22 80 443; do
    nc -zv $TARGET $port
done

# Aliases (add to ~/.bashrc)
alias ll='ls -la'
alias ports='ss -tulpn'
alias myip='curl -s api.ipify.org'
```

---

## Practical Security Scenarios

```bash
# Who's currently logged in?
who
w
last                        # login history
lastb                       # failed login attempts

# Recent file changes (possible indicator of compromise)
find / -newer /tmp/reference -type f 2>/dev/null
find /home -mtime -1        # files changed in last day

# Check for unusual SUID binaries (compare against baseline)
find / -perm -4000 2>/dev/null | tee /tmp/suid_list.txt

# Check cron jobs for backdoors
for user in $(cut -f1 -d: /etc/passwd); do 
    crontab -u $user -l 2>/dev/null | grep -v "^#"
done
cat /etc/crontab
ls /etc/cron.d/

# Check network connections for suspicious outbound
ss -tp | grep ESTABLISHED | awk '{print $4}' | cut -d: -f1 | sort -u

# Look for world-writable files in PATH
for dir in $(echo $PATH | tr ':' ' '); do
    find $dir -perm -002 2>/dev/null
done

# Check /tmp for executable files (malware staging area)
find /tmp /var/tmp -executable -type f 2>/dev/null
```

---

## Summary

| Category | Key commands |
|----------|-------------|
| Navigation | `cd`, `ls -la`, `find`, `locate` |
| Viewing | `cat`, `less`, `tail -f`, `xxd` |
| Searching | `grep -r`, `find`, `locate` |
| Text processing | `awk`, `sed`, `cut`, `sort`, `uniq` |
| Networking | `ip`, `ss`, `dig`, `nc`, `curl` |
| Processes | `ps`, `top`, `lsof`, `kill` |
| Permissions | `chmod`, `chown`, `sudo -l` |

---

## Exercises

1. Find all SUID binaries on your system. Research 3 of them — could any be used for privilege escalation?
2. Write a one-liner that reads `/var/log/auth.log` and shows the top 10 IPs that failed SSH login attempts.
3. Write a script that checks if a list of ports are open on a given host, without using nmap.
4. Find all PHP files in `/var/www` that contain `eval(` — a common webshell indicator.
5. Use `proc` to find the full command line and open files of a running nginx process.
