# Chapter 39: VaultDB — Storage Engine: Reading and Writing Pages

The storage engine is the bottom layer of VaultDB — it reads and writes pages to disk files. Nothing else cares how data is physically stored; everything else talks to this layer.

## Table of Contents

1. The Disk Manager
2. The Database File Format
3. Page Allocation and Free List
4. Serializing Rows to Pages
5. The Table Heap
6. Exercises

---

## 1. The Disk Manager

The disk manager owns one file per database. Pages are read and written by PageID.

```go
// storage/disk.go
package storage

import (
    "fmt"
    "os"
    "sync"
)

// DiskManager reads and writes fixed-size pages to a single file
type DiskManager struct {
    mu       sync.RWMutex
    file     *os.File
    filename string
    numPages uint64
}

func NewDiskManager(filename string) (*DiskManager, error) {
    f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return nil, fmt.Errorf("open db file: %w", err)
    }

    info, err := f.Stat()
    if err != nil {
        return nil, err
    }

    dm := &DiskManager{
        file:     f,
        filename: filename,
        numPages: uint64(info.Size()) / PageSize,
    }

    // Page 0 is reserved for the database header
    if dm.numPages == 0 {
        if err := dm.initHeader(); err != nil {
            return nil, err
        }
    }

    return dm, nil
}

func (dm *DiskManager) ReadPage(id PageID) (*Page, error) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()

    if uint64(id) >= dm.numPages {
        return nil, fmt.Errorf("page %d out of range (have %d pages)", id, dm.numPages)
    }

    page := new(Page)
    offset := int64(id) * PageSize
    n, err := dm.file.ReadAt(page[:], offset)
    if err != nil {
        return nil, fmt.Errorf("read page %d: %w", id, err)
    }
    if n != PageSize {
        return nil, fmt.Errorf("short read: got %d bytes, want %d", n, PageSize)
    }
    return page, nil
}

func (dm *DiskManager) WritePage(id PageID, page *Page) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()

    offset := int64(id) * PageSize
    n, err := dm.file.WriteAt(page[:], offset)
    if err != nil {
        return fmt.Errorf("write page %d: %w", id, err)
    }
    if n != PageSize {
        return fmt.Errorf("short write: wrote %d bytes, want %d", n, PageSize)
    }
    return nil
}

// AllocatePage adds a new page at the end of the file and returns its ID
func (dm *DiskManager) AllocatePage() (PageID, error) {
    dm.mu.Lock()
    defer dm.mu.Unlock()

    id := PageID(dm.numPages)
    dm.numPages++

    // Write a zero page to extend the file
    page := new(Page)
    offset := int64(id) * PageSize
    _, err := dm.file.WriteAt(page[:], offset)
    if err != nil {
        dm.numPages--
        return InvalidPageID, fmt.Errorf("allocate page: %w", err)
    }
    return id, nil
}

func (dm *DiskManager) Sync() error {
    return dm.file.Sync()
}

func (dm *DiskManager) Close() error {
    dm.file.Sync()
    return dm.file.Close()
}

func (dm *DiskManager) NumPages() uint64 {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    return dm.numPages
}
```

---

## 2. The Database File Format

Page 0 is the **database header page** — it stores metadata about all tables.

