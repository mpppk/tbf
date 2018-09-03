# tbf: CLI for Tech Book Festival
tbfは[技術書典](https://techbookfest.org)のサークル情報を取得するためのCLIツールです。

![](https://i.gyazo.com/8bd958b53fdc3e140f5bbe6b354c8194.gif)


# Usage
## fzf/pecoで絞り込んだサークル詳細ページを表示する

```
tbf list | fzf -m | awk '{print $1}' | xargs tbf describe | jq -r .DetailURL | xargs open
```

## fzf/pecoで絞り込んだサークルのサイトを表示する

```
tbf list | fzf -m | awk '{print $1}' | xargs tbf describe | jq -r .WebURL | xargs open
```
