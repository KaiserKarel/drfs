# DRFS: Drive-Reply File System

## Storing your files in ordered reply threads.

### But why?

Well, although there are limits on the total amount of data stored in Google Drive, comment threads and replies do not 
contribute to this total. One can upload their entire music library for glacier storage at no cost whatsoever.

### Should I use this?

No. This is merely an experiment and a fun little project. Google will definitely notice the usage of replies as storage
and at least send you a warning, at most ban your account.

### How does it work?

Currently this is a WIP. The main module contains the file primites. 
`package os` implements functions which mimick the standard library 
`os` functions. `package recovery` contains helpers for reindexingand 
recovering files.

### Tests

Maybe more later. `package e2e` contains assorted tests.

### Bugs

Currently many. drfs can correctly upload and download `testdata/lorem_short.txt`. However `lorem.txt` results in some errors (16 missing bytes)
 
### Speed

The rate limits of Drive should allow for an upload speed of approximately 
400kB/s. Currently DRFS wastes 50% of API calls by updating file indexes. 
This can be reduced to about 2% by having a single global file index and updating 
this after every batch write.

Each (service) account has a rate limit of 10% of the project rate limit. Using multiple
accounts thus increases the amount of API calls by a factor 10.

