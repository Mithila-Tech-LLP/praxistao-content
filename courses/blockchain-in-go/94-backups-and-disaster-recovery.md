# Chapter 94: Backups and Disaster Recovery

Every chapter so far in this volume has assumed the happy path: the VM stays up, the disk stays healthy, the containers keep running. Real disks fail. Real VMs get accidentally deleted from a cloud dashboard. A `docker compose down -v` typed in the wrong terminal window can erase a node's entire BoltDB file in half a second. This chapter builds the safety net for all of that: scheduled backups of the data that actually matters, a restore procedure you have personally run at least once before you ever need it for real, and a chaos test that proves — rather than assumes — that a "restored" node genuinely rejoins the network and catches back up.

## Table of Contents

1. [What Actually Needs Backing Up](#1-what-actually-needs-backing-up)
2. [The 3-2-1 Backup Principle](#2-the-3-2-1-backup-principle)
3. [Backing Up BoltDB Safely, While Live](#3-backing-up-boltdb-safely-while-live)
4. [The Backup Script](#4-the-backup-script)
5. [Scheduling Backups with Cron](#5-scheduling-backups-with-cron)
6. [Backing Up Wallet Files](#6-backing-up-wallet-files)
7. [Shipping Backups Off the Server](#7-shipping-backups-off-the-server)
8. [The Restore Procedure](#8-the-restore-procedure)
9. [The Chaos Test](#9-the-chaos-test)
10. [Running the Chaos Test](#10-running-the-chaos-test)
11. [Why an Untested Backup Is Not a Backup](#11-why-an-untested-backup-is-not-a-backup)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. What Actually Needs Backing Up

Before writing any scripts, it is worth being precise about what would actually be lost if the Chapter 88 VM's disk died right now, because not everything is equally irreplaceable:

- **Each node's BoltDB file** (`/app/data`, from Chapter 54 onward) — the entire chain history, the UTXO index, and any locally cached state. This is technically *recoverable* from peers via Chapter 49's synchronization, as long as at least one other honest node in the network still has it — but if you are running the network's only nodes (an early testnet, before Chapter 90's public participants join), losing every node's data at once means losing the chain permanently.
- **The faucet's wallet file** (Chapter 90) — this is genuinely irreplaceable. A private key is not derived from anything else in the system; if it is lost, whatever funds it controlled are gone forever, exactly the way losing a real Bitcoin wallet's seed phrase means losing the coins permanently, with no customer support line to call.
- **Grafana's dashboard configuration and Prometheus's historical metrics** (Chapter 91) — nice to have, genuinely useful for spotting long-term trends, but their loss would be an inconvenience, not a catastrophe; dashboards can be rebuilt from the JSON committed to your repository, and metric history simply starts over.
- **Everything else** — the `Dockerfile`s, `docker-compose.yml`, the `Caddyfile`, the GoChain source itself — already lives in Git (Chapter 92's CI/CD assumes exactly this), and is therefore already backed up, for free, by every `git push` you have ever run.

This chapter focuses its effort where it matters most: the BoltDB data directories and the wallet files, in that priority order.

---

## 2. The 3-2-1 Backup Principle

A widely used rule of thumb for backup strategy, worth internalizing once and reusing for the rest of your career, not just this project: keep **3** copies of your data, on **2** different types of storage media, with **1** copy stored somewhere physically or geographically separate from the original.

```
   3 copies                2 storage types              1 offsite
  +-------------+         +-----------------+          +--------------+
  | 1. the live  |         | disk on the VM   |          | a location    |
  |    data on   |         | (fast, but same  |          | that survives  |
  |    the VM    |  --->   | machine failure  |  --->    | the VM/data    |
  | 2. a local   |         | destroys both    |          | center itself  |
  |    backup    |         | copies at once)  |          | disappearing   |
  |    file      |         +-----------------+          +--------------+
  | 3. an        |
  |    offsite   |
  |    copy      |
  +-------------+
```

Applied to GoChain specifically: the live BoltDB file on the VM is copy one; a compressed backup archive saved to a *different* disk (or at minimum a different directory tree, ideally an entirely separate storage volume) on the same VM is copy two; and a copy uploaded to object storage (or simply downloaded to your own laptop) run by an entirely different provider is copy three — the one that survives even a total, permanent loss of the VM and everything on it.

---

## 3. Backing Up BoltDB Safely, While Live

The naive approach — `cp` the BoltDB file while the node process is still running and actively writing to it — is unsafe. BoltDB, like any real database, may have a write in progress at the exact instant you copy the file, producing a backup that is truncated or internally inconsistent, the same category of problem Chapter 53 raised about flat files and crash-safety in the first place.

BoltDB's own Go API includes exactly the tool for this: `DB.View`, wrapped around `tx.WriteTo`, takes a **consistent, point-in-time snapshot** of the database from inside a read transaction, guaranteeing the resulting bytes represent the database exactly as it existed at one single, coherent moment — never a half-written, in-progress state:

```go
// gochain/storage/backup.go
//
// Backup writes a consistent, point-in-time snapshot of the BoltDB
// database to w. Because it runs inside a read-only transaction,
// concurrent writers to the live database are completely unaffected -
// they simply continue writing to a version of the data this snapshot
// will not see, rather than being blocked or corrupted.
package storage

import (
	"io"

	bolt "go.etcd.io/bbolt"
)

func Backup(db *bolt.DB, w io.Writer) error {
	return db.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(w)
		return err
	})
}
```

This is the exact mechanism a small companion CLI, `gochain-backup`, wraps into a standalone tool usable from a shell script without needing to embed this logic in the node process itself:

```go
// cmd/gochain-backup/main.go
package main

import (
	"flag"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
	"github.com/you/gochain/storage"
)

func main() {
	dbPath := flag.String("db", "/app/data/gochain.db", "path to the live BoltDB file")
	outPath := flag.String("out", "", "path to write the backup snapshot to")
	flag.Parse()

	if *outPath == "" {
		log.Fatal("gochain-backup: --out is required")
	}

	// Opened read-only and with a short timeout: if the live node has
	// the file locked in a way that would block, fail fast rather than
	// hanging a scheduled cron job indefinitely.
	db, err := bolt.Open(*dbPath, 0444, &bolt.Options{ReadOnly: true, Timeout: 5e9})
	if err != nil {
		log.Fatalf("gochain-backup: failed to open db: %v", err)
	}
	defer db.Close()

	out, err := os.Create(*outPath)
	if err != nil {
		log.Fatalf("gochain-backup: failed to create output file: %v", err)
	}
	defer out.Close()

	if err := storage.Backup(db, out); err != nil {
		log.Fatalf("gochain-backup: backup failed: %v", err)
	}
	log.Printf("gochain-backup: wrote consistent snapshot to %s", *outPath)
}
```

Opening the database with `ReadOnly: true` from a *second* process, alongside the live node process that already has it open for writing, is safe specifically because BoltDB supports multiple concurrent readers alongside a single writer — the same single-writer/multiple-reader model briefly mentioned back in Chapter 54's comparison of BoltDB against Badger.

---

## 4. The Backup Script

A shell script ties the pieces together: run `gochain-backup` for each node's data directory, compress the result, and stamp it with a timestamp so multiple generations of backups can coexist without overwriting each other:

```bash
#!/usr/bin/env bash
# /opt/gochain/scripts/backup.sh
#
# Snapshots every node's BoltDB file and the faucet wallet, compresses
# each into a timestamped archive, and prunes anything older than 14
# days so the backup directory does not grow forever.
set -euo pipefail

BACKUP_DIR="/opt/gochain/backups"
TIMESTAMP="$(date +%Y%m%d-%H%M%S)"
RETAIN_DAYS=14

mkdir -p "$BACKUP_DIR"

# Each node's data lives in its own named Docker volume (Chapter 87);
# `docker compose exec` runs gochain-backup *inside* the running
# container, where the volume is actually mounted, and streams the
# result out to the host's filesystem.
for node in node1 node2 node3; do
  snapshot="${BACKUP_DIR}/${node}-${TIMESTAMP}.db"
  docker compose exec -T "$node" \
    gochain-backup --db /app/data/gochain.db --out /tmp/snapshot.db
  docker compose cp "${node}:/tmp/snapshot.db" "$snapshot"
  gzip "$snapshot"
  echo "backed up ${node} -> ${snapshot}.gz"
done

# The faucet's wallet file is small and simple enough to copy directly -
# it is a single encrypted JSON file (Chapter 40), not a live database,
# so there is no consistency concern to worry about here.
cp /opt/gochain/faucet-wallet.json \
   "${BACKUP_DIR}/faucet-wallet-${TIMESTAMP}.json"

# Prune backups older than RETAIN_DAYS, so this directory does not
# silently consume the entire disk over months of daily cron runs.
find "$BACKUP_DIR" -type f -mtime "+${RETAIN_DAYS}" -delete

echo "backup complete: ${TIMESTAMP}"
```

`set -euo pipefail` is a good habit for every operational shell script from here on: `-e` exits immediately on any command failure rather than plowing ahead with a broken backup, `-u` treats an unset variable as an error rather than silently substituting an empty string, and `-o pipefail` makes a pipeline fail if *any* stage of it fails, not just the last one.

---

## 5. Scheduling Backups with Cron

**Cron** is the standard Linux tool for running a command on a recurring schedule, configured through a **crontab** (a small text file of schedule-plus-command lines). Add one line to run the backup script every night:

```bash
crontab -e
```

```cron
# minute hour day month weekday   command
# Run every day at 3:15 AM server time - chosen deliberately off the
# hour, since a huge number of scheduled jobs across the internet fire
# exactly on the hour, and staggering slightly avoids competing for
# disk I/O with anything else on a shared host.
15 3 * * * /opt/gochain/scripts/backup.sh >> /opt/gochain/backups/backup.log 2>&1
```

Redirecting both stdout and stderr (`>> ... 2>&1`) into a log file is essential for a cron job specifically — cron jobs run with no terminal attached, so any output that is not explicitly redirected is silently discarded, which is exactly how a failing backup script can go unnoticed for months until the day you actually need a backup and discover there are none.

---

## 6. Backing Up Wallet Files

Section 4's script already copies the faucet's wallet file, but the same discipline applies to any wallet involved in running the testnet — recall Chapter 40's point that a wallet file is encrypted at rest with a password-derived key, so the backup itself does not need additional encryption on top of that (the file is already safe to store, as long as the password used to encrypt it is not stored alongside it). What the backup absolutely must preserve is the wallet file **and**, separately and never in the same location, the password (or, for an HD wallet from Chapter 38, the seed phrase) needed to decrypt it — losing either one alone is as good as losing both.

A reasonable split: the encrypted wallet file goes into the same backup rotation as everything else in this chapter; the password or seed phrase is written down and stored somewhere entirely separate — a password manager, a physical safe, whatever your personal threat model calls for — precisely because a backup script that stored both together would defeat the entire purpose of encrypting the wallet file in the first place.

---

## 7. Shipping Backups Off the Server

A backup that only ever lives on the same VM's disk satisfies zero of the 3-2-1 principle's "different storage" and "offsite" requirements — if that disk fails, the live data and every backup of it disappear together. The simplest fix, needing no new infrastructure: sync the backup directory to an object storage bucket (AWS S3, DigitalOcean Spaces, or similar — all speak a compatible-enough API that the same `aws s3` or `s3cmd`-style tooling works across providers) as the last step of the same script:

```bash
# Appended to backup.sh, after the prune step:

# rclone (or aws s3 sync / s3cmd sync, depending on your provider) is a
# general-purpose tool for syncing a local directory to remote object
# storage - only new or changed files are actually uploaded each run.
rclone sync "$BACKUP_DIR" "remote:gochain-backups" --log-level INFO
```

At this point, all three copies from Section 2's principle exist: the live data on the VM, the local timestamped archive in `/opt/gochain/backups`, and the offsite copy in a bucket owned by a provider entirely independent of whichever one hosts your VM.

---

## 8. The Restore Procedure

A backup is only useful if you know, precisely and without improvising under pressure, how to turn it back into a running node. Documented here as the exact steps to actually follow:

1. **Stop the affected node**, so nothing is trying to write to its data directory while you replace it: `docker compose stop node2`.
2. **Locate the most recent good backup**: `ls -lt /opt/gochain/backups/node2-*.db.gz | head -1`.
3. **Decompress it** into a temporary location: `gunzip -c node2-20260415-031500.db.gz > /tmp/restored.db`.
4. **Replace the node's live data file** inside its Docker volume:
   ```bash
   docker compose cp /tmp/restored.db node2:/app/data/gochain.db
   ```
5. **Restart the node**: `docker compose start node2`.
6. **Watch its logs** to confirm it starts up cleanly against the restored file, then catches up on anything mined since the backup was taken: `docker compose logs -f node2`.

Step 6 matters more than it looks: a restored node is, by definition, slightly behind the live chain (it holds whatever height existed at backup time, not the current tip), so a correct restore is not "the node has all the data back" but "the node has enough data back to resume, and Chapter 49's synchronization logic pulls in the rest" — which is exactly what the chaos test in the next section verifies rather than assumes.

---

## 9. The Chaos Test

**Chaos engineering** is the practice of deliberately injecting failure into a system, on purpose, specifically to verify that the safety mechanisms you believe exist actually work — rather than trusting that they would work if a real failure ever happened. Netflix's well-known "Chaos Monkey," which randomly terminates production servers to prove the rest of the system tolerates it, is the most famous example of this idea in the wild. This chapter's chaos test is a small, deliberately destructive script that does exactly this to one GoChain node, on purpose, so you find out today — not during a real incident — whether your backup and restore procedure genuinely works.

```bash
#!/usr/bin/env bash
# /opt/gochain/scripts/chaos-test.sh
#
# Deliberately destroys node2's live data, restores it from the most
# recent backup, and verifies it rejoins the network and catches up -
# proving the backup/restore procedure works, rather than assuming it
# does because nobody has ever needed it yet.
set -euo pipefail

echo "== recording pre-chaos state =="
BEFORE_HEIGHT=$(curl -s http://localhost:8080/chain/height | jq .height)
echo "node1 (reference) height before: ${BEFORE_HEIGHT}"

echo "== taking a fresh backup of node2 =="
/opt/gochain/scripts/backup.sh

echo "== simulating disaster: killing node2 and destroying its volume =="
docker compose kill node2
docker compose rm -f node2
docker volume rm gochain_node2-data

echo "== recreating node2 from a clean container (empty volume) =="
docker compose up -d node2
sleep 5

echo "== restoring the backup taken moments ago =="
LATEST=$(ls -t /opt/gochain/backups/node2-*.db.gz | head -1)
gunzip -c "$LATEST" > /tmp/restored.db
docker compose stop node2
docker compose cp /tmp/restored.db node2:/app/data/gochain.db
docker compose start node2

echo "== waiting for node2 to rejoin and sync =="
sleep 15
AFTER_HEIGHT=$(curl -s http://localhost:8081/chain/height | jq .height)
CURRENT_HEIGHT=$(curl -s http://localhost:8080/chain/height | jq .height)

echo "node1 (reference) current height: ${CURRENT_HEIGHT}"
echo "node2 (restored)  height:         ${AFTER_HEIGHT}"

if [ "$AFTER_HEIGHT" -eq "$CURRENT_HEIGHT" ]; then
  echo "PASS: node2 fully recovered and caught up to the live chain."
else
  echo "FAIL: node2 did not converge - investigate before trusting this backup process."
  exit 1
fi
```

---

## 10. Running the Chaos Test

```bash
chmod +x /opt/gochain/scripts/chaos-test.sh
/opt/gochain/scripts/chaos-test.sh
```

Expected output on a healthy setup:

```
== recording pre-chaos state ==
node1 (reference) height before: 4213
== taking a fresh backup of node2 ==
backed up node2 -> /opt/gochain/backups/node2-20260415-142200.db.gz
== simulating disaster: killing node2 and destroying its volume ==
== recreating node2 from a clean container (empty volume) ==
== restoring the backup taken moments ago ==
== waiting for node2 to rejoin and sync ==
node1 (reference) current height: 4215
node2 (restored)  height:         4215
PASS: node2 fully recovered and caught up to the live chain.
```

Notice the reference height moved from 4213 to 4215 during the test itself — mining did not stop while you were busy destroying and restoring `node2`, which is precisely the point: a real disaster does not politely pause the rest of the network while you fix one machine, and the test above only passes if `node2` catches up to wherever the chain actually is *by the time it finishes recovering*, not to some frozen snapshot of where the chain used to be.

---

## 11. Why an Untested Backup Is Not a Backup

The uncomfortable but important truth this chapter is built around: a backup file sitting untested on a disk somewhere is a *hope*, not a guarantee. Restore procedures rot silently — a script that worked perfectly six months ago can fail today because a file path changed, a container name changed, a tool used in the script was upgraded and its flags changed, or the backup itself has quietly been failing for weeks without anyone noticing (exactly the failure mode Section 5's warning about cron output redirection exists to prevent). The only way to actually know your backups work is to have restored from one recently, deliberately, under conditions you controlled — which is exactly what Sections 9-10 just had you do.

A reasonable operational habit going forward: run the chaos test on a fixed schedule (monthly is a common, reasonable cadence for a small testnet), not just once while reading this chapter, and treat a failing chaos test with exactly the urgency you would treat a failing test in Chapter 92's CI pipeline — because in both cases, the entire purpose of the check is to catch a broken safety net before it is the only thing standing between you and a real, permanent loss.

---

## Summary

- What actually matters most to back up, in priority order: each node's BoltDB data, then wallet files (irreplaceable if lost), then Grafana/Prometheus state (merely inconvenient if lost).
- The 3-2-1 principle — 3 copies, 2 storage media, 1 offsite — is a durable rule of thumb worth applying beyond just this project.
- BoltDB's `tx.WriteTo`, run inside a read-only transaction via `DB.View`, produces a consistent point-in-time snapshot safely, even while the live node continues writing.
- A cron-scheduled shell script snapshots every node, compresses and timestamps the result, prunes old backups, and syncs a copy to offsite object storage.
- Wallet files are already encrypted at rest (Chapter 40); the backup script must never store the decryption password or seed phrase alongside the encrypted file itself.
- The documented restore procedure is: stop the node, decompress the latest good backup, copy it into the container's data path, restart, and watch it catch up via Chapter 49's synchronization.
- A **chaos test** deliberately destroys a running node's data and proves, with an automated pass/fail check comparing heights, that restore-and-resync genuinely works — rather than assuming it would.
- An untested backup is a hope, not a guarantee; running the chaos test on a recurring schedule is what actually keeps that guarantee true over time.

---

## Exercises

### Easy

1. Implement `gochain-backup` and run it manually against one node's BoltDB file while the node is still running. Confirm the resulting snapshot file is non-empty and roughly the same size as the live database.

2. Add the cron entry from Section 5, then temporarily change the schedule to run one minute from now instead of 3:15 AM, and confirm a real backup file appears in `/opt/gochain/backups` without you running the script manually.

3. Intentionally break the backup script (for example, point `--db` at a path that does not exist) and confirm the cron job's log file captures the resulting error message, proving the `2>&1` redirection from Section 5 actually works.

### Medium

4. Implement the offsite sync step from Section 7 using either `rclone` or your cloud provider's own CLI, pointed at a real object storage bucket. Delete a local backup file, then restore it from the offsite copy, confirming the round trip works end to end.

5. Follow the restore procedure in Section 8 manually, step by step, against a real node in your testnet (not the automated chaos-test script), and time how long the entire process takes from "node is down" to "node is back and caught up." Write down the exact commands you ran, in order, as your own personal runbook.

6. Modify the chaos test to target the faucet's wallet file instead of a node's BoltDB file: back it up, delete it, restore it, and confirm the faucet can still sign a valid payout transaction afterward (proving the restored wallet's private key still works, not just that a file of the right size exists).

### Hard

7. Run the full chaos test from Sections 9-10 against your own deployed testnet, capturing the complete before/after height comparison output, and report how long the restore-and-catch-up process took relative to how many blocks `node2` had to sync (compare this against Chapter 87's Exercise 7, which measured a similar recovery after a hard `kill` without any data loss at all).

8. Extend the chaos test to simulate two nodes failing simultaneously (`node2` and `node3`), restored from backups taken at different times (one an hour older than the other), and verify both correctly converge to the same final chain height. Explain what would happen if you restored a backup that was *older* than the last block your restore procedure assumes was "caught up."

9. Design (and, if practical, implement) an automated alert that fires if the backup script has not successfully completed within the last 25 hours — a real safety net for the exact silent-failure scenario Section 11 warns about. Explain what signal you chose to monitor (a heartbeat file's timestamp, a Prometheus metric pushed at the end of the script, an external "dead man's switch" service) and why you rejected the alternatives.
