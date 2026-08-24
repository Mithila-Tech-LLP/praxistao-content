# Chapter 40: Low-Level Design — Parking Lot, Elevator & Library System

LLD interviews test your ability to translate real-world systems into clean, extensible object-oriented (or interface-based) code. The focus is on classes/interfaces, relationships, and extensibility — not implementation details.

## Table of Contents

1. [How to Approach LLD Interviews](#1-how-to-approach-lld-interviews)
2. [Parking Lot Design](#2-parking-lot-design)
3. [Elevator System Design](#3-elevator-system-design)
4. [Library Management System](#4-library-management-system)
5. [Summary](#summary)

---

## 1. How to Approach LLD Interviews

```
Step 1: Clarify requirements (5 min)
  "What types of vehicles? What fee structure? Multiple floors?"
  
Step 2: Identify core entities/classes (5 min)
  List all the nouns: ParkingLot, Floor, Spot, Vehicle, Ticket, Payment
  
Step 3: Identify relationships (5 min)
  ParkingLot HAS-MANY Floor
  Floor HAS-MANY Spot  
  Ticket BELONGS-TO Spot + Vehicle
  
Step 4: Define interfaces first, then implementations (15 min)
  Start with behavior: what does each entity DO?
  
Step 5: Show one complete flow (10 min)
  "Let me trace through parking a car: vehicle arrives → find spot → issue ticket..."

Key principles:
  - Every decision: explain the why and the trade-off
  - Extensibility: "I'm using an interface here so we can add new vehicle types later"
  - Don't over-engineer: match complexity to requirements
```

---

## 2. Parking Lot Design

### Requirements
```
- Multiple floors
- Different vehicle sizes: motorcycle, car, large (truck/bus)
- Different spot sizes: small, medium, large
- Issue ticket on entry, calculate fee on exit
- Find available spot for vehicle type
- Vehicle: small fits any spot, car fits medium/large, truck fits large only
```

### Entities & Relationships

```
ParkingLot
  └── []Floor
       └── []ParkingSpot
            └── Vehicle (when occupied)

Ticket
  ├── Vehicle
  ├── ParkingSpot
  ├── EntryTime
  └── ExitTime

Payment
  ├── Ticket
  └── Amount
```

### Go Implementation

```go
package parkinglot

import (
    "errors"
    "sync"
    "time"
)

// Vehicle types
type VehicleType string
const (
    Motorcycle VehicleType = "motorcycle"
    Car        VehicleType = "car"
    Truck      VehicleType = "truck"
)

// Spot sizes
type SpotSize string
const (
    SmallSpot  SpotSize = "small"
    MediumSpot SpotSize = "medium"
    LargeSpot  SpotSize = "large"
)

type Vehicle struct {
    LicensePlate string
    Type         VehicleType
}

func (v *Vehicle) FitsInSpot(spot SpotSize) bool {
    switch v.Type {
    case Motorcycle:
        return true // fits anywhere
    case Car:
        return spot == MediumSpot || spot == LargeSpot
    case Truck:
        return spot == LargeSpot
    }
    return false
}

type ParkingSpot struct {
    ID       string
    Size     SpotSize
    FloorNum int
    vehicle  *Vehicle
    mu       sync.Mutex
}

func (s *ParkingSpot) IsAvailable() bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.vehicle == nil
}

func (s *ParkingSpot) Park(v *Vehicle) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.vehicle != nil {
        return errors.New("spot already occupied")
    }
    s.vehicle = v
    return nil
}

func (s *ParkingSpot) Unpark() *Vehicle {
    s.mu.Lock()
    defer s.mu.Unlock()
    v := s.vehicle
    s.vehicle = nil
    return v
}

type Ticket struct {
    ID        string
    Vehicle   *Vehicle
    Spot      *ParkingSpot
    EntryTime time.Time
    ExitTime  *time.Time
}

// FeeCalculator interface — easy to swap different pricing models
type FeeCalculator interface {
    Calculate(ticket *Ticket) float64
}

type HourlyFeeCalculator struct {
    RatePerHour map[VehicleType]float64
}

func (c *HourlyFeeCalculator) Calculate(ticket *Ticket) float64 {
    if ticket.ExitTime == nil {
        return 0
    }
    duration := ticket.ExitTime.Sub(ticket.EntryTime).Hours()
    rate := c.RatePerHour[ticket.Vehicle.Type]
    return duration * rate
}

type Floor struct {
    Number int
    spots  []*ParkingSpot
}

func (f *Floor) FindAvailableSpot(v *Vehicle) *ParkingSpot {
    for _, spot := range f.spots {
        if spot.IsAvailable() && v.FitsInSpot(spot.Size) {
            return spot
        }
    }
    return nil
}

type ParkingLot struct {
    Name      string
    floors    []*Floor
    tickets   map[string]*Ticket // ticketID → ticket
    calculator FeeCalculator
    mu        sync.RWMutex
}

func (lot *ParkingLot) Park(vehicle *Vehicle) (*Ticket, error) {
    lot.mu.Lock()
    defer lot.mu.Unlock()
    
    // Find available spot across all floors
    var foundSpot *ParkingSpot
    for _, floor := range lot.floors {
        if spot := floor.FindAvailableSpot(vehicle); spot != nil {
            foundSpot = spot
            break
        }
    }
    
    if foundSpot == nil {
        return nil, errors.New("no available spot for vehicle type")
    }
    
    if err := foundSpot.Park(vehicle); err != nil {
        return nil, err
    }
    
    ticket := &Ticket{
        ID:        generateTicketID(),
        Vehicle:   vehicle,
        Spot:      foundSpot,
        EntryTime: time.Now(),
    }
    lot.tickets[ticket.ID] = ticket
    return ticket, nil
}

func (lot *ParkingLot) Exit(ticketID string) (float64, error) {
    lot.mu.Lock()
    defer lot.mu.Unlock()
    
    ticket, exists := lot.tickets[ticketID]
    if !exists {
        return 0, errors.New("ticket not found")
    }
    
    now := time.Now()
    ticket.ExitTime = &now
    ticket.Spot.Unpark()
    
    fee := lot.calculator.Calculate(ticket)
    delete(lot.tickets, ticketID)
    
    return fee, nil
}

func (lot *ParkingLot) AvailableSpots() map[SpotSize]int {
    lot.mu.RLock()
    defer lot.mu.RUnlock()
    
    counts := map[SpotSize]int{SmallSpot: 0, MediumSpot: 0, LargeSpot: 0}
    for _, floor := range lot.floors {
        for _, spot := range floor.spots {
            if spot.IsAvailable() {
                counts[spot.Size]++
            }
        }
    }
    return counts
}
```

---

## 3. Elevator System Design

### Requirements
```
- Multiple elevators, multiple floors
- Requests: inside (destination floor) or outside (floor + direction)
- Elevator states: idle, moving up, moving down, stopped
- Algorithm: scan algorithm (SCAN/LOOK) — service requests in one direction, then reverse
```

### Core Entities

```go
type Direction string
const (
    Up   Direction = "up"
    Down Direction = "down"
    Idle Direction = "idle"
)

type Request struct {
    Floor     int
    Direction Direction // for external requests
    Type      string    // "internal" or "external"
}

type Elevator struct {
    ID          int
    CurrentFloor int
    Direction   Direction
    Destinations map[int]bool // floors to stop at
    mu          sync.Mutex
}

func (e *Elevator) AddDestination(floor int) {
    e.mu.Lock()
    defer e.mu.Unlock()
    e.Destinations[floor] = true
}

// SCAN algorithm: service all requests in current direction,
// then reverse when no more in that direction
func (e *Elevator) NextFloor() int {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    if len(e.Destinations) == 0 {
        e.Direction = Idle
        return e.CurrentFloor
    }
    
    if e.Direction == Up {
        // Find smallest floor above current
        next := -1
        for floor := range e.Destinations {
            if floor > e.CurrentFloor && (next == -1 || floor < next) {
                next = floor
            }
        }
        if next != -1 { return next }
        
        // No floors above, reverse direction
        e.Direction = Down
    }
    
    if e.Direction == Down {
        // Find largest floor below current
        next := -1
        for floor := range e.Destinations {
            if floor < e.CurrentFloor && (next == -1 || floor > next) {
                next = floor
            }
        }
        if next != -1 { return next }
        
        // No floors below, reverse
        e.Direction = Up
    }
    
    return e.CurrentFloor
}

type ElevatorSystem struct {
    elevators []*Elevator
}

// Dispatch: assign request to the best elevator
// Simple strategy: assign to elevator with minimum additional travel
func (s *ElevatorSystem) Dispatch(req Request) *Elevator {
    var bestElevator *Elevator
    bestCost := 1<<31 - 1 // max int
    
    for _, e := range s.elevators {
        cost := s.estimateCost(e, req)
        if cost < bestCost {
            bestCost = cost
            bestElevator = e
        }
    }
    
    if bestElevator != nil {
        bestElevator.AddDestination(req.Floor)
    }
    return bestElevator
}

func (s *ElevatorSystem) estimateCost(e *Elevator, req Request) int {
    // Simple heuristic: distance from current position
    dist := req.Floor - e.CurrentFloor
    if dist < 0 { dist = -dist }
    
    // Penalize elevators moving in opposite direction
    if e.Direction == Up && req.Floor < e.CurrentFloor {
        dist += 10 // higher penalty
    }
    if e.Direction == Down && req.Floor > e.CurrentFloor {
        dist += 10
    }
    return dist
}
```

---

## 4. Library Management System

### Requirements
```
- Books (physical copies), Members, Librarians
- Members can borrow books (max 5 at once, 30-day period)
- Librarians can add/remove books, manage members
- Search books by title, author, ISBN
- Late return fine calculation
- Reservation: reserve a book when all copies are checked out
```

### Core Entities

```go
type Book struct {
    ISBN      string
    Title     string
    Author    string
    Copies    []*BookCopy
}

type BookCopy struct {
    ID     string
    ISBN   string
    Status string // "available", "borrowed", "reserved"
}

type Member struct {
    ID           string
    Name         string
    Email        string
    Borrowings   []*Borrowing
    Reservations []*Reservation
}

const (
    MaxBorrowingsPerMember = 5
    LoanPeriodDays         = 30
    FinePerDay             = 0.25 // $0.25 per day late
)

type Borrowing struct {
    ID         string
    MemberID   string
    BookCopyID string
    DueDate    time.Time
    ReturnDate *time.Time
}

func (b *Borrowing) Fine() float64 {
    if b.ReturnDate == nil {
        // Still borrowed
        if time.Now().After(b.DueDate) {
            days := time.Since(b.DueDate).Hours() / 24
            return days * FinePerDay
        }
        return 0
    }
    if b.ReturnDate.After(b.DueDate) {
        days := b.ReturnDate.Sub(b.DueDate).Hours() / 24
        return days * FinePerDay
    }
    return 0
}

type Reservation struct {
    ID         string
    MemberID   string
    ISBN       string
    ReservedAt time.Time
    Status     string // "active", "fulfilled", "expired"
}

// Repository interfaces for testability
type BookRepository interface {
    FindByISBN(isbn string) (*Book, error)
    FindByTitle(title string) ([]*Book, error)
    FindByAuthor(author string) ([]*Book, error)
    Save(book *Book) error
}

type MemberRepository interface {
    FindByID(id string) (*Member, error)
    Save(member *Member) error
}

type LibraryService struct {
    books   BookRepository
    members MemberRepository
    mu      sync.Mutex
}

func (svc *LibraryService) BorrowBook(memberID, isbn string) (*Borrowing, error) {
    svc.mu.Lock()
    defer svc.mu.Unlock()
    
    member, err := svc.members.FindByID(memberID)
    if err != nil { return nil, err }
    
    // Check member's borrowing limit
    activeBorrowings := 0
    for _, b := range member.Borrowings {
        if b.ReturnDate == nil { activeBorrowings++ }
    }
    if activeBorrowings >= MaxBorrowingsPerMember {
        return nil, errors.New("borrowing limit reached")
    }
    
    // Find available copy
    book, err := svc.books.FindByISBN(isbn)
    if err != nil { return nil, err }
    
    var availableCopy *BookCopy
    for _, copy := range book.Copies {
        if copy.Status == "available" {
            availableCopy = copy
            break
        }
    }
    
    if availableCopy == nil {
        return nil, errors.New("no copies available — consider reserving")
    }
    
    availableCopy.Status = "borrowed"
    
    borrowing := &Borrowing{
        ID:         generateID(),
        MemberID:   memberID,
        BookCopyID: availableCopy.ID,
        DueDate:    time.Now().AddDate(0, 0, LoanPeriodDays),
    }
    
    member.Borrowings = append(member.Borrowings, borrowing)
    svc.members.Save(member)
    svc.books.Save(book)
    
    return borrowing, nil
}
```

---

## Summary

- **LLD approach:** clarify → identify entities → define relationships → interfaces first → trace one flow.
- **Parking Lot:** ParkingLot has floors, floors have spots, spots hold vehicles. FeeCalculator interface for extensible pricing.
- **Elevator:** SCAN algorithm (service current direction, reverse when done). Dispatch to elevator with minimum estimated cost.
- **Library:** Book has copies, Member has borrowings. Repository interfaces for data access abstraction.
- Always explain the interface choice: "FeeCalculator is an interface so we can add dynamic pricing without changing the parking lot logic."
- Thread safety: use mutexes for shared state (spots being parked, ticket creation).
