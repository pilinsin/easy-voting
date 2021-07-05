# EasyVoting
This is an Online Voting App based on IPFS.<br>
Blockchain is not used.<br>

# Requirement
[go-ipfs](https://github.com/ipfs/go-ipfs)  
gui

# Usage
<img alt="system_process" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/system_process.png"><br>
## Online Voter Registration
<img alt="system_process" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/registration.png"><br>
(User)<br>
Generate a RSA key pair.<br>
The private key is stored locally.<br>
The public key is added to IPFS, and then publish its CID with arbitrary key.<br>
Register an email address and the IPNS address for a Manager's server.<br>

## Voting Setup
<img alt="system_process" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/voting_setup.png"><br>

(Manager)<br>

Obtain the list of email addresses and registration IPNS addresses from the server.<br> 

For each user, process the following.<br>
1. Obtain the user public key from the registration IPNS address.<br>
2. Generate an userID.<br>
3. Generate a KeyFile.<br>

```Go
userID := util.GenUniqueID(30,6)
KeyFile := ipfs.KeyFileGenerate()
```

4. Encode the userID and the KeyFile with the user public key.<br>
5. Send an email include the encoded userID and KeyFile to the user.<br>
<br>

Calculate voting IPNS addresses corresponding to the KeyFiles.<br>
Generate a manager's RSA key pair.<br>
Generate a votingID.<br>

```Go
votingID := util.GenUniqueID(30,30)
```
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
  votingIPNSAddrs []string  
}  
type Candidate struct{  
  url      string  
  group    string  
  groupURL string  
}  
```

## Voting
<img alt="system_process" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/voting.png"><br>

(User)<br>
Obtain VotingInfo.<br>
Obtain the encoded userID and KeyFile from the email.<br>
Decode the userID and KeyFile with the user private key.<br>

Calculate a voting IPNS address corresponding to the KeyFile.<br>
Verify the address with the voting IPNS addresses.<br>

Reflect the votingType on a voting form.<br>
Generate a voting data.<br>

```Go
type VoteInt map[string]int
votingData := map[string]int{userID: vote}
//or  
//type VoteBool map[string]bool
//votingData := map[string]bool{userID: vote}  
```

Encode the voting data with the manager's public key.<br> 
Add the encoded voting data to IPFS and publish to the voting IPNS.<br>

## Counting Setup
<img alt="system_process" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/counting_setup.png"><br>

(Manager)<br>
Obtain VotingInfo.<br>
Collect the encoded voting data from the voting IPNS addresses.<br>
Decode them with the manager's private key.<br>
Concatenate them.<br>
Add the whole voting data to IPFS and announce its CID.<br>
   
## Counting
<img alt="system_process" src="https://github.com/m-vlanbdg2ln52gla/EasyVoting/blob/main/images/counting.png"><br>

(User)<br>
Obtain the whole voting data.<br>
Check own voting data.<br>
Tally them.<br>

# Voting Type
This supports the following types:  
[Single](https://en.wikipedia.org/wiki/Single_transferable_vote)  
[Block](https://en.wikipedia.org/wiki/Multiple_non-transferable_vote)  
[Approval](https://en.wikipedia.org/wiki/Approval_voting)  
[Range](https://en.wikipedia.org/wiki/Score_voting)  
[Cumulative](https://en.wikipedia.org/wiki/Cumulative_voting)  
[Preference](https://en.wikipedia.org/wiki/Ranked_voting)  


# TODO
Registration process  
Counting process  
GUI part  

# Support


ETH Address: 0x81f5877EFC75906230849205ce11387C119bd9d8
