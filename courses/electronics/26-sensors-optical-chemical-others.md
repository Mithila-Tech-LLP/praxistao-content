# Chapter 26: Optical, Chemical, and Other Sensors

## 26.1 Light-Dependent Resistor (LDR)

### Working Principle

An LDR (also called photoresistor or CdS cell) is a passive component whose **resistance decreases** with increasing light intensity:

```
CdS (Cadmium Sulfide) semiconductor:
  In darkness: high resistance (1 MΩ - 10 MΩ)
  In bright light: low resistance (100 Ω - 1 kΩ)

Photoconductivity:
  Photons with energy > bandgap → create electron-hole pairs
  More photons → more carriers → lower resistance

Response: Not linear, follows power law:
  R = C × E^(-γ)
  Where E = illuminance (lux), C and γ are material constants
  γ ≈ 0.7-0.9 for CdS
```

**LDR Circuit:**
```
VCC (3.3V or 5V)
      │
   [10kΩ] ← Fixed resistor (voltage divider)
      │
      ├── Vout → ADC
      │
    [LDR] ← Resistance varies with light
      │
     GND

Vout = VCC × R_LDR / (R_fixed + R_LDR)

Bright light: R_LDR low → Vout low
Dark: R_LDR high → Vout high
```

**LDR characteristics:**
- Response time: slow (10-100 ms) — can't measure fast-switching light
- Spectral response: peak at 540-560 nm (yellow-green, similar to human eye)
- RoHS concern: CdS is restricted in EU/UK (cadmium is toxic) — alternatives: PbS, InGaAs
- Cost: $0.20-$1.00

**Arduino code:**
```cpp
int ldr_raw = analogRead(A0);
float voltage = ldr_raw * (5.0 / 1023.0);
float resistance = 10000.0 * (5.0 - voltage) / voltage;  // R_LDR from divider
// Map to lux using calibration or lookup table
```

**Applications:**
- Street light automatic on/off (turns on when dark)
- Night light
- Automatic display brightness
- Photography light meter
- Burglar alarm (light beam detection)

---

## 26.2 Photodiode and Phototransistor

### Photodiode

A reverse-biased PN junction that generates current when illuminated:

```
Photodiode physics:
  Photon absorption in depletion region → electron-hole pair
  Electric field sweeps carriers apart → photocurrent

I_ph = R_λ × P_opt (A)

Where:
  R_λ = responsivity (A/W, wavelength-dependent)
  P_opt = optical power (W)

Typical responsivity:
  Silicon @ 850 nm: ~0.6 A/W
  Silicon @ 630 nm (red): ~0.4 A/W
  Silicon @ 400 nm (UV): ~0.1 A/W
```

**Photodiode modes:**

**Photovoltaic mode (zero-bias):**
```
+  [photodiode] −
     │
    [Load R]
     │
    GND

Generates voltage (like solar cell)
Very linear, best for precision light measurement
Current: pA to µA
```

**Photoconductive mode (reverse bias):**
```
VCC ──────── Vout ──── [TIA op-amp] ──── output
             │
          [photodiode]
             │ (cathode to VCC for reverse bias)
            GND

Faster response (lower junction capacitance)
Used for high-speed optical receivers
Dark current increases with bias
```

**Transimpedance Amplifier (TIA):**
```
         Rf (1 MΩ for low light, 1 kΩ for bright)
    ┌────/\/\/────┐
    │            │
I_ph→ [−] op-amp ├── Vout = -I_ph × Rf
      [+]        │
       │         │
      GND

I_ph = 1 µA, Rf = 1 MΩ → Vout = 1V
I_ph = 1 mA, Rf = 1 kΩ → Vout = 1V
```

**Phototransistor:**
- Photodiode + built-in transistor amplification
- 100-1000× more sensitive than photodiode
- Slower (base junction capacitance)
- Used in: optocouplers, TV remote receivers, light barriers

---

## 26.3 Color Sensor — TCS34725

Measures color components (R, G, B, Clear):

**Working principle:**
```
4 photodiodes with different color filters:
  R filter: passes red light
  G filter: passes green light
  B filter: passes blue light
  Clear: no filter (total light)

Each has dedicated 16-bit ADC
Ratio of R:G:B → color identification
```

