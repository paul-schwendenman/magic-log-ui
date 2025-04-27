package jqfilter

import (
	"fmt"
	"log"

	"github.com/itchyny/gojq"
)

var compiled *gojq.Code

func Init(query string) {
	if query == "" {
		return
	}

	parsed, err := gojq.Parse(query)
	if err != nil {
		log.Fatalf("❌ Failed to parse jq query: %v", err)
	}
	compiled, err = gojq.Compile(parsed)
	if err != nil {
		log.Fatalf("❌ Failed to compile jq query: %v", err)
	}
	log.Printf("🔀 JQ filter enabled: %s", query)
}

func Apply(entry map[string]string) map[string]string {
	if compiled == nil {
		return entry
	}

	generic := make(map[string]interface{}, len(entry))
	for k, v := range entry {
		generic[k] = v
	}

	iter := compiled.Run(generic)
	v, ok := iter.Next()
	if !ok {
		log.Printf("❌ jq query returned no result")
		return entry
	}
	if err, ok := v.(error); ok {
		log.Printf("❌ jq query error: %v", err)
		return entry
	}

	mapped, ok := v.(map[string]interface{})
	if !ok {
		mapped = map[string]interface{}{"value": v}
	}

	newEntry := make(map[string]string, len(mapped))
	for k, v := range mapped {
		newEntry[k] = fmt.Sprintf("%v", v)
	}

	return newEntry
}
