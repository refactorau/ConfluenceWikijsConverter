# Quick and dirty Confluence Export to WikiJS converter

Reads the exported confluence file and...
* copies the file as is into the wikijs import folder, or
* tries to work out if the file is part of a hierarchical structure, in which case...
  * detects where in the folder structure the file should go.
  * removes any breadcrumbs embedded into the html file
  * updates any links to attachments

## Usage
Optionally build a binary using `go build .`

Alternatively compile and run using `go run .`

If you have 
a) a folder of exported confluence files in the folder "FromConfluence",
b) a folder that wikijs is pointed to in the Storage->Local File System section, e.g. /opt/wikijs,
d) a root folder within that directory you want to import your files to wikijs, e.g. MyConfluence,
e) you have cloned this repository to the "ConfluenceWikijsConverter" folder, then
f) run the following command:


```
cd ConfluenceWikijsConverter
go run ../FromConfluence /opt/wikijs/MyConfluence
```

Then within wikijs admin...
a) Storage -> Local File System
b) Ensure the "Path" is set to the correct location and the local file system setting is enabled
c) Ensure that the Apply button has been clicked if either of those settings were modfied.
d) Click the "Import Everything" button and watch the logs for any errors.

## License
MIT
