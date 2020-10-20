# Logbus

**In development, do not use it in production.**

logbus is a log processing pipeline that from one `input datasource` like message queue,
to `multiple output channel` like stdout, files, common store db (MongoDB, MySQL, ES).

You can implement your own input/output interface.

**Note**: logbus is not support distribute.

## Feature

- Input source
  + [X] Redis List
- Output channel
  + [X] Stdout
  + [X] MongoDB
- Your can custom a transformer callback function processing log decide finally output log

## Intro

logbus define a standard log format `StdLog`, I use [logrusredis-hook](https://github.com/zhangjie2012/logrusredis-hook) output log to redis
`LIST` by `StdLogWash`, but all of this is not necessary.

logbus a pipeline framework, the core task is `Read` and `Write`, checkout `Serve` code.
**Only** you need do it inplement your owne `input`, `output`, `transformer`
(of course, you can use default, or PR a new).

For data processing, it provide:

- `DefaultTransformer` do nothing, all log passed.
- `StatLogTransformer` only pass which log has a valid `StateId`

## Usage

```
go get github.com/zhangjie2012/logbus
```

For example code `example/main.go`, log from redis `LIST` and to stdout/MongoDB, It's a real scenes for me.
You can implement your `input`, `output`, `transformer` and call `Serve` built up.
