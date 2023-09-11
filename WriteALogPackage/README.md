# How logs work

A log is an append-only sequence of records. You append records to the end
of the log, and you typically read top to bottom, oldest to newest—similar to
running tail -f on a file. You can log any data. When you append a record to a log, the log assigns the record a unique and
sequential offset number that acts like the ID for that record. A log is like a
table that always orders the records by time and indexes each record by its
offset and time created.


### Segment
- We don't have infinitive space &rarr; can't append to them same file forever.
- Split the log into a list of segments.
- Log grows too big &rarr; free up disk space by deleting old segment we have already processed in background process.
- Special segment is active segment (a segment we actively write to). When we filled the active segment, we can create new active segment.

### Index
- Each segment comprises a store file(data) and index file.
- Index file is where we index each record in the store file.
- Requires 2 fields: offset and stored position of the record.
- Read file: Get entry from index file &rarr; read the record at that position.
- Index files are small enough that we can memory-map them and make operations on the file
as fast as operating on in-memory data.

## Build a log
- Record: the data stored in our log.
- Store: the file we store records in.
- Index: the file we store index entries in.
- Segment: the abstraction that ties a store and an index together.
- Log: the abstraction that ties all the segments together.