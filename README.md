# driveignore

This simple cli tools works **together with** the [google drive sync](https://www.google.com/drive/download/) allowing you to have .driveignore files. Driveignore uses hardlinks, meaning no files duplicates, no repetitive cli calls, and blazing fast 'upload' times.

## installing

- make sure you have [golang](https://golang.org/dl/) installed
- run the `go get github.com/shilangyu/driveignore` command

Done! You will now have `driveignore` as a command in your terminal.

## how to use

You can get all the help about each command by using the `--help` (`-h`) flag.

- Create a folder and add it to your downloaded google drive watch list
- Create a `.driveignore` the same way you would a `.gitignore` in the root of a directory you wish to sync
- run `driveignore upload [path to your folder from step 1]`. The current working directory will be cloned to the drive folder with respect to the `.driveignore` blacklist

And you're done! Google drive will take care of the rest, which is syncing the files to the cloud. Once a file has been uploaded through `driveignore upload` you wont have to upload it again, google drive will listen to changes because the 'uploaded' files are hardlinks.

## global vs local .driveignore

You can create a global `.driveignore` using the `driveignore global [path to global driveignore]`, that way if you want to upload a directory without a `.driveignore` the global one will be used. You can also force a merge of local and global `.driveignore` during upload using the `--mergeIgnores` flag.