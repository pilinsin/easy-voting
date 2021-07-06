# EasyVoting
IPFSとGUIライブラリを使用したオンライン投票アプリです。  
ブロックチェーンは使用しません。

# Requirement
[go-ipfs](https://github.com/ipfs/go-ipfs)  
gui

# Usage
<img alt="system_process" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/system_process.png"><br>
## Online Voter Registration
<img alt="registration" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/registration.png"><br>

(ユーザー)  
RSA鍵生成を行い、秘密鍵はローカルに保存します。  
公開鍵をIPFSにaddし、任意のIPNSキーでpublishします。  
メールアドレスとpublishしたIPNSのアドレスをサーバーに登録します。  

## Voting Setup
<img alt="voting_setup" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/voting_setup.png"><br>

(マネージャー)  

サーバーからメールアドレスと登録用IPNSアドレスのリストを取得します。  

各ユーザーに対して以下の処理を行います。  
1. 登録用IPNSアドレスからユーザー公開鍵を取得  
2. ユーザーidの生成  
3. KeyFileの生成  

```Go
userID := util.GenUniqueID(30,6)
KeyFile := ipfs.KeyFileGenerate()
```

4. ユーザーidとKeyFileをユーザー公開鍵で暗号化  
5. 暗号化したユーザーidとKeyFileをユーザーにメールで送信  
  

全ユーザーのKeyPairに対応する投票用IPNSアドレスをリスト化します。  
マネージャー公開鍵&秘密鍵を生成します。  
投票idを生成します。  

```Go
votingID := util.GenUniqueID(30,30)
```

VotingInfoを生成してIPFSにaddし、そのCIDを公表します。  

```Go
type VotingInfo struct{  
  votingID        string   
  manPubKey       rsa.PublicKey  
  begin           string  
  end             string  
  votingType      string  
  candidates      map[string]Candidate  
  votingIPNSAddrs []string  
}  
type Candidate struct{  
  url      string  
  group    string  
  groupURL string  
}  
```


## Voting
<img alt="voting" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/voting.png"><br>

(ユーザー)  
VotingInfoを取得します。  
メールから暗号化ユーザーidと暗号化KeyFileを取得します。  
ローカルに保存しておいたユーザー秘密鍵でユーザーidとKeyFileを取得します。  
KeyFileを入力し対応する投票用IPNSのアドレスを求めます。  
そのアドレスを投票用IPNSアドレスリストと比較してログイン認証を行います。  
投票方式を投票フォームに反映させます。  
投票データを生成します。  

```Go
type VoteInt map[string]int
votingData := map[string]VoteInt{userID: vote}
//or  
//type VoteBool map[string]bool
//votingData := map[string]VoteBool{userID: vote}  
```

投票データをマネージャー公開鍵で暗号化します。  
IPFSにaddしてKeyFileを用いて投票用IPNSにpublishします。  

## Counting Setup
<img alt="counting_setup" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/counting_setup.png"><br>

(マネージャー)  
VotingInfoを取得します。  
投票用IPNSアドレスリストから暗号化投票データを収集します。  
マネージャー秘密鍵で投票データを取得します。  
投票データ全体をIPFSにaddし、そのCIDを公表します。  
   
## Counting
<img alt="counting" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/counting.png"><br>

(ユーザー)  
投票データ全体を取得します。  
自身のユーザーidから投票内容を確認します。  
投票データリストを集計して投票結果を取得します。  

# Voting Type
以下の投票方式に対応  
・[単記投票](https://ja.wikipedia.org/wiki/%E5%8D%98%E8%A8%98%E7%A7%BB%E8%AD%B2%E5%BC%8F%E6%8A%95%E7%A5%A8)  
・[連記投票](https://ja.wikipedia.org/wiki/%E9%80%A3%E8%A8%98%E6%8A%95%E7%A5%A8)  
・[認定投票](https://ja.wikipedia.org/wiki/%E8%AA%8D%E5%AE%9A%E6%8A%95%E7%A5%A8)  
・[範囲投票](https://ja.wikipedia.org/wiki/%E6%8E%A1%E7%82%B9%E6%8A%95%E7%A5%A8)  
・[累積投票](https://ja.wikipedia.org/wiki/%E7%B4%AF%E7%A9%8D%E6%8A%95%E7%A5%A8)  
・[選好投票](https://ja.wikipedia.org/wiki/%E9%81%B8%E5%A5%BD%E6%8A%95%E7%A5%A8)  


# TODO
・登録処理  
・集計処理  
・GUI  
・IPNSを使用する処理の高速化

# Support
どの組織にも属さずフリーランスで開発しています。  
投票という制度をより身近で容易かつ公正なものにするため、無料でリリースする予定です。  
開発の継続のため、皆様のご支援が必要です。  
どうかよろしくお願いいたします。

ETH Address: 0x81f5877EFC75906230849205ce11387C119bd9d8
