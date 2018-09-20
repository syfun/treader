package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/syfun/treader/pkg/client"
)

var book = flag.String("book", "", "书籍名称")

func readFromStdin(prompt string) string {
	fmt.Println(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimRight(text, "\n")
}

func main() {
	flag.Parse()
	rst, _, err := client.Search(context.Background(), *book, 0, 5)
	if err != nil {
		fmt.Println(err)
		return
	}
	if rst.Total == 0 {
		fmt.Printf("没有搜索到 '%v' 相关记录", *book)
		return
	}
	s := fmt.Sprintf("\n\n'%v' 搜索结果如下:\n\n总数: %v\n\n", *book, rst.Total)

	for i, bk := range rst.Books {
		s += fmt.Sprintf("编号: %v\n书名: %v\n作者: %v\n简介: %v\n\n", i, bk.Title, bk.Author, bk.ShortIntro)
	}

	fmt.Println(s)

	text := readFromStdin("请选择编号: ")
	i, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println(err)
		return
	}
	bk := rst.Books[i]
	tocs, _, err := client.ListTocs(context.Background(), bk.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(tocs) == 0 {
		fmt.Println("未找到书源")
	}
	s = "书源: \n"
	for i, toc := range tocs {
		s += fmt.Sprintf("编号: %v\n名称: %v\n最新章节: %v\n\n", i, toc.Name, toc.LastChapter)
	}
	fmt.Println(s)

	text = readFromStdin("请选择编号: ")
	i, _ = strconv.Atoi(text)
	toc := tocs[i]
	t, _, err := client.ListChapters(context.Background(), toc.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(t.Chapters) == 0 {
		fmt.Println("没有章节")
		return
	}

	f, _ := os.Create(fmt.Sprintf("./store/%s.txt", bk.Title))
	defer f.Close()
	for _, chapter := range t.Chapters {
		rst, _, err := client.GetChapter(context.Background(), chapter.Link)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !rst.Success {
			return
		}
		fmt.Println(chapter.Title)
		fmt.Fprintf(f, "%v\n\n%v\n\n", chapter.Title, rst.Chapter.Body)
	}
}
