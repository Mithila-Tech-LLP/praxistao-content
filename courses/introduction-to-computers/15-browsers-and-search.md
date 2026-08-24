# Chapter 15: Browsers and Search Engines

> **"A browser is your window to the web. A search engine is the librarian who helps you find what you need. Both are so deeply embedded in daily life that most people have never thought about how they actually work."**

---

## Table of Contents

1. [What Is a Browser?](#1-what-is-a-browser)
2. [Inside a Browser — How It Works](#2-inside-a-browser--how-it-works)
3. [Browser Features You Should Know](#3-browser-features-you-should-know)
4. [What Is a Search Engine?](#4-what-is-a-search-engine)
5. [How Google Finds Answers](#5-how-google-finds-answers)
6. [Search Tips — Finding Things Better](#6-search-tips--finding-things-better)
7. [Cookies, Cache, and Browsing History](#7-cookies-cache-and-browsing-history)
8. [Browser Privacy and Incognito Mode](#8-browser-privacy-and-incognito-mode)
9. [Summary](#summary)

---

## 1. What Is a Browser?

A **web browser** is an app that retrieves web pages and displays them to you.

```
Popular browsers (2026):
  Chrome (Google)         ~65% of all web traffic
  Safari (Apple)          ~19% (default on iPhone, iPad, Mac)
  Edge (Microsoft)        ~5%  (default on Windows)
  Firefox (Mozilla)       ~3%  (open source, privacy-focused)
  Samsung Internet        ~3%  (default on Samsung phones)
  Brave                   ~2%  (strong privacy, ad-blocking)
```

All browsers can show you any website. The differences are:
- Speed (how fast they load pages)
- Privacy (how much data they collect)
- Battery use (Safari is very efficient on Apple devices)
- Extensions/add-ons (Chrome has the most)
- Default integration (Safari knows your Apple ID, Chrome knows your Google account)

---

## 2. Inside a Browser — How It Works

A browser has several important "engines":

```
Rendering Engine:
  Reads HTML + CSS and draws the page on screen.
  Chrome/Edge/Brave: Blink engine
  Safari: WebKit engine
  Firefox: Gecko engine
  
  This is why the same website looks slightly different on Safari vs Chrome —
  different engines interpret some things slightly differently.

JavaScript Engine:
  Runs JavaScript code in web pages.
  Chrome: V8 (also used in Node.js)
  Safari: JavaScriptCore
  Firefox: SpiderMonkey
  
  A faster JS engine = faster web apps.

Network stack:
  Handles HTTP/HTTPS requests, DNS lookup, caching.
  
Security sandbox:
  Each browser tab runs in its own isolated process.
  If one tab crashes or is attacked, others are protected.
```

---

## 3. Browser Features You Should Know

```
Tabs:
  Open multiple websites at once.
  Cmd+T (Mac) / Ctrl+T (Windows) = new tab
  Cmd+W / Ctrl+W = close tab
  Cmd+Shift+T / Ctrl+Shift+T = reopen last closed tab

Bookmarks / Favorites:
  Save websites to visit again.
  Cmd+D / Ctrl+D = bookmark current page
  
Browser history:
  Every site you've visited is logged.
  Cmd+Y / Ctrl+H = open history
  
Extensions / Add-ons:
  Small programs that modify browser behavior.
  Popular ones:
    uBlock Origin  → block ads
    1Password      → fill passwords automatically
    Grammarly      → check grammar in text fields
    React DevTools → for web developers
    
Address bar (also called "omnibox" in Chrome):
  Type a URL → go to that website
  Type words → search Google (or your default search engine)
  
Developer Tools (F12):
  See the HTML, CSS, JavaScript of any website
  Inspect any element, see network requests, debug JavaScript
  Essential for web developers
  
Keyboard shortcuts:
  Ctrl/Cmd+L  → jump to address bar
  F5/Cmd+R    → refresh page
  Ctrl/Cmd+F  → find text on page
  Ctrl/Cmd++  → zoom in
  Ctrl/Cmd+0  → reset zoom
  Backspace/Cmd+[ → go back
```

---

## 4. What Is a Search Engine?

A **search engine** is a service that indexes the web and lets you search it.

```
Major search engines:
  Google      → 91% of global searches (the dominant search engine)
  Bing        → 3% (Microsoft's search engine)
  Yahoo       → 1%
  DuckDuckGo  → privacy-focused, doesn't track you
  Baidu       → dominant in China
  Yandex      → dominant in Russia
  
Search engines are not the same as browsers.
  Browser: Chrome, Firefox, Safari — apps that display web pages
  Search engine: Google, Bing — services that help you find web pages
  
  Chrome is a browser. Google is a search engine.
  You can use Chrome with Bing as your search engine.
  You can use Firefox with Google as your search engine.
```

---

## 5. How Google Finds Answers

Google knows about billions of web pages. How?

```
Step 1: Crawling
  Google runs "Googlebot" — software that visits web pages automatically.
  Starting from known URLs, it follows every link to discover new pages.
  Does this continuously — billions of pages per day.
  
Step 2: Indexing
  Every crawled page is analyzed:
  What words are on it? What topics? What links does it have?
  All this is stored in Google's index — essentially a giant database.
  Index size: hundreds of billions of web pages.
  
Step 3: Ranking (PageRank + ML)
  When you search "how do penguins stay warm":
  Google looks for pages matching these words in its index.
  Then ranks them by relevance + quality:
    - How many other good websites link to this page? (PageRank)
    - Does it match the search intent?
    - Is it a trusted, authoritative source?
    - Is it recent (for news topics)?
    - Does it load fast?
    - Is it mobile-friendly?
  Machine learning models have made this much smarter over time.
  
Step 4: Results
  Top 10 results displayed in milliseconds.
  Featured snippets, knowledge panels, maps, images, news — all from different systems.
```

---

## 6. Search Tips — Finding Things Better

Most people use Google at maybe 20% efficiency. Here's how to search smarter:

```
Exact phrase — use quotes:
  "climate change 2024 report"
  Only shows pages containing this exact phrase.

Exclude words — use minus:
  jaguar speed -car
  Finds jaguar (animal) speed, not Jaguar car speed.

Search a specific site:
  site:wikipedia.org black holes
  Only shows Wikipedia pages about black holes.

File type:
  filetype:pdf climate report
  Only shows PDF files.

Within a time range:
  Use "Tools" → "Any time" → "Past year" for recent results.

Definitions:
  define:ubiquitous
  Shows definition directly.

Calculator:
  Type math directly: "2^32", "15% of 450", "speed of light in km/h"

Unit conversion:
  "100 USD to EUR", "5 miles to km", "98.6 F to C"

Weather:
  "weather London", "weather this weekend"

Quick answers:
  "how old is Obama", "when did WWII end", "capital of Thailand"
```

---

## 7. Cookies, Cache, and Browsing History

**Cookies:**
```
A small text file that websites save in your browser.

Uses:
  Login:    Website saves "user is logged in" → so you stay logged in
  Cart:     Shopping site saves your cart items
  Settings: Website remembers you prefer dark mode
  Tracking: Ad networks track which websites you visit (→ targeted ads)
  
Cookies are NOT viruses. They're just text files.
But tracking cookies are a real privacy concern.
```

**Cache:**
```
Your browser saves copies of websites you've visited.
Next time you visit, it loads from local storage (fast) instead of downloading again.

If a website looks "broken" or "old":
  Clear the cache: Ctrl+Shift+Delete (Windows) / Cmd+Shift+Delete (Mac)
  This forces a fresh download of the page.
```

**Browsing History:**
```
Your browser keeps a log of every page you've visited.
Stored locally on your device.
Useful for: finding a page you visited before but forgot the URL.
Concern: if someone has access to your computer, they can see your history.
```

---

## 8. Browser Privacy and Incognito Mode

**Incognito / Private browsing mode:**
```
What it DOES:
  ✓ Doesn't save browsing history to your device
  ✓ Doesn't save cookies after the window closes
  ✓ Doesn't save typed passwords
  ✓ Good for: logging into a different account, gift shopping without
    spoiling search history, using a shared/public computer

What it DOES NOT do:
  ✗ Does NOT hide your IP from websites you visit
  ✗ Does NOT prevent your ISP from seeing what sites you visit
  ✗ Does NOT make you anonymous online
  ✗ Does NOT prevent tracking by the sites themselves
  
For real privacy: use a VPN, or Tor Browser.
```

**Third-party cookies:**
Many websites embed trackers from advertising companies.
When you visit site A, site B (the ad network) also gets notified.
Over time, the ad network builds a profile of your interests.
This is how you search for "running shoes" and see running shoe ads for weeks.

Modern browsers (Safari, Firefox, Brave) block many of these trackers by default. Chrome is rolling out similar protections.

---

## Summary

| Concept | What It Means |
|---------|--------------|
| Browser | App that retrieves and displays web pages |
| Rendering engine | Converts HTML/CSS into visual display |
| Search engine | Service that indexes the web and finds pages |
| Cookie | Small text file websites save in your browser |
| Cache | Saved copies of websites for faster reloading |
| Incognito mode | Doesn't save history/cookies locally (not fully anonymous) |
| Extensions | Add-ons that enhance browser functionality |

**Now you can navigate the web and find anything. Next: how do you communicate with others online?**
