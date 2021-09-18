# img2sx: Image 2 MSX converter

Simple image converter that converts images between the following formats:

* PNG
* SC2 (applies resize)
* TO DO: more formats

It's also a Go library to handle SC2 files as `image.Image` interface implementors.

## Usage

```
img2sx -i <input_file> -opt <options> <output_file.png>
```

options. Comma-separated list:

* sc2 output:
```
crop: if source is larger than the image, it crops the destination image (default)
stretch: if source and dest file size are diferent, crops
aspect: if source and dest sizes are different, rescales keeping the aspect ratio
```

## sh2 format

Shrinked sc2 format: it merges into one all the equal tiles. All the tables and names are optional,
so you can have a single patterns table or no names table at all.

The format is:

* Header (4 bytes): `SH2\n`
* Indexes
    * Pattern table 1:
        - 1 byte, offset: end of header 
        - number of tiles (8 bytes each) belonging to the pattern table 1 and color table 1
    * Pattern tables 2 and 3:
        - 1 byte, offset: end of the previous pattern tables
        - number of tiles belonging to pattern table 2 or 3 and color tables 2 or 3
        - if `0`, it is a copy of Pattern table 1
    * Enabling name tables:
        - 1 byte, offset: end of pattern table 3
        - If `0`, there is no name tables
        - If `!0`, it should read 256 bytes for each of the three name
          tables.
* Byte contents:
    - patterns for pattern table 1
    - colors for pattern table 1
    - 

Example: hex dump of a very simple file with 2 tiles per pattern table:

```
53 48 32 0A     <-- SH2\n header
02 02 02        <-- each pattern table has 2 tiles
01              <-- name tables are enabled
00 01 02 03 04 05 06 07    <-- data of pattern table 1 starts here (2 tiles * 8 patterns)
08 09 0A 0B 0C 0D 0E 0F
10 11 12 13 14 15 16 17    <-- data of color table starts here (2 tiles x 8 patterns)
18 19 1A 1B 1C 1D 1E 1F
20 21 22 23 24 25 26 27    <-- data of pattern table 2
28 29 2A 2B 2C 2D 2E 2F
30 31 32 33 34 35 36 37    <-- data of color table 2
38 39 3A 3B 3C 3D 3E 3F
40 41 42 43 44 45 46 47    <-- data of pattern table 3
48 49 4A 4B 4C 4D 4E 4F
50 51 52 53 54 55 56 57    <-- data of color table 3
58 59 5A 5B 5C 5D 5E 5F
60 61 62 63 64 65 66 67    <-- name tables start here (3 * 256 bytes)
68 69 6A 6B 6C 6D 6E 6F
70 71 72....
```



## API

TO DO