```go
// storage/header.go
package storage

import (
    "encoding/binary"
    "fmt"
)

// Database header (stored in page 0)
// Layout:
// Bytes  0- 7: magic number ("VAULTDB\x01")
// Bytes  8-15: number of tables
// Bytes 16-23: first free list page ID (InvalidPageID if none)
// Bytes 24-...: table catalog entries (variable number)

var MagicNumber = [8]byte{'V', 'A', 'U', 'L', 'T', 'D', 'B', 0x01}

const (
    maxTableNameLen  = 64
    maxColumnNameLen = 32
    maxColumns       = 16
)

type ColumnDef struct {
    Name   string
    TypeID TypeID
}

type TableDef struct {
    Name       string
    RootPageID PageID // root page of the B-Tree for this table
    Columns    []ColumnDef
}

type Catalog struct {
    Tables  []TableDef
    FreePage PageID
}

func (dm *DiskManager) initHeader() error {
    page := new(Page)
    copy(page[:8], MagicNumber[:])
    binary.BigEndian.PutUint64(page[8:16], 0) // 0 tables
    binary.BigEndian.PutUint64(page[16:24], uint64(InvalidPageID))

    dm.numPages = 1
    offset := int64(0)
    _, err := dm.file.WriteAt(page[:], offset)
    return err
}

func (dm *DiskManager) ReadCatalog() (*Catalog, error) {
    page, err := dm.ReadPage(0)
    if err != nil {
        return nil, err
    }

    if [8]byte(page[:8]) != MagicNumber {
        return nil, fmt.Errorf("invalid database file: bad magic number")
    }

    catalog := &Catalog{
        FreePage: PageID(binary.BigEndian.Uint64(page[16:24])),
    }

    numTables := binary.BigEndian.Uint64(page[8:16])
    offset := 24

    for i := 0; i < int(numTables); i++ {
        var tbl TableDef

        nameLen := int(page[offset])
        offset++
        tbl.Name = string(page[offset : offset+nameLen])
        offset += nameLen

        tbl.RootPageID = PageID(binary.BigEndian.Uint64(page[offset:]))
        offset += 8

        numCols := int(page[offset])
        offset++

        for j := 0; j < numCols; j++ {
            var col ColumnDef
            colNameLen := int(page[offset])
            offset++
            col.Name = string(page[offset : offset+colNameLen])
            offset += colNameLen
            col.TypeID = TypeID(page[offset])
            offset++
            tbl.Columns = append(tbl.Columns, col)
        }

        catalog.Tables = append(catalog.Tables, tbl)
    }

    return catalog, nil
}

func (dm *DiskManager) WriteCatalog(catalog *Catalog) error {
    page := new(Page)
    copy(page[:8], MagicNumber[:])

    binary.BigEndian.PutUint64(page[8:16], uint64(len(catalog.Tables)))
    binary.BigEndian.PutUint64(page[16:24], uint64(catalog.FreePage))

    offset := 24
    for _, tbl := range catalog.Tables {
        page[offset] = byte(len(tbl.Name))
        offset++
        copy(page[offset:], tbl.Name)
        offset += len(tbl.Name)
        binary.BigEndian.PutUint64(page[offset:], uint64(tbl.RootPageID))
        offset += 8
        page[offset] = byte(len(tbl.Columns))
        offset++
        for _, col := range tbl.Columns {
            page[offset] = byte(len(col.Name))
            offset++
            copy(page[offset:], col.Name)
            offset += len(col.Name)
            page[offset] = byte(col.TypeID)
            offset++
        }
    }

    return dm.WritePage(0, page)
}
```

---

## 3. Serializing Rows to Pages

A **row** is a sequence of values (one per column). We need to encode and decode rows to/from bytes.

```go
// storage/row.go
package storage

import (
    "encoding/binary"
    "fmt"
)

// Row is a list of values, one per column
type Row []Value

// Encode serializes a row to bytes
// Format:
//   For each value:
//   - 1 byte: type ID (TypeNull = value is NULL)
//   - if not null: 4 bytes length + length bytes of data
func EncodeRow(row Row) []byte {
    var buf []byte
    for _, v := range row {
        buf = append(buf, byte(v.Type))
        if v.Type != TypeNull {
            lenBuf := make([]byte, 4)
            binary.BigEndian.PutUint32(lenBuf, uint32(len(v.Data)))
            buf = append(buf, lenBuf...)
            buf = append(buf, v.Data...)
        }
    }
    return buf
}

// DecodeRow deserializes bytes back into a row
func DecodeRow(data []byte, numCols int) (Row, error) {
    row := make(Row, numCols)
    offset := 0
    for i := 0; i < numCols; i++ {
        if offset >= len(data) {
            return nil, fmt.Errorf("unexpected end of data at column %d", i)
        }
        typeID := TypeID(data[offset])
        offset++
        row[i].Type = typeID
        if typeID == TypeNull {
            continue
        }
        if offset+4 > len(data) {
            return nil, fmt.Errorf("truncated length at column %d", i)
        }
        length := int(binary.BigEndian.Uint32(data[offset:]))
        offset += 4
        if offset+length > len(data) {
            return nil, fmt.Errorf("truncated value at column %d", i)
        }
        row[i].Data = make([]byte, length)
        copy(row[i].Data, data[offset:offset+length])
        offset += length
    }
    return row, nil
}
```

---

## 4. The Table Heap

A **heap** is the simplest way to store rows: append new rows to pages, scan sequentially to read.

