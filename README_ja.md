# EasyVoting
IPFSとGUIライブラリを使用したオンライン投票アプリです。  
ブロックチェーンは使用しません。

## Features
* 匿名投票
* 投票のやり直し
* 投票結果の検証
* 投票結果の集計


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
投票idを生成します。  
投票用IPNSアドレスマップを生成します。  
検証鍵マップを生成します。  

```Go
votingID := util.GenUniqueID(30,30)
var votingIPNSAddrs map[string]string
var verfKeys        map[string]rsa.PublicKey
```

サーバーからメールアドレスと登録用IPNSアドレスのリストを取得します。  

各ユーザーに対して以下の処理を行います。  
1. 登録用IPNSアドレスからユーザー公開鍵を取得  
2. ユーザーidの生成  
3. keyFileの生成  

```Go
userID := util.GenUniqueID(30,6)
keyFile := ipfs.GenKeyFile()
```

4. RSA鍵生成で署名鍵と検証鍵を生成  
5. ユーザーidとkeyFileと署名鍵をユーザー公開鍵で暗号化  
6. 暗号化したユーザーidとkeyFileと署名鍵をユーザーにメールで送信  
7. keyFileに対応する投票用IPNSアドレスを求める
8. 投票idとユーザーidのハッシュ値を得る

```Go
hash := util.Hash(userID, votingID)
```

9. ハッシュ値をキー、投票用IPNSアドレスを値としてマップに追加する  
10. ハッシュ値をキー、検証鍵を値としてマップに追加する  

```Go
votingIPNSAddrs[hash] = addr
verfKeys[hash] = verfKey
```
 
マネージャー公開鍵&秘密鍵を生成します。  
VotingInfoを生成してIPFSにaddし、そのCIDを公表します。  

```Go
type VotingInfo struct{  
  votingID        string   
  manPubKey       rsa.PublicKey  
  begin           string  
  end             string  
  votingType      string  
  candidates      map[string]Candidate  
  votingIPNSAddrs map[string]string
  verfKeys        map[string]rsa.PublicKey 
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
メールから暗号化ユーザーidと暗号化keyFileと暗号化署名鍵を取得します。  
ローカルに保存しておいたユーザー秘密鍵でユーザーidとkeyFileと署名鍵を取得します。  
keyFileを入力し対応する投票用IPNSのアドレスを求めます。  
そのアドレスを投票用IPNSアドレスマップと比較してログイン認証を行います。  
投票方式を投票フォームに反映させます。  
署名鍵で署名を付与した投票データを生成します。   

```Go
type VoteInt map[string]int
type VotingData struct{
  data VoteInt
  enc []byte
}
votingData := voting.GenVotingData(voteInt)
```

投票データをマネージャー公開鍵で暗号化します。  
IPFSにaddしてkeyFileを用いて投票用IPNSにpublishします。  

## Counting Setup
<img alt="counting_setup" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/counting_setup.png"><br>

(マネージャー)  
VotingInfoを取得します。  
投票用IPNSアドレスマップから暗号化投票データを収集します。  
マネージャー秘密鍵で投票データを取得します。  
IPNSアドレスマップのキーを用いて投票データを値とするマップを生成します。

```Go
var votingDataMap map[string]VotingData
for k, v := range votingIPNSAddrs{
  encVotingData := Get(v)
  mvd := rsa.Decrypt(encVotingData, manPriKey)
  votingData := voting.UnmarshalVotingData(mvd)
  votingDataMap[k] = votingData
}
```

投票データマップをIPFSにaddし、そのCIDを公表します。  
   
## Counting
<img alt="counting" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/counting.png"><br>

(ユーザー)  
投票データマップを取得します。  
VotingInfoから検証鍵マップを取得します。  
投票データマップと検証鍵マップから任意の投票データを検証します。  
投票データを集計して投票結果を取得します。  

# Voting Type
以下の投票方式に対応  
* [単記投票](https://ja.wikipedia.org/wiki/%E5%8D%98%E8%A8%98%E7%A7%BB%E8%AD%B2%E5%BC%8F%E6%8A%95%E7%A5%A8)  
* [連記投票](https://ja.wikipedia.org/wiki/%E9%80%A3%E8%A8%98%E6%8A%95%E7%A5%A8)  
* [認定投票](https://ja.wikipedia.org/wiki/%E8%AA%8D%E5%AE%9A%E6%8A%95%E7%A5%A8)  
* [範囲投票](https://ja.wikipedia.org/wiki/%E6%8E%A1%E7%82%B9%E6%8A%95%E7%A5%A8)  
* [累積投票](https://ja.wikipedia.org/wiki/%E7%B4%AF%E7%A9%8D%E6%8A%95%E7%A5%A8)  
* [選好投票](https://ja.wikipedia.org/wiki/%E9%81%B8%E5%A5%BD%E6%8A%95%E7%A5%A8)  


# TODO
* 登録処理  
* 集計処理  
* GUI  
* IPNSを使用する処理の高速化

# Support
どの組織にも属さずフリーランスで開発しています。  
投票という制度をより身近で容易かつ公正なものにするため、無料でリリースする予定です。  
開発の継続のため、皆様のご支援が必要です。  
どうかよろしくお願いいたします。

ETH Address:  0x81f5877EFC75906230849205ce11387C119bd9d8
