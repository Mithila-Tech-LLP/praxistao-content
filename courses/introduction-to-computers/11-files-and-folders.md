# Chapter 11: Files and Folders — How Computers Organize Everything

> **"Without a file system, your hard drive would be a junkyard — data piled everywhere with no way to find anything. Files and folders are the filing system that makes storage usable."**

---

## Table of Contents

1. [What Is a File?](#1-what-is-a-file)
2. [What Is a Folder?](#2-what-is-a-folder)
3. [The Folder Tree — How Everything Is Organized](#3-the-folder-tree--how-everything-is-organized)
4. [File Paths — The Address of a File](#4-file-paths--the-address-of-a-file)
5. [Important System Folders](#5-important-system-folders)
6. [File Properties — Size, Date, Type](#6-file-properties--size-date-type)
7. [Searching for Files](#7-searching-for-files)
8. [The Recycle Bin / Trash](#8-the-recycle-bin--trash)
9. [Summary](#summary)

---

## 1. What Is a File?

A **file** is a collection of data with a name, stored on your drive.

```
Examples of files:
  photo_of_cat.jpg      → your cat photo
  homework.docx         → your Word document
  favourite_song.mp3    → a music track
  game_save.dat         → your game's saved progress
  notes.txt             → a text file you wrote
  
A file has:
  Name:      "photo_of_cat"
  Extension: ".jpg"  (tells the OS what type of file it is)
  Contents:  the actual data (pixels, text, audio, etc.)
  Metadata:  size, date created, date modified, permissions
```

Files are the fundamental unit of storage. Everything on your computer — apps, photos, documents, operating system itself — is stored as files.

---

## 2. What Is a Folder?

A **folder** (also called a **directory**) is a container that holds files — or other folders.

Folders don't actually store data. They're just a way to organize files. Think of them as labeled boxes that hold files.

```
Without folders:
  
  Your entire computer's 500GB drive is one flat list:
  
  photo_of_cat.jpg
  photo_of_dog.jpg
  homework_january.docx
  homework_february.docx
  game_save.dat
  system_config.bin
  browser.exe
  ...
  
  500,000 files, no organization. Impossible to find anything.
  
With folders:
  
  Photos/
    Pets/
      photo_of_cat.jpg
      photo_of_dog.jpg
    Holidays/
      beach_2024.jpg
  Documents/
    Homework/
      january.docx
      february.docx
  Games/
    Save/
      game_save.dat
```

Folders can contain other folders — creating a hierarchy that can be as deep as you need.

---

## 3. The Folder Tree — How Everything Is Organized

Every OS organizes its entire drive as a tree structure starting from a "root."

**Windows:**
```
C:\ (root — the C: drive)
├── Windows\           (OS files — don't touch)
│   ├── System32\
│   └── SysWOW64\
├── Program Files\     (installed apps)
│   ├── Chrome\
│   └── Spotify\
├── Users\
│   └── YourName\      (your personal folder)
│       ├── Desktop\
│       ├── Documents\
│       ├── Downloads\
│       ├── Music\
│       ├── Pictures\
│       └── Videos\
└── Temp\              (temporary files)
```

**macOS:**
```
/ (root)
├── Applications/      (installed apps)
├── System/            (macOS files — don't touch)
├── Library/           (system settings)
├── Users/
│   └── yourname/      (your home folder, ~)
│       ├── Desktop/
│       ├── Documents/
│       ├── Downloads/
│       ├── Music/
│       ├── Pictures/
│       └── Movies/
└── tmp/               (temporary files)
```

**The key rule:** Your personal files live inside your **home folder** (`C:\Users\Name` on Windows, `/Users/name` on Mac). The OS and app files live elsewhere and you usually don't need to touch them.

---

## 4. File Paths — The Address of a File

A **path** is the complete address of a file, showing every folder you pass through to reach it.

```
Windows path:
  C:\Users\John\Documents\Homework\january.docx
  
  C:\          → root of C drive
  Users\       → Users folder
  John\        → John's home folder
  Documents\   → Documents subfolder
  Homework\    → Homework subfolder
  january.docx → the file itself
  
Mac / Linux path:
  /Users/john/Documents/Homework/january.docx
  
  /            → root
  Users/       → Users folder
  john/        → john's home folder
  Documents/   → Documents subfolder
  Homework/    → Homework subfolder
  january.docx → the file itself

Windows uses backslash (\), Mac/Linux use forward slash (/)
```

Paths are important when you're programming or using the command line. You'll use them constantly if you become a developer.

---

## 5. Important System Folders

**Folders on Windows you should know:**
```
C:\Windows\System32\    Core operating system files. Never delete.
C:\Program Files\       Where 64-bit apps are installed.
C:\Users\Name\AppData\  Hidden folder — app settings, cache.
C:\Users\Name\Desktop\  Everything on your desktop is here.
C:\Users\Name\Downloads\ Downloaded files land here.
```

**Folders on Mac you should know:**
```
/Applications/          Your apps.
/Library/               System-wide settings and extensions.
~/Library/              Your personal settings (~ means your home).
~/Desktop/              Your desktop files.
~/Downloads/            Your downloaded files.
```

---

## 6. File Properties — Size, Date, Type

Right-clicking a file and selecting "Properties" (Windows) or "Get Info" (Mac) shows:

```
Name:           january.docx
Type:           Microsoft Word Document
Size:           24.3 KB (kilobytes)
Location:       C:\Users\John\Documents\Homework
Created:        January 5, 2024 at 9:12 AM
Modified:       January 20, 2024 at 3:47 PM
Opened:         March 10, 2024 at 8:00 AM
Permissions:    Read & Write (for your user)
```

**File size matters when:**
- Your storage is filling up (large files take lots of space)
- Sending files by email (most email limits ~25MB attachments)
- Uploading to the internet (large files = slow upload)

---

## 7. Searching for Files

```
Windows: Win+S or click the search bar
  Type the filename (or part of it)
  Can search inside documents too

Mac: Cmd+Space (Spotlight)
  Type the filename, also searches inside documents
  Can find apps, settings, calculator, calendar events

Phone:
  iOS: Swipe down from home screen → search bar
  Android: Swipe down or use Google widget
```

**Tips for finding files:**
- If you don't remember the name, search by date: "files modified last week"
- If you know the content: "documents containing 'budget report'"
- Keep files organized in folders — future you will be grateful

---

## 8. The Recycle Bin / Trash

```
When you delete a file:
  It goes to the Recycle Bin (Windows) or Trash (Mac/iOS)
  It hasn't actually been deleted yet.
  You can restore it anytime.
  
When you empty the Recycle Bin / Trash:
  The file's entry is removed from the file system
  The space is marked as "available" for new data
  But the data is still physically on the disk until overwritten
  
  (This is why "deleted" files can sometimes be recovered by
  forensic software — they may still be on the disk physically)
  
Permanent deletion:
  Windows: Shift+Delete (skips Recycle Bin)
  Mac: Option+Cmd+Delete
  Still technically recoverable until overwritten
  
Secure deletion:
  Software that overwrites the data with random 0s and 1s
  Makes recovery impossible
  Important for privacy before selling/recycling a device
```

---

## Summary

| Concept | What It Means |
|---------|--------------|
| File | A named collection of data on storage |
| Folder / Directory | A container that organizes files and other folders |
| File extension | The .suffix that tells the OS what type of file it is |
| Path | The full address of a file (every folder from root to file) |
| Root | The top-level of the folder tree (C:\ on Windows, / on Mac/Linux) |
| Home folder | Your personal folder where your files live |
| Recycle Bin / Trash | Temporary holding area for deleted files |

**Files and folders organize storage. Next: the applications and programs that you actually use every day.**
