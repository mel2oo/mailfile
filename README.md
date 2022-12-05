# mailfile

mailfile是一个邮件文件解析库，支持解析eml/msg邮件文件。



## Installation

```
go get -u github.com/mel2oo/mailfile
```



## Example

### MSG:

```
msg, err := msg.New("testdata/complete.msg")
if err != nil {
	return
}

eml.Format().Output()
```



### EML:

```
eml, err := eml.New("testdata/2.eml")
if err != nil {
	return
}

eml.Format().Output()
```

