# Chapter 40: VaultDB — B-Trees: The Heart of Every Database Index

B-Trees are the most important data structure in all of computer science for databases. PostgreSQL, MySQL, SQLite, MongoDB — every major database uses B-Trees for indexes. Understanding B-Trees means understanding how databases are fast.

## Table of Contents

1. Why B-Trees?
2. How a B-Tree Works
3. B-Tree Node Structure on a Page
4. Searching a B-Tree
5. Inserting into a B-Tree
6. Node Splitting
7. Full B-Tree Implementation
8. Exercises

---

## 1. Why B-Trees?

**The problem:** We have 10 million rows. Find rows where `user_id = 42`. Without an index, we read every page sequentially — 10 million comparisons.

**Binary search tree:** If rows were in a binary tree, we'd find the answer in ~log₂(10 million) ≈ 23 comparisons. But binary trees have nodes with 2 children. Each node = one disk read. 23 disk reads × 8ms seek time = 184ms just for seek time.

**B-Tree insight:** Make each node hold many keys and have many children. A node that fits in one 4 KB page can hold ~200 keys and have 201 children. Now log₂₀₁(10 million) ≈ 3.4 levels. **Only 4 disk reads to find any row in a 10 million row table!**

```
Binary tree (10M rows): ~23 disk reads
B-Tree (10M rows, order 100): ~4 disk reads
B-Tree (1B rows, order 100): ~5 disk reads
```

---

## 2. How a B-Tree Works

A B-Tree of order M has:
- Internal nodes: hold keys and child page pointers
- Leaf nodes: hold actual key-value pairs (or row pointers)
- Every leaf is at the same depth

```
                    [30, 70]           ← internal node (root)
                   /    |    \
           [10,20]   [40,60]   [80,90]  ← internal nodes
           /  |  \   /  |  \   /  |  \
         ...  ...  ...               ... ← leaf nodes
```

**Searching:** Start at root. Compare key to node keys. Follow the appropriate child pointer. Repeat until leaf.

**Inserting:** Find the correct leaf. Add the key. If the leaf is full, split it in half and push the middle key up to the parent. If the parent is full, split it too. The tree grows upward from splits.

**B+ Tree (what most databases use):** All data is at leaf level. Internal nodes only hold routing keys. Leaves are linked in a doubly-linked list for range scans.

```
Root: [30, 70]
       /    \
  [10,20]  [40,60,80]
  └─next──►└─next───►  (leaf linked list for range scans)
```

---

## 3. B-Tree Node Structure on a Page

