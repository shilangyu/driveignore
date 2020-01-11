# driveignore

[![](https://github.com/shilangyu/driveignore/workflows/ci/badge.svg)](https://github.com/shilangyu/driveignore/actions)

This simple cli tools works **together with** the [google drive sync](https://www.google.com/drive/download/) allowing you to have .driveignore files. Driveignore uses hardlinks, meaning no files duplicates, no repetitive cli calls, and blazing fast 'upload' speeds.

## installing

- grab an executable from the [release tab](https://github.com/shilangyu/driveignore/releases)
- add it to `PATH`
- install [google drive sync](https://www.google.com/drive/download/)

or

- install [golang](https://golang.org/dl/)
- run the `go get github.com/shilangyu/driveignore` command
- install [google drive sync](https://www.google.com/drive/download/)

Done! You will now have `driveignore` as a command in your terminal.

## how to use

You can get all the help about each command by using the `--help` (`-h`) flag.

- Create an empty folder and add it to the google drive watch list
- Create a `.driveignore` the same way you would a `.gitignore` in the root of a directory you wish to sync
- run `driveignore upload [path to your folder from step 1]`. The current working directory will be cloned to the drive folder with respect to the `.driveignore` blacklist

And you're done! Google drive will take care of the rest, which is syncing the files to the cloud. Once a file has been uploaded through `driveignore upload` you wont have to upload it again, google drive will listen to changes because the 'uploaded' files are hardlinks.

## global vs local .driveignore

You can create a global `.driveignore` using the `driveignore global` (it will print the path to it), that way if you want to upload a directory without a `.driveignore` the global one will be used. You can also force a merge of local and global `.driveignore` during upload using the `--mergeIgnores` flag.

## help output

```
This simple cli allows you to have .driveignore(s)
It will look for a .driveignore, ignore the specified files
and make a hard link of your files to your drivesync folder
meaning no files duplicates, and no repetitive cli calls.

Usage:
  driveignore [command]

Available Commands:
  clean       Cleans your drive sync folder from old files
  diff        Compares your directory with the drive one
  global      Get the path to your global .driveignore
  help        Help about any command
  unify       Unifies 2 directories where input is the source
  upload      Upload a directory to your drive folder

Flags:
  -h, --help      help for driveignore
      --verbose   Prints out whats happening

Use "driveignore [command] --help" for more information about a command.
```

## linux? 😢

There is no drive sync for linux:

> There is no Drive app for Linux at this time. Please use Drive on the web and on your mobile devices.
