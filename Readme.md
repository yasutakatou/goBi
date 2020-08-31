# goBi

**語尾を自動に変換してくれるツール (Go言語製)**.

![yannsu](https://github.com/yasutakatou/goBi/blob/pic/yannsu.gif)

# できたもの

[ポートフォリオや個人開発で使えそうなアイデア](https://qiita.com/MasatoraAtarashi/items/eec4642fe1e6ce79304d)

![2](https://github.com/yasutakatou/goBi/blob/pic/2.png)

こちらの記事に書かれた大量のアイディアに触発されて作りました。

### v0.2 途中まで入力したものを、訂正して違う語尾ルールへ変換できるようにしました

![yannsu2](https://github.com/yasutakatou/goBi/blob/pic/yannsu2.gif)

# 環境

動作環境は**Windows系OS**のみです。<br>
キー入力している**ウィンドゥが変わったら、それまで入力してたカウンターをリセットする**処理にWindowsのAPIを使用しているためです。<br>
途中まで入力したのが、他の窓で暴発しないようにしたためで、それでも良いなら簡単に他OSへ移植できます。<br>

# インストール

パス通った所に置きたい場合はこちら

```
go get github.com/yasutakatou/goBi
```

バイナリ作りたいときはこっちです。

```
git clone https://github.com/yasutakatou/goBi
cd goBi
go build goBi.go
```

[バイナリをダウンロードして即使いたいならこっち](https://github.com/yasutakatou/goBi/releases).<br>

# アンインストール

リポジトリのフォルダとバイナリを消す

# 使い方

バイナリに**変換前**と、**変換したい単語**を"@"(デフォ)で区切って渡してください

```
>goBi.exe desu@deyannsu masita@tadeyannsu 190@yannsu
```

こう指定すれば”です”、”ました”、”。”が入力されると”やんす系”に変換します。<br>
数字一個の指定は対応しているASCIIコードが一個来たら変換するモードです

# オプション

```
Usage of goBi.exe:
  -debug
        [-debug=debug mode (true is enable)]
  -del int
        [-del=string delete key] (default 8)
  -split string
        [-split=string for split] (default "@")
  -zenkaku
        [-zenkaku=zenkaku mode (true is enable)] (default true)
```

-debugはデバッグモード。色々出力されます。<br>
-delは入力を消す操作です。デフォはバックスペースが割当たってます。<br>
-splitは変換前後の区切り文字。デフォは"@"です。<br>
-zenkakuは全角モード。つまり”です”なら4回キーを叩くけど、削るのは2文字なので半分消すモードです。もし英語で使うならfalseにしてください。<br>


# LICENSE

GPL-3.0 License

