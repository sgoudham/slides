package navigation

import (
	"strconv"
)

type repeatableFunc func(slide, totalSlides int) int

// State tracks the current buffer, page, and total number of slides
type State struct {
	Buffer        string
	Page          int
	Section       int
	TotalSlides   int
	TotalSections int
}

// Navigate receives the current State and keyPress, and returns the new State.
func Navigate(state State, keyPress string) State {
	switch keyPress {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		newBuffer := keyPress

		if bufferIsNumeric(state.Buffer) {
			newBuffer = state.Buffer + keyPress
		}

		return State{
			Buffer:        newBuffer,
			Page:          state.Page,
			Section:       state.Section,
			TotalSlides:   state.TotalSlides,
			TotalSections: state.TotalSections,
		}
	case "g":
		switch state.Buffer {
		case "g":
			return State{
				Page:          0,
				Section:       0,
				TotalSlides:   state.TotalSlides,
				TotalSections: state.TotalSections,
			}
		default:
			return State{
				Buffer:        "g",
				Page:          state.Page,
				Section:       state.Section,
				TotalSlides:   state.TotalSlides,
				TotalSections: state.TotalSections,
			}
		}
	case "G":
		targetSlide := state.TotalSlides - 1
		if bufferIsNumeric(state.Buffer) {
			targetSlide = navigateSlide(state.Buffer, state.TotalSlides)
		}

		return State{
			Page:        targetSlide,
			TotalSlides: state.TotalSlides,
		}
	case " ", "j", "right", "l", "enter", "pgdown":
		return State{
			Page:        navigateNext(state, state.Page),
			TotalSlides: state.TotalSlides,
		}
	case "k", "left", "h", "pgup", "N":
		return State{
			Page:        navigatePrevious(state, state.Page),
			TotalSlides: state.TotalSlides,
		}
	case "up", "p":
		return State{
			Page:    state.Page,
			Section: navigatePrevious(state, state.Section),
		}
	case "down", "n":
		return State{
			Page:    state.Page,
			Section: navigateNext(state, state.Section),
		}
	default:
		return State{
			Page:        state.Page,
			TotalSlides: state.TotalSlides,
		}
	}
}

func bufferIsNumeric(buffer string) bool {
	_, err := strconv.Atoi(buffer)
	return err == nil
}

func navigateNext(state State, pageOrSection int) int {
	return repeatableAction(func(slide, totalSlides int) int {
		if slide < totalSlides-1 {
			return slide + 1
		}

		return totalSlides - 1
	}, state, pageOrSection)
}

func navigateSlide(buffer string, totalSlides int) int {
	destinationSlide, _ := strconv.Atoi(buffer)
	destinationSlide--

	if destinationSlide > totalSlides-1 {
		return totalSlides - 1
	}

	if destinationSlide < 0 {
		return 0
	}

	return destinationSlide
}

func navigatePrevious(state State, pageOrSection int) int {
	return repeatableAction(func(slide, totalSlides int) int {
		if slide > 0 {
			return slide - 1
		}

		return slide
	}, state, pageOrSection)
}

func repeatableAction(fn repeatableFunc, state State, pageOrSection int) int {
	if !bufferIsNumeric(state.Buffer) {
		return fn(state.Page, state.TotalSlides)
	}

	repeat, _ := strconv.Atoi(state.Buffer)
	page := state.Page

	if repeat == 0 {
		// This is how behaviour works in Vim, so following principle of least astonishment.
		return fn(pageOrSection, state.TotalSlides)
	}

	for i := 0; i < repeat; i++ {
		page = fn(page, state.TotalSlides)
	}

	return page
}
