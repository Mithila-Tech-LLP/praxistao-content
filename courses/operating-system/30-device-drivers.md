# Chapter 30: Device Drivers — Extending the Kernel

> **"A device driver is a bridge between the abstract world of the OS (files, bytes, requests) and the concrete world of hardware (registers, DMA, interrupts, timing). Without drivers, hardware would be completely opaque. With them, a USB webcam appears as `/dev/video0` and the rest of the OS never needs to know it has sensors and analog-to-digital converters inside."**

---

## Table of Contents

1. [What Is a Device Driver?](#1-what-is-a-device-driver)
2. [The Driver Model](#2-the-driver-model)
3. [Character vs. Block vs. Network Devices](#3-character-vs-block-vs-network-devices)
4. [Character Device Drivers — Step by Step](#4-character-device-drivers--step-by-step)
5. [Block Device Drivers](#5-block-device-drivers)
6. [Kernel Modules — Loadable Drivers](#6-kernel-modules--loadable-drivers)
7. [DMA — Direct Memory Access](#7-dma--direct-memory-access)
8. [The Device Tree and ACPI](#8-the-device-tree-and-acpi)
9. [Driver Development — Minimal Example](#9-driver-development--minimal-example)
10. [Summary](#summary)

---

## 1. What Is a Device Driver?

**A device driver is kernel code that:**
1. Initializes hardware (configures registers, sets operating mode)
2. Provides a standard VFS interface (`open`, `read`, `write`, `ioctl`, `close`)
3. Handles IRQs from the device
4. Manages DMA transfers
5. Reports errors

**Why drivers run in kernel space:**
- Drivers access hardware registers (I/O ports and MMIO regions) — requires Ring 0
- Drivers handle IRQs — requires direct interrupt handling
- Drivers share kernel data structures (device queues, page tables) — requires kernel access
- A crashed driver can crash the whole OS (this is why 80% of Windows crashes were driver bugs in the Windows XP era)

**Microkernel alternative:**
In a microkernel OS (QNX, Minix), drivers run as user-space processes. A buggy driver crashes only itself, not the whole OS. The cost: extra IPC for every I/O operation.

---

## 2. The Driver Model

Linux's device driver model is built around three abstractions:

**Bus:**
The hardware interconnect that devices plug into:
```
pci:    PCIe/PCI bus (graphics cards, network cards, SSDs)
usb:    USB bus (keyboards, mice, drives, cameras)
i2c:    I2C bus (sensors, EEPROM, embedded devices)
spi:    SPI bus (flash memory, displays)
platform: On-chip devices without a proper bus (ARM SoC peripherals)
acpi:   ACPI namespace (x86 system devices)
```

**Device:**
A specific piece of hardware on a bus:
```c
struct device {
    struct device        *parent;    // bus controller
    struct device_type   *type;
    struct bus_type      *bus;
    struct device_driver *driver;    // which driver manages this device
    void                 *platform_data;  // device-specific data
    struct device_node   *of_node;   // Device Tree node (embedded)
    // ...
};
```

**Driver:**
Code that knows how to talk to a specific class of hardware:
```c
struct device_driver {
    const char         *name;
    struct bus_type    *bus;
    int (*probe)(struct device *dev);   // called when device is found
    void (*remove)(struct device *dev); // called when device is removed
    // ...
};
```

**Binding:**
The kernel matches devices to drivers by comparing identifiers:
```
PCI device: vendor 0x8086, device 0x1234 (Intel NIC)
  → search driver table for any driver claiming (0x8086, 0x1234)
  → found: e1000 driver
  → call e1000_probe() to initialize the device

USB device: vendor 0x045e, product 0x00cb (Microsoft keyboard)
  → search USB driver table
  → found: hid-keyboard driver
  → call hid_probe() to set up the device
```

---

## 3. Character vs. Block vs. Network Devices

**Character devices:**
- Data accessed as a stream of bytes (no seeking, no blocks)
- Examples: keyboard, mouse, serial port, terminal, `/dev/random`, `/dev/null`
- Interface: `read()`, `write()`, `ioctl()`, `poll()`
- Major/minor numbers identify them in `/dev/`

```bash
ls -l /dev/tty0 /dev/null /dev/urandom
# crw--w---- 1 root tty  4, 0 ... /dev/tty0     # char device, major 4, minor 0
# crw-rw-rw- 1 root root 1, 3 ... /dev/null     # char device, major 1, minor 3
# crw-rw-rw- 1 root root 1, 9 ... /dev/urandom  # char device, major 1, minor 9
```

**Block devices:**
- Data accessed in fixed-size blocks (512B or 4096B)
- Support random access (seek to any block)
- Examples: HDD, SSD, NVMe, USB drives, loop devices
- Requests go through the **block I/O layer** for scheduling and merging

```bash
ls -l /dev/sda
# brw-rw---- 1 root disk 8, 0 ... /dev/sda  # block device, major 8, minor 0
```

**Network devices:**
- Not accessed as files (no `/dev/eth0` to read/write)
- Applications use sockets; kernel routes packets through network devices
- Interface: `netif_receive_skb()` (receive), `dev_queue_xmit()` (transmit)

```bash
ip link show eth0   # network interfaces aren't in /dev/
# 2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP
```

**ioctl (I/O Control):**
For device-specific operations that don't fit read/write:
```c
// User space:
int fd = open("/dev/sda", O_RDONLY);
struct hd_geometry geom;
ioctl(fd, HDIO_GETGEO, &geom);  // get disk geometry
printf("Heads: %d, Sectors: %d, Cylinders: %d\n",
       geom.heads, geom.sectors, geom.cylinders);

// Other ioctl examples:
ioctl(tty_fd, TIOCGWINSZ, &ws);    // get terminal window size
ioctl(sock_fd, SIOCGIFADDR, &ifr); // get network interface IP address
ioctl(drm_fd, DRM_IOCTL_MODE_GETRESOURCES, &res); // get display resources
```

---

## 4. Character Device Drivers — Step by Step

Let's trace how a character device driver works, using a keyboard as an example:

**1. Driver registration:**
```c
static const struct file_operations kbd_fops = {
    .read    = kbd_read,
    .write   = kbd_write,  // typically returns -EINVAL for keyboards
    .open    = kbd_open,
    .release = kbd_close,
    .poll    = kbd_poll,   // for select()/poll()
    .ioctl   = kbd_ioctl,
};

static int __init kbd_init(void) {
    // Register character device with major number 4:
    register_chrdev(4, "tty", &kbd_fops);
    
    // Request IRQ1 (keyboard):
    request_irq(1, kbd_irq_handler, IRQF_SHARED, "keyboard", NULL);
    
    return 0;
}
```

**2. Hardware interrupt fires (key pressed):**
```c
// Top half — minimal work:
static irqreturn_t kbd_irq_handler(int irq, void *dev_id) {
    // Read scan code from keyboard controller (I/O port 0x60):
    uint8_t scancode = inb(0x60);
    
    // Add to a circular buffer (very fast):
    kbd_buf[kbd_head] = scancode;
    kbd_head = (kbd_head + 1) % KBD_BUF_SIZE;
    
    // Wake up any process waiting for keyboard input:
    wake_up_interruptible(&kbd_wait_queue);
    
    return IRQ_HANDLED;
}
```

**3. User calls read():**
```c
// Bottom: user process reads /dev/tty
static ssize_t kbd_read(struct file *file, char __user *buf,
                        size_t count, loff_t *ppos) {
    // If buffer empty, wait for a keypress:
    wait_event_interruptible(kbd_wait_queue, kbd_tail != kbd_head);
    
    if (signal_pending(current))
        return -ERESTARTSYS;
    
    // Convert scan code to ASCII character:
    uint8_t scancode = kbd_buf[kbd_tail];
    kbd_tail = (kbd_tail + 1) % KBD_BUF_SIZE;
    char ascii = scancode_to_ascii(scancode);
    
    // Copy to user space:
    if (copy_to_user(buf, &ascii, 1))
        return -EFAULT;
    
    return 1;  // 1 byte returned
}
```

**Key ideas:**
- `inb()`: read from I/O port (port I/O)
- `wait_event_interruptible()`: put process to sleep until condition is true
- `wake_up_interruptible()`: wake process from sleep
- `copy_to_user()`: safe copy from kernel to user space (validates pointer)

---

## 5. Block Device Drivers

Block devices go through an extra layer — the **block I/O (BIO) layer**:

```
Application: read(fd, buf, 4096)
       ↓
VFS / File System (ext4)
       ↓
Block I/O Layer (merges requests, schedules I/O order)
       ↓
Block Device Driver (sends commands to hardware)
       ↓
Hardware (disk, SSD, NVMe)
       ↓
DMA (hardware writes directly to memory)
       ↓
Interrupt (notifies CPU that I/O is complete)
       ↓
Block layer processes completion, wakes up file system
       ↓
File system returns data to application
```

**Block request merging:**
The block layer collects multiple small requests and merges adjacent ones:
```
Process A: read block 100
Process B: read block 101
Process C: read block 102
→ Merged into: read blocks 100-102 in one command (3× the throughput)
```

**I/O schedulers:**
Linux supports pluggable I/O schedulers for block devices:
```bash
cat /sys/block/sda/queue/scheduler
# [mq-deadline] kyber bfq none

# For HDDs: BFQ (Budget Fair Queueing) — prioritizes latency-sensitive I/O
echo bfq > /sys/block/sda/queue/scheduler

# For NVMe SSDs: none (no scheduling needed — NVMe handles it in hardware)
echo none > /sys/block/nvme0n1/queue/scheduler
```

---

## 6. Kernel Modules — Loadable Drivers

Instead of compiling all drivers into the kernel, Linux supports **loadable kernel modules (LKMs)** — `.ko` files that can be loaded/unloaded at runtime.

```bash
# List loaded modules:
lsmod
# Module                  Size  Used by
# e1000e               286720  0       ← Intel Ethernet driver
# bluetooth            757760  56 btusb
# usbcore              290816  4 xhci_hcd,btusb

# Load a module:
modprobe e1000e      # loads with dependencies resolved automatically
insmod e1000e.ko     # load specific .ko file (no dependency resolution)

# Unload a module:
modprobe -r e1000e   # removes module + unused dependencies
rmmod e1000e         # force remove

# Get module information:
modinfo e1000e
# description:    Intel(R) PRO/1000 Network Driver
# author:         Intel Corporation
# license:        GPL v2
# srcversion:     ...

# Module parameters:
modprobe e1000e InterruptThrottleRate=3000
# passes "InterruptThrottleRate" parameter to e1000e_init
```

**Module structure:**
```c
#include <linux/module.h>
#include <linux/init.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Your Name");
MODULE_DESCRIPTION("My example driver");

static int __init my_driver_init(void) {
    printk(KERN_INFO "my_driver: loaded\n");
    // register device, request IRQ, etc.
    return 0;
}

static void __exit my_driver_exit(void) {
    printk(KERN_INFO "my_driver: unloaded\n");
    // release resources, unregister device
}

module_init(my_driver_init);
module_exit(my_driver_exit);
```

**`/sys/module/`:**
```bash
ls /sys/module/e1000e/
# drivers/  holders/  initstate  notes/  parameters/  sections/  taint

cat /sys/module/e1000e/parameters/InterruptThrottleRate
# 3
```

---

## 7. DMA — Direct Memory Access

**The problem without DMA:**
To read 1MB from disk:
```
Without DMA:
  Disk controller reads 1 byte → interrupts CPU
  CPU copies byte from I/O port to RAM
  Repeat 1,048,576 times
  = 1 million CPU interruptions to transfer 1MB → unacceptable!
```

**DMA (Direct Memory Access):**
A hardware DMA controller can transfer data directly between a device and RAM **without CPU involvement**:

```
With DMA:
  1. CPU tells DMA controller:
     - Source: disk sector N
     - Destination: physical RAM address 0x10000000
     - Length: 4096 bytes
  2. DMA controller handles the transfer autonomously:
     - Reads from disk controller
     - Writes to RAM via memory bus (CPU bus is not busy)
  3. When done: DMA controller raises an interrupt
  4. CPU processes interrupt: "DMA transfer complete"
  5. Data is in RAM, ready to use

CPU was free to do other work during the entire transfer!
```

**DMA in Linux:**
```c
// Allocate DMA-coherent memory (CPU + device can both access):
dma_addr_t dma_handle;
void *cpu_addr = dma_alloc_coherent(dev, size, &dma_handle, GFP_KERNEL);
// cpu_addr: virtual address (for CPU to read the data after DMA completes)
// dma_handle: physical/bus address (for device to write to)

// Program device with dma_handle (tell it where to write):
writel(dma_handle, device_base + DMA_ADDR_REGISTER);
writel(size, device_base + DMA_LENGTH_REGISTER);
writel(DMA_START, device_base + DMA_CONTROL_REGISTER);

// In IRQ handler (DMA complete):
data_ready = 1;
wake_up(&wait_queue);

// CPU reads result via cpu_addr:
memcpy(user_buffer, cpu_addr, size);

// Free DMA memory:
dma_free_coherent(dev, size, cpu_addr, dma_handle);
```

**IOMMU:**
On modern systems, an IOMMU (I/O Memory Management Unit) protects RAM from rogue DMA:
- Device can only write to pages explicitly mapped in the IOMMU
- Prevents DMA attacks: a malicious PCI device can't overwrite arbitrary kernel memory
- Required for secure virtualization (prevents guest VMs from attacking hypervisor via DMA)

---

## 8. The Device Tree and ACPI

**How does the kernel know what hardware exists?**

**ACPI (Advanced Configuration and Power Interface) — x86:**
On x86 PCs, firmware (UEFI/BIOS) provides ACPI tables describing hardware:
```bash
# List ACPI tables:
ls /sys/firmware/acpi/tables/
# DSDT  FADT  HPET  MADT  MCFG  SSDT  ...

# MADT: Multiple APIC Description Table — lists CPUs and APIC IDs
# DSDT: Differentiated System Description Table — full hardware description
```

**Device Tree — ARM/embedded:**
ARM systems (phones, Raspberry Pi, embedded boards) use a **Device Tree Blob (DTB)** — a binary file describing hardware topology passed by the bootloader to the kernel:

```
/dts-v1/;
/ {
    cpus {
        cpu@0 {
            compatible = "arm,cortex-a72";
            device_type = "cpu";
        };
    };
    
    memory@0 {
        device_type = "memory";
        reg = <0x0 0x00000000 0x0 0x40000000>;  /* 1GB at address 0 */
    };
    
    uart0: uart@ff000000 {
        compatible = "brcm,bcm2835-aux-uart";
        reg = <0xff000000 0x40>;
        interrupts = <0x1 0x1d 0x4>;
    };
    
    ethernet@7ef12000 {
        compatible = "brcm,bcm2711-genet-v5";
        reg = <0x7ef12000 0x4000>;
    };
};
```

The kernel reads the Device Tree, creates `platform_device` objects for each node, and binds the matching drivers (based on the `compatible` string).

---

## 9. Driver Development — Minimal Example

A minimal `/dev/hello` character driver:

```c
#include <linux/module.h>
#include <linux/fs.h>
#include <linux/uaccess.h>

#define DEVICE_NAME "hello"
#define MAJOR_NUM 240

static const char hello_msg[] = "Hello from kernel driver!\n";

static ssize_t hello_read(struct file *file, char __user *buf,
                           size_t count, loff_t *ppos)
{
    size_t len = min(count, sizeof(hello_msg) - (size_t)*ppos);
    if (len == 0)
        return 0;  // EOF
    
    if (copy_to_user(buf, hello_msg + *ppos, len))
        return -EFAULT;
    
    *ppos += len;
    return len;
}

static const struct file_operations hello_fops = {
    .owner = THIS_MODULE,
    .read  = hello_read,
};

static int __init hello_init(void) {
    int ret = register_chrdev(MAJOR_NUM, DEVICE_NAME, &hello_fops);
    if (ret < 0) {
        printk(KERN_ALERT "hello: failed to register device\n");
        return ret;
    }
    printk(KERN_INFO "hello: registered with major %d\n", MAJOR_NUM);
    return 0;
}

static void __exit hello_exit(void) {
    unregister_chrdev(MAJOR_NUM, DEVICE_NAME);
    printk(KERN_INFO "hello: unloaded\n");
}

module_init(hello_init);
module_exit(hello_exit);
MODULE_LICENSE("GPL");
```

```bash
# Build, load, and use:
make -C /lib/modules/$(uname -r)/build M=$(pwd) modules
insmod hello.ko
mknod /dev/hello c 240 0      # create device file
cat /dev/hello
# Hello from kernel driver!
rmmod hello
```

---

## Summary

| Concept | Description |
|---------|------------|
| Device driver | Kernel code bridging OS abstractions and hardware |
| Driver model | Bus → device → driver hierarchy; probe/remove lifecycle |
| Char device | Stream-based; read/write byte by byte; keyboard, serial, /dev/null |
| Block device | Random-access fixed-size blocks; disks, SSDs; goes through block layer |
| Network device | Not a file; accessed via sockets; packets, not bytes |
| ioctl | Device-specific commands beyond read/write |
| Major number | Identifies the driver class (8 = SCSI/SATA) |
| Minor number | Identifies the specific device within a driver class |
| Kernel module | Loadable `.ko` file; driver loaded at runtime without kernel recompile |
| DMA | Hardware transfer device → RAM without CPU; generates interrupt when done |
| IOMMU | Hardware protection: limits which memory a device can DMA to |
| Device Tree | ARM/embedded hardware description passed by bootloader to kernel |
| ACPI | x86 firmware tables describing hardware topology to the kernel |
| copy_to_user | Safe copy from kernel buffer to user-space pointer (validates address) |
| request_irq | Register an ISR for a specific IRQ line |
| Top half | Fast IRQ handler; bottom half defers slower work |
