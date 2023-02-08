# aster-nft

## アプリ機能概要

アップロードされた画像を100分割し、抽選を行うことができるWebアプリケーション。

管理者：利用者、出品者情報の閲覧・管理。画像のアップロード。画像情報の閲覧。

利用者：画像の抽選

出品者：画像のアップロード

## 使用技術

Go、MySQL、HTML、CSS、JavaScript


## ディレクトリの役割

templates ・・・　HTMLテンプレート

view　・・・　HTMLファイル作成

db ・・・　データべース操作

handler　・・・　リクエスト処理、dbとの橋渡し

lib　・・・　画像処理など


## ER図

https://drive.google.com/file/d/1y-7Kb0TTKKOD1zcEFYzRh0gF9hG_ED_g/view?usp=share_link


## 技術的ポイント

画像リストから、抽選済み部がグレースケール化された画像を確認できます。

カーソルを合わせると、抽選者の情報などが閲覧できます。

画像データや情報を非同期通信で取得しています。

## 事前準備

### MySQL

#### MySQLのインストール

ブラウザから以下のURLにアクセスし、端末のOSに合わせたインストーラをダウンロードし、インストールする。

https://dev.mysql.com/downloads/mysql/

※バージョンは8.0を使用

#### DBの作成

以下を実行し、データベースを作成する。

```
CREATE DATABASE NFT;
```

#### ユーザの作成

上記で設定したDBに対し、作業ユーザを設定する

参照：【MySQL入門】ユーザー作成の方法を解説！8.0からの変更点も紹介　https://www.sejuku.net/blog/82303　

## 起動方法

### ローカル環境

環境変数を設定する
```
export DBUSER=（データベースに設定したユーザ)
export DBPASS=（ユーザに設定したパスワード)
```

main.goの存在するディレクトリで以下を実行する。

```
go run .
```

### 検証環境

環境変数を設定する
```
export DBUSER=aster
export DBPASS=kym7izpzt2tV
```

nft_site/に移動し、以下を実行する。

```
./nft-site
```

## アクセス

ブラウザから以下のURLにアクセスする。

ローカル　http://localhost:8080/admin/login　

検証環境　http://3.114.104.27:8000/admin/login　


## ソースコードの注意点
### 管理者アカウントの認証

管理者サイトの認証情報はファイルで管理しています

dataディレクトリにaccnt.txtとpswd.txtとして保管されています。

パスワードの暗号化にはsha512という規格のハッシュ関数を使用しています。

利用者/出品者サイトには、consumers/sellersテーブルにアカウント情報を登録することで、ログインできるようになります

### テンプレートの仕様

各サイトのトップページ（下参照）にアクセスした際、そのサイトに必要なテンプレートファイルが一気に読み込まれます (Parsefiles view.go)

管理者：/usrList

利用者：/lottery

出品者：/upload

上記のページを経由せずに別のページURLに移動すると、画面が表示できません。(ExecuteTemplateError handler.go)

### 画像処理

出品側と抽選側で画像処理の流れが異なります

#### 出品側の画像処理

画像をアップロードした際、オリジナル画像がuploaded/に保存されます

100分割された画像はout/original/ (画像名) /に保存されます。

画像リストから画像を表示するときは、1枚画と抽選済み部の情報がサーバ側から送られ、フロントエンドでグレースケール化処理を行います。(image.js)

#### 抽選側の画像処理

保有画像一覧から画像を表示した際、サーバ側で一部グレースケールの画像が生成されます。(lib.go)

lib.goでは、out/ (画像名) /に保存された分割画像をもとに一枚の画像を生成し、out直下に保存します。

その画像をsrcでmyImage.htmlに埋めこみ、表示します

### ローカル環境から検証環境への対応

#### <mail.go> http.ListenAndServe(":8000", nil)

 :8080　→　:8000

#### <db.go> ConnectDb()のcfgのAddrの値

localhost　→　orangebot.cluster-czickfmhh6ua.ap-northeast-1.rds.amazonaws.com

#### <*.js> APIの向き

localhost:8080 → 3.114.104.27:8000