**TCS34725 Specifications:**
| Parameter     | Value               |
|--------------|---------------------|
| Color channels| RGBC (4-channel)    |
| ADC          | 16-bit each         |
| Interface    | I2C (0x29)          |
| Voltage      | 3.3V                |
| Integration time | 2.4 ms - 614 ms |
| Gain         | 1×, 4×, 16×, 60×   |
| IR blocking  | Built-in filter     |

**Color temperature and lux:**
```c
#include <Adafruit_TCS34725.h>
Adafruit_TCS34725 tcs = Adafruit_TCS34725(TCS34725_INTEGRATIONTIME_700MS,
                                            TCS34725_GAIN_1X);

void loop() {
    uint16_t r, g, b, c;
    tcs.getRawData(&r, &g, &b, &c);

    uint16_t colorTemp = tcs.calculateColorTemperature(r, g, b);
    uint16_t lux = tcs.calculateLux(r, g, b);

    Serial.print("Color Temp: "); Serial.print(colorTemp); Serial.println(" K");
    Serial.print("Lux: "); Serial.println(lux);
}
```

**CCT (Color Correlated Temperature):**
- 2700K: Warm white (incandescent)
- 4000K: Cool white (office fluorescent)
- 6500K: Daylight
- Used in: smart lighting control, white balance in cameras

**Applications:**
- Sorting by color (industrial automation)
- Food quality inspection
- Color matching (painting, printing)
- Plant growth lighting control

---

## 26.4 Optical Encoder

### Rotary Encoder — Position and Speed Detection

```
Incremental encoder:
  Slotted disk with N slots attached to rotating shaft
  IR LED shines through slots
  Phototransistor detects pulses

  With 2 tracks (A and B, 90° offset):
     Direction detection possible:
       A leads B → clockwise
       B leads A → counter-clockwise

    ┌─┐   ┌─┐   ┌─┐   ┌─┐
A:  │ │   │ │   │ │   │ │
    ┘ └───┘ └───┘ └───┘ └
      ┌─┐   ┌─┐   ┌─┐
B:    │ │   │ │   │ │
    ──┘ └───┘ └───┘ └────
     (B follows A by 90°)
                        ← clockwise rotation
```

**Resolution:**
- 600 PPR (Pulses Per Revolution) encoder on 100 RPM motor
- 600 × 100 RPM / 60 = 1000 Hz pulse frequency
- Timer input capture measures pulse period → RPM

**Quadrature decoding:**
```c
// State machine decoding (software):
int encoder_value = 0;
uint8_t prev_state = 0;

void EXTI_IRQHandler(void) {
    uint8_t curr_state = (READ_A() << 1) | READ_B();
    // State transitions: 00→01→11→10→00 = CCW
    //                    00→10→11→01→00 = CW
    if ((prev_state == 0b00 && curr_state == 0b01) ||
        (prev_state == 0b01 && curr_state == 0b11)) encoder_value++;
    if ((prev_state == 0b00 && curr_state == 0b10) ||
        (prev_state == 0b10 && curr_state == 0b11)) encoder_value--;
    prev_state = curr_state;
}
```

**Absolute encoder:**
- Each position has unique binary code (Gray code)
- No need to count pulses from home position
- Knows absolute position even after power cycle
- More expensive

**Used in:** CNC machines, servo motors, industrial robots, 3D printers

---

## 26.5 Photoelectric Sensors (Industrial)

Used in factories for object detection:

**Through-beam:**
```
Emitter ─ IR beam ─ Receiver
If beam broken → object detected
Range: up to 60 m
Most reliable (separate emitter and receiver)
```

**Retro-reflective:**
```
Combined emitter+receiver ─ beam ─ Reflector
If beam broken → object detected
Easier to install (only one side needs wiring)
```

**Diffuse:**
```
Combined unit sends beam
Detects reflection from object itself
Short range, depends on target reflectivity
```

---

## 26.6 Gas Sensors — MQ Series

### MQ-2 — Combustible Gas Sensor

**Working principle (metal oxide semiconductors):**
```
SnO₂ (tin dioxide) heated to ~300°C by internal heater coil:
  In clean air: SnO₂ has high resistance (barrier at grain boundaries)
  Gas adsorbs onto surface → oxygen replaced → electron donors
  → Resistance decreases!

RS/R0 ratio: reference ratio in clean air
Lower ratio → higher gas concentration
```

**MQ-2 Detects:**
- LPG (liquefied petroleum gas)
- Methane (CH₄)
- Hydrogen (H₂)
- Smoke
- Alcohol
- Carbon monoxide (partially)