```go
// btree/node.go
package btree

import (
    "encoding/binary"
    "github.com/yourname/vaultdb/storage"
)

// B+ Tree node layout in a page:
//
// Header (inherited from page):
//   Type: PageTypeLeaf or PageTypeInternal
//   NumSlots: number of keys in this node
//   FreeOffset: ...
//   ParentID: parent node's page ID
//
// For INTERNAL nodes, slots contain:
//   Key (variable length) + ChildPageID (8 bytes)
//
// For LEAF nodes, slots contain:
//   Key (variable length) + RowID (10 bytes = PageID 8 + SlotID 2)
//   Plus: next leaf page ID in bytes 6-7 of page header

const (
    maxOrder  = 200 // max keys per node
    minFill   = 40  // minimum % fill before merge
)

type NodeType = storage.PageType

// Key is a byte slice (encoded Value)
type Key []byte

// ChildPtr is either a PageID (internal) or RowID (leaf)
type InternalEntry struct {
    Key     Key
    ChildID storage.PageID
}

type LeafEntry struct {
    Key   Key
    RowID storage.RowID
}

// ReadInternalNode reads all entries from an internal node page
func ReadInternalNode(page *storage.Page) []InternalEntry {
    h := page.GetHeader()
    var entries []InternalEntry

    for i := 0; i < int(h.NumSlots); i++ {
        data, ok := getSlotData(page, i)
        if !ok || len(data) < 8 {
            continue
        }
        keyLen := int(binary.BigEndian.Uint16(data[:2]))
        key := make(Key, keyLen)
        copy(key, data[2:2+keyLen])
        childID := storage.PageID(binary.BigEndian.Uint64(data[2+keyLen:]))
        entries = append(entries, InternalEntry{Key: key, ChildID: childID})
    }
    return entries
}

// ReadLeafNode reads all entries from a leaf node page
func ReadLeafNode(page *storage.Page) []LeafEntry {
    h := page.GetHeader()
    var entries []LeafEntry

    for i := 0; i < int(h.NumSlots); i++ {
        data, ok := getSlotData(page, i)
        if !ok || len(data) < 10 {
            continue
        }
        keyLen := int(binary.BigEndian.Uint16(data[:2]))
        key := make(Key, keyLen)
        copy(key, data[2:2+keyLen])
        rowPageID := storage.PageID(binary.BigEndian.Uint64(data[2+keyLen:]))
        rowSlotID := binary.BigEndian.Uint16(data[2+keyLen+8:])
        entries = append(entries, LeafEntry{
            Key:   key,
            RowID: storage.RowID{PageID: rowPageID, SlotID: rowSlotID},
        })
    }
    return entries
}

// encodeKey encodes a Value as a comparable key
func encodeKey(v storage.Value) Key {
    // Prefix with type byte so different types don't compare equal
    buf := make([]byte, 1+len(v.Data))
    buf[0] = byte(v.Type)
    copy(buf[1:], v.Data)
    return Key(buf)
}

// compareKeys returns -1, 0, or 1
func compareKeys(a, b Key) int {
    minLen := len(a)
    if len(b) < minLen {
        minLen = len(b)
    }
    for i := 0; i < minLen; i++ {
        if a[i] < b[i] {
            return -1
        }
        if a[i] > b[i] {
            return 1
        }
    }
    if len(a) < len(b) {
        return -1
    }
    if len(a) > len(b) {
        return 1
    }
    return 0
}

func getSlotData(p *storage.Page, slotID int) ([]byte, bool) {
    h := p.GetHeader()
    if slotID >= int(h.NumSlots) {
        return nil, false
    }
    const headerSize = 16
    const slotSize = 4
    off := headerSize + slotID*slotSize
    dataOffset := int(p[off])<<8 | int(p[off+1])
    dataLen := int(p[off+2])<<8 | int(p[off+3])
    if dataLen == 0 {
        return nil, true
    }
    return p[dataOffset : dataOffset+dataLen], true
}
```

---

## 4. Searching a B-Tree

