package commands

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func setShiftCodesCommand() {
	// example ping command
	command := Command{Name: "show-me-the-shift-codes", Description: "Shows you the latest shift code"}

	command.SetCommandAction(func(s *discordgo.Session, m *discordgo.MessageCreate, _ argsInterface) error {
		resp, err := soup.Get("http://orcz.com/Borderlands_2:_Golden_Key")
		if err != nil {
			panic(err)
		}
		doc := soup.HTMLParse(resp)

		rows := doc.Find("table").FindAll("tr")

		shiftCode := ShiftCode{Reward: "", Date: "", Status: "", PCCode: "", PSCode: "", XBCode: ""}
		// rows[1:] skips the header
		for _, row := range rows[1:] {
			columns := row.FindAll("td")

			status := parseStatus(columns[3])

			if status != "Works" {
				continue
			}

			// source := parseSource(columns[0])
			reward := parseReward(columns[1])
			date := parseDate(columns[2])
			pcCode := parseCode(columns[4])
			psCode := parseCode(columns[5])
			xbCode := parseCode(columns[6])

			shiftCode = ShiftCode{Reward: reward, Date: date, Status: status, PCCode: pcCode, PSCode: psCode, XBCode: xbCode}
		}

		s.ChannelMessageSend(m.ChannelID, shiftCode.ToString())
		return nil
	}, nil)
	Commands[command.Name] = &command
}

// ShiftCode ...
type ShiftCode struct {
	Reward string
	Date   string
	Status string
	PCCode string
	PSCode string
	XBCode string
}

// ToString converts object to a string
func (c *ShiftCode) ToString() string {
	return fmt.Sprintf(
		"Reward: %s\nDate: %s\nStatus: %s\nPC: %s\nPS: %s\nXB: %s",
		c.Reward, c.Date, c.Status, c.PCCode, c.PSCode, c.XBCode,
	)
}

func parseReward(column soup.Root) string {
	return strings.TrimSpace(column.Text())
}

func parseDate(column soup.Root) string {
	return strings.TrimSpace(column.Text())
}

func parseStatus(column soup.Root) string {
	if strings.Contains(column.Text(), "Works") {
		return "Works"
	}
	return "Expired"
}

func parseCode(column soup.Root) string {
	return strings.TrimSpace(column.FullText())
}
