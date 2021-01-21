# delta

A Difference Algorithm based on "An O(ND)Difference Algorithm and Its Variations" by Eugene W. Myers.

## Install

`$ go install github.com/stuartmscott/delta/cmd/delta`

## Usage

```
$ echo "Hello World" > file1
$ echo "Hi Earth" > file2
$ delta file1 file2
@ 1
- 101 e
- 108 l
- 108 l
- 111 o
+ 105 i
@ 3
-  87 W
- 111 o
+  69 E
+  97 a
@ 6
- 108 l
- 100 d
+ 116 t
+ 104 h
```
