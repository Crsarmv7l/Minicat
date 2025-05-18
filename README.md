WIP Go implementation of Radiolib's CC1101 functionality specifically for my minicat project, 
a long running attempt to create a Yardstick One clone with an MCU and the CC1101.

The library protion that my program uses will be submitted to Go drivers once I am satisfied with it and flesh out more functions.

Like my implementation with Radiolib, this version also suffers some bitshifts and shortening when making large 
transmissions.

Please see here for discussion of that issue:
https://github.com/jgromes/RadioLib/discussions/1138



**Radiolib:**

All code here should be considered derivative work of Radiolib except for the code I have committed to radiolib
(including TX FIFO refills, autosetRXBandwidth, PQT changes and more).

Numerous liberties have been taken with this code base like removing the module abstraction present in radiolib, the combining of functions and other modifications.