**Specifications:**
| Parameter     | Value               |
|--------------|---------------------|
| Heater voltage| 5V ±0.1V            |
| Heater current| 150 mA (0.75W)      |
| Load resistance| Adjustable 10-47 kΩ|
| Preheat time | 20 seconds min      |
| Operating temp| -10 to +50°C       |

**Circuit:**
```
VCC (5V)
    │
   [5V Heater circuit] (H pins)

VCC (5V)
    │
   [RL=10kΩ]  ← Load resistor (set sensitivity)
    │
    ├── Vout → ADC
    │
   [MQ Sensor] (A pins)
    │
   GND

Vout = VCC × RL / (Rs + RL)
Where Rs = sensor resistance (decreases with gas)
Lower Rs (more gas) → lower Vout? No, wait:
VCC → Rs → Vout → RL → GND
Vout = VCC × RL / (Rs + RL)
Rs low (gas present) → RL/(low+RL) → higher Vout! YES
```

**Warm-up required!**
```cpp
void setup() {
    Serial.println("Warming up MQ-2...");
    delay(20000);  // 20 second warm-up required!
    // Read baseline
    float clean_air = analogRead(A0);
    Serial.println("Ready!");
}
```

### MQ Sensor Family

| Sensor | Primary Detection       | Other         |
|--------|------------------------|---------------|
| MQ-2   | LPG, Propane, Hydrogen | Methane, Smoke|
| MQ-3   | Alcohol (ethanol)      | Benzene, CH4  |
| MQ-4   | Methane (natural gas)  | Other gases   |
| MQ-5   | LPG, Natural gas       |               |
| MQ-6   | LPG, Butane            |               |
| MQ-7   | Carbon monoxide (CO)   |               |
| MQ-8   | Hydrogen (H₂)          |               |
| MQ-9   | Carbon monoxide + LPG  |               |
| MQ-131 | Ozone (O₃)             |               |
| MQ-135 | Air quality (NH₃, NOx, alcohol, CO₂ roughly) |   |
| MQ-136 | Hydrogen sulfide (H₂S) |               |
| MQ-137 | Ammonia (NH₃)          |               |

**Limitations of MQ sensors:**
- Cross-sensitivity: responds to multiple gases
- Requires calibration in fresh air and known gas concentrations
- Long warm-up time
- High power consumption (heater!)
- Not suitable for precise concentration measurement
- Good for: alarm systems (detect threshold), not analysis

### ENS160 / CCS811 — Digital TVOC and CO₂ Sensor

More precise digital sensors for air quality:

**CCS811:**
- Measures: eCO₂ (estimated CO₂, 400-8192 ppm), TVOC (Total Volatile Organic Compounds, 0-1187 ppb)
- Interface: I2C (0x5A or 0x5B)
- Supply: 3.3V
- Burn-in: 48 hours for accurate readings!
- Baseline: 20 minutes for stable readings
- Note: Measures TVOC and *estimates* CO₂ based on VOC level — not true CO₂!

**True CO₂ sensors (NDIR):**
- SCD30, SCD40/41 (Sensirion): Real CO₂ via NDIR (Non-Dispersive Infrared)
- CO₂ absorbs 4.26 µm infrared specifically
- Accurate, stable, but expensive ($30-50)
- Used in: HVAC, ventilation control, indoor air quality monitoring

---

## 26.7 pH Sensor

### Working Principle

pH measures hydrogen ion (H⁺) concentration in solution:

```
pH = -log₁₀[H⁺]

pH 0-14 scale:
  pH 0-6:  Acid  (H⁺ rich)
  pH 7:    Neutral (pure water)
  pH 8-14: Base/Alkaline (OH⁻ rich)

pH electrode (glass electrode):
  Reference half-cell + sensing half-cell
  Nernst equation:
  E = E₀ + (RT/F) × ln([H⁺])
    = E₀ - 0.0592 × pH (at 25°C)

Sensitivity: 59.16 mV/pH unit at 25°C
At pH 7: 0 mV (with proper calibration)
At pH 4: +177 mV
At pH 10: -177 mV
```

**Temperature effect:**
```
Slope changes with temperature!
At 20°C: 58.16 mV/pH
At 25°C: 59.16 mV/pH
At 37°C: 61.54 mV/pH

Must temperature compensate for accuracy!
pH meter with built-in temp sensor (PT100 or NTC)
```

