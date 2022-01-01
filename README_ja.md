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
(マネージャー)  
セットアップページから登録ページを生成して遷移します。  
登録ページのユーザー用cidを公開します。  
登録処理スイッチをonにして待機します。  

(ユーザー)  
登録ページのcidを入力して遷移します。  
userDataを入力して登録します。  
UserInfoが出力されるので、コピーして保持しておきます。  

## Voting
<!--
<img alt="voting" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/voting.png"><br>
-->
(マネージャー)  
セットアップページから登録ページcid(ユーザー用)を含む必要情報を入力して投票ページに遷移します。  
投票ページcid(ユーザー用)を公開します。
userDataを入力することで、そのユーザーが登録済みかどうか確認可能です。  
投票終了まで待機します。  
投票終了後、resultMapを生成します。  

(ユーザー)  
投票ページcidを入力して投票ページに遷移します。  
登録時に取得したuserInfoを入力します。  
投票フォームから投票内容を入力して投票します。   
投票終了後、マネージャーによって生成されたresultMapを用いて検証・集計を行います。  

<!--
```Go
type VoteInt map[string]int
type VotingData struct{
  data VoteInt
  enc []byte
}
votingData := voting.GenVotingData(voteInt)
```
-->
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




