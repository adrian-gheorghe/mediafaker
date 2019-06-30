# mediafaker

[![CircleCI](https://circleci.com/gh/adrian-gheorghe/mediafaker.svg?style=svg)](https://circleci.com/gh/adrian-gheorghe/mediafaker)

mediafaker is a small utility written in go that creates a fake version of a source directory passed to it. mediafaker is dedicated to the developers that need to work on local environments without the need to copy static resources from production environments.

The path can be local or over SSH provided correct credentials are passed to the utility.

# usage
the mediafaker binary can be executed in 3 different ways

## local
```sh
mediafaker --source sourceDirPath --destination destDirPath
```
## ssh

## from moni tree export
Using https://github.com/adrian-gheorghe/moni you can create a json tree export of the directory you want to fake. 
This json file can be passed into filefaker to create a fake directory 

```json
[
    {
        "Path": "./testdata",
        "Type": "directory",
        "Mode": "drwxr-xr-x",
        "Size": 288,
        "Modtime": "2019-06-15 12:15:50.422928417 +0100 BST",
        "Sum": "",
        "MediaType": "",
        "Content": "",
        "ImageInfo": {
            "Width": 0,
            "Height": 0,
            "PixelInfo": null,
            "BlockWidth": 0,
            "BlockHeight": 0
        },
        "Children": null
    },
    {
        "Path": "testdata/.DS_Store",
        "Type": "file",
        "Mode": "-rw-r--r--",
        "Size": 8196,
        "Modtime": "2019-06-15 12:15:50.423694836 +0100 BST",
        "Sum": "46e1c0012e80786a970b3f8569222b51",
        "MediaType": "application/octet-stream",
        "Content": "",
        "ImageInfo": {
            "Width": 0,
            "Height": 0,
            "PixelInfo": null,
            "BlockWidth": 0,
            "BlockHeight": 0
        },
        "Children": null
    },
    {
        "Path": "testdata/a",
        "Type": "directory",
        "Mode": "drwxr-xr-x",
        "Size": 160,
        "Modtime": "2019-01-27 16:45:34.039871793 +0000 GMT",
        "Sum": "",
        "MediaType": "",
        "Content": "",
        "ImageInfo": {
            "Width": 0,
            "Height": 0,
            "PixelInfo": null,
            "BlockWidth": 0,
            "BlockHeight": 0
        },
        "Children": null
    },
    {
        "Path": "testdata/a/ac",
        "Type": "directory",
        "Mode": "drwxr-xr-x",
        "Size": 96,
        "Modtime": "2019-01-27 16:45:33.875895413 +0000 GMT",
        "Sum": "",
        "MediaType": "",
        "Content": "",
        "ImageInfo": {
            "Width": 0,
            "Height": 0,
            "PixelInfo": null,
            "BlockWidth": 0,
            "BlockHeight": 0
        },
        "Children": null
    },
    {
        "Path": "testdata/a/ac/acd.txt",
        "Type": "file",
        "Mode": "-rw-r--r--",
        "Size": 0,
        "Modtime": "2019-01-27 16:45:33.875876063 +0000 GMT",
        "Sum": "d41d8cd98f00b204e9800998ecf8427e",
        "MediaType": "text/plain; charset=utf-8",
        "Content": "",
        "ImageInfo": {
            "Width": 0,
            "Height": 0,
            "PixelInfo": null,
            "BlockWidth": 0,
            "BlockHeight": 0
        },
        "Children": null
    }
]
```

## Image Pixelation
One of the main features of mediafaker is it will read image files and create a pixelated version of the original image keeping the same dimensions and main colors. This saves time and precious disk space while also allowing you to have a drop in replacement for your images




filefaker will proceed to create an exact replica of the source directory, by creating mock files. The mock files created try to follow the originals as much as possible.


## TODOS
- pdf and doc files
- csv files
- load from export
- ssh deploy moni and retrieve info


Files use their initial names
Directory depth parameter
File Age ignore parameter
Documents created are valid, but empty.