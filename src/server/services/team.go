package services

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// 팀 상태 저장소 (전역 변수로 관리)
var TeamSessions = make(map[string][]string) // messageID -> participants

// TeamService 구조체
type TeamService struct{}

// 새로운 TeamService 인스턴스 생성
func NewTeamService() *TeamService {
	return &TeamService{}
}

// 랜덤 팀 구성 (핵심 비즈니스 로직)
func (ts *TeamService) CreateRandomTeams(participants []string) ([]string, []string) {
	if len(participants) != 6 {
		return nil, nil
	}

	rand.Seed(time.Now().UnixNano())

	shuffled := make([]string, len(participants))
	copy(shuffled, participants)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:3], shuffled[3:]
}

// 참가자 추가 (중복 체크 및 인원 제한 로직)
func (ts *TeamService) AddParticipant(messageID, participant string) bool {
	participants := TeamSessions[messageID]

	// 중복 확인
	for _, p := range participants {
		if p == participant {
			return false // 이미 존재
		}
	}

	// 인원 수 확인
	if len(participants) >= 6 {
		return false // 인원 초과
	}

	TeamSessions[messageID] = append(participants, participant)
	return true
}

// 여러 참가자 추가 (콤마 분리 파싱 로직)
func (ts *TeamService) AddMultipleParticipants(messageID string, input string, separator string) int {
	participants := TeamSessions[messageID]
	nicknames := strings.Split(input, separator)
	addedCount := 0

	for _, nickname := range nicknames {
		nickname = strings.TrimSpace(nickname)
		if nickname != "" && len(participants) < 6 {
			// 중복 확인
			exists := false
			for _, p := range participants {
				if p == nickname {
					exists = true
					break
				}
			}
			if !exists {
				participants = append(participants, nickname)
				addedCount++
			}
		}
	}

	TeamSessions[messageID] = participants
	return addedCount
}

func (ts *TeamService) CreateTeamResultEmbed(team1, team2 []string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "🏙️ City-Bot 팀 구성 결과 🏙️",
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "🔵 A팀", Value: strings.Join(team1, "\n"), Inline: true},
			{Name: "🔴 B팀", Value: strings.Join(team2, "\n"), Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "3 vs 3 팀이 구성되었습니다! 아래 버튼으로 다시 구성할 수 있습니다.",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
