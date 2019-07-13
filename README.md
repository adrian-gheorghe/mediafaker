# mediafaker

[![CircleCI](https://circleci.com/gh/adrian-gheorghe/mediafaker.svg?style=svg)](https://circleci.com/gh/adrian-gheorghe/mediafaker)

mediafaker is a small utility written in go that creates a fake version of a source directory passed to it. mediafaker is dedicated to the developers that need to work on local environments without the need to copy static resources from production environments.

The path can be local or over SSH provided correct credentials are passed to the utility.

# Manual
```bash
This utility creates a fake simplified version of a directory tree of your chosing, making it easier to work locally on legacy projects that have a large media asset folder.

Usage:
  mediafaker [command]

Available Commands:
  help        Help about any command
  local       runs mediafaker on a local directory
  ssh         runs mediafaker on a remote tree file accessible via ssh
  url         runs mediafaker from a tree json stored remotely

Flags:
  -d, --destination string   Local Destination directory path where mediafaker should store the files
  -e, --extcopy strings      List of extensions that should be copied automatically
  -h, --help                 help for mediafaker
  -j, --jsonlog              Change logger format to json
  -m, --maxcopy int          Maximum Size(in bytes) a file should have to be copied automatically if it cannot be faked (default 30000)
      --version              version for mediafaker

Use "mediafaker [command] --help" for more information about a command.
```

## Installation
```bash
wget -O https://raw.githubusercontent.com/adrian-gheorghe/mediafaker/master/install.sh | bash
```
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/adrian-gheorghe/mediafaker/master/install.sh)"
```

## Download
Download latest from the releases page: https://github.com/adrian-gheorghe/mediafaker/releases.

# Usage
mediafaker can be used in 2 ways, to either fake a local path or a remote url / ssh path.

## Local
The simplest way is when the path you want to fake is on the same host as the destination path. Then you can run mediafaker local

```sh
mediafaker local \
    --source="/opt/media" \
    --destination="/home/project/public/fake/destination"
```
## Remote
When the source and destination paths are on different hosts mediafaker uses https://github.com/adrian-gheorghe/moni in order to generate a json representation of the directory you want to fake.

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

### URL
You can use moni https://github.com/adrian-gheorghe/moni or any other tool to generate a json output of the following form. mediafaker can use this to generate a fake version of the tree in your desired source.

```bash
moni url \
    --source="http://example.org/path/to/moni.output.json" \
    --destination="/home/project/public/fake/destination"
```

This flow can easily be automated to give you access to a fresh abstracted version of the tree on demand. 
You can get moni to run on your server periodically and create a new version of the output.json file as often as you like.

You can run mediafaker locally or on your dev/stage environments in order to have an up to date faked version of the public assets in production.

### SSH
mediafaker has a ssh client built in. When calling the ssh command, mediafaker attempts to:
- connect to your remote host using the credentials provided
- download moni from github
- run moni to generate a new output.json file
- downloads the compressed output file 
- removes moni and the output tree from your remote host
- fakes the output tree locally

```bash
mediafaker ssh \
    --source "/var/www/html/public" \
    --destination "/home/project/public/fake/destination" \
    --ssh-host "22.22.22.22" \
    --ssh-user "user" \
    --ssh-key "/home/.ssh/id_rsa_private_key_to_use"
```

## Image Pixelation
One of the main features of mediafaker is that it will read image files and create a pixelated version of the original image keeping the same dimensions and main colors. This saves time and precious disk space while also allowing you to have a drop in replacement for your images

filefaker will proceed to create an exact replica of the source directory, by creating mock files. The mock files created try to follow the originals as much as possible.

## Documents
- Xlsx support added using the https://github.com/tealeg/xlsx package 
- Pdf support added using the https://github.com/johnfercher/maroto package
- Docx, Pptx, CSV to be added
- Audio and Video file support to be added soon (mp3, wav, mp4, avi)
- Json and CSV file support to be added soon

## Features not yet implemented
- directory/file permissions
- fake from mediatype
- fake file and directory names
- directory depth parameter
- file age ignore parameter
