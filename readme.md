![](https://i.gyazo.com/792e6681b4a0ae5ec3414fb3529e3e4f.png)

[![CircleCI](https://circleci.com/gh/mpppk/tbf.svg?style=svg)](https://circleci.com/gh/mpppk/tbf)
[![codebeat badge](https://codebeat.co/badges/2cd1f4de-1e7d-4da3-900d-1bcb013c9448)](https://codebeat.co/projects/github-com-mpppk-tbf-master)
[![codecov](https://codecov.io/gh/mpppk/tbf/branch/master/graph/badge.svg)](https://codecov.io/gh/mpppk/tbf)

tbfは[技術書典](https://techbookfest.org)のサークル情報を取得/表示するためのCLIツールです。  
pecoやfzfなどのfuzzy finderと組み合わせることでサークル情報を絞り込むことができます。

![](https://i.gyazo.com/8bd958b53fdc3e140f5bbe6b354c8194.gif)

# Installation

## binary
Download from [GitHub Releases](https://github.com/mpppk/tbf/releases)

## brew

[TODO]

## From source

```Shell
$ go get github.com/mpppk/tbf
```


# Usage
## tbf list
技術書典ウェブサイトをクロールした結果のcsvやURLから、サークル情報を表示します。 
デフォルトでは[https://raw.githubusercontent.com/mpppk/tbf/master/data/latest_circles.csv](https://raw.githubusercontent.com/mpppk/tbf/master/data/latest_circles.csv
)からcsvを取得します。csvはローカルにキャッシュされますが、変更があった場合は再度ダウンロードします。

```
$ tbf list | head -n5
csv file will be downloaded becase checksums are different between meta(2413623400) and local file(2781509532)
new csv file is downloaded from https://raw.githubusercontent.com/mpppk/tbf/master/data/latest_circles.csv to latest_circles.csv
か46 トゲトゲ団（トゲトゲダン） by トゲトゲ 【ソフトウェア全般】 : ゲームエンジ ン(UnrealEngine4)
か77 ナナナナロク（ナナナナロク） by 776 【ソフトウェア全般】 : ストリーミング処理と可視化（予定）※VagrantとElasticsearch(Kibana)を軸としたTwitterデータ取得
け08 TY製作所（ティーワイセイサクジョ） by 吉野 【科学技術】 : 3Dプリンター及び レーザー加工機の取り扱いや造形物について機械ごとにまとめた漫画本あるいは解説本とグッズ
こ40 Firebase Japan User Group（ファイヤーベースジャパンユーザーグループ） by Firebase Japan User Group 【ソフトウェア全般】 : Firebaseについて
か08 Route 312（ルートサンイチニ） by mzsm 【ソフトウェア全般】 : Python製Webフ レームワークDjangoのTips集とか
```

### Tips: fuzzy finderで絞り込んだサークル詳細ページをブラウザで表示する
あらかじめpeco/fzfなどのfuzzy finderとjqをインストールしておく必要があります。

#### mac + peco (ctrl + spaceで複数選択)

```
tbf list | peco | awk '{print $1}' | xargs tbf describe | jq -r .DetailURL | xargs open
```

#### mac + fzf (tabで複数選択)

```
tbf list | fzf -m | awk '{print $1}' | xargs tbf describe | jq -r .DetailURL | xargs open
```

#### linux + peco (ctrl + spaceで複数選択)

```
tbf list | peco | awk '{print $1}' | xargs tbf describe | jq -r .DetailURL | xargs xdg-open
```

#### linux + fzf (tabで複数選択)

```
tbf list | fzf -m | awk '{print $1}' | xargs tbf describe | jq -r .DetailURL | xargs xdg-open
```

#### windows + peco (ctrl + spaceで複数選択)

```
TODO
```

#### windows + fzf (tabで複数選択)

```
TODO
```

#### docker(URLの表示のみ)

```
$ docker run -it --rm mpppk/tbf list_circle_url
https:/techbookfest.org/event/tbf05/circle/47030001
https:/techbookfest.org/event/tbf05/circle/35050003
https:/techbookfest.org/event/tbf05/circle/45010003
```

### Tips: fuzzy finderで絞り込んだサークルのサイトをブラウザで表示する

#### mac + peco (ctrl + spaceで複数選択)

```
tbf list | peco | awk '{print $1}' | xargs tbf describe | jq -r .WebURL | xargs open
```

#### mac + fzf (tabで複数選択)

```
tbf list | fzf -m | awk '{print $1}' | xargs tbf describe | jq -r .WebURL | xargs open
```

## tbf describe
引数として与えられたスペースに配置されたサークル情報をjsonで出力します。

```
$ tbf describe あ01 あ02
{"DetailURL":"https:/techbookfest.org/event/tbf05/circle/24830001","Space":"あ01","Name":"毬栗ロ マン（イガグリロマン）","Penname":"いっこう","Genre":"ソフトウェア全般","ImageURL":"https://lh3.googleusercontent.com/3HYptOqYpzH0-ZaaG55rG7vk1COYse9e6tcZBX5DlsAilF_67wwOXVVB7oVb5-mRHC1z6Z5QODOVkazRp9-5kQ","WebURL":"","GenreFreeFormat":"WebXR(Webブラウザで実現するVR並びにAR)の解説本、これまで頒布してきたものが基礎的な内容だったので今回はより具体的な内容を想定しています。"}
{"DetailURL":"https:/techbookfest.org/event/tbf05/circle/28360002","Space":"あ02","Name":"いしだ け（イシダケ）","Penname":"t_ishida,コンドウアヤ","Genre":"ソフトウェア全般","ImageURL":"https://lh3.googleusercontent.com/DT5P_6OcobmPWDzVnu1loCAt_DDrcmQ8P2Y5hE3RWoRb6Fx-4dcuA7U3oPP3yQyAXr3FzH-6Jc8_iI5Z_1Pp","WebURL":"http://www.dezapatan.com","GenreFreeFormat":"体系的なプログラミング制作 を目指してPHPで緩く解説しています"}

```

## tbf crawl
 chromeを起動して技術書典ウェブサイトからサークル情報をクローリングし、csvとして保存します。  
`tbf list`ではデフォルトでクロール済みのcsvをHTTP経由で取得するため、通常このコマンドを実行する必要はありません。
また、技術書典ウェブサイトへ継続的にリクエストを送ることになるので、利用には注意してください。

```
$ tbf crawl
# → ウェブサイトをクローリングし、結果をcircles.csvという名前で保存する
```
