# DRFS: Drive-Reply File System

## Storing your files in ordered reply threads.

### But why?

Well, although there are limits on the total amount of data stored in Google Drive, comment threads and replies do not 
contribute to this total. One can upload their entire music library for glacier storage at no cost whatsoever.

### Should I use this?

No. This is merely an experiment and a fun little project. Google will definitely notice the usage of replies as storage
and at least send you a warning, at most ban your account.

### How does it work?

Each branch contains a different version of GDFS. You are looking at `v1`, which is I mainly used to explore the Drive 
API and as a prototype. 

The usages at Drive are determined by total data stored; not files created nor associated data stored (such as 
file history, comments and replies). We can thus "cheat" their storage policy by putting all our data in comment threads.

Drive provides the following APIs to interact with files, comments and replies:

#### API Overview
*warning: oversimplification*
```
file.create     (...)         -> file{fileID}
file.get        (fileID)      -> file{...}
file.update     (file)        -> ok|error 

comment.create  (...)         -> comment{commentID}
comment.get     (commentID)   -> comment{...}
comment.update  (commentID)   -> ok|error

reply.create    (...)         -> reply{replyID}
reply.get       (replyID)     -> reply{...}
reply.update    (replyID)     -> ok|error 
```

These IDs are strings and non-sequential.

The `v1` implementation is quite naive. A `gdfs.File` is a file-comment association. Each comment reply contains at most
4096 bytes encoded as UTF-8, which is the maximum allowed size. Each file write creates a new reply. This model implements 
the following `io` interfaces:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

However due to the usage of a single comment thread per file; implementing [io.WriterAt](https://golang.org/pkg/io/#WriterAt) 
is quite difficult and inefficient. gdfs.Files cannot be concurrently written to, which is a requirement for `WriterAt`.

*The only option I see is, to lock reply creation and have calls to WriteAt update existing comments. Which I implement in v2* 