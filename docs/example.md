# Bitcoin

Bitcoin (abbreviation: BTC or XBT; sign: ₿) is a protocol which implements a highly available, public, permanent, and decentralized ledger. In order to add to the ledger, a user must prove they control an entry in the ledger. The protocol specifies that the entry indicates an amount of a token, bitcoin with a minuscule b. The user can update the ledger, assigning some of their bitcoin to another entry in the ledger. Because the token has characteristics of money, it can be thought of as a digital currency.

## Units and divisibility

The unit of account of the bitcoin system is the bitcoin. Currency codes for representing bitcoin are BTC and XBT. Its Unicode character is ₿. One bitcoin is divisible to eight decimal places. Units for smaller amounts of bitcoin are the millibitcoin (mBTC), equal to 1⁄1000 bitcoin, and the satoshi (sat), which is the smallest possible division, and named in homage to bitcoin's creator, representing 1⁄100000000 (one hundred millionth) bitcoin. 100,000 satoshis are one mBTC. 

## Blockchain

The bitcoin blockchain is a public ledger that records bitcoin transactions. It is implemented as a chain of blocks, each block containing a cryptographic hash of the previous block up to the genesis block in the chain. A network of communicating nodes running bitcoin software maintains the blockchain. Transactions of the form payer X sends Y bitcoins to payee Z are broadcast to this network using readily available software applications. 

Network nodes can validate transactions, add them to their copy of the ledger, and then broadcast these ledger additions to other nodes. To achieve independent verification of the chain of ownership, each network node stores its own copy of the blockchain. At varying intervals of time averaging to every 10 minutes, a new group of accepted transactions, called a block, is created, added to the blockchain, and quickly published to all nodes, without requiring central oversight. This allows bitcoin software to determine when a particular bitcoin was spent, which is needed to prevent double-spending. A conventional ledger records the transfers of actual bills or promissory notes that exist apart from it, but as a digital ledger, bitcoins only exist by virtue of the blockchain; they are represented by the unspent outputs of transactions.

Individual blocks, public addresses, and transactions within blocks can be examined using a blockchain explorer.

## Supply

Every 10 minutes, the successful miner finding the new block is allowed by the rest of the network to collect for themselves all transaction fees from transactions they included in the block, as well as a predetermined reward of newly created bitcoins. As of 11 May 2020, this reward was ₿6.25 in newly created bitcoins per block. To claim this reward, a special transaction called a coinbase is included in the block, with the miner as the payee. All bitcoins in existence have been created through this type of transaction. The bitcoin protocol specifies that the reward for adding a block will be reduced by half every 210,000 blocks (approximately every four years), until ₿21 million are generated. Assuming the protocol is not changed and the 10 minute average block creation time remains constant, the last new bitcoin would be generated around the year 2140. After that, a successful miner would be rewarded by transaction fees only.

## Wallets

A wallet stores the information necessary to transact bitcoins. While wallets are often described as a place to hold or store bitcoins, due to the nature of the system, bitcoins are inseparable from the blockchain transaction ledger. A wallet is more correctly defined as something that "stores the digital credentials for your bitcoin holdings" and allows one to access (and spend) them. glossary  Bitcoin uses public-key cryptography, in which two cryptographic keys, one public and one private, are generated. At its most basic, a wallet is a collection of these keys. 

### Software wallets

The first wallet program, simply named Bitcoin, and sometimes referred to as the Satoshi client, was released in 2009 by Satoshi Nakamoto as open-source software. In version 0.5 the client moved from the wxWidgets user interface toolkit to Qt, and the whole bundle was referred to as Bitcoin-Qt. After the release of version 0.9, the software bundle was renamed Bitcoin Core to distinguish itself from the underlying network. Bitcoin Core is, perhaps, the best known implementation or client. Alternative clients (forks of Bitcoin Core) exist, such as Bitcoin XT, Bitcoin Unlimited, and Parity Bitcoin.

### Cold storage

Wallet software is targeted by hackers because of the lucrative potential for stealing bitcoins. A technique called "cold storage" keeps private keys out of reach of hackers; this is accomplished by keeping private keys offline at all times  by generating them on a device that is not connected to the internet. The credentials necessary to spend bitcoins can be stored offline in a number of different ways, from specialized hardware wallets to simple paper printouts of the private key.

### Hardware wallets

A hardware wallet is a computer peripheral that signs transactions as requested by the user. These devices store private keys and carry out signing and encryption internally, and do not share any sensitive information with the host computer except already signed (and thus unalterable) transactions. Because hardware wallets never expose their private keys, even computers that may be compromised by malware do not have a vector to access or steal them.

### Paper wallets

A paper wallet is created with a keypair generated on a computer with no internet connection; the private key is written or printed onto the paper and then erased from the computer. The paper wallet can then be stored in a safe physical location for later retrieval.