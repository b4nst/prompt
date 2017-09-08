Prompt
======

Prompt a keyboard interactive menu on your cli, and echo the selected option. 

Usage
-----

```shell
prompt Message option1 option2 ...
```
Navigate with <kbd>↑</kbd>, <kbd>↓</kbd>, <kbd>←</kbd>, <kbd>→ </kbd>, <kbd>TAB</kbd>, <kbd>SHIFT</kbd> + <kbd>TAB</kbd>, <kbd>SPACE</kbd>, <kbd>BACKSPACE</kbd>     
Validate or Cancel with <kbd>ENTER</kbd> or <kbd>ESC</kbd>


Menu is displayed in **Stderr**, selected option to **Stdout**

![screenshot](screenshot.gif)

Install
-------

```shell
go get github.com/bastienar/prompt
```