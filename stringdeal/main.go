package main

import (
	"bytes"
	"fmt"
)

// intsToString is like fmt.Sprint(values) but adds commas.
func intsToString(values []int) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteByte(']')
	return buf.String()
}

func comma(s string) string {
	var buf bytes.Buffer
	for i, j := len(s)-1, 0; i >= 0; i-- {
		buf.WriteByte(s[i])
		j++
		if j%3 == 0 && j != 0 {
			buf.WriteByte(',')
		}
	}
	// reverse the string
	content := buf.Bytes()
	for i, j := int(len(s)/2), 0; j <= i; j++ {
		content[j], content[len(content)-j-1] = content[len(content)-j-1], content[j]
	}
	buf.Reset()
	buf.Write(content)
	return buf.String()
}

func main() {
	fmt.Println(intsToString([]int{1, 2, 3})) // "[1, 2, 3]"
	fmt.Println(comma("1234567890"))          // "1,234,567,890"
}
