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
Get a mCfgCid (registration Manager Config CID).  
Input the mCfgCid and transition to the registration manager page.  
Get the rCfgCid (Registration Config CID) from the page and open it.  
Turn the registrate switch on, and then wait.    

(User)  
Input the rCfgCid and transition to the registration page.  
Input a userData and registrate.  

```Go
var userData []string
```

Copy and keep a userIdentity output.  
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
Get a mCfgCid (voting Manager Config CID).  
Input the mCfgCid and transition to the voting manager page.  
Get the vCfgCid (Voting Config CID) from the page and open it.  
When input a userData for a user, the user can be verified.    
After the voting time finished, generate a resultMap.    

(User)  
Input the vCfgCid and transition to the voting page.  
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