**Arduino pH module:**
```c
// pH electrode → high-impedance op-amp buffer → ADC
// Common module: SEN0161 (DFRobot) or similar

float voltage = analogRead(A0) * (5.0 / 1023.0);
float pH = 7.0 + ((2.5 - voltage) / 0.18);  // Calibration formula
// Note: Exact formula depends on electrode calibration!
```

**Calibration (2-point):**
1. Rinse electrode, dip in pH 7.0 buffer solution
2. Adjust trim pot until reading shows 7.00
3. Rinse, dip in pH 4.0 (or 10.0) buffer
4. Record slope
5. Apply two-point correction in firmware

**pH electrode maintenance:**
- Must be kept moist (store in KCl solution or pH 7 buffer)
- Drying out permanently damages the glass membrane!
- Reference junction must not be clogged
- Lifespan: 1-2 years with proper care

---

## 26.8 Soil Moisture Sensor

Two main types for garden/agricultural IoT:

### Resistive Soil Moisture Sensor

```
Two electrodes in soil
Resistance decreases with moisture:
  Dry soil: 100 kΩ - ∞
  Wet soil: 1 kΩ - 10 kΩ

Problems:
  - Electrolytic corrosion (electrodes degrade fast!)
  - Susceptible to soil mineral content (salinity affects reading)
  - DC voltage causes electrolysis

Cheap FC-28 module: $0.50-1.00, expect to replace every few months
```

### Capacitive Soil Moisture Sensor

```
Two interleaved electrodes (capacitor plates)
Soil is dielectric between plates:
  Dry soil: lower dielectric constant ε_r ≈ 3-5
  Wet soil: water ε_r = 80! → much higher capacitance

Capacitance measured by oscillator frequency:
  Higher moisture → higher C → lower frequency

No DC current → no corrosion → much more durable!

VEML7700/STEMMA QT Capacitive Soil Sensor: I2C, $8-15
Chirp Sensor (I2C): temperature, moisture, light in one
```

**Capacitive sensor calibration:**
```c
// Map raw ADC reading to percentage:
const int dry_value = 520;   // Read in air (0% moisture)
const int wet_value = 250;   // Read in water (100% moisture)

int raw = analogRead(A0);
int moisture_pct = map(raw, dry_value, wet_value, 0, 100);
moisture_pct = constrain(moisture_pct, 0, 100);
```

---

## 26.9 Current and Voltage Sensors

### ACS712 — Current Sensor (Hall Effect)

```
Current flows through copper conductor on chip
Copper current creates magnetic field
Hall effect sensor in chip measures field
Output: Voltage proportional to current

ACS712-05B: ±5A range, 185 mV/A sensitivity
ACS712-20A: ±20A range, 100 mV/A
ACS712-30A: ±30A range, 66 mV/A

Output at 0A: VCC/2 = 2.5V (bipolar, measures AC and DC!)
At +5A: 2.5 + 5 × 0.185 = 3.425V
At -5A: 2.5 - 5 × 0.185 = 1.575V

Accuracy: ±1.5% of full-scale (too noisy for precision)
```

**INA219 — Bidirectional Current/Power Monitor (I2C):**
```
Shunt resistor in series with load
INA219 measures:
  - Shunt voltage (±320 mV)
  - Bus voltage (0-32V)
  - Current = Vshunt / R_shunt
  - Power = Bus_V × Current

Interface: I2C (0x40-0x4F, configurable)
Resolution: 0.8 mA (default), up to 0.1 mA
Accuracy: 1% typical

Used in: Battery monitoring, solar charge controllers, power management
```

```c
#include <Adafruit_INA219.h>
Adafruit_INA219 ina219;

void setup() {
    ina219.begin();
    ina219.setCalibration_32V_2A();  // 32V, 2A range
}

void loop() {
    float shuntvoltage = ina219.getShuntVoltage_mV();
    float busvoltage = ina219.getBusVoltage_V();
    float current_mA = ina219.getCurrent_mA();
    float power_mW = ina219.getPower_mW();

    Serial.print("Bus: "); Serial.print(busvoltage); Serial.println(" V");
    Serial.print("Current: "); Serial.print(current_mA); Serial.println(" mA");
}
```

---

## 26.10 Flex and Force Sensors

### Flex Sensor

A resistive sensor that changes resistance when bent:

