package main

import (
	"fmt"
	"github.com/BaseMax/go-clipboard-plus/pkg/transform"
)

func main() {
	fmt.Println("=== go-clipboard-plus Transformation Demo ===")
	fmt.Println()

	// Test Base64
	fmt.Println("1. Base64 Encoding/Decoding:")
	input1 := "hello world"
	encoded, _ := transform.Base64Encode(input1)
	fmt.Printf("   Input: %s\n", input1)
	fmt.Printf("   Encoded: %s\n", encoded)
	decoded, _ := transform.Base64Decode(encoded)
	fmt.Printf("   Decoded: %s\n", decoded)
	fmt.Println()

	// Test case transformations
	fmt.Println("2. Case Transformations:")
	input2 := "hello world"
	upper, _ := transform.Upper(input2)
	lower, _ := transform.Lower("HELLO WORLD")
	title, _ := transform.Title("hello world")
	fmt.Printf("   Original: %s\n", input2)
	fmt.Printf("   Upper: %s\n", upper)
	fmt.Printf("   Lower: %s\n", lower)
	fmt.Printf("   Title: %s\n", title)
	fmt.Println()

	// Test JSON
	fmt.Println("3. JSON Formatting:")
	input3 := `{"name":"test","age":30,"city":"NYC","active":true}`
	formatted, _ := transform.FormatJSON(input3)
	fmt.Printf("   Original: %s\n", input3)
	fmt.Printf("   Formatted:\n%s\n", formatted)
	fmt.Println()

	// Test URL encoding
	fmt.Println("4. URL Encoding/Decoding:")
	input4 := "hello world & test=value"
	urlEncoded, _ := transform.URLEncode(input4)
	urlDecoded, _ := transform.URLDecode(urlEncoded)
	fmt.Printf("   Original: %s\n", input4)
	fmt.Printf("   Encoded: %s\n", urlEncoded)
	fmt.Printf("   Decoded: %s\n", urlDecoded)
	fmt.Println()

	// Test reverse
	fmt.Println("5. Text Reversal:")
	input5 := "hello world"
	reversed, _ := transform.Reverse(input5)
	fmt.Printf("   Original: %s\n", input5)
	fmt.Printf("   Reversed: %s\n", reversed)
	fmt.Println()

	// Test trim
	fmt.Println("6. Whitespace Trimming:")
	input6 := "   hello world   \n"
	trimmed, _ := transform.Trim(input6)
	fmt.Printf("   Original: '%s'\n", input6)
	fmt.Printf("   Trimmed: '%s'\n", trimmed)
	fmt.Println()

	fmt.Println("Demo complete!")
}
