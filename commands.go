package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
)

// Globas Flags
var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "ENV_CONFIRM",
		Name:   "confirm",

		Usage: "",
	},
	cli.StringFlag{
		EnvVar: "ENV_F",
		Name:   "F",
		Usage:  "",
	},
}

//User : Block user
type User struct {
	id         int64
	screenName string
}

func doBlock(c *cli.Context) {
	api := doOauth()
	defer api.Close()

	var successBlock []anaconda.User
	var failedBlock []User

	users, err := readUsers(c)

	if err != nil {
		log.Fatal(err)
	}

	statusChan := make(chan string)
	for _, user := range users {
		go func(user User) {
			twitterUser, err := blockUser(user, api)
			statusChan <- user.screenName
			if err != nil {
				// ブロックに失敗した時
				failedBlock = append(failedBlock, user)
			} else {
				// ブロックできた時
				successBlock = append(successBlock, twitterUser)
			}
		}(user)
	}

	for i := 0; i < len(users); i++ {
		fmt.Println("Blocking...", <-statusChan)
	}

	showBlockedList(successBlock, failedBlock)
}

func readUsers(c *cli.Context) (users []User, err error) {
	// 引数がひとつもなければ標準入力を読み込む
	if len(c.Args()) < 1 {
		scanner := bufio.NewScanner(os.Stdin)

		var input string

		for scanner.Scan() {
			input += scanner.Text() + " "
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		delimiter := c.String("F")

		// 区切り文字が指定されていなかったら ',' にする
		if delimiter == "" {
			delimiter = " "
		}

		// 空白を指定された区切り文字に置換する
		input = strings.Replace(input, " ", delimiter, -1)

		inputSlice := strings.Split(input, delimiter)

		for _, input := range inputSlice {

			var u User

			screenName := input
			var userID int64
			if n, err := strconv.ParseInt(input, 10, 64); err == nil {
				userID = n
			}

			u.screenName = screenName
			u.id = userID

			users = append(users, u)

		}

	} else {
		for i := 0; i < len(c.Args()); i++ {
			var u User
			str := c.Args()[i]
			// 余計な空白を取り除く
			str = strings.TrimSpace(str)
			// @が先頭についていたら取り除く
			str = strings.TrimLeft(str, "@")
			userID, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				userID = 0
			}

			// 両方にセットする
			u.screenName = str
			u.id = userID

			users = append(users, u)
		}
	}

	fmt.Println("Loaded Block Lists...")

	return
}

func blockUser(user User, api *anaconda.TwitterApi) (twitterUser anaconda.User, err error) {

	// userIDが0かuserIDの長さが10でないと気はscreenNameでブロックする
	if user.id == 0 || len(user.screenName) != 10 {
		// screen_name の時
		if twitterUser, err = api.BlockUser(user.screenName, nil); err == nil {
			return twitterUser, nil
		}
	} else {

		// Id の時
		if _, err := api.BlockUserId(user.id, nil); err == nil {
			return twitterUser, nil
		}
	}

	return
}

func showBlockedList(successBlock []anaconda.User, failedBlock []User) (err error) {

	red := ansi.ColorCode("red")
	blue := ansi.ColorCode("blue")
	reset := ansi.ColorCode("reset")

	if len(successBlock) > 0 {
		fmt.Println(red, "Blocked", reset)
		for idx, user := range successBlock {
			fmt.Println(idx+1, ":", user.Name, "(@", user.ScreenName, ")")
		}
	}

	if len(failedBlock) > 0 {
		fmt.Println(blue, "Failed Block", reset)
		for idx, user := range failedBlock {
			if user.id == 0 {
				fmt.Println(idx+1, ":", user.screenName)
			} else {
				fmt.Println(idx+1, ":", user.id)
			}
		}
	}

	return
}