```
Construction:
  Conductive carbon on polyimide (Kapton) substrate
  Flat: ~10 kΩ (no bend)
  Bent 45°: ~15 kΩ
  Bent 90°: ~30 kΩ

Spectra Symbol 2.2" flex sensor: 10-110 kΩ range
Spectra Symbol 4.5": 10-110 kΩ (more sensitive)

Circuit: Voltage divider with 10-47 kΩ fixed resistor
```

**Applications:**
- VR gloves (finger angle detection)
- Robot gripper feedback
- Wearable gesture detection
- Automotive seat occupancy detection

### Force Sensitive Resistor (FSR)

Resistance decreases as force increases:

```
FSR materials: Polymer thick film between conductive traces
  No force: very high resistance (>1 MΩ)
  Light touch (20g): ~100 kΩ
  Heavy press (1kg): ~1 kΩ
  Max force (10kg): ~100 Ω

Force range: 20g - 10 kg
Repeatability: ±15% (not for precision measurements!)
Response time: <5 ms

Interlink FSR402 (round, 12.5mm active area): most common
```

**Circuit:**
```
VCC ──[10kΩ]──── Vout ──── ADC
                   │
                  FSR
                   │
                  GND

More force → lower FSR → higher Vout
```

**Calibration (FSR is highly non-linear):**
```c
// Empirical calibration:
float voltage = analogRead(A0) * (3.3 / 1023.0);
float resistance = 10000.0 * (3.3 - voltage) / voltage;  // R = Rfixed × (Vcc/Vout - 1)
float force_N = 1.0 / (resistance / 1000.0);  // Approximate, needs proper calibration
```

---

## 26.11 Ultrasonic Fluid Level Sensor

Same principle as HC-SR04 but for liquid level:

**MaxSonar WR Series (waterproof):**
- Acoustic beam points downward into tank
- Measures time from sensor to liquid surface
- Outputs: UART, PWM, or analog voltage
- Range: 20-765 cm
- Used in: Water tanks, fuel tanks, waste water

**Pressure-based level (continuous):**
```
Pressure sensor at bottom of tank:
P = ρ × g × h

Where:
  ρ = fluid density (1000 kg/m³ for water)
  g = 9.81 m/s²
  h = height of fluid

h = P / (ρ × g)
1 meter of water = 9810 Pa = 0.0981 bar = 1.422 PSI

Use: Waterproof pressure transducer (4-20 mA output)
MSP300 series: Stainless steel, 4-20 mA, 0-1 to 0-100 bar
```

---

## 26.12 Flow Sensors

### Hall-Effect Water Flow Sensor (YF-S201)

```
Paddlewheel rotor inside pipe housing
As water flows, rotor spins
Magnet on rotor triggers Hall effect sensor
Pulses per liter: 450 pulses/L (YF-S201)

Flow rate calculation:
frequency = pulses per second
flow_L_per_min = frequency / 7.5  (YF-S201 calibration)

void flow_interrupt_handler(void) {
    pulse_count++;  // Count pulses in ISR
}

// Every second:
float freq = pulse_count;  // Pulses per second
float flow = freq / 7.5;   // L/min
volume += flow / 60.0;     // Total liters
pulse_count = 0;
```

**Range:** 1-30 L/min, ±10% accuracy
**Applications:** Water meters, aquaponics, coolant monitoring

---

## 26.13 Geiger-Müller Counter (Radiation)

For nuclear radiation detection:

**Working principle:**
```
Geiger-Müller tube: gas-filled tube with high voltage (~400-600V)
Ionizing particle enters → ionizes gas → avalanche discharge → pulse!

Pulse counted per unit time → CPM (Counts Per Minute)
Converted to radiation units:
  For background radiation: 25 CPM ≈ 0.25 µSv/hr (typical background)
  Dose rate: µSv/hr or mR/hr

Tube types: SBM-20 (γ, β), SBM-19 (γ only), SI-8B (β only)
```

**DIY Geiger counter:**

```mermaid
flowchart LR
    HV["High Voltage Generator\n(400-600V from 5V boost converter)"]
    GM["GM Tube"]
    ST["Schmitt Trigger"]
    MCU["MCU GPIO interrupt\nCount pulses → CPM"]
    DISP["Display µSv/hr"]
    HV -->|"biases"| GM
    GM -->|"pulse output"| ST
    ST --> MCU --> DISP
```

---

## 26.14 Sound / Microphone Sensors

### Electret Microphone Module (KY-038, KY-037)

