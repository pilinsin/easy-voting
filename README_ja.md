# EasyVoting
[I2P](https://geti2p.net/en/)と[Fyne](https://fyne.io/)を使用したオンライン投票アプリです。  
ブロックチェーンは使用しません。

## Features
* 匿名投票
* 投票期間内ならば何度でも投票可能
* 誰でも投票結果の集計が可能

## Usage
### I2P Setup
[I2P](https://github.com/i2p/i2p.i2p)をインストールして起動してください。  
起動後、[SAM](https://geti2p.net/en/docs/api/samv3)が有効になっているか確認し、無効ならば有効にしてください。  
### Bootstrap
(登録マネージャー)  
セットアップページにて、自身のI2Pブートストラップを生成すると、そのアドレスが出力されます。  
自身のブートストラップを生成するか他のブートストラップのアドレスを入力することでブートストラップのリストが得られ、そのアドレスが出力されます。
以下、ブートストラップのリストのアドレスをbaddrsとします。
### Registration
(登録マネージャー)  
セットアップページにて、
- 登録ページのタイトル
- 登録ユーザーデータセット(csv)
- baddrs

を入力します。  
登録ユーザーデータセットはユーザーの登録に必要なデータをまとめたcsvファイルで、  

| label 0 | label 1 | ... | label M |
| --- | --- | --- | --- |
| user 00 | user 01 | ... | user 0M |
| user 10 | user 11 | ... | user 1M |
| ... | ... | ... | ... |
| user N0 | user N1 | ... | user NM |

という形式です。

入力後、Submitボタンを押すとrCfgAddrとmanIdentityが出力されるので、それらをロードページフォームに入力して登録ページに遷移します。  
rCfgAddrを公開して登録が完了するまで待機します。

(ユーザー)  
rCfgAddrを入力して登録ページに遷移します。  
登録に必要なデータを入力して登録を行います。  
userIdentityが出力されるので、それをコピーして保持しておきます。  

### Voting
(投票マネージャー)  
投票セットアップページにて、
- 投票ページのタイトル
- 開始時刻
- 終了時刻
- タイムゾーン(Location)
- rCfgAddr
- 認証者数(投票時刻を証明するのに必要な件数)
- 候補者情報
  - 画像
  - 氏名
  - グループ名
  - URL
- 投票パラメータ
  - 最小票数
  - 最大票数
  - 合計票数
- 投票方式

を入力します。  
入力後、Submitボタンを押すとvCfgAddrとmanIdentityが生成されるので、ロードページフォームに入力して投票ページに遷移します。  
vCfgAddrを公開して投票終了まで待機します。  
投票終了後にマネージャーも投票ボタンを押すことで投票データの復号鍵を公開します。  

(ユーザー)  
vCfgAddrとuserIdentityを入力して投票ページに遷移します。  
userIdentityを入力しない場合、投票結果の集計のみが可能です。   

投票フォームから投票データを生成して投票を行います。  
投票マネージャーが投票データの複合鍵を公開していれば、自身の投票データの確認と投票結果の集計が可能です。　　

# Voting Type
以下の投票方式に対応しています。  
* [単記投票](https://ja.m.wikipedia.org/wiki/%E5%8D%98%E8%A8%98%E7%A7%BB%E8%AD%B2%E5%BC%8F%E6%8A%95%E7%A5%A8)  
* [連記投票](https://ja.m.wikipedia.org/wiki/%E9%80%A3%E8%A8%98%E6%8A%95%E7%A5%A8)  
* [承認投票](https://ja.m.wikipedia.org/wiki/%E8%AA%8D%E5%AE%9A%E6%8A%95%E7%A5%A8)  
* [範囲投票](https://ja.m.wikipedia.org/wiki/%E6%8E%A1%E7%82%B9%E6%8A%95%E7%A5%A8)  
* [累積投票](https://ja.m.wikipedia.org/wiki/%E7%B4%AF%E7%A9%8D%E6%8A%95%E7%A5%A8)  
* [選好投票](https://ja.m.wikipedia.org/wiki/%E9%81%B8%E5%A5%BD%E6%8A%95%E7%A5%A8)  


# Support
どの組織にも属さずフリーランスで開発しています。  
このシステムは投票という制度をより身近で簡単で公平なものにすることを目的としています。  
開発継続のため、ご支援ご協力をお願いいたします。  

- BitCoin (BTC)
bc1qu0zl5z4zgvx2ar3zgdgmt3thl3fnt0ujvzdx9e
- Ethereum (ETH)
0x81f5877EFC75906230849205ce11387C119bd9d8
- Tron (TRX)
TCc7D7thmW4egbiUEk2uH3Y21shfbjVNvn
- Monero (XMR)
49y6hymbjLqf1LRrGoARqGNxD95UeHtpGbfYmutrLZaWhfFwefPHkDUiKkab3aCNBv36xAUu4VQus1V1g8hhYWrhLemRjPt
- Zcash (ZEC)
t1diehqpgftGp9dvEMKcAoCUxZnGgodcU96
- Basic Attention Token (BAT)
0xe83D64a10256aE37d3039344fE49ec9D1d75dd5c
- FileCoin (FIL)
f1mhulmnu4apv3thlsnmjw3nigl5hzcgozfabpsyi
