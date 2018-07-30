# qrss-grabbers-db
Code to reconcile QRSS Grabber databases

This program fetches 3 databases from the web.
See the end of main.go for those sources.

There are two outputs:
*   A complete merged list goes to stdout (and is captured at OUTPUT.txt).
*   A simple merged CSV file is written to MERGED.txt.

In the stdout, labels like [Andy] or [Henry] show what that database had different.
The first database to have a field is not marked.
