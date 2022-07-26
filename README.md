# aster-nft

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

※ _db.go_ 24行目記載

#### 認証設定

_db.go_ 20-23行目を、MySQLインストール時の設定に合わせる。

## 起動方法

### 起動

_main.go_ の存在するディレクトリで以下を実行する。

```
go run main.go
```

### アクセス

ブラウザから以下のURLにアクセスする。

http://localhost:8080/
