# xferspdy

This project aims to provide few things
* Generic command line utilities to diff and patch binary files
* Pluggable storage of files (Files can stored in S3 instead of disk, for example)
* Provide a server and client components to facilitate diff transfer
* P2P File updates

The codebase is almost entirely in golang and the core algorithms and ideas are as represented here Rsync Algorithm https://rsync.samba.org/tech_report/node2.html
