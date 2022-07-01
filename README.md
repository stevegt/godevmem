# godevmem

This project is a Go library and CLI based on a rough translation of
the devmem utility for reading/writing to /dev/mem on a Linux machine.
This technique is often used for GPIO and PRU access on Raspberry Pi
or Beaglebone SBCs.

The version of devmem this translation is based on is the devmem2
project at https://github.com/rcn-ee/devmem2.  I chose this one
because I know it to work on a Beaglebone Blue for PRU shared memory
access; I haven't used or tested it myself in any other way.  

**WARNING: /dev/mem access can brick devices or cause mechanical
movement when motors and actuators are involved; use at your own
risk.**


