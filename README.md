# byfiri: a file system

the name comes from the words **by**te, because it is in-memory file system with a byte array in its core, **fi**le since it is a file system, and **ri** just because it sounds coolish.

to start it run **go build** and after that **byfiry.exe {size} {size abbreviation}**, ex. **byfiri.exe 1 GB**

## CLI Commands:

-   **open** {name}

expects file or directory name as parameter. If name is file prints content, if it is directory: enters the directory

-   **make** {f/d} {name}

with option f: creates file with given name and starts reading user input that will be the content of the file. Exiting is performed by typing ':q' or ':Q' on new line
with option d: creates empty directory with given name

-   **del** {name}

deletes file or directory with given name. Directory deletion deletes all directory contents recursively

-   **edit** {name}

edits file (but not directory) content. User input handled in the same way as make command

-   **ren** {old name} {new name}

renames file or directory

-   **list**

shows all files contained in a directory

-   **tree**

prints the contents of the current directory tree

-   **exit**

exits byfiri