```
Electret microphone capsule:
  Pre-charged diaphragm (permanent electret)
  Capacitor microphone + built-in JFET buffer
  Frequency response: 20 Hz - 20 kHz

Module contains:
  Electret mic + LM393 comparator (for threshold trigger)
  A0: Analog output (amplified mic signal, 0-5V swings)
  D0: Digital output (HIGH when sound exceeds threshold, set by pot)

Sensitivity: -40 to -44 dBV (mV output per Pa sound pressure)
```

**Sound level measurement:**
```c
// Calculate RMS of sound:
const int SAMPLE_TIME = 50;  // ms
long startMillis = millis();
int peakToPeak = 0;
int signalMin = 1024;
int signalMax = 0;

while (millis() - startMillis < SAMPLE_TIME) {
    int sample = analogRead(A0);
    if (sample > signalMax) signalMax = sample;
    if (sample < signalMin) signalMin = sample;
}
peakToPeak = signalMax - signalMin;
float db = map(peakToPeak, 0, 1023, 49, 100);  // Approximate dB
```

### I2S MEMS Microphones (Digital)

Better quality for audio processing:

**SPH0645LM4H (Adafruit I2S Mic):**
```
Output: I2S digital audio
SNR: 65 dB
Frequency response: 100 Hz - 10 kHz
Interface: I2S (BCK, WS, SD)
Supply: 3.3V
Used in: Smart speakers (basic), voice recognition, audio logging

With ESP32 I2S:
i2s_config_t i2s_config = {
    .mode = I2S_MODE_MASTER | I2S_MODE_RX | I2S_MODE_PDM,
    .sample_rate = 16000,
    .bits_per_sample = I2S_BITS_PER_SAMPLE_16BIT,
    .channel_format = I2S_CHANNEL_FMT_ONLY_LEFT,
};
```

---

## 26.15 Ultrasonic Fingerprint Sensor

Used in modern smartphones (under-display fingerprint):

```
Qualcomm 3D Sonic Sensor (in Galaxy S10+, S20, etc.):
  Ultrasonic pulses penetrate AMOLED display
  Reflect differently from fingerprint ridges vs valleys
  3D fingerprint map (unlike optical which is 2D)
  Works with wet fingers, some conditions optical fails

Technology:
  PMUT (Piezoelectric Micromachined Ultrasonic Transducer)
  128×128 PMUT array
  Each element transmits and receives 40 MHz ultrasound
  Image reconstruction algorithm → fingerprint
  Spoof-resistant (3D structure, not 2D image)
```

---

## 26.16 Summary — Optical, Chemical, and Other Sensors

| Sensor Type        | Principle           | Range/Range      | Interface   | Key Use               |
|-------------------|--------------------|-----------------|-----------|-----------------------|
| LDR               | Photoconductivity   | 1 kΩ - 1 MΩ    | Analog R  | Light level (simple)  |
| Photodiode        | PN junction         | fA - mA         | Analog I  | Optical receiver      |
| Phototransistor   | Photodiode + BJT    | µA - mA         | Analog    | Object detection      |
| TCS34725          | RGBC filters + ADC  | 0.1 - 10k lux  | I2C       | Color measurement     |
| Rotary encoder    | IR + slots          | unlimited count | Digital   | Position/speed        |
| MQ-2 (gas)        | Metal oxide semi.   | 300-10,000 ppm  | Analog    | Gas/smoke alarm       |
| CCS811            | Metal oxide semi.   | CO₂/TVOC        | I2C       | Air quality           |
| SCD40             | NDIR optical        | 400-2000 ppm CO₂| I2C       | True CO₂              |
| pH electrode      | Electrochemical     | 0-14 pH         | Analog mV | Water quality         |
| Capacitive soil   | Capacitance         | 0-100% moisture | Analog    | Garden automation     |
| ACS712            | Hall effect         | ±5/20/30 A      | Analog    | Current monitoring    |
| INA219            | Shunt voltage       | ±3.2A, 0-32V   | I2C       | Power monitoring      |
| Flex sensor       | Piezoresistive      | 0-90° bend      | Analog R  | Gesture control       |
| FSR               | Pressure resistance | 20g - 10 kg     | Analog R  | Touch/weight          |
| YF-S201           | Hall pulse          | 1-30 L/min      | Digital   | Water flow            |
| Electret mic      | Capacitive+JFET     | 20 Hz - 20 kHz  | Analog    | Sound detection       |
| I2S MEMS mic      | MEMS capacitive     | 100 Hz - 10 kHz | I2S       | Voice, audio          |
