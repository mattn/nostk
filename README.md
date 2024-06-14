nostk
========
Implementing a CLI client to use [Nostr Protocol](https://github.com/nostr-protocol/nostr).

### Develop Environment
* Ubuntu 23.04 and later
* Go Language 1.22.4 and later

### Features
* Initializing the nostk environment
* Generating a key pair
* Edit your contact list
* Edit custom emoji list
* Edit relay list
* Publish relay list
* Edit profile
* Publish profile
* Display home timeline ([kind 1](https://github.com/nostr-protocol/nips/blob/master/01.md#kinds))
* Display your's note ([kind 1](https://github.com/nostr-protocol/nips/blob/master/01.md#kinds))

### ToDo
* Support content warning (WIP)
* Mention to any user
* Publishing a message with message citations
* Log viewer
* any more

### Requirements
* [nbd-wtf / go-nostr](https://github.com/nbd-wtf/go-nostr)
* Some kind of text editor
* Setting $EDITOR environment variable

### Setup
#### Install tools
1. Install [git](https://www.git-scm.com/)
2. Install [golang](https://go.dev/)

#### Install nostk:
##### Windows
```command.com
SETX EDITOR=<Text editor's full path name>
go install github.com/mitsugu/nostk@HEAD
```

##### Ubuntu and maybe other distribution
For bash
```bash
echo 'export EDITOR=vim' >> ~/.bashrc
go install github.com/mitsugu/nostk@HEAD
```

#### Placement of config.json
1. Download [config.json](https://raw.githubusercontent.com/mitsugu/nostk/main/config.json)
2. Move config.json to "$HOME/.nostk" directory
3. Adjust defaultReadNo, multiplierReadRelayWaitTime, and defaultContentWarning in config.json to your liking.

#### Setting nostk:
1. nostk init (must)
2. nostk genkey (must)
3. nostk editRelays (must)
4. nostk editContacts (must)
5. nostk editProfile (should \*)
6. nostk pubProfile (should \*)
7. nostk editEmoji (Optional)
8. nostk pubRelays (Optional)

\* Unless there is a special reason, it is recommended to use a web app such as [nostter](https://nostter.app/home) instead of nostk.

### Usage
#### Display help documanets
``` bash
nostk help

nostk -h

nostk --help

nostk
```

#### Init nostk
``` bash
nostk init
```

#### Ganerate Key Pair
``` bash
nostk genkey
```

#### Edit contact list
``` bash
nostk editContacts
```

#### Edit custom emoji
``` bash
nostk editEmoji
```

#### Edit relay list
``` bash
nostk editRelays
```

#### Publish relay list
``` bash
nostk pubRelays
```

#### Edit profile
``` bash
nostk editProfile
```

#### Publish profile
``` bash
nostk pubProfile
```

#### Publish message
``` bash
nostk pubMessage <text message>

nostk pubMessage < (ps)

(ps) | nostk pubMessage
```

#### Display home timeline (kind 1)
``` bash
nostk catHome [number]
```

#### Display your's note (kind 1)
``` bash
nostk catSelf [number]
```

#### Display specified event
``` bash
nostk catEvent <Hex event id>
```

#### Remove specified event (Test implementation)
``` bash
nostk removeEvent <Hex event id> [reason]
```

#### Custom emoji reaction (Test implementation)
``` bash
nostk emojiReaction <Hex event id> <Hex pubkey> <short code of custom emoji>
```

### About content warning note
  The catHome subcommand does not directly display notes with content warnings.  

  The corresponding note will be printed to indicate that it is a content warning note, the reason will be displayed if a reason is set, and the event ID of the note will also be displayed.  
To display a content warning note, run the catEvent subcommand by specifying the note's Event ID in hex.  

