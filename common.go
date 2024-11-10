package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"log"
	"regexp"
	"strings"
)

const (
	subCommandName=1
)

/*
nostr.Tags Data analysis related functions.

WHY WAS IT WRITTEN?
For checking nostr.Tags for undesired data, etc.
*/
// class definition
type Tags struct{}
/* Tags.hasPrefix {{{

WHY WAS IT WRITTEN?
The purpose is to check whether nostr.Tags contains a private key,
but it is also versatile.
*/
func (r Tags)hasPrefix(tgs nostr.Tags, pref string) bool {
	for _, tg := range tgs {
		for _, s := range tg {
			if ret := strings.HasPrefix(s, pref); ret == true {
				return ret
			}
		}
	}
	return false
}

// }}}

/* checkTags {{{
WHAT'S THIS?
Checks whether the tag is the target of processing of the specified kind.
*/
type ChkTblMap map[int][]string
func NewChkTblMap() ChkTblMap {
	return ChkTblMap{
		1:     {"content-warning", "client", "e", "emoji", "expiration", "p", "q", "r", "t"},
		6:     {"e", "p"},
		10000: {"e", "p", "t", "word"},
		10001: {"e"},
		30315: {"d", "emoji", "expiration", "r"},
	}
}
func (r ChkTblMap)sliceToMap(kind int) map[string]struct{} {
	m := make(map[string]struct{})
	for _, v := range r[kind]{
		m[v] = struct{}{}
	}
	return m
}
func (r ChkTblMap)contains(kind int, target string) bool {
	m := r.sliceToMap(kind)
	_, exists := m[target]
	return exists
}
func checkTags(kind int, tgs nostr.Tags) error {
	list := NewChkTblMap()
	for _, tg := range tgs {
		if result := list.contains(kind, tg[indexTagName]); result != true {
			log.Printf("kind : %v, tagName : %v\n", kind, tg[indexTagName])
			return errors.New("Inclusion of invalid tag in specified kind")
		}
	}
	return nil
}

// }}}

/* subcommand's kind getter {{{
*/
type SubCmdKindTbl map[string]int
func NewSubCmdKindTbl() SubCmdKindTbl {
	return SubCmdKindTbl{
		"pubMessage": 1,
		"pubMessageTo": 1,
		//"emojiReaction": 6,
	}
}

func (r SubCmdKindTbl) hasSubcommand(args []string) error {
	if subCommandName >= len(args) {
		return fmt.Errorf("index %d out of range in args", subCommandName)
	}

	cmd := args[subCommandName]
	if _, exists := r[cmd]; !exists {
		return fmt.Errorf("Not supported subcommand %s", cmd)
	}
	return nil
}
func (r SubCmdKindTbl) get(args []string) (int, error) {
	if err := r.hasSubcommand(args); err != nil {
		return 0, err
	}
	return r[args[subCommandName]], nil
}

// }}}

/* subcommand's json builder {{{
*/
type ConvArgsTagsTbl map[string]map[int][]string
func NewConvArgsTagsTbl() ConvArgsTagsTbl {
	return ConvArgsTagsTbl{
		"pubMessage":{
			3:{"content-warning",},
		},
		"pubMessageTo":{
			3:{"p"},
		},
		/*
		"emojiReaction":{
			6:{"e","p","k","emoji"},
		},
		*/
	}
}

