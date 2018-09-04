# tbf: CLI for Tech Book Festival
[![CircleCI](https://circleci.com/gh/mpppk/tbf.svg?style=svg)](https://circleci.com/gh/mpppk/tbf)
[![codebeat badge](https://codebeat.co/badges/2cd1f4de-1e7d-4da3-900d-1bcb013c9448)](https://codebeat.co/projects/github-com-mpppk-tbf-master)
[![codecov](https://codecov.io/gh/mpppk/tbf/branch/master/graph/badge.svg)](https://codecov.io/gh/mpppk/tbf)

tbfは[技術書典](https://techbookfest.org)のサークル情報を取得するためのCLIツールです。

![](https://i.gyazo.com/8bd958b53fdc3e140f5bbe6b354c8194.gif)


# Usage
## fzf/pecoで絞り込んだサークル詳細ページをブラウザで表示する

```
tbf list | fzf -m | awk '{print $1}' | xargs tbf describe | jq -r .DetailURL | xargs open
```

## fzf/pecoで絞り込んだサークルのサイトをブラウザで表示する

```
tbf list | fzf -m | awk '{print $1}' | xargs tbf describe | jq -r .WebURL | xargs open
```
