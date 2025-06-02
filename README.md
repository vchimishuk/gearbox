**gearbox** is a non-interactive console client for transmission-daemon.

### Documentation
Comparing to `transmission-remote(1)` Gearbox aim to be an easy to use console
[Transmission](https://transmissionbt.com) client with clean an simple
interface. Gearbox's goal is not to implement all the functionality supported by
`transmission-remote(1)` but implement only the necessary one focusing on
keeping interface clean and simple. Another feature of Gearbox is that default
configuration can be stored in file. For instance, default list of columns
printed by `list` command can be stored in configuration file to prevent
entering it every time when using `list` command.

```
$ gearbox -H
gearbox is a non-interactive client for transmission-daemon

usage: gearbox [-H] [-h host] [-p port] command [opt]... [arg]...

commands:
  add [-S] [-l label[,label...]] file...
  delete [-d] id...
  edit [-l label[,label...]] id...
  info [-fi] id...
  list [-ar] [-c column] [-n count] [-s column]
  start id...
  stats
  stop id...
```
See `gearbox(1)` man page for more information.

### Examples
```
$ gearbox list -a
STATUS     SIZE      URATE   RATIO  NAME
seeding   27 GB    90 kB/s    5.90  Rainbow Discography
seeding  5.3 GB  2344 kB/s   19.72  Joan Jett & The Blackhearts - Discography
seeding  2.9 GB   475 kB/s   87.15  Red Rider - 1984
seeding   35 GB    24 kB/s  120.86  Loudness
seeding  2.1 GB     0 kB/s   77.99  X-Wild
```
```
$ gearbox add foo.torrent
```
```
$ gearbox info 1539
           Id: 1539
         Name: REVOLUTION OS
       Labels: video,documental
       Status: seeding
         Size: 3.7 GB
   Downloaded: 3.7 GB
     Uploaded: 77 GB
        Ratio: 20.95
Download rate: 0 kB/s
  Upload rate: 0 kB/s
  Last active: 23 Apr 25 18:11 +0300
        Added: 11 Apr 21 17:44 +0300
      Created: 27 Sep 08 00:20 +0300
      Comment: https://rutracker.org/forum/viewtopic.php?t=1132547
```

### Build and run
The app can build and run using standard `go` command.
```
$ go build
$ ./gearbox
```
