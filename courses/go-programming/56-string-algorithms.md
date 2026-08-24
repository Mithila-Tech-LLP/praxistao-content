# Chapter 38: String Algorithms

Strings are everywhere in software — parsing, searching, pattern matching, compression. The naive approach to pattern searching is O(nm), but KMP and Rabin-Karp reduce it to O(n+m). Understanding these algorithms builds intuition for the Go `strings` and `bytes` packages and for interview problems involving substrings, subsequences, and anagrams.

## Table of Contents

1. [String Fundamentals in Go](#1-string-fundamentals-in-go)
2. [Pattern Matching — KMP](#2-pattern-matching--kmp)
3. [Rabin-Karp (Rolling Hash)](#3-rabin-karp-rolling-hash)
4. [Sliding Window Problems](#4-sliding-window-problems)
5. [Two Pointer Techniques](#5-two-pointer-techniques)
6. [Advanced String Problems](#6-advanced-string-problems)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. String Fundamentals in Go

```go
// Strings are immutable byte slices in Go.
s := "Hello, 世界"

// Length in bytes vs runes:
fmt.Println(len(s))           // 13 (bytes)
fmt.Println(len([]rune(s)))   // 9 (Unicode code points)

// Iterate bytes (wrong for multi-byte chars):
for i := 0; i < len(s); i++ {
    fmt.Printf("%d:%c ", i, s[i])
}

// Iterate runes (correct for Unicode):
for i, r := range s {
    fmt.Printf("%d:%c ", i, r)  // i is BYTE index, r is the rune
}

// String <-> []byte and []rune conversions:
bytes := []byte(s)        // O(n) copy
runes := []rune(s)        // O(n) copy
s2 := string(bytes)       // O(n) copy back
s3 := string(runes)       // O(n) copy back

// Efficient string building (avoid O(n²) concatenation):
var sb strings.Builder
for i := 0; i < 1000; i++ {
    sb.WriteString("hello")
    sb.WriteByte(' ')
}
result := sb.String()

// Common stdlib:
import "strings"
strings.Contains(s, "ell")       // true
strings.HasPrefix(s, "Hel")      // true
strings.Index(s, "llo")          // 2
strings.Split(s, ", ")           // ["Hello" "世界"]
strings.Join([]string{"a","b"}, "-")  // "a-b"
strings.TrimSpace("  hi  ")      // "hi"
strings.ToLower("ABC")           // "abc"
strings.Replace("aaa", "a", "b", 2)  // "bba"
strings.Count("aaa", "a")        // 3
strings.Repeat("ab", 3)          // "ababab"
```

### String Comparison and Sorting
```go
// Anagram check — frequency count:
func IsAnagram(s, t string) bool {
    if len(s) != len(t) { return false }
    count := [26]int{}
    for i := range s {
        count[s[i]-'a']++
        count[t[i]-'a']--
    }
    return count == [26]int{}
}

// Palindrome check:
func IsPalindrome(s string) bool {
    runes := []rune(s)
    for lo, hi := 0, len(runes)-1; lo < hi; lo, hi = lo+1, hi-1 {
        if runes[lo] != runes[hi] { return false }
    }
    return true
}
```

---

## 2. Pattern Matching — KMP

Naive search: O(nm). For each position in text, compare pattern character by character.

**KMP** avoids redundant comparisons by precomputing a failure function that tells us how far back to jump when a mismatch occurs.

```go
// KMP (Knuth-Morris-Pratt) — O(n+m) time, O(m) space
//
// Failure function (partial match table): lps[i] = length of longest proper
// prefix of pattern[0..i] that is also a suffix.
// Example: pattern = "AAACAAAA"
//          lps    = [0,1,2,0,1,2,3,3]

func buildLPS(pattern string) []int {
    m := len(pattern)
    lps := make([]int, m)
    length, i := 0, 1  // lps[0] is always 0

    for i < m {
        if pattern[i] == pattern[length] {
            length++
            lps[i] = length
            i++
        } else {
            if length != 0 {
                length = lps[length-1]  // Don't increment i
            } else {
                lps[i] = 0
                i++
            }
        }
    }
    return lps
}

func KMPSearch(text, pattern string) []int {
    n, m := len(text), len(pattern)
    if m == 0 { return nil }

    lps := buildLPS(pattern)
    matches := []int{}
    i, j := 0, 0  // Pointers into text and pattern

    for i < n {
        if text[i] == pattern[j] {
            i++; j++
        }
        if j == m {
            matches = append(matches, i-j)  // Found match at i-j
            j = lps[j-1]                    // Continue searching
        } else if i < n && text[i] != pattern[j] {
            if j != 0 {
                j = lps[j-1]  // Skip already-matched prefix
            } else {
                i++
            }
        }
    }
    return matches
}

// Example:
// text = "AABAACAADAABAABA", pattern = "AABA"
// KMPSearch returns [0, 9, 12]
```

---

## 3. Rabin-Karp (Rolling Hash)

Rabin-Karp uses a hash of the pattern and slides a hash window over the text. When hashes match, verify with string comparison (to handle collisions).

```go
// Rabin-Karp — O(n+m) average, O(nm) worst case (hash collisions)
// Excellent for: multiple pattern search, 2D pattern matching

const (
    base = 31
    mod  = 1_000_000_007
)

func RabinKarpSearch(text, pattern string) []int {
    n, m := len(text), len(pattern)
    if m > n { return nil }

    // Compute base^(m-1) mod p:
    pow := 1
    for i := 0; i < m-1; i++ {
        pow = pow * base % mod
    }

    // Compute hash of pattern and first window:
    patHash, winHash := 0, 0
    for i := 0; i < m; i++ {
        patHash = (patHash*base + int(pattern[i]-'a'+1)) % mod
        winHash = (winHash*base + int(text[i]-'a'+1)) % mod
    }

    matches := []int{}
    for i := 0; i <= n-m; i++ {
        if winHash == patHash && text[i:i+m] == pattern {
            matches = append(matches, i)
        }
        if i < n-m {
            // Roll the hash: remove leading char, add next char
            winHash = (winHash - int(text[i]-'a'+1)*pow%mod + mod) % mod
            winHash = (winHash*base + int(text[i+m]-'a'+1)) % mod
        }
    }
    return matches
}
```

### Z-Algorithm — fast prefix matching
```go
// Z array: Z[i] = length of longest substring starting at i that matches a prefix of s.
func buildZ(s string) []int {
    n := len(s)
    z := make([]int, n)
    z[0] = n
    l, r := 0, 0

    for i := 1; i < n; i++ {
        if i < r {
            z[i] = min(r-i, z[i-l])
        }
        for i+z[i] < n && s[z[i]] == s[i+z[i]] {
            z[i]++
        }
        if i+z[i] > r {
            l, r = i, i+z[i]
        }
    }
    return z
}

func ZSearch(text, pattern string) []int {
    s := pattern + "$" + text  // $ is separator not in alphabet
    z := buildZ(s)
    m := len(pattern)
    matches := []int{}

    for i := m + 1; i < len(s); i++ {
        if z[i] == m {
            matches = append(matches, i-m-1)
        }
    }
    return matches
}
```

---

## 4. Sliding Window Problems

The sliding window technique maintains a window `[left, right]` and expands/contracts it based on a condition.

### Longest Substring Without Repeating Characters
```go
func LengthOfLongestSubstring(s string) int {
    lastSeen := make(map[byte]int)
    maxLen, left := 0, 0

    for right := 0; right < len(s); right++ {
        c := s[right]
        if idx, ok := lastSeen[c]; ok && idx >= left {
            left = idx + 1  // Shrink window past the repeated character
        }
        lastSeen[c] = right
        if w := right - left + 1; w > maxLen { maxLen = w }
    }
    return maxLen
}
```

### Minimum Window Substring
```go
// MinWindow: smallest window in s containing all characters of t.
func MinWindow(s, t string) string {
    if len(s) < len(t) { return "" }

    need := make(map[byte]int)
    for i := range t { need[t[i]]++ }

    have, required := 0, len(need)
    window := make(map[byte]int)
    left := 0
    minLen, minLeft := len(s)+1, 0

    for right := 0; right < len(s); right++ {
        c := s[right]
        window[c]++
        if need[c] > 0 && window[c] == need[c] {
            have++
        }
        for have == required {
            if right-left+1 < minLen {
                minLen = right - left + 1
                minLeft = left
            }
            lc := s[left]
            window[lc]--
            if need[lc] > 0 && window[lc] < need[lc] {
                have--
            }
            left++
        }
    }
    if minLen == len(s)+1 { return "" }
    return s[minLeft : minLeft+minLen]
}
```

### Find All Anagrams in a String
```go
// FindAnagrams: all start indices where s's anagram exists in text.
func FindAnagrams(s, p string) []int {
    if len(s) < len(p) { return nil }

    pCount := [26]int{}
    wCount := [26]int{}
    for i := range p { pCount[p[i]-'a']++ }
    for i := 0; i < len(p); i++ { wCount[s[i]-'a']++ }

    result := []int{}
    if pCount == wCount { result = append(result, 0) }

    for right := len(p); right < len(s); right++ {
        wCount[s[right]-'a']++
        wCount[s[right-len(p)]-'a']--
        if pCount == wCount { result = append(result, right-len(p)+1) }
    }
    return result
}
```

---

## 5. Two Pointer Techniques

### Valid Palindrome (ignoring non-alphanumeric)
```go
func IsValidPalindrome(s string) bool {
    lo, hi := 0, len(s)-1
    for lo < hi {
        for lo < hi && !isAlphaNum(s[lo]) { lo++ }
        for lo < hi && !isAlphaNum(s[hi]) { hi-- }
        if lo < hi {
            if toLower(s[lo]) != toLower(s[hi]) { return false }
            lo++; hi--
        }
    }
    return true
}

func isAlphaNum(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func toLower(c byte) byte {
    if c >= 'A' && c <= 'Z' { return c + 32 }
    return c
}
```

### Reverse Words in a String
```go
func ReverseWords(s string) string {
    words := strings.Fields(s)  // Split on any whitespace, trim
    lo, hi := 0, len(words)-1
    for lo < hi {
        words[lo], words[hi] = words[hi], words[lo]
        lo++; hi--
    }
    return strings.Join(words, " ")
}
```

---

## 6. Advanced String Problems

### Longest Common Prefix
```go
func LongestCommonPrefix(strs []string) string {
    if len(strs) == 0 { return "" }
    prefix := strs[0]
    for _, s := range strs[1:] {
        for !strings.HasPrefix(s, prefix) {
            prefix = prefix[:len(prefix)-1]
            if prefix == "" { return "" }
        }
    }
    return prefix
}
```

### String Compression
```go
// Compress consecutive duplicate characters: "aabcccccaaa" → "a2b1c5a3"
func Compress(s string) string {
    if len(s) == 0 { return s }
    var sb strings.Builder
    count := 1

    for i := 1; i < len(s); i++ {
        if s[i] == s[i-1] {
            count++
        } else {
            sb.WriteByte(s[i-1])
            if count > 1 { sb.WriteString(strconv.Itoa(count)) }
            count = 1
        }
    }
    sb.WriteByte(s[len(s)-1])
    if count > 1 { sb.WriteString(strconv.Itoa(count)) }

    result := sb.String()
    if len(result) >= len(s) { return s }
    return result
}
```

### Decode Ways
```go
// DecodeWays: "226" can decode as "BZ"(2,26), "VF"(22,6), "BBF"(2,2,6) → 3 ways
// DP + string processing:
func NumDecodings(s string) int {
    n := len(s)
    if n == 0 || s[0] == '0' { return 0 }

    dp := make([]int, n+1)
    dp[0] = 1  // Empty string — one way
    dp[1] = 1  // Single non-zero digit — one way

    for i := 2; i <= n; i++ {
        oneDigit := int(s[i-1] - '0')
        twoDigit := int(s[i-2]-'0')*10 + int(s[i-1]-'0')

        if oneDigit >= 1 { dp[i] += dp[i-1] }
        if twoDigit >= 10 && twoDigit <= 26 { dp[i] += dp[i-2] }
    }
    return dp[n]
}
```

### Regular Expression Matching (DP)
```go
// IsMatch: '.' matches any single char, '*' matches zero or more of preceding.
// dp[i][j] = p[:j] matches s[:i]
func IsMatch(s, p string) bool {
    m, n := len(s), len(p)
    dp := make([][]bool, m+1)
    for i := range dp { dp[i] = make([]bool, n+1) }
    dp[0][0] = true

    // Pattern like "a*b*c*" can match empty string:
    for j := 1; j <= n; j++ {
        if p[j-1] == '*' && j >= 2 { dp[0][j] = dp[0][j-2] }
    }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if p[j-1] == '*' {
                dp[i][j] = dp[i][j-2]  // Use '*' as zero occurrences
                if p[j-2] == '.' || p[j-2] == s[i-1] {
                    dp[i][j] = dp[i][j] || dp[i-1][j]  // Use one more occurrence
                }
            } else if p[j-1] == '.' || p[j-1] == s[i-1] {
                dp[i][j] = dp[i-1][j-1]
            }
        }
    }
    return dp[m][n]
}
```

---

## Summary

- **Go strings**: immutable bytes, iterate with `range` for Unicode, use `strings.Builder` for concatenation
- **KMP**: O(n+m) pattern search via failure function (prefix = suffix table); avoid re-scanning matched prefix
- **Rabin-Karp**: rolling hash O(n+m) average; good for multiple pattern search
- **Z-algorithm**: Z[i] = length of longest prefix match starting at i; concatenate pattern + "$" + text
- **Sliding window**: two pointers `[left, right]`, expand right, contract left based on constraint
- **Classic problems**: minimum window substring, all anagrams, longest non-repeating substring
- **DP on strings**: LCS, edit distance, regex matching, decode ways — all use 2D table or 1D rolling

---

## Exercises

### Easy
1. Implement `ReverseString(s []byte)` in-place. Then `ReverseStringByWords(s string) string` (reverse each word individually, not the order of words). Then `ReverseVowels(s string) string` — two pointers, swap only when both point to vowels.
2. Implement `IsomorphicStrings(s, t string) bool` — each character in s maps to exactly one in t and vice versa. Use two maps (s→t and t→s). Verify: `"egg","add"` → true; `"foo","bar"` → false.
3. Implement `WordPattern(pattern, s string) bool` — `pattern = "abba"`, `s = "dog cat cat dog"` → true. Same bijection as isomorphic strings but between pattern chars and words.

### Medium
4. **Group Anagrams**: given a list of strings, group them by anagram. Two strings are anagrams if sorted they're equal. Use `sort.Slice` to sort characters. Return `[][]string`. Bonus: can you use a frequency count array as map key instead of sorting?
5. **Longest Repeating Character Replacement**: given string s and integer k, find the length of the longest substring where you can replace at most k characters to make all characters the same. Sliding window: maintain max frequency count in window; window is valid if `windowLen - maxFreq <= k`. Verify: `s="AABABBA", k=1` → 4.
6. **Implement `strStr`**: implement `strings.Index` yourself using KMP. Write a test that compares your implementation against `strings.Index` for 1000 random (text, pattern) pairs. The test should be deterministic using a seeded `rand.New(rand.NewSource(42))`.

### Hard
7. **Shortest Palindrome**: given string s, add minimum characters to the front to make it a palindrome. Key insight: find the longest palindromic prefix using KMP. Build `s + "#" + reverse(s)`, compute the KMP failure function, and the last value tells you the length of the longest palindromic prefix. Verify: `"aacecaaa"` → `"aaacecaaa"`.
8. **Substring with Concatenation of All Words**: given string s and a list of words (all same length), find all starting indices where a concatenation of all words (in any order) appears. Sliding window over word-sized chunks, maintain a word-frequency map. Time: O(n * totalLen). Verify: `s="barfoothefoobarman", words=["foo","bar"]` → `[0,9]`.
