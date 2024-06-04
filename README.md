# bolter

View BoltDB files in your terminal.
This is a fork of [hasit/bolter](https://github.com/hasit/bolter), since it is has been unmaintained for around a year.
This branch fixes a bug, comments the code, and aims to make he codebase cleaner.

![List all items](assets/viewbucket.gif)

## Install

```
$ go install github.com/vanillaiice/bolter/v2@latest
```

## Usage

```
NAME:
   bolter - view boltdb files interactively in your terminal

USAGE:
    [global options] command [command options] [arguments...]

VERSION:
   2.0.3

AUTHORS:
   Hasit Mistry <hasitnm@gmailcom>
   vanillaiice <vanillaiice1@proton.me>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file FILE, -f FILE  load boltdb FILE
   --no-values           do not print values (use if values are huge and/or not printable)
   --more more           use more to print all listings
   --help, -h            show help
   --version, -v         print the version

COPYRIGHT:
   (c) 2024 Hasit Mistry, vanillaiice
```

### List all buckets

```
$ bolter -f emails.db
+---------------------------+
|          BUCKETS          |
+---------------------------+
| john@doe.com              |
| jane@roe.com              |
| sample@example.com        |
| test@test.com             |
+---------------------------+
```

### List all items in bucket

```
$ bolter -f emails.db -b john@doe.com
Bucket: john@doe.com
+---------------+---------------------+
|      KEY      |        VALUE        |
+---------------+---------------------+
| emailLastSent |                     |
| subLocation   |                     |
| subTag        |                     |
| userActive    | true                |
| userCreatedOn | 2016-10-28 07:21:49 |
| userEmail     | john@doe.com        |
| userFirstName | John                |
| userLastName  | Doe                 |
+---------------+---------------------+
```

### Nested buckets

You can easily list all items in a nested bucket:

```
$ bolter -f my.db
+-----------+
|  BUCKETS  |
+-----------+
|   root    |
+-----------+

$ bolter -f my.db -b root
Bucket: root
+---------+---------+
|   KEY   |  VALUE  |
+---------+---------+
| nested* |         |
+---------+---------+

* means the key ('nested' in this case) is a bucket.

$ bolter -f my.db -b root.nested
Bucket: root.nested
+---------+---------+
|   KEY   |  VALUE  |
+---------+---------+
|  mykey  | myvalue |
+---------+---------+
```

### Machine friendly output

```
$ bolter -f emails.db -m
john@doe.com
jane@roe.com
sample@example.com
test@test.com

$ bolter -f emails.db -b john@doe.com -m
emailLastSent=
subLocation=
subTag=
userActive=true
userCreatedOn=2016-10-28 07:21:49
userEmail=john@doe.com
userFirstName=John
userLastName=Doe
nested-bucket*=
```

## Contribute

Feel free to ask questions, post issues and open pull requests on github.
When contributing, make sure to format your code with `gofmt`.
