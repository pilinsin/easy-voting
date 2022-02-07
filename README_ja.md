# EasyVoting
[IPFS](https://ipfs.io/)と[Fyne](https://fyne.io/)を使用したオンライン投票アプリです。  
ブロックチェーンは使用しません。

[日本語版IPFS解説サイト](https://ipfs-book.decentralized-web.jp/)
## Features
* 匿名投票
* 投票期間内ならば何度でも投票可能
* 誰でも全ての投票結果の検証が可能
* 誰でも投票結果の集計が可能


# Requirement
[go-ipfs](https://github.com/ipfs/go-ipfs)  
[fyne](https://github.com/fyne-io/fyne)

# Usage
<!--
<img alt="system_process" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/system_process.png"><br>
-->
## Registration
<!--
<img alt="registration" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/registration.png"><br>
-->
(登録マネージャー)  
セットアップページにて、情報を入力してmCfgCid (registration Manager Config CID)を得ます。  
mCfgCidを入力して登録マネージャーページに遷移します。  
ページに表示されるrCfgCid (Registration Config CID)を取得して公開します。  
登録処理スイッチをオンにして待機します。  

(ユーザー)  
rCfgCidを入力して登録ページに遷移します。  
userDataを入力して登録を行います。  

```Go
var userData []string
```

userIdentityが出力されるので、それをコピーして保持しておきます。  
```Go
type UserIdentity struct{
  userHash UserHash
  userPriKey *ecies.PriKey
  userSignKey *ed25519.SignKey
  rKeyFile *ipfs.KeyFile
}
type UserHash string
```

## Voting
<!--
<img alt="voting" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/voting.png"><br>
-->
(投票マネージャー)  
投票セットアップページにて、rCfgCidを含む幾つかの情報を入力し、mCfgCid (voting Manager Config CID)を得ます。   
mCfgCidを入力して投票マネージャーページに遷移します。   
vCfgCid (Voting Config CID)を取得して公開します。  
userDataを入力することで、そのユーザーが登録されているかどうか検証が可能です。  
投票期間終了後、resultMapを生成します。   

(ユーザー)  
vCfgCidを入力して投票ページに遷移します。  
userIdentityを含む幾つかの情報を入力して投票を行います。  
投票期間終了後、resultMapを用いて検証と集計が可能です。  


# Voting Type
以下の投票方式に対応しています。  
* [単記投票](https://ja.m.wikipedia.org/wiki/%E5%8D%98%E8%A8%98%E7%A7%BB%E8%AD%B2%E5%BC%8F%E6%8A%95%E7%A5%A8)  
* [連記投票](https://ja.m.wikipedia.org/wiki/%E9%80%A3%E8%A8%98%E6%8A%95%E7%A5%A8)  
* [承認投票](https://ja.m.wikipedia.org/wiki/%E8%AA%8D%E5%AE%9A%E6%8A%95%E7%A5%A8)  
* [範囲投票](https://ja.m.wikipedia.org/wiki/%E6%8E%A1%E7%82%B9%E6%8A%95%E7%A5%A8)  
* [累積投票](https://ja.m.wikipedia.org/wiki/%E7%B4%AF%E7%A9%8D%E6%8A%95%E7%A5%A8)  
* [選好投票](https://ja.m.wikipedia.org/wiki/%E9%81%B8%E5%A5%BD%E6%8A%95%E7%A5%A8)  


# TODO
* GUIデザイン
* バグ修正
* 登録時の新規ユーザーへの対応

# Support
どの組織にも属さずフリーランスで開発しています。  
このシステムは投票という制度をより身近で簡単で公平なものにすることを目的としています。  
開発継続のため、ご支援ご協力をお願いいたします。  

Ethereum Address: 0x81f5877EFC75906230849205ce11387C119bd9d8




