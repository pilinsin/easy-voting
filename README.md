# EasyVoting
This is an Online Voting App based on [IPFS](https://ipfs.io/).<br>
Blockchain is not used.<br>

## Feature
* Anonymity  
* Revote  
* Confirmation of every voting data by anyone  
* Counting by anyone


# Requirement
[go-ipfs](https://github.com/ipfs/go-ipfs)  
gui

# Process Flow
<img alt="system_process" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/system_process.png"><br>
## Online Voter Registration
<img alt="registration" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/registration.png"><br>
(User)<br>
Generate a RSA key pair.<br>
The private key is stored locally.<br>
The public key is added to IPFS, and then publish its CID with arbitrary key.<br>
Register an email address and the IPNS address for a Manager's server.<br>

## Voting Setup
<img alt="voting_setup" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/voting_setup.png"><br>

(Manager)<br>
Generate a votingID.<br>
Make voting IPNS address map and verification key map.<br>
```Go
votingID := util.GenUniqueID(30,30)

var votingIPNSAddrs map[string]string
var verfKeys        map[string]rsa.PublicKey
```

Obtain the list of email addresses and registration IPNS addresses from the server.<br> 

For each user, process the following.<br>
1. Obtain the user public key from the registration IPNS address.<br>
2. Generate an userID.<br>
3. Generate a keyFile.<br>

```Go
userID := util.GenUniqueID(30,6)
keyFile := ipfs.GenKeyFile()
```
4. Generate RSA signKey and verfKey.<br>
5. Encode the userID, keyFile and signKey with the user public key.<br>
6. Send an email include the encoded userID, keyFile and signKey to the user.<br>
7. Calculate voting IPNS address corresponding to the keyFile.<br>
8. Calculate a hash.<br>

```Go
hash := util.Hash(userID, votingID)
```

9. Append the IPNS address to the voting IPNS address map
10. Append the verification key to the key map.<br>

```Go
votingIPNSAddrs[hash] = addr
verfKeys[hash] = verfKey
```

<br>
Generate a manager's RSA key pair.<br>
The private key is stored locally.<br>
<br>
Add VotingInfo to IPFS and announce its CID.<br>

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

(User)<br>
Obtain VotingInfo.<br>
Obtain the encoded userID, keyFile and signKey from the email.<br>
Decode the userID, keyFile and signKey with the user private key.<br>

Calculate a voting IPNS address corresponding to the keyFile.<br>
Verify the address with the voting IPNS addresses.<br>

Reflect the votingType on a voting form.<br>
Generate a voting data with a signature.<br>

```Go
type VoteInt map[string]int
type VotingData struct{
  data VoteInt
  enc []byte
}
votingData := voting.GenVotingData(voteInt)
```

Encode the voting data with the manager's public key.<br> 
Add the encoded voting data to IPFS and publish to the voting IPNS.<br>

## Counting Setup
<img alt="counting_setup" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/counting_setup.png"><br>

(Manager)<br>
Obtain VotingInfo.<br>
Collect the encoded voting data from the voting IPNS address map.<br>
Decode them with the manager's private key.<br>
Concatenate them.<br>
```Go
var votingDataMap map[string]VotingData
for k, v := range votingIPNSAddrs{
  encVotingData := Get(v)
  mvd := rsa.Decrypt(encVotingData, manPriKey)
  votingData := voting.UnmarshalVotingData(mvd)
  votingDataMap[k] = votingData
}
```

Add the votingDataMap to IPFS and announce its CID.<br>
   
## Counting
<img alt="counting" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/counting.png"><br>

(User)<br>
Obtain the verfKeys from VotingInfo.<br>
Obtain the votingDataMap.<br>
Using the verfKeys and votingDataMap, check arbitrary voting data.<br>
Tally them.<br>


# Usage
## Registration
(User)  
Input an arbitrary IPNS key, and then get a registration IPNS Address.<br>
Register the address and an email address for a Manager's server.<br>

## Voting Setup
(Manager)  
Obtain the list of email addresses and registration IPNS addresses from the server.<br> 
Input the following informations from a form of the app:
* Beginning time for voting
* End time for voting
* Voting type
* Candidate informations 

A CID of a VotingInfo is output, so announce it.<br>

## Voting
(User)  
Input the CID of VotingInfo.<br>
Receive an email include encoded userID, keyFile and signKey.<br>
Input them.<br>
Vote from a form.<br>

## Counting Setup
(Manager)  
Input the CID of VotingInfo.<br>
A CID of a VotingDataMap is output, so announce it.<br>

## Counting
(User)  
Input the CID of VotingInfo.<br>
Input the CID of VotingDataMap.<br>
If you want to verify your own vote, input your encoded userID.<br>
The verification result is output.<br>
The voting result is output.<br>
 

# Voting Type
This supports the following types:  
* [Single](https://en.wikipedia.org/wiki/Single_transferable_vote)  
* [Block](https://en.wikipedia.org/wiki/Multiple_non-transferable_vote)  
* [Approval](https://en.wikipedia.org/wiki/Approval_voting)  
* [Range](https://en.wikipedia.org/wiki/Score_voting)  
* [Cumulative](https://en.wikipedia.org/wiki/Cumulative_voting)  
* [Preference](https://en.wikipedia.org/wiki/Ranked_voting)  


# TODO
* Registration process  
* Counting process  
* GUI   
* Acceleration of IPNS processes 

# Support
I develop it in freelance.<br>
I am going to release it free to make voting more common, and easy and fair.<br>
Your support lets development continue.<br>

ETH Address: 0x81f5877EFC75906230849205ce11387C119bd9d8