```go
// storage/heap.go
package storage

import "fmt"

// RowID uniquely identifies a row: (pageID, slotID)
type RowID struct {
    PageID PageID
    SlotID uint16
}

// Heap stores rows across a chain of pages
type Heap struct {
    dm       *DiskManager
    rootPage PageID
}

func NewHeap(dm *DiskManager, rootPage PageID) *Heap {
    return &Heap{dm: dm, rootPage: rootPage}
}

// InsertRow appends a row to the heap
func (h *Heap) InsertRow(row Row) (RowID, error) {
    data := EncodeRow(row)

    // Find a page with enough space
    pageID, page, err := h.findPageWithSpace(len(data) + slotSize)
    if err != nil {
        return RowID{}, err
    }

    slotID, err := insertIntoPage(page, data)
    if err != nil {
        return RowID{}, err
    }

    if err := h.dm.WritePage(pageID, page); err != nil {
        return RowID{}, err
    }

    return RowID{PageID: pageID, SlotID: uint16(slotID)}, nil
}

// ScanAll returns all rows in the heap
func (h *Heap) ScanAll(numCols int) ([]Row, error) {
    var rows []Row
    pageID := h.rootPage

    for pageID != InvalidPageID {
        page, err := h.dm.ReadPage(pageID)
        if err != nil {
            return nil, err
        }

        h2 := page.GetHeader()
        for slot := 0; slot < int(h2.NumSlots); slot++ {
            data, ok := getSlotData(page, slot)
            if !ok || len(data) == 0 {
                continue // deleted slot
            }
            row, err := DecodeRow(data, numCols)
            if err != nil {
                return nil, fmt.Errorf("decode row page=%d slot=%d: %w", pageID, slot, err)
            }
            rows = append(rows, row)
        }

        pageID = nextPage(page)
    }
    return rows, nil
}

// findPageWithSpace finds or allocates a page with at least `need` bytes free
func (h *Heap) findPageWithSpace(need int) (PageID, *Page, error) {
    pageID := h.rootPage
    for pageID != InvalidPageID {
        page, err := h.dm.ReadPage(pageID)
        if err != nil {
            return 0, nil, err
        }
        if page.FreeSpace() >= need {
            return pageID, page, nil
        }
        pageID = nextPage(page)
    }

    // Allocate a new page
    newID, err := h.dm.AllocatePage()
    if err != nil {
        return 0, nil, err
    }
    newPage := new(Page)
    newPage.Initialize(PageTypeLeaf)

    // Link to previous page (update last page's "next" pointer)
    // Simplified: store next page ID in bytes 6-7 of header
    // (In production, use a proper page header field)

    return newID, newPage, nil
}

// insertIntoPage adds data to a page using the slot mechanism
func insertIntoPage(p *Page, data []byte) (int, error) {
    h := p.GetHeader()

    slotTableEnd := headerSize + int(h.NumSlots)*slotSize
    dataStart := int(h.FreeOffset) - len(data)

    if dataStart < slotTableEnd {
        return 0, fmt.Errorf("page full")
    }

    // Write data at the bottom
    copy(p[dataStart:], data)

    // Write slot entry at the top
    slotOffset := headerSize + int(h.NumSlots)*slotSize
    p[slotOffset] = byte(uint16(dataStart) >> 8)
    p[slotOffset+1] = byte(uint16(dataStart))
    p[slotOffset+2] = byte(uint16(len(data)) >> 8)
    p[slotOffset+3] = byte(uint16(len(data)))

    h.NumSlots++
    h.FreeOffset = uint16(dataStart)
    p.SetHeader(h)

    return int(h.NumSlots) - 1, nil
}

func getSlotData(p *Page, slotID int) ([]byte, bool) {
    h := p.GetHeader()
    if slotID >= int(h.NumSlots) {
        return nil, false
    }
    off := headerSize + slotID*slotSize
    dataOffset := int(p[off])<<8 | int(p[off+1])
    dataLen := int(p[off+2])<<8 | int(p[off+3])
    if dataLen == 0 {
        return nil, true // deleted
    }
    return p[dataOffset : dataOffset+dataLen], true
}

func nextPage(p *Page) PageID {
    // Bytes 6-7 of header: next page ID (0 = none)
    id := PageID(uint16(p[6])<<8 | uint16(p[7]))
    if id == 0 {
        return InvalidPageID
    }
    return id
}
```

---

## Summary

- The disk manager reads and writes fixed-size pages to a file using `ReadAt`/`WriteAt` for thread-safe I/O.
- Page 0 stores the database catalog (table names, column definitions, root page IDs).
- Rows are serialized to bytes using a simple type-length-value encoding.
- The heap is the simplest table storage: append rows to pages, scan sequentially to read.
- Next we'll add a B-Tree on top of the heap for fast indexed lookups.

### Exercises

**Easy:** Write a test that creates a `DiskManager`, allocates 3 pages, writes a known byte pattern to each, reads them back, and verifies the content matches.

**Medium:** Implement `DeleteRow(rid RowID)` on the `Heap`: mark the slot as deleted (set length = 0) and write the page back. Verify that `ScanAll` skips deleted rows.

**Hard:** Implement `UpdateRow(rid RowID, row Row)`: if the new row fits in the same slot, update in place. If it's larger, delete the old slot and insert as a new row (updating any index entries would be the caller's responsibility).
