package content_type

import (
  "encoding/json"
  "strings"
)

func Detect(result string) string {
  trimmed := strings.TrimSpace(result)
  if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
    var js json.RawMessage
    if json.Unmarshal([]byte(trimmed), &js) == nil {
      return "json"
    }
  }
  if strings.HasPrefix(trimmed, "<!DOCTYPE html") || strings.HasPrefix(trimmed, "<html") {
    return "html"
  }
  if strings.HasPrefix(trimmed, "<?xml") || (strings.HasPrefix(trimmed, "<") && strings.HasSuffix(trimmed, ">")) {
    return "xml"
  }
  return "text"
}