func (r ConvArgsTagsTbl) hasSubcommand(args []string) error {
	if subCommandName >= len(args) {
		return fmt.Errorf("index %d out of range in args", subCommandName)
	}

	cmd := args[subCommandName]
	if _, exists := r[cmd]; !exists {
		return fmt.Errorf("Not supported subcommand %s", cmd)
	}
	return nil
}
type RawArg struct {
	Kind	int			`json:"kind"`
	Content	string		`json:"content"`
	Tags	nostr.Tags	`json:"tags"`
}
func buildJson(args []string) (string, error) {
	var ret	RawArg
	list := NewConvArgsTagsTbl()
	kindList := NewSubCmdKindTbl()
	if err := list.hasSubcommand(args); err != nil {
		return "", err
	}
	for i := range args {
		switch i {
			case 0:		// nostk
				continue
			case 1:		// subcommand name
				if tmpKind, err := kindList.get(args); err != nil {
					return "", err
				} else {
					ret.Kind = tmpKind
				}
			case 2:		// content
				ret.Content = args[i]
			default:	// tags
				if err := addTags(args, i, &ret.Tags); err != nil {
					return "", err
				}
		}
	}

	strJson, err := json.Marshal(ret)
	if err != nil {
		return "",err
	}
	return string(strJson), nil
}

// }}}

/* addTags {{{
*/
func addTags(args []string, index int, tgs *nostr.Tags) error {
	t := []string{}
	list := NewConvArgsTagsTbl()

	if err := list.hasSubcommand(args); err != nil {
		return err
	}
	t = append(t, list[args[subCommandName]][index][0])
	t = append(t, args[index])
	*tgs = append(*tgs, t)
	return nil
}

// }}}

/* is64HexString {{{ */
func is64HexString(s string) bool {
	if len(s) != 64 {
		return false
	}
	match, _ := regexp.MatchString("^[a-fA-F0-9]{64}$", s)
	return match
}

// }}}

/* isHexString {{{
*/
func isHexString(s string) bool {
	match, _ := regexp.MatchString("^[a-fA-F0-9]+$", s)
	return match
}

// }}}

/* getPrefixInString {{{
*/
func getPrefixInString(str string) (string, error) {
	pref, _, err := nip19.Decode(str)
	if err != nil {
		return "", err
	}
	return pref, nil
}

// }}}

/* excludeHashtagsParsign {{{
*/
func excludeHashtagsParsign(src string) (string, error) {
	const strexp = `(?:^|\s)([#﹟＃][^#﹟＃]\S*[#﹟＃]\S*)`
	re, err := regexp.Compile(strexp)
	if err != nil {
		return "", err
	}
	result := re.ReplaceAllString(src, "")
	return result, nil
}

// }}}

/* setHashTags {{{
*/
func setHashTags(buf string, tgs *nostr.Tags) {
	const strexp = `(?:^|\s)([#﹟＃][^#\s﹟＃]+[^\s|$])`

	re := regexp.MustCompile(strexp)
	matches := re.FindAllString(buf, -1)
	for i := range matches {
		t := ExTag{}
		t.addTagName("t")
		rtmp := regexp.MustCompile(`[\s﹟＃#]`)
		result := rtmp.ReplaceAllString(matches[i], "")
		t.addTagValue(result)
		*tgs = append(*tgs, t.getNostrTag())
	}
}

// }}}

/* setContentWarning {{{
*/
func setContentWarning(r string, tgs *nostr.Tags) {
	const CWTag = "content-warning"
	var t []string
	t = nil
	t = append(t, CWTag)
	t = append(t, r)
	*tgs = append(*tgs, t)
}

// }}}

/* setPerson {{{
*/
func setPerson(p string, tgs *nostr.Tags) {
	const PTag = "p"
	var t []string
	t = nil
	t = append(t, PTag)
	t = append(t, p)
	*tgs = append(*tgs, t)
}

// }}}

/* containsNsec1 {{{
*/
func containsNsec1(text string) bool {
	pattern := `nsec1[a-zA-Z0-9]{58}`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(text, -1)

	for _, match := range matches {
		alphanumericPart := match[5:]
		if !regexp.MustCompile(`nsec1`).MatchString(alphanumericPart) {
			return true
		}
	}

	return false
}

// }}}

/* containsHsec1 {{{
*/
/*
func containsHsec1(text string) bool {
	pattern := `[a-zA-Z0-9]{64}`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(text, -1)

	for _, match := range matches {
		if _, err := nip19.EncodePrivateKey(match); err == nil {
			return true
		}
	}
	return false
}
*/

// }}}