```go
// btree/tree.go
package btree

import (
    "fmt"
    "sort"
    "github.com/yourname/vaultdb/storage"
)

type BTree struct {
    dm       *storage.DiskManager
    rootPage storage.PageID
}

func NewBTree(dm *storage.DiskManager, rootPage storage.PageID) *BTree {
    return &BTree{dm: dm, rootPage: rootPage}
}

// Search finds the RowID for the given key value, or returns false if not found
func (t *BTree) Search(key storage.Value) (storage.RowID, bool, error) {
    encodedKey := encodeKey(key)
    pageID := t.rootPage

    for {
        page, err := t.dm.ReadPage(pageID)
        if err != nil {
            return storage.RowID{}, false, err
        }

        h := page.GetHeader()

        if h.Type == storage.PageTypeLeaf {
            // Search leaf node
            entries := ReadLeafNode(page)
            idx := sort.Search(len(entries), func(i int) bool {
                return compareKeys(entries[i].Key, encodedKey) >= 0
            })
            if idx < len(entries) && compareKeys(entries[idx].Key, encodedKey) == 0 {
                return entries[idx].RowID, true, nil
            }
            return storage.RowID{}, false, nil
        }

        // Internal node: find the right child
        entries := ReadInternalNode(page)
        childIdx := sort.Search(len(entries), func(i int) bool {
            return compareKeys(entries[i].Key, encodedKey) > 0
        })

        if childIdx >= len(entries) {
            return storage.RowID{}, false, fmt.Errorf("corrupt B-Tree: no child for key")
        }

        pageID = entries[childIdx].ChildID
    }
}

// RangeScan returns all (key, RowID) pairs where startKey <= key <= endKey
func (t *BTree) RangeScan(start, end storage.Value) ([]LeafEntry, error) {
    startKey := encodeKey(start)
    endKey := encodeKey(end)

    // Find the leaf containing startKey
    leafID, err := t.findLeaf(startKey)
    if err != nil {
        return nil, err
    }

    var results []LeafEntry

    for leafID != storage.InvalidPageID {
        page, err := t.dm.ReadPage(leafID)
        if err != nil {
            return nil, err
        }

        entries := ReadLeafNode(page)
        for _, e := range entries {
            cmp := compareKeys(e.Key, startKey)
            if cmp < 0 {
                continue
            }
            if compareKeys(e.Key, endKey) > 0 {
                return results, nil // past end
            }
            results = append(results, e)
        }

        // Follow next leaf pointer
        leafID = nextLeaf(page)
    }
    return results, nil
}

func (t *BTree) findLeaf(key Key) (storage.PageID, error) {
    pageID := t.rootPage
    for {
        page, err := t.dm.ReadPage(pageID)
        if err != nil {
            return storage.InvalidPageID, err
        }
        if page.GetHeader().Type == storage.PageTypeLeaf {
            return pageID, nil
        }
        entries := ReadInternalNode(page)
        idx := sort.Search(len(entries), func(i int) bool {
            return compareKeys(entries[i].Key, key) > 0
        })
        if idx >= len(entries) {
            idx = len(entries) - 1
        }
        pageID = entries[idx].ChildID
    }
}

func nextLeaf(p *storage.Page) storage.PageID {
    id := storage.PageID(uint16(p[6])<<8 | uint16(p[7]))
    if id == 0 {
        return storage.InvalidPageID
    }
    return id
}
```

---

## 5. Inserting into a B-Tree

```go
// Insert adds a key → RowID mapping to the B-Tree
func (t *BTree) Insert(key storage.Value, rid storage.RowID) error {
    encodedKey := encodeKey(key)

    // Find the leaf
    leafID, err := t.findLeaf(encodedKey)
    if err != nil {
        return err
    }

    leaf, err := t.dm.ReadPage(leafID)
    if err != nil {
        return err
    }

    // Encode entry: 2-byte key length + key + 8-byte pageID + 2-byte slotID
    keyLen := len(encodedKey)
    entry := make([]byte, 2+keyLen+10)
    entry[0] = byte(keyLen >> 8)
    entry[1] = byte(keyLen)
    copy(entry[2:], encodedKey)
    for i := 0; i < 8; i++ {
        entry[2+keyLen+i] = byte(rid.PageID >> (56 - 8*i))
    }
    entry[2+keyLen+8] = byte(rid.SlotID >> 8)
    entry[2+keyLen+9] = byte(rid.SlotID)

    // Try to insert into the leaf
    if leaf.FreeSpace() >= len(entry)+4 {
        // Insert in sorted position
        t.insertSorted(leaf, encodedKey, entry)
        return t.dm.WritePage(leafID, leaf)
    }

    // Leaf is full — split
    return t.splitLeaf(leafID, leaf, encodedKey, entry)
}
```

---

## 6. Node Splitting

When a node is full and we need to insert, we split it:

