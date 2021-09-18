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


## API

TO DO
