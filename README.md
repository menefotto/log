# Log

[![GoDoc](https://godoc.org/github.com/wind85/log?status.svg)](https://godoc.org/github.com/wind85/log)
[![Build Status](https://travis-ci.org/wind85/log.svg?branch=master)](https://travis-ci.org/wind85/log)
[![Coverage Status](https://coveralls.io/repos/github/wind85/log/badge.svg?branch=master)](https://coveralls.io/github/wind85/log?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

### Log package

This package is a think wrapper around the standard library package syslog. Since
logging is a blocking operation can become quite expensive, specially in log heavy 
applications. This package is born in order to reduce to minum or totally avoid, 
blocking in fact offers two main functions Log and SyncLog.
The first of the two offers no blocking but does not grant that the message will
be logged, SyncLog should instead assure the message will be logged, if case or a
local logger ( that is a file logger ) is the syslog case can't be forced to sync.
It provides the following:


-New(prefix,logfilename) func that given a prefix string, and file name (optional) 
 construct a Logger, it can fails if for any reason the file can't be opened or
 created or there is a failure during the syslog instantiation.
 

-The other provided functions are Log() and SyncLog() which sends the message, the 
 first one send the message in a non blocking manner and the second one as well,
 though the second one syncs to disk in case of file logger assuring the log is sent.
 
-Close() closes the logger, it's important to do so because otherwise we would leak 
 the goroutine who does the job in the background.

#### How to use it
It's pretty simple
```
  l := log.New("[MYAPPNAME]","") // syslog
  l := log.New("[MYAPPNAME]","local.log") // local file log
  // then in order to log
  l.Log("My message") // logs optmistically
  l.SyncLog("My message") tries to sync to disk in case it's logger with local file
  // very important don't forget to close it once done
  l.Close() 
```

### Important
Due to a design decision this logger is an optimistic logger, that means you can 
afford to lose log messages, however to minimize the risk when you close the logger 
in the edge case that your program immediately terminates, please introduce a 
artificial sleep of at least 1 second.


#### Philosophy
This software is developed following the "mantra" keep it simple, stupid or better 
known as KISS. Something so simple like a cache with auto eviction should not required 
over engineered solutions. Though it provides most of the functionality needed by 
generic configuration files, and most important of all meaning full error messages.

#### Disclaimer
This software in alpha quality, don't use it in a production environment, it's not even completed.

#### Thank You Notes
None.
