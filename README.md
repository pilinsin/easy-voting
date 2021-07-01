# EasyVoting
IPFSとGUIライブラリを使用したオンライン投票アプリです。  


# Usage
管理者(マネージャー)と投票者(ユーザー)が存在します。
## Registration Setup
<マネージャー>  
まずオンライン投票を利用するユーザーの登録が必要です。  
アプリ上で入力した情報をサーバーに送信して登録するので、  
使用するサーバーのアドレスを入力しIPFSにadd、そのpathを公表します。  

## Online Voter Registration
<ユーザー>  
RSA鍵生成を行い、秘密鍵はローカルに保存します。  
公開鍵をIPFSにaddし、任意のIPNSkeyでpublishします。  
本人確認用に個人情報を入力し、ハッシュ化します。
個人情報ハッシュとメールアドレスとpublishしたIPNSのアドレスをサーバーに送信して登録します。  

## Voting Setup
<マネージャー>  
投票idを生成します。  
サーバーからメールアドレスと登録用IPNSアドレスのリストを取得します。  

登録用IPNSアドレスからユーザー公開鍵を取得します。  
ユーザーidを生成します。  
hash(投票id+ユーザーid)を投票用IPNSのキー名として、KeyPairを生成します。    
ユーザーidとKeyPairをユーザー公開鍵で暗号化してメールします。  

全ユーザーのKeyPairに対応する投票用IPNSアドレスをリスト化します。  
マネージャー公開鍵&秘密鍵を生成します。  

```
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

としてVotingInfoをIPFSにaddし、そのpathを公表します。  

## Voting
<ユーザー>  
VotingInfoのpathを入力してデータを取得します。  
メールから暗号化ユーザーidと暗号化KeyPairを取得します。  
ローカルに保存しておいたユーザー秘密鍵でユーザーidとKeyPairを取得します。  
KeyPairを入力し対応する投票用IPNSのアドレスを求め、投票用IPNSアドレスリストと比較してログイン認証を行います。  
投票方式が反映された投票フォームから投票内容を生成します。  
ユーザーidを入力してユーザーidをキー、投票内容を値とする投票データを生成します。  
投票データをマネージャー公開鍵で暗号化し、IPFSにaddしてKeyPairを用いて投票用IPNSにpublishします。  

## Counting Setup
<マネージャー>  
投票終了時刻経過後に処理を行います。  
VotingInfoを取得します。  
投票用IPNSアドレスリストから暗号化投票データを収集します。  
マネージャー秘密鍵で投票データを取得します。  
全ユーザーの投票データを纏めてリスト化します。  
投票データリストをIPFSにaddし、そのpathを公表します。  
   
## Counting
<ユーザー>  
投票データリストを取得します。  
自身のユーザーidから投票内容を確認します。  
投票データリストを集計して投票結果を取得します。  

# Voting Type
・単記投票  
・連記投票  
・認定投票  
・範囲投票  
・累積投票  
・選好投票  


# TODO


