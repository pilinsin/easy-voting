# EasyVoting
This is an Online Voting App based on [I2P](https://geti2p.net/en/) and [Fyne](https://fyne.io/).<br>
Blockchain is not used.<br>

## Feature
* Anonymous voting  
* Revote  
* Counting by anyone  


## Usage
### Setup I2P
Install and run [I2P](https://github.com/i2p/i2p.i2p).  
Check if [SAM](https://geti2p.net/en/docs/api/samv3) is enabled, and if not enable it.  
### Bootstrap
##### Manager  
If you generate your own I2P bootstrap on the setup page, its address is output.  
Enter the address or others, and then an address of a list of the addresses will be output.  
In this section, the address is "baddrs".
### Registration
##### Manager  
On the setup page, enter the following information:
- Title of the registration page
- User Data Set
- baddrs

The user data set is a csv file containing the data required for user registration.

| label 0 | label 1 | ... | label M |
| --- | --- | --- | --- |
| user 00 | user 01 | ... | user 0M |
| user 10 | user 11 | ... | user 1M |
| ... | ... | ... | ... |
| user N0 | user N1 | ... | user NM |

After entering the information, press the Submit button to output rCfgAddr and manIdentity, which are then entered into the load page form to move to the registration page.  
Publish rCfgAddr and wait until registration is complete.

##### User  
Enter rCfgAddr to go to the registration page.  
Enter the data necessary for registration.   
The userIdentity will be output, so copy and keep it.  

### Voting
##### Manager  
On the setup page, enter the following information:
- Voting Page Title
- Start Time
- End Time
- Time Zone (Location)
- rCfgAddr
- Number of Verificators(number needed to verify voting time)
- Candidate Information
  - Image
  - Name
  - Group
  - URL
- Voting Parameters
  - Minimum Number of Votes
  - Maximum Number of Votes
  - Total Votes
- Voting Type

After entering the information, press the Submit button to generate the vCfgAddr and manIdentity, which will be entered into the load page form to move to the voting page.  
Publish vCfgAddr and wait until voting close.  
After voting is completed, press the Vote button to release the decryption key for the voting data.

##### User  
Enter vCfgAddr and userIdentity to go to the voting page.  
If you do not enter userIdentity, you can only tally the voting results.  

Voting is done via the voting form.  
If the voting manager has published the decryption key for your voting data, you can check your voting data and tally the results of all votes.

## Voting Type
This supports the following types:  
* [Single](https://en.wikipedia.org/wiki/Single_transferable_vote)  
* [Block](https://en.wikipedia.org/wiki/Multiple_non-transferable_vote)  
* [Approval](https://en.wikipedia.org/wiki/Approval_voting)  
* [Range](https://en.wikipedia.org/wiki/Score_voting)  
* [Cumulative](https://en.wikipedia.org/wiki/Cumulative_voting)  
* [Preference](https://en.wikipedia.org/wiki/Ranked_voting)  

## Support
I develop it in freelance.<br>
I am going to release it free to make voting more common, easy and fair.<br>
Your support lets development continue.<br>

- BitCoin (BTC)
bc1qu0zl5z4zgvx2ar3zgdgmt3thl3fnt0ujvzdx9e
- Ethereum (ETH)
0x81f5877EFC75906230849205ce11387C119bd9d8
- Tron (TRX)
TCc7D7thmW4egbiUEk2uH3Y21shfbjVNvn
- Monero (XMR)
49y6hymbjLqf1LRrGoARqGNxD95UeHtpGbfYmutrLZaWhfFwefPHkDUiKkab3aCNBv36xAUu4VQus1V1g8hhYWrhLemRjPt
- Zcash (ZEC)
t1diehqpgftGp9dvEMKcAoCUxZnGgodcU96
- Basic Attention Token (BAT)
0xe83D64a10256aE37d3039344fE49ec9D1d75dd5c
- FileCoin (FIL)
f1mhulmnu4apv3thlsnmjw3nigl5hzcgozfabpsyi
