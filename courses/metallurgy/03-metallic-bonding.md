# Chapter 03: Metallic Bonding — Why Metals Are Special

> **"The metallic bond is one of the most democratic bonds in nature. Every atom contributes its valence electrons to a common pool, and every atom benefits from that pool. It is this collective arrangement that gives metals their unique combination of strength, ductility, and conductivity."**

---

## Table of Contents

1. [The Three Types of Primary Bonds](#1-the-three-types-of-primary-bonds)
2. [The Metallic Bond — Electron Sea Model](#2-the-metallic-bond--electron-sea-model)
3. [Why Metals Conduct Electricity and Heat](#3-why-metals-conduct-electricity-and-heat)
4. [Why Metals Are Ductile — The Key Insight](#4-why-metals-are-ductile--the-key-insight)
5. [Why Metals Have Luster](#5-why-metals-have-luster)
6. [Metallic Bond Strength and Melting Point](#6-metallic-bond-strength-and-melting-point)
7. [Band Theory — A More Accurate Picture](#7-band-theory--a-more-accurate-picture)
8. [How Alloying Affects the Bond](#8-how-alloying-affects-the-bond)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. The Three Types of Primary Bonds

Before understanding metallic bonding, let's place it in context. There are three primary (strong) bond types:

### Ionic Bond
One atom transfers electrons to another. Both become ions — one positive (cation), one negative (anion). They attract electrostatically.

Example: NaCl. Na loses one electron → Na⁺. Cl gains it → Cl⁻. Strong bond, but **non-directional**.

Properties: hard, brittle, high melting point, electrically insulating (ions can't move), dissolves in polar solvents.

### Covalent Bond
Two atoms **share** electrons. The electrons are localized between the two atoms, forming a directional bond.

Example: Diamond (C-C bonds). Extremely strong, but totally brittle — breaking one bond means the structure fails.

Properties: very hard, brittle, generally insulating (electrons are fixed), high melting point.

### Metallic Bond
Atoms release their valence electrons into a **shared electron cloud** (the "electron sea"). The positively charged ion cores sit in this sea of electrons and are held together by the electrostatic attraction between the positive cores and the negative electron cloud.

Properties: **conductors** (electrons move freely), **ductile** (bond is non-directional), shiny, and — crucially — **mixed strength and ductility simultaneously**, which neither ionic nor covalent bonds provide.

| Bond Type | Electron Sharing | Directionality | Typical Materials | Ductile? | Conductive? |
|-----------|-----------------|----------------|-------------------|----------|-------------|
| Ionic | Transferred | Non-directional | NaCl, MgO, Al₂O₃ | No | No (solid) |
| Covalent | Shared, localized | Directional | Diamond, SiC, Si₃N₄ | No | No |
| Metallic | Shared, delocalized | Non-directional | Fe, Al, Ni, Cu | **Yes** | **Yes** |

---

## 2. The Metallic Bond — Electron Sea Model

The **Drude electron sea model** (1900) is simple but captures the essential physics:

```
   +  +  +  +  +  +  +  +
  Fe Fe Fe Fe Fe Fe Fe Fe
   +  +  +  +  +  +  +  +
  Fe Fe Fe Fe Fe Fe Fe Fe
   +  +  +  +  +  +  +  +
  Fe Fe Fe Fe Fe Fe Fe Fe
   +  +  +  +  +  +  +  +

   ~~~~~~~~~ electron sea ~~~~~~~~~~~
   Each Fe atom: [Ar] 3d⁶ 4s² → releases 2 electrons
   Ion cores: Fe²⁺, each sitting in the sea
```

**How it works:**
1. Each metal atom releases its valence electrons.  
   Iron releases 2 electrons (from 4s²); aluminum releases 3 (from 3s² 3p¹); copper releases 1 (from 4s¹).
2. These electrons become **delocalized** — they don't belong to any one atom but roam freely throughout the entire crystal.
3. The remaining ion cores (positively charged) sit in a regular array.
4. The **attraction** between the mobile electron cloud and the positive ion cores holds the structure together.

This is fundamentally different from ionic bonding: in NaCl, the electron is permanently with the chloride ion. In iron, the electron is shared by all 10²³ iron atoms simultaneously.

### Why the Electron Sea is Special

**Non-directional:** The bond has no preferred direction. An Fe atom is attracted to all the electrons around it equally. This means atoms can **slide past each other** while still being in the electron sea — this is why metals are ductile.

Compare to a covalent bond: if you try to move atoms past each other in diamond, you break the specific directed C-C bonds → shatters. In iron, the electron sea rearranges around the new positions → deforms without breaking.

---

## 3. Why Metals Conduct Electricity and Heat

### Electrical Conductivity

When you apply a voltage across a metal, the free electrons — already mobile — drift toward the positive terminal. This drift of charge is electric current.

```
No voltage:   → ← → ↑ ↓ → ← ↑   (random thermal motion)
With voltage: → → → → → → → →   (net drift, superimposed on thermal)
```

Conductivity depends on:
- **Number of free electrons per atom**: Cu (1), Al (3), Fe (2) — more electrons → better conductor
- **Mean free path**: how far an electron travels before scattering off a defect or vibrating atom
- **Temperature**: higher T → more atom vibration → more electron scattering → **resistivity increases with temperature** (opposite of semiconductors)

**The best metal conductors (at room temperature):**

| Metal | Electrical Conductivity (MS/m) | Notes |
|-------|-------------------------------|-------|
| Silver | 63.0 | Best, but expensive |
| Copper | 59.6 | Standard; 1 free e⁻/atom |
| Gold | 45.2 | Expensive; used for contacts |
| Aluminum | 37.7 | Light; used for power lines |
| Nickel | 14.3 | Magnetic; lower conductivity |
| Iron | 10.0 | Many scattering sites (d electrons) |
| Stainless (316) | 1.4 | Very low; alloying disrupts sea |

Note: stainless steel conducts 40× worse than copper — all the Cr and Ni disruptions scatter electrons.

### Thermal Conductivity

The same free electrons that carry charge also carry heat energy. When you heat one end of a copper rod, the electrons on that end gain kinetic energy and quickly spread it through the entire material.

This is why: **good electrical conductors are also good thermal conductors** (the Wiedemann-Franz law: κ/σT ≈ constant for metals).

**Thermal conductivity (W/m·K):**
- Silver: 429; Copper: 401; Aluminum: 237; Iron: 80; Nickel: 91; Inconel 718: 11.4

Again, the superalloy has terrible thermal conductivity — 35× worse than copper. This matters enormously for turbine blade design: the alloy resists heat transfer, so you need active cooling channels inside the blade.

---

## 4. Why Metals Are Ductile — The Key Insight

This is the most important consequence of metallic bonding. Understanding it properly sets up Chapter 12 (Dislocations).

**The Problem with Ionic and Covalent Bonds:**

Consider an ionic crystal. Each ion is surrounded by ions of opposite charge. If you apply a shear force and shift one layer by one atom spacing:

```
Before:   + - + - + -      After shift:   + - + - + -
          - + - + - +                     - + - + - +
          + - + - + -  →                  + - + - + -
          ↑                                    ↑
          shift this plane by one step
```

After shifting, **like charges** end up next to each other (+/+ and -/-). Repulsion is enormous. The crystal shatters. This is why table salt shatters when you hit it.

**In a Metal:**

Apply the same shear force to an iron crystal. Shift one layer by one atom spacing:

```
Before:   Fe Fe Fe Fe         After:   Fe Fe Fe Fe
          Fe Fe Fe Fe    →             Fe Fe Fe Fe
          Fe Fe Fe Fe                  Fe Fe Fe Fe
```

After shifting, the atoms are still in identical chemical environments — surrounded by other iron atoms and bathed in the same electron sea. The electron sea **instantly rearranges** around the new configuration. The energy is virtually the same before and after. The crystal can deform without the bond breaking.

This is the atomic origin of metallic ductility:
1. The electron sea is isotropic (same in all directions)
2. Atoms can move without encountering a repulsion spike
3. Deformation is progressive — not sudden catastrophic failure

> Real metals deform via **dislocations** (Chapter 12) which makes plastic deformation even easier than this simple-shift model suggests. But the fundamental reason dislocations can move is the non-directional metallic bond.

---

## 5. Why Metals Have Luster

When light (electromagnetic radiation) hits a metal surface, the free electrons oscillate in response to the oscillating electric field of the light wave. These oscillating electrons re-emit light in the visible range.

The result: metals reflect most visible wavelengths → they look shiny ("metallic luster").

Most metals reflect all visible wavelengths equally → they look silver-grey (iron, aluminum, chromium, nickel).

**Colored metals are exceptions:**
- **Gold** (Au): reflects red and yellow more than blue/violet → golden color. Due to relativistic effects on its 5d/6s electrons.
- **Copper** (Cu): reflects red and orange more than blue → copper/orange color. Due to its specific band structure near the Fermi level.
- **Brass** (Cu-Zn): yellow color, tunable by Zn content.

---

## 6. Metallic Bond Strength and Melting Point

The strength of the metallic bond correlates with:
1. **Number of valence electrons contributed** — more electrons → stronger attraction to ion cores
2. **Size of the ion core** — smaller core → electrons are closer to positive charge → stronger bond
3. **d-electron involvement** — transition metals with partially filled d shells have d-d bonding contributions that enormously increase bond strength

This is why the **melting points of transition metals peak in the middle groups** (Group VIB: Cr, Mo, W):

```
Period 4 melting points (°C):
Ca(839) → Sc(1541) → Ti(1668) → V(1910) → Cr(1907) → Mn(1246) → Fe(1538) → Co(1495) → Ni(1455) → Cu(1085)
                                                ↑
                                            Peak at Cr (Group VIB)
                                            Half-filled d shell
                                            Strongest d-d bonding
```

The peak occurs at the **half-filled d shell** (5 d electrons), where exchange interaction is maximized.

**Period 5 and 6** follow the same pattern but peak even higher:
- Mo (Period 5, Group VIB): 2623°C
- W (Period 6, Group VIB): 3422°C — the highest melting metal

This is why tungsten is used in incandescent light bulb filaments and in tools that cut at high speed. Its metallic bond is the strongest known among metals.

---

## 7. Band Theory — A More Accurate Picture

The electron sea model is a simplification. The more accurate picture is **energy band theory**:

When N atoms come together to form a crystal, their discrete atomic energy levels broaden into **bands** of very closely spaced energy levels:

```
Isolated atom        In a solid
    ___                ═══════ conduction band (empty or partially filled)
    ___         →         ← band gap (no allowed energies)
    ___                ═══════ valence band (filled)
```

In a **metal**, the valence band and conduction band **overlap** — there is no gap. Electrons at the Fermi level are free to move into nearby empty states → conduction, ductility.

In an **insulator**, the gap is large (>3 eV). No electrons can jump at room temperature → no conduction.

In a **semiconductor** (Si, Ge), the gap is small (1–2 eV) — some electrons thermally jump → limited conduction, tunable by doping.

### Why the d-Band Matters

For transition metals like Ni, Fe, Co, the **d-band** partially overlaps the s-band. The d-band is narrower (d electrons are more localized) but holds up to 10 electrons. The degree of d-band filling and its position relative to the Fermi level controls:
- Catalytic reactivity
- Magnetic moment
- High-temperature strength (partially-filled d bands increase bond strength)

Modern superalloy design increasingly relies on **d-electron theory** to predict alloying effects computationally (DFT calculations, Chapter 64).

---

## 8. How Alloying Affects the Bond

When you add a second element to a metal:

### Solid Solution Strengthening (preview of Ch 13)
Foreign atoms (solute) disturb the regular electron sea. This creates local **strain fields** and **electronic perturbations** that scatter the electrons → increases electrical resistivity, and more importantly, impedes dislocation motion → stronger alloy.

Example: Pure copper has yield strength ~70 MPa. Add 30% Zn → brass: yield strength ~200 MPa. The Zn atoms distort the electron sea → dislocations can't move as easily.

### Intermetallic Compounds
When two metals have very different electronegativities, the bonding changes character from purely metallic toward covalent or ionic. The compound is ordered (specific atom arrangements) and typically **hard and brittle**.

Example: Ni₃Al (gamma prime, γ′) — the key strengthening phase in superalloys.
- Ni and Al have very different sizes and some electronegativity difference
- Ni₃Al forms an ordered L1₂ structure (Al atoms at cube corners, Ni at face centers)
- The ordered structure means dislocations must overcome **antiphase boundary energy** to move
- This makes γ′ extraordinarily resistant to dislocation motion at high temperature

The γ′ hardening mechanism that powers jet engines is fundamentally a consequence of the partial ionic-covalent character introduced into the metallic bond by adding aluminum to nickel. We'll explore this in detail in Chapter 44.

---

## Summary

- Three primary bonds: **ionic** (electron transfer, brittle), **covalent** (electron sharing, brittle), **metallic** (electron sea, ductile + conductive).
- The **metallic bond** = positive ion cores immersed in a delocalized electron sea.
- **Non-directional** → ductility: atoms can slide past each other without breaking bonds; the electron sea rearranges.
- **Free electrons** → electrical and thermal conductivity (good conductors are also good thermal conductors — Wiedemann-Franz law).
- **Free electrons oscillating with light** → metallic luster; colored metals (Au, Cu) are exceptions due to band structure.
- **Bond strength** peaks at half-filled d shells (Cr, Mo, W) → highest melting points; d-d bonding is critical.
- **Band theory** is the accurate quantum picture: metals have overlapping valence and conduction bands — no gap.
- **Alloying** disrupts the electron sea → increases strength but decreases conductivity. Ordered intermetallics (like Ni₃Al = γ′) combine metallic and covalent character → exceptional high-temperature strength.

**Next chapter:** The metallic bond explains individual atoms. But metals are crystalline — the atoms pack into regular, repeating 3D arrays. The geometry of that packing determines density, slip systems, and ultimately much of the mechanical behavior.

---

## Exercises

1. Explain in your own words why copper conducts electricity but table salt does not, using only bonding concepts.
2. Why does stainless steel (Inconel 718) have 35× lower thermal conductivity than copper, despite both being metals?
3. Copper has electrical conductivity 59.6 MS/m. You add 30% Zn to make brass. Predict: does conductivity increase, decrease, or stay the same? Why? (Look up actual value to check.)
4. Why does tungsten have a higher melting point than iron, even though both are Group VIII neighbors in the periodic table? (Hint: compare their d-electron counts.)
5. In the electron sea model, what happens to bond character when you add aluminum (EN = 1.61) to nickel (EN = 1.91)? Does this make the alloy harder or softer at high temperature? Why?