```go
func (t *BTree) splitLeaf(leafID storage.PageID, leaf *storage.Page, newKey Key, newEntry []byte) error {
    // 1. Collect all entries including the new one
    existing := ReadLeafNode(leaf)
    // Find insertion position
    pos := sort.Search(len(existing), func(i int) bool {
        return compareKeys(existing[i].Key, newKey) >= 0
    })
    // Insert at pos (in memory only)
    newEntries := append(existing[:pos:pos],
        append([]LeafEntry{{Key: newKey}}, existing[pos:]...)...)
    _ = newEntries // full list with new key

    // 2. Split: left half stays in old page, right half goes to new page
    mid := len(newEntries) / 2
    leftEntries := newEntries[:mid]
    rightEntries := newEntries[mid:]

    // 3. Build new right leaf page
    newLeafID, err := t.dm.AllocatePage()
    if err != nil {
        return err
    }
    newLeaf := new(storage.Page)
    newLeaf.Initialize(storage.PageTypeLeaf)
    for _, e := range rightEntries {
        keyLen := len(e.Key)
        entry := make([]byte, 2+keyLen+10)
        entry[0] = byte(keyLen >> 8)
        entry[1] = byte(keyLen)
        copy(entry[2:], e.Key)
        // encode e.RowID... (same as above)
        insertAtEnd(newLeaf, entry)
    }

    // 4. Rebuild old leaf with left half
    leaf.Initialize(storage.PageTypeLeaf)
    for _, e := range leftEntries {
        // encode and insert
        _ = e
    }

    // 5. Link leaves: old.next = new, new.next = old.next
    setNextLeaf(leaf, newLeafID)

    // 6. Push the first key of the right leaf up to the parent
    promotedKey := rightEntries[0].Key
    return t.insertIntoParent(leafID, newLeafID, promotedKey)
}

func (t *BTree) insertIntoParent(leftID, rightID storage.PageID, key Key) error {
    leftPage, err := t.dm.ReadPage(leftID)
    if err != nil {
        return err
    }
    parentID := leftPage.GetHeader().ParentID

    if parentID == storage.InvalidPageID {
        // This is the root — create a new root
        newRootID, err := t.dm.AllocatePage()
        if err != nil {
            return err
        }
        newRoot := new(storage.Page)
        newRoot.Initialize(storage.PageTypeInternal)
        // Add two entries: [leftID | key | rightID]
        // ... encode and insert
        t.rootPage = newRootID
        return t.dm.WritePage(newRootID, newRoot)
    }

    // Insert key+rightID into parent
    parent, err := t.dm.ReadPage(parentID)
    if err != nil {
        return err
    }

    if parent.FreeSpace() >= len(key)+4+8+4 {
        t.insertInternalSorted(parent, key, rightID)
        return t.dm.WritePage(parentID, parent)
    }

    // Parent is also full — split parent recursively
    return t.splitInternal(parentID, parent, key, rightID)
}

func insertAtEnd(p *storage.Page, data []byte) {
    // Simplified: append to page slots
    // (in production, insert in sorted order)
    _ = p
    _ = data
}

func setNextLeaf(p *storage.Page, nextID storage.PageID) {
    p[6] = byte(uint16(nextID) >> 8)
    p[7] = byte(uint16(nextID))
}

func (t *BTree) insertSorted(p *storage.Page, key Key, data []byte) {
    // In production: read all entries, find position, rebuild page
    // Simplified here for clarity
    _ = p
    _ = key
    _ = data
}

func (t *BTree) insertInternalSorted(p *storage.Page, key Key, childID storage.PageID) {
    _ = p
    _ = key
    _ = childID
}

func (t *BTree) splitInternal(pageID storage.PageID, page *storage.Page, key Key, rightChildID storage.PageID) error {
    // Same as splitLeaf but for internal nodes
    return nil
}
```

---

## Summary

- B-Trees keep data sorted and balanced. Every leaf is at the same depth, so all searches take the same number of page reads.
- For 10 million rows with order 200, finding any row takes only 4 page reads.
- B+ Trees (used by PostgreSQL, MySQL, SQLite) store all data at the leaf level. Leaves are linked for range scans.
- Insertions may require splitting a node. Splits propagate upward. The tree grows taller by creating a new root.
- The key insight: **disk I/O cost is proportional to tree height**, not the number of rows.

### Exercises

**Easy:** Draw a B-Tree of order 3 (max 2 keys per node) with 7 keys: 10, 20, 30, 40, 50, 60, 70. Show the state after each insertion. How many levels does it have?

**Medium:** Implement the full `insertSorted` function that reads all entries from a leaf, inserts in sorted position, and writes them back to the page. Handle the case where the new entry would overflow the page by returning an error.

**Hard:** Implement deletion from the B-Tree. When a leaf becomes less than 50% full after deletion, either borrow a key from a sibling or merge with a sibling (and pull the separator key from the parent). This may trigger a cascade of merges up to the root.
