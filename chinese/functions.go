package chinese

import (
	"log"
	"regexp"
	"strings"
	"time"
	"fmt"
	"math/rand"
	"hash/crc32"

	"github.com/huichen/sego"
	"gopkg.in/redis.v3"
)

func (this *Chinese) Learn(text string, id int64) {
	text = strings.Trim(filter(text), " \n")

	if this.debug {
		log.Printf("filtered: %s\n", text)
	}

	segments := this.seg.Segment([]byte(text))

	if len(segments) == 0 {
		return
	}

	if this.debug {
		logTokens(segments)
	}

	tags := make([]string, len(segments))
	words := make([]string, len(segments))
	unrecognized := 0

	for i, s := range segments {
		tag := s.Token().Pos()
		word := s.Token().Text()

		tag = customTag(word, tag)

		tags[i] = tag
		words[i] = word

		if (tag == "x") || (tag == "eng") || isCustomTag(tag) {
			unrecognized += 1
		}
	}

	if unrecognized >= int(float32(len(words)) * 0.6) {
		// To many unrecognized words!
		if this.debug {
			log.Println("Too many unrecognized tags!")
		}

		return
	}

	for i, word := range words {
		tag := tags[i]

		if isCustomTag(tag) {
			continue
		}

		// Output
		log.Printf("{%d}: %s -> %s\n", i, word, tag)
		addMember(this.redis, fmt.Sprintf("chn%dword%s", id, tag), word)

		// Store the coexisting status
		// A word "w" of tag "tag[j]" coexists with "word"
		for j, w := range words {
			addMember(this.redis, fmt.Sprintf("chn%d%scoexist%s", id, tags[j], word), w)
		}
	}

	model := strings.Join(tags, " ")
	log.Printf("model: %s\n", model)
	addMember(this.redis, fmt.Sprintf("chn%dmodels", id), model)
}

func (this *Chinese) Speak(id int64) string {
	// Random seed
	rand.Seed(time.Now().Unix())

	model := randMember(this.redis, fmt.Sprintf("chn%dmodels", id))

	if model == "" {
		return ""
	}

	if this.debug {
		log.Printf("model = %s", model)
	}

	mo := strings.Split(model, " ")

	tagType := -1
	sentence := ""
	for _, m := range mo {
		word := ""
		if isCustomTag(m) {
			word = customUntag(m, &tagType)
		} else {
			word = randMember(this.redis, fmt.Sprintf("chn%dword%s", id, m))
		}

		if word == "" {
			continue
		}

		sentence += word
	}

	return sentence
}

func (this *Chinese) Answer(text string, id int64) string {
	text = strings.Trim(filter(text), " \n")

	// Set text as seed of random
	rand.Seed(int64(crc32.ChecksumIEEE([]byte(text))))

	model := randMember(this.redis, fmt.Sprintf("chn%dmodels", id))

	if model == "" {
		return ""
	}

	segments := this.seg.Segment([]byte(text))

	sentence := ""
	tags := strings.Split(model, " ")
	tagType := -1
	for _, t := range tags {
		word := ""
		if isCustomTag(t) {
			word = customUntag(t, &tagType)
		} else {
			for _, s := range segments {
				word = randMember(this.redis, fmt.Sprintf("chn%d%scoexist%s", id, t, s.Token().Text()))

				if word != "" {

					if this.debug {
						log.Printf("Got word %s for %s", word, t)
					}

					break
				}
			}

			if word == "" {
				if this.debug {
					log.Printf("Falling back for tag %s", t)
				}

				word = randMember(this.redis, fmt.Sprintf("chn%dword%s", id, t))
			}
		}

		if word != "" {
			sentence += word
		}
	}

	return sentence
}

func filter(text string) string {
	text = filterReg(text, `^([[(<].*? ?[\])>] )+`)
	text = filterReg(text, `([^<]*>|[^<>]*<\/)(([a-z][0-9a-z]*:)\/\/[a-z0-9&#=.\/\-?_]+)`)
	text = filterReg(text, `^(\S+, ?)*\S+: `)
	text = filterReg(text, `((\/|\@)[a-zA-Z0-9]*) `)
	return text
}

func filterReg(text string, reg string) string {
	re := regexp.MustCompile(reg)
	return re.ReplaceAllString(text, "")
}

// Scope start
var startTags []string = []string {
	"(", "[", "{", "（", "［", "【", "《", "『", "｢", "「", "‘", "“"}
// Scope end
var endTags []string = []string {
	")", "]", "}", "）", "］", "】", "》", "』", "｣", "」", "’", "”"}
// Balanced tags
var balTags []string = []string {
	"`", "'", "\""}
// Literials
var litTags []string = []string {
	".", ",", "?", "!", ";", "。", "，", "；", "？", "！", "…"}

func contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}

	return false
}

func customTag(word string, tag string) string {
	if tag == "x" {
		// Only process unknown tags
		if contains(startTags, word) {
			tag = "__my_start"
		} else if contains(endTags, word) {
			tag = "__my_end"
		} else if contains(balTags, word) {
			tag = "__my_bal"
		} else if contains(litTags, word) {
			tag = "__my_lit_" + word
		}
	}

	return tag
}

func customUntag(tag string, tagType *int) string {
	if tag == "__my_start" {
		(*tagType) = rand.Intn(len(startTags))
		return startTags[(*tagType)]
	} else if tag == "__my_end" {
		t := (*tagType)
		if t == -1 {
			t = rand.Intn(len(endTags))
		}
		(*tagType) = -1
		return endTags[t]
	} else if tag == "__my_bal" {
		t := rand.Intn(len(balTags))
		return balTags[t]
	} else if strings.HasPrefix(tag, "__my_lit") {
		return tag[9:]
	} else {
		// Should never reach here
		return ""
	}
}

func isCustomTag(tag string) bool {
	return strings.HasPrefix(tag, "__my_")
}

func weightedRandom(max int64) int64 {
	total := (1 + max) * max / 2

	if total <= 1 {
		return -1
	}

	r := rand.Int63n(total)
	var t int64 = 0
	var i int64 = 0

	for i = 0; i < max; i++ {
		t += i

		if t >= r {
			return i
		}
	}

	return -1
}

func addMember(c *redis.Client, setName string, member string) {
	exists, _ := c.Exists(setName).Result()

	score, err2 := c.ZScore(setName, member).Result()
	if err2 == redis.Nil {
		score = 1
	} else {
		score += 1
	}

	err := c.ZAdd(setName, redis.Z {
		Score: score,
		Member: member,
	}).Err()

	if err != nil {
		panic(err)
	}

	if !exists {
		// Add expiration
		log.Printf("Adding expiration to %s", setName)

		err = c.Expire(setName, 48 * time.Hour).Err()

		if err != nil {
			panic(err)
		}
	}
}

func randMember(c *redis.Client, setName string) string {
	max, err1 := c.ZCard(setName).Result()

	if err1 == redis.Nil {
		return ""
	} else if err1 != nil {
		panic(err1)
	}

	index := weightedRandom(max)

	m, err2 := c.ZRange(setName, index, index).Result()

	if (err2 == redis.Nil) || (len(m) == 0) {
		return ""
	} else if err2 != nil {
		panic(err2)
	}

	return m[0]
}

func logTokens(segments []sego.Segment) {
	for _, s := range segments {
		log.Printf("%s:%s", s.Token().Text(), s.Token().Pos())
	}
}
