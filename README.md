# EasyVoting
This is an Online Voting App based on IPFS.<br>
Blockchain is not used.<br>


# Usage
There are an Election Manager (Manager) and Voters (User).<br>
## Online Voter Registration
(User)<br>
Generate a RSA key pair.<br>
The private key is stored locally.<br>
The public key is added to IPFS, and then publish its CID with arbitrary key.<br>
Register an email address and the IPNS address for a Manager's server.<br>

## Voting Setup
(Manager)<br>
Generate a votingID.<br>
Obtain the list of email addresses and registration IPNS addresses from the server.<br> 
Obtain the user public key from the registration IPNS addresses.<br>
<br>
For each user:<br>
Generate an userID.<br>
Generate a KeyPair.<br>

```
KeyPair := util.Hash(votingID + userID)
``` 

Encode the userID and the KeyPair with the user public key.<br>
Send an email include the encoded userID and KeyPair to the user.<br>
<br>

Calculate voting IPNS addresses corresponding to the KeyPairs.<br>
Generate a manager's RSA key pair.<br>
<br>
Add VotingInfo to IPFS and announce its CID.<br>

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

## Voting
(User)<br>
Obtain VotingInfo.<br>
Obtain the encoded userID and KeyPair from the email.<br>
Decode the userID and KeyPair with the user private key.<br>

Calculate a voting IPNS address corresponding to the KeyPair.<br>
Verify the address with the voting IPNS addresses.<br>

Reflect the votingType on a voting form.<br>
Generate a voting data.<br>

```
votingData := map[string]int{userID: num}
//or  
//votingData := map[string]bool{userID: flag}  

```

Encode the voting data with the manager's public key.<br> 
Add the encoded voting data to IPFS and publish to the voting IPNS.<br>

## Counting Setup
(Manager)<br>
Obtain VotingInfo.<br>
Collect the encoded voting data from the voting IPNS addresses.<br>
Decode them with the manager's private key.<br>
Add the whole voting data to IPFS and announce its CID.<br>
   
## Counting
(User)<br>
Obtain the whole voting data.<br>
Check own voting data.<br>
Tally them.<br>

# Voting Type
•Single  
•Block  
•Approval  
•Range  
•Cumulative  
•Preference  


# TODO


