package main

import (
	_ "image/png"
	"time"
	"github.com/gen2brain/raylib-go/raylib"
	"regexp"
)

type MessageItem_t struct {
	Ownership bool
	Message string
	Position_x int
	Position_y int
	ArrivalTime time.Time
}

var (
	max_messages_nb = 15
	message_duration = 10 * time.Second
	message_font_size = 15
	message_max_len = 20
	message_color = rl.Red
	message_color_ours = rl.Blue
	message_padding = 5
	)




/*                +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+
                  | Messages auxiliary functions |
                  +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+ */
func SplitMessage(m string, length int) []string {
	split := make([]string, 0)
	if len(m) < 2 {
		return split
	}
	m = m[1:]
	for {
		if len(m) <= length {
			return append(split, m)
		} else {
			split = append(split, m[:length])
			m = m[length:]
		}
	}
}

func FilterMessages(messages []MessageItem_t) []MessageItem_t {
	res := make([]MessageItem_t, 0)
	for i := 0 ; i < len(messages) ; i ++ {
		if (time.Since(messages[i].ArrivalTime) < message_duration) {
			res = append(res, messages[i])
		}
	}
	return res
}

func organizeMessages(messages []MessageItem_t) {
	// The messages are printed on the bottom left corner of the screen
	for i := 0 ; i < len(messages) ; i ++ {
		messages[i].Position_x = 0
		messages[i].Position_y = rl.GetScreenHeight() - (i+2) * (message_font_size + message_padding)
	}
}

func NewMessageItem(current []MessageItem_t, new_message string, config Configuration_t) ([]MessageItem_t) {
	//ownership := has_send > 0
	regex := regexp.MustCompile(`(?P<pseudo>.*):.*`)
	pseudo := regex.FindStringSubmatch(new_message)[1]
	ownership := pseudo == config.Pseudo
	if len(current) == max_messages_nb {
		// On enlève le premier!
		return append(current[1:], (MessageItem_t {ownership, new_message, 0, 0, time.Now()}))
	}
	return append(current, (MessageItem_t {ownership, new_message, 0, 0, time.Now()}))
}

func isPrintable(key int32) (bool) {
	return key >= 32 && key <= 126
}

