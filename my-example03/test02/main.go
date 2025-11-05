package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 簡單的檔案讀取和替換
func copyFileWithReplace(inputFile, outputFile string, keywords []string, replacement []string) error {
	// 開啟輸入檔案
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	// 建立輸出檔案
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)

	lineCount := 0
	replaceCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		var checkReplace bool = true
		for _, keyword := range keywords {
			// 如果包含關鍵字，就替換
			if strings.Contains(line, keyword) {
				//line = replacement
				for _, rep := range replacement {
					// 寫入行
					newline := fmt.Sprintf(`                	<li>%s</l>`, rep)
					_, err := writer.WriteString(newline + "\n")
					if err != nil {
						return err
					}
					replaceCount++
					fmt.Printf("第 %d 行已替換\n", lineCount)
				}
				checkReplace = false
				continue //跳過
			}
		}
		if checkReplace {
			// 寫入行
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}

	}

	writer.Flush()

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("完成！處理 %d 行，替換 %d 處\n", lineCount, replaceCount)
	return nil
}

// 多關鍵字替換版本
func copyFileWithMultipleKeywords(inputFile, outputFile string, keywords []string, replacement string) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)

	lineCount := 0
	replaceCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		originalLine := line

		// 檢查所有關鍵字
		for _, keyword := range keywords {
			if strings.Contains(line, keyword) {
				line = replacement
				replaceCount++
				fmt.Printf("第 %d 行找到 '%s': %s\n", lineCount, keyword, strings.TrimSpace(originalLine))
				break
			}
		}

		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	writer.Flush()

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("多關鍵字處理完成！總共 %d 行，替換 %d 處\n", lineCount, replaceCount)
	return nil
}

func main() {
	// 使用範例 1：單一關鍵字
	fmt.Println("=== 單一關鍵字處理 ===")

	keywords := []string{"<!--####0###-->", "<!--####1###-->"}
	wkeywords := []string{"ai1", "ai2", "ai3", "ai4"}
	err := copyFileWithReplace("example.txt", "result.html", keywords, wkeywords)
	if err != nil {
		fmt.Printf("錯誤: %v\n", err)
	}

	//使用範例 2：多關鍵字
	// fmt.Println("\n=== 多關鍵字處理 ===")
	// keywords := []string{"錯誤", "bug", "error", "fixme"}
	// err = copyFileWithMultipleKeywords("source.txt", "result2.txt", keywords, "**[已修正]**")
	// if err != nil {
	// 	fmt.Printf("錯誤: %v\n", err)
	// }

}
