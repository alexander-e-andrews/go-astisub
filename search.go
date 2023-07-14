package astisub

import (
	"fmt"
	"strings"
	"time"
)

type TimePoint struct {
	StartAt time.Duration
	EndAt   time.Duration
}

// So the captions appear twice, so we need to remove all those duplicates
func CleanOverlappingTimePoints(){

}

func (s *Subtitles) SimpleSearchSubtitles(searchTerm string) (finalSpeakers map[string][]TimeStamps) {
	// Given a search team, going to try to find the times where those words are

	// Should be all from one voice

	return SearchThrough(s.Items, searchTerm)
}

func SearchThrough(items []*Item, searchTerm string) (finalSpeakers map[string][]TimeStamps) {
	searchTerm = strings.ToLower(searchTerm)
	finalSpeakers = make(map[string][]TimeStamps)
	sim := make(SpeakerItemMap)
	// I am going to assume the items come in chronological order
	// Build the full set of all the lines, this definitely would be memory intensive on a whole movie
	for _, it := range items {
		speakerLines := it.extractLines()
		for speaker, lin := range speakerLines {
			lines, ok := sim[speaker]
			sl := SpeakerLine{
				Line:    strings.ToLower(lin),
				StartAt: it.StartAt,
				EndAt:   it.EndAt,
			}
			if !ok {
				sim[speaker] = []SpeakerLine{sl}
			} else {
				sim[speaker] = append(lines, sl)
			}
		}
	}

	for speaker, ummkay := range sim {
		speakerPossibleTimes := searchSpeakerLines(searchTerm, ummkay)
		finalSpeakers[speaker] = speakerPossibleTimes
	}
	return
}

type SpeakerLine struct {
	Line    string
	StartAt time.Duration
	EndAt   time.Duration
}

type TimeStamps struct {
	StartAt time.Duration
	EndAt   time.Duration
}

// The [speaker] line
type SpeakerLineMap map[string]string

type SpeakerItemMap map[string][]SpeakerLine

func searchSpeakerLines(searchLine string, lines []SpeakerLine) (matchTimes []TimeStamps) {
	searchLineRune := []rune(searchLine)
	searchLineRuneIndex := 0

	// so if there is a startAt then we have begun searching
	// if there is an endAt, we have found a match
	var startAt *time.Duration
	var endAt *time.Duration

	for _, line := range lines {
		//lineRuneIndex := 0
		lineRune := []rune(line.Line)
		// We exit the loop if we hit the end of the line, or we have gotten through our search term
		// Cant exit early in case the term is reapeated in this same line
		for x := 0; x < len(lineRune); x++ {

			if lineRune[x] == searchLineRune[searchLineRuneIndex] {
				if startAt == nil {
					startAt = &line.StartAt
				}
				searchLineRuneIndex += 1
			} else {
				x += searchLineRuneIndex
				searchLineRuneIndex = 0
				startAt = nil
			}

			if searchLineRuneIndex == len(searchLineRune) {
				endAt = &line.EndAt
				ts := TimeStamps{}
				if startAt != nil {
					ts = TimeStamps{
						StartAt: *startAt,
						EndAt:   *endAt,
					}
				} else {
					fmt.Printf("startAt null for line: %s\n", line.Line)
				}

				searchLineRuneIndex = 0
				matchTimes = append(matchTimes, ts)
				startAt = nil
				endAt = nil
			}
		}

	}

	return
}

// Set to true to have every line be spoken by same person
var noSpeaker = true

// Given an Item, we will pull out all its text, should probably do it per speaker
func (it *Item) extractLines() (speakerLines SpeakerLineMap) {
	speakerLines = make(SpeakerLineMap)
	var speaker string
	for _, ln := range it.Lines {
		if !noSpeaker {
			speaker = ln.VoiceName
		}

		for _, ln2 := range ln.Items {
			line := ln2.Text
			fullLine, ok := speakerLines[speaker]
			if !ok {
				speakerLines[speaker] = line
			} else {
				speakerLines[speaker] = fullLine + " " + line
			}
		}
	}

	return
}
