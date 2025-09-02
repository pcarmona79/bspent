# bspent - BSP files utility

Allows you to extract or parse the entities text inside a Q2 BSP map (.bsp) or an entities file (.ent).

Information and error messages will be printed to stderr and the entities definition will always be printed to stdout. This allows you to redirect the stdout to a file.

To build you'll need Go installed in your system, then run in a console:

``` sh
make
```

Or, if you're in Windows, type

``` batchfile
make windows
```

This will create a ``bin/`` directory that will contain the bspent executable (bspent.exe when building for Windows).

## Usage:

```
bspent <-p|-x> <filename>
  -p: Parse entities inside BSP file.
  -x: Writes to standard output the entities of the BSP file.
  <filename> must be a .bsp or .ent file.
```

To extract:

``` sh
./bspent -x file.bsp > file.ent
```

To parse:

``` sh
./bspent -p file.bsp > file.ent
```

You also can parse several files adding them to the command line:

``` sh
./bspent -p file1.bsp file2.bsp file3.bsp
```
