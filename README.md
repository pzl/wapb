What, Another Paste bin?
========================

`wapb` - Yep, another pastebin. I hate this. There are so many clones, it's stupid.


At least it's a small utility, so it's rebuilding a _wheel_ instead of an _entire car_.


I primarily use this as a paste-able handoff utility from phone to computer. I no longer want to email myself links. Paste from phone, fetch from computer. Without an app (so, browser-accessible from the phone). I looked at a solid [list of clones](https://github.com/awesome-selfhosted/awesome-selfhosted#pastebins), but found that the intersection of features I needed was not met by any one tool. Those features being:


- Single Binary distributable, without runtime requirements (goodbye python, ruby, php)
- Persistent storage
- Ability to set self-destroy timers
- Web UI
- Listing / discovery


**Close Contenders**

- [mkaczanowski/pastebin](https://github.com/mkaczanowski/pastebin) meets a lot of this, but is missing file upload, and listing entries.
- [raftario/filite](https://github.com/raftario/filite) also meets much of this, but without timers. It's also going through a rewrite by it's author at time of writing.
- [prologic/pastebin](https://github.com/prologic/pastebin) uses entirely in-memory caching for storage. Reboot will lose it all

**Nice-to-have features, but not mission-critical**

- File upload
- API


**NON Targeted features**

This is intended to be on a local network deployment, in a semi-trusted environment. So some features are out-of-scope:

- Encryption of contents
- Login / Authorization