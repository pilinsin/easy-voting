# EasyVoting
This is an Online Voting App based on IPFS.  
Blockchain is not used.  


# Usage
There are an Election Manager (Manager) and Voters (User).  
## Online Voter Registration
<User>  
Generate a RSA key pair.  
The private key is stored locally.  
The public key is added to IPFS, and then publish its CID with arbitrary key.  
Register an email address and the IPNS address for a Manager's server.  

## Voting Setup
<Manager>  
Generate a votingID.  
Obtain the list of email addresses and registration IPNS addresses from the server.  
Obtain the user public key from the registration IPNS addresses.  

For each user:  
Generate an userID.
Generate a KeyPair.

```
KeyPair := util.Hash(votingID + userID)
``` 

Encode the userID and the KeyPair with the user public key.  
Send an email include the encoded userID and KeyPair to the user.  
  

Calculate voting IPNS addresses corresponding to the KeyPairs.  
Generate a manager's RSA key pair.  

Add VotingInfo to IPFS and announce its CID.  

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
<User> 
Obtain VotingInfo.  
Obtain the encoded userID and KeyPair from the email.   
Decode the userID and KeyPair with the user private key.  

Calculate a voting IPNS address corresponding to the KeyPair.  
Verify the address with the voting IPNS addresses.    

Reflect the votingType on a voting form.  
Generate a voting data.

```
votingData := map[string]Vote{userID: vote}
```

Encode the voting data with the manager's public key.    
Add the encoded voting data to IPFS and publish to the voting IPNS.  

## Counting Setup
<Manager>  
Obtain VotingInfo.  
Collect the encoded voting data from the voting IPNS addresses.  
Decode them with the manager's private key.  
Add the whole voting data to IPFS and announce its CID.   
   
## Counting
<User>  
Obtain the whole voting data.  
Check own voting data.  
Tally them.  

# Voting Type
•Single  
•Block  
•Approval  
•Range  
•Cumulative  
•Preference  


# TODO


