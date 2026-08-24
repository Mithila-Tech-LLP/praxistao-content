# Chapter 14: What Is the World Wide Web?

> **"The web is a global library where every book links to every other book, and every person on Earth can both read from it and add to it, for free. Nothing like it has existed in human history."**

---

## Table of Contents

1. [The Web in Plain English](#1-the-web-in-plain-english)
2. [Web Pages, Websites, and Web Apps](#2-web-pages-websites-and-web-apps)
3. [HTML — The Structure of a Web Page](#3-html--the-structure-of-a-web-page)
4. [CSS — Making It Look Good](#4-css--making-it-look-good)
5. [JavaScript — Making It Interactive](#5-javascript--making-it-interactive)
6. [URLs — The Address of a Web Page](#6-urls--the-address-of-a-web-page)
7. [How a Server Serves a Web Page](#7-how-a-server-serves-a-web-page)
8. [HTTPS — Secure Web](#8-https--secure-web)
9. [The Web Has Changed Everything](#9-the-web-has-changed-everything)
10. [Summary](#summary)

---

## 1. The Web in Plain English

The **World Wide Web** (just "the web" for short) is a system for sharing and linking documents over the internet.

```
Before the web (before 1991):
  Information was locked in books, libraries, proprietary databases.
  Sharing required fax machines, postal mail, or calling someone.
  Finding information required knowing where to look.
  
After the web:
  Any information can be published to a URL accessible worldwide.
  Web pages link to other web pages (hyperlinks).
  A search engine can find anything in milliseconds.
  Anyone with internet access can publish.
```

**Three key ideas make the web work:**
1. **HTML** — a language for writing structured documents
2. **HTTP** — a protocol for transferring those documents
3. **URLs** — addresses that uniquely identify each document

Tim Berners-Lee invented all three at CERN (Switzerland) in 1989–1991. His original goal: make it easier for scientists to share research papers. What he created changed the world.

---

## 2. Web Pages, Websites, and Web Apps

```
Web page:
  A single document on the web.
  Has its own URL.
  Contains text, images, links.
  Example: this very page about the web (hypothetically)
  
Website:
  A collection of related web pages.
  Share a domain name (google.com, wikipedia.org).
  Usually maintained by one organization.
  Example: Wikipedia is a website. Each article is a web page.
  
Web app:
  A website that behaves like an application.
  Does things beyond displaying information.
  Example: Gmail (compose, send, receive email),
           Google Docs (edit documents),
           YouTube (stream video, comment, like).
  The line between website and web app is blurry.
```

---

## 3. HTML — The Structure of a Web Page

**HTML** (HyperText Markup Language) is the language used to write web pages. It's not a programming language — it's a **markup language** that tells a browser what's on the page.

```html
<!-- This is what HTML looks like: -->
<html>
  <head>
    <title>My First Web Page</title>
  </head>
  <body>
    <h1>Hello, World!</h1>
    <p>This is a paragraph of text.</p>
    <a href="https://google.com">Click here to visit Google</a>
    <img src="cat.jpg" alt="A photo of my cat">
  </body>
</html>
```

```
Reading the HTML tags:
  <h1>...</h1>  → Heading (big text)
  <p>...</p>    → Paragraph
  <a href="">   → Link (hyperlink to another page)
  <img src="">  → Image
  <div>...</div>→ Container (groups things together)
  
Tags always come in pairs: <tag> starts, </tag> ends.
```

**What HTML does NOT do:**
HTML only says WHAT is on the page. Not how it looks. The heading is a heading — but how big? What color? Where positioned? That's CSS.

---

## 4. CSS — Making It Look Good

**CSS** (Cascading Style Sheets) controls how HTML looks.

```css
/* This is CSS: */
h1 {
  color: blue;
  font-size: 48px;
  font-family: Arial;
}

p {
  color: black;
  line-height: 1.6;
  max-width: 600px;
}

body {
  background-color: white;
  margin: 20px;
}
```

```
HTML is the skeleton.
CSS is the clothes, hair, and makeup.

Same HTML + different CSS = completely different looking website.
This is why you can have a "dark mode" and "light mode":
Same content, different CSS applied.

Every website you've ever seen: all the colors, fonts, layouts,
spacing, animations — all CSS.
```

CSS is what makes Google's homepage white and minimal, while a gaming site looks dark and dramatic — the underlying HTML structure may be similar.

---

## 5. JavaScript — Making It Interactive

**JavaScript** (JS) is the programming language of the web. It makes pages interactive and dynamic.

```
Without JavaScript:
  Web pages are static — like a printed pamphlet.
  You can read, you can click links, but nothing reacts.
  
With JavaScript:
  Click a button → form appears without reloading the page
  Type in search box → results appear instantly as you type
  Maps → drag to move, pinch to zoom
  Notifications → "You have a new message"
  Infinite scroll → new content loads as you scroll down
  Image slider → photos rotate automatically
  
Examples of heavy JavaScript apps:
  Google Maps     → a full mapping application in the browser
  Gmail           → email client that runs entirely in JavaScript
  Figma           → professional design tool in the browser
  Google Docs     → word processor in the browser (no installation)
```

Together:
- HTML = structure (skeleton)
- CSS = style (appearance)
- JavaScript = behavior (interactivity)

Every modern website uses all three.

---

## 6. URLs — The Address of a Web Page

A **URL** (Uniform Resource Locator) is the address of any resource on the web.

```
https://www.example.com:443/blog/post?id=42&lang=en#comments
  │      │   │          │   │         │               │
  │      │   │          │   │         │               └─ Fragment (jump to section)
  │      │   │          │   │         └─ Query string (extra parameters)
  │      │   │          │   └─ Path (which page/file)
  │      │   │          └─ Port (usually hidden, 443 = default HTTPS)
  │      │   └─ Domain name
  │      └─ Subdomain (www = world wide web, but can be "blog." "shop." etc.)
  └─ Protocol (https = secure, http = not secure)

Common examples:
  https://google.com               → Google homepage
  https://en.wikipedia.org/wiki/Cat → Wikipedia article on cats
  https://github.com/username/repo → A GitHub repository
  https://youtube.com/watch?v=xxxxx → Specific YouTube video
```

---

## 7. How a Server Serves a Web Page

```
Step 1: You type https://wikipedia.org in your browser

Step 2: DNS lookup
  "What's the IP address for wikipedia.org?"
  Answer: 208.80.154.224

Step 3: Your browser connects to Wikipedia's server

Step 4: Browser sends HTTP request:
  GET /wiki/Computer HTTP/2
  Host: en.wikipedia.org

Step 5: Wikipedia's server:
  Looks up the "Computer" article in its database
  Constructs an HTML page with the article content
  Sends it back

Step 6: Browser receives HTML, then:
  Fetches CSS (for styling)
  Fetches JavaScript (for interactivity)
  Fetches images and other media

Step 7: Browser renders (draws) the complete page on your screen

Total time: ~200-500ms (less than half a second)
```

---

## 8. HTTPS — Secure Web

```
HTTP:  Data sent as plain text. Anyone watching the network can read it.
HTTPS: Data is encrypted. Even if intercepted, it's unreadable.

The "S" = Secure. Uses TLS (Transport Layer Security) encryption.

How to spot it:
  🔒 padlock in the browser address bar = HTTPS
  URL starts with https:// = encrypted
  
Why it matters:
  Without HTTPS: You type your bank password.
    → Travels across internet as plain text.
    → Anyone on the same Wi-Fi can read it.
    
  With HTTPS: You type your bank password.
    → Encrypted to random-looking garbage immediately.
    → Travels across internet.
    → Only the bank's server can decrypt it.
    
RULE: Never enter a password or credit card on a site without 🔒

Today (2026): ~90%+ of web traffic is HTTPS.
Browsers warn you loudly when a site is HTTP-only.
```

---

## 9. The Web Has Changed Everything

```
Before the Web (pre-1991):
  Information: in books, libraries, expert heads — hard to find
  Shopping: physical stores, mail-order catalogs
  Music: CDs, record stores
  TV: scheduled broadcasts, can't pause
  News: newspapers, TV news at specific times
  Communication: phone calls, postal mail, fax
  Education: classrooms, expensive courses
  Work: mostly in person, physical documents
  
After the Web (post-1991):
  Information: any fact, instantly, free (Wikipedia, Google)
  Shopping: everything, delivered (Amazon, everywhere)
  Music: any song, instantly (Spotify, YouTube)
  TV: any show, any time, from anywhere (Netflix, YouTube)
  News: constant, instant, from every angle (Twitter, websites)
  Communication: free, instant, video, global (email, WhatsApp)
  Education: any course, free or cheap (Khan Academy, YouTube, Coursera)
  Work: remote work, global teams, digital documents
```

The web is arguably the most transformative technology in human history — beyond electricity, beyond the printing press, in how rapidly it changed daily life globally.

---

## Summary

| Concept | What It Means |
|---------|--------------|
| Web | A system of linked documents on the internet |
| Web page | A single document with a URL |
| Website | A collection of related web pages |
| HTML | Markup language — what's on the page |
| CSS | Stylesheet language — how the page looks |
| JavaScript | Programming language — how the page behaves |
| URL | The address of a web page |
| HTTP/HTTPS | Protocol for transferring web pages (S = secure/encrypted) |
| DNS | Translates domain names to IP addresses |

**Now you know what the web is. Next: how do you navigate it? Meet browsers and search engines.**
