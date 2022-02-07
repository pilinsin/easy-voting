# EasyVoting
This is an Online Voting App based on [IPFS](https://ipfs.io/) and [Fyne](https://fyne.io/).<br>
Blockchain is not used.<br>

## Feature
* Anonymous voting  
* Revote  
* Confirmation of every voting data by anyone  
* Counting by anyone


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
(Registration Manager)  
In the registration setup page, input some informations.  
Get a rCfgCid (Registration Config CID) and a ManIdentity.  
The ManIdentity is stored locally with a key input in the setup page.   
Input the rCfgCid and key, and then transition to the registration manager page.  
Turn the registrate switch on, and then wait.    

```Go
type ManIdentity struct{
  rPriKey *ecies.PriKey
  hnmKeyFile *ipfs.KeyFile
}
```

(User)  
Input the rCfgCid and transition to the registration page.  
Input a userData and registrate.  

```Go
var userData []string
```

A UserIdentity is output and stored locally.  
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
(Manager)  
In the voting setup page, input some informations include the rCfgCid.  
Get a vCfgCid (Voting Config CID) and a ManIdentity.  
The ManIdentity is stored locally with a key input in the voting setup page.  
Input the vCfgCid and key, and then transition to the voting manager page.  
When input a userData for a user, the user can be verified.    
After the voting time finished, generate a resultMap.    
```Go
type ManIdentity struct{
  manPriKey *ecies.PriKey
  resMapKeyFile *ipfs.KeyFile
}
```

(User)  
Input the vCfgCid and the key, and then transition to the voting page.  
Input some informations for voting include the userIdentity.  
After the voting time finished, verify and count from the resultMap.   

# Voting Type
This supports the following types:  
* [Single](https://en.wikipedia.org/wiki/Single_transferable_vote)  
* [Block](https://en.wikipedia.org/wiki/Multiple_non-transferable_vote)  
* [Approval](https://en.wikipedia.org/wiki/Approval_voting)  
* [Range](https://en.wikipedia.org/wiki/Score_voting)  
* [Cumulative](https://en.wikipedia.org/wiki/Cumulative_voting)  
* [Preference](https://en.wikipedia.org/wiki/Ranked_voting)  


# TODO
* GUI design
* Bug fix  
* Improvement of the registration process


# Support
I develop it in freelance.<br>
I am going to release it free to make voting more common, and easy and fair.<br>
Your support lets development continue.<br>

Ethereum Address: 0x81f5877EFC75906230849205ce11387C119bd9d8
