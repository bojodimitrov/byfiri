package cli

const openHelp = `open

Opens both file and directory

open [file name] - Prints contents of a file
open [path\file name] - Prints the content of a file that is located behind provided path
open .. - Changes current directory to the parent one
open [directory name] - Changes current directory 
open [path] - Changes current directory according to the specified path
`

const makeHelp = `make

Creates both file and directory

make d [directory name] - Creates a directory within the current one
make f [file name] - Creates file within the current directory. After issuing command the user is prompted to type the file contents:
	>some content|
	It is a multiline input. In order to quit input ':q' or ':Q' needs to be typed on a new line
`

const deleteHelp = `del

Delete both file or directory

del [directory/file name] - Deletes file or directory. If it is directory, all the contents of the directory are deleted as well.
`

const editHelp = `edit

Edits only file content

edit [file name] - Edit the content of a file by overwriting it. After issuing command the user is prompted to type the file contents:
	>some content|
	It is a multiline input. In order to quit input ':q' or ':Q' needs to be typed on a new line
`

const renameHelp = `ren

Renames both file and directory

ren [directory/file name] - Renames the given directory or file
`

const listHelp = `list

Prints the content in current directory
`

const treeHelp = `tree

Prints tree of all files and directories contained in current one
`

const exitHelp = `exit

Closes byfiri
`